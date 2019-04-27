package main

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
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
	}

	db, err := sqlx.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	contr := NewController(db)

	r := mux.NewRouter()
	r.HandleFunc("/authenticate", contr.AuthenticateHandler).Methods("POST")

	subs := r.PathPrefix("/subscriptions").Subrouter()
	subs.Use(contr.AuthMiddleware)
	subs.HandleFunc("", contr.ListSubsHandler).Methods("GET")
	subs.HandleFunc("/{name}", contr.AddSubHandler).Methods("PUT")
	subs.HandleFunc("/{name}", contr.RemoveSubHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
