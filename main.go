package main

import (
	cors "db_practice/config"
	"db_practice/handlers"
	"db_practice/internal/db"
	"db_practice/internal/users"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting the application")
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DEV_CONN_STR")

	if connStr == "" {
		panic("Connection string not found in configuration")
	}

	udb, err := db.NewDB(connStr, false, "")
	if err != nil {
		log.Fatalf("FAILURE OPENING DATABASE CONNECTION: %v", err)
	}
	defer udb.Close()

	uClient := users.NewUsersClient(udb)

	// Setup the HTTP server and router
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the practice API")
	})

	uHandler := handlers.NewUsersHandler(uClient)
	router.HandleFunc("/users/create", uHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{email}", uHandler.GetUserByEmail).Methods("GET")

	handler := cors.SetCORS(router)

	port := 8000
	fmt.Printf("Server is running on :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Println("Application started successfully")
}
