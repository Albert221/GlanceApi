package server

import (
	"github.com/Albert221/ReddigramApi/controller"
	"github.com/Albert221/ReddigramApi/mysql"
	"github.com/Albert221/ReddigramApi/reddit"
	"github.com/Albert221/ReddigramApi/repository"
	"github.com/gbrlsnchs/jwt/v3"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
	"time"
)

type Server struct {
	router  *mux.Router
	port    string
	db      *sqlx.DB
	jwtAlgo jwt.Algorithm

	redditRepo       repository.RedditRepository
	subscriptionRepo repository.SubscriptionRepository

	authContr         *controller.AuthController
	subscriptionContr *controller.SubscriptionController
	suggestionContr   *controller.SuggestionController
}

type Options struct {
	Port   string
	Secret string

	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
}

func NewServer(options Options) *Server {
	srv := &Server{
		router:  mux.NewRouter(),
		port:    options.Port,
		db:      setupDb(options),
		jwtAlgo: jwt.NewHS256([]byte(options.Secret)),
	}

	// Setup repositories
	srv.redditRepo = reddit.NewRepository()
	srv.subscriptionRepo = mysql.NewSubscriptionRepository(srv.db)

	// Setup controllers
	srv.authContr = controller.NewAuthController(srv.jwtAlgo, srv.redditRepo)
	srv.suggestionContr = controller.NewSuggestionController()
	srv.subscriptionContr = controller.NewSubscriptionController(srv.subscriptionRepo)

	srv.setupRoutes()

	return srv
}

func setupDb(options Options) *sqlx.DB {
	dsn := mysqlDriver.Config{
		User:   options.DbUser,
		Passwd: options.DbPassword,
		Net:    "tcp",
		Addr:   options.DbHost,
		DBName: options.DbName,

		Collation: "utf8mb4_unicode_ci",

		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sqlx.Open("mysql", dsn.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}

func createRateLimiterMiddleware() mux.MiddlewareFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	}

	rlim := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(rlim)

	return middleware.Handler
}

func (s *Server) setupRoutes() {
	r := s.router
	r.Use(createRateLimiterMiddleware(), s.authContr.AuthenticationMiddleware)
	// Auth
	r.HandleFunc("/authenticate", s.authContr.AuthenticateHandler).Methods("POST")
	// Suggestions
	r.HandleFunc("/suggested-subreddits", s.suggestionContr.SuggestedSubredditsHandler).Methods("POST")
	// Subscriptions
	r.HandleFunc("/subscriptions", s.subscriptionContr.ListHandler).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", s.subscriptionContr.SubscribeHandler).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", s.subscriptionContr.UnsubscribeHandler).Methods("DELETE")
}

func (s *Server) Listen() error {
	return http.ListenAndServe(":"+s.port, s.router)
}
