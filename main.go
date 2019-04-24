package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/authenticate", AuthenticateHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
