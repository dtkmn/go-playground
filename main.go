package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Handler for fetching resources
func getResourceHandler(w http.ResponseWriter, r *http.Request) {
	resources, err := GetResources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resources)
}

// Handler for adding a new resource
func addResourceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add resource")
	var resource Resource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := AddResource(resource.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Resource added successfully")
	w.WriteHeader(http.StatusCreated)
}

// Handler for sending a message to Kafka
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sendMessage(message.Message)
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Check if DATABASE_URL is loaded
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	log.Println("DATABASE_URL:", connStr)

	// Initialize the router
	r := mux.NewRouter()

	// Define your routes
	r.HandleFunc("/api/resource", getResourceHandler).Methods("GET")
	r.HandleFunc("/api/resource", addResourceHandler).Methods("POST")
	r.HandleFunc("/api/send-message", sendMessageHandler).Methods("POST")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	if err := server.Shutdown(nil); err != nil {
		log.Fatal(err)
	}
}
