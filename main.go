package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Albert221/ReddigramApi/handlers"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")

	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}

	contr := handlers.NewController(db, secret)

	r := mux.NewRouter()
	r.Use(createRateLimiterMiddleware())

	r.HandleFunc("/suggested-subreddits", contr.SuggestedSubredditsHandler).Methods("POST")

	r.HandleFunc("/authenticate", contr.AuthenticateHandler).Methods("POST")

	subs := r.PathPrefix("/subscriptions").Subrouter()
	subs.Use(contr.AuthMiddleware)
	subs.HandleFunc("", contr.ListSubsHandler).Methods("GET")
	subs.HandleFunc("/{id}", contr.AddSubHandler).Methods("PUT")
	subs.HandleFunc("/{id}", contr.RemoveSubHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func createDB() (*sqlx.DB, error) {
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbName := os.Getenv("DBNAME")

	dsn := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbHost,
		DBName:               dbName,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
		ParseTime:            true,
	}

	return sqlx.Open("mysql", dsn.FormatDSN())
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
