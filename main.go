package main

import (
	"fmt"
	"github.com/Albert221/ReddigramApi/server"
	"log"
	"os"
)

func main() {
	srv := server.NewServer(server.Options{
		Port:   os.Getenv("PORT"),
		Secret: os.Getenv("SECRET"),

		DbHost:     os.Getenv("DBHOST"),
		DbUser:     os.Getenv("DBUSER"),
		DbPassword: os.Getenv("DBPASSWORD"),
		DbName:     os.Getenv("DBNAME"),
	})

	fmt.Println("Server started listening...")
	log.Fatal(srv.Listen())
}
