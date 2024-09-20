package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"log"
	"net/http"
	"sync"

	"kyri56xcaesar/discord_bots_app/guild/user"
)

var next_user_ID, next_bot_ID = 1, 1

var mu sync.Mutex
var members []user.Member
var bots []user.Bot

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	certFile := os.Getenv("CERTFILE")
	keyFile := os.Getenv("KEYFILE")
	http_port := os.Getenv("HTTP_PORT")
	https_port := os.Getenv("HTTPS_PORT")
	ip := os.Getenv("IP")

	fmt.Printf("ip: %v, certFile_path: %v\nkeyFile_path: %v\nhttp_port: %v\nhttps_port: %v\n", ip, certFile, keyFile, http_port, https_port)

	if ip == "" || (https_port == "" && http_port == "") || certFile == "" || keyFile == "" {
		log.Fatalf("Required environment variables are missing")
	}

	r := mux.NewRouter()

	// Root handler for health check
	r.HandleFunc("/", rootHandler)

	// Subrouter for /guild
	guildRouter := r.PathPrefix("/guild").Subrouter()
	guildRouter.HandleFunc("/{type}", userHandler).Methods("GET", "POST")

	// Enable CORS for all routes
	corsObj := handlers.AllowedOrigins([]string{"*"}) // Allow all origins
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	// Start the server
	log.Println("Server listening on port " + http_port)
	log.Fatal(http.ListenAndServe(":"+http_port, handlers.CORS(corsObj, headersOk, methodsOk)(r)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	resourceType := vars["type"]

	switch resourceType {
	case "members":
		if r.Method == "GET" {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(members)

		} else if r.Method == "POST" {

			var newMembers []user.Member
			err := json.NewDecoder(r.Body).Decode(&newMembers)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for i := range newMembers {
				newMembers[i].ID = next_user_ID
				next_user_ID++
				members = append(members, newMembers[i])
				// fmt.Printf("The member is now: %+v\n", newMembers[i])
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newMembers)
		}
	case "bots":
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bots)
		} else if r.Method == "POST" {

			var newBots []user.Bot
			err := json.NewDecoder(r.Body).Decode(&newBots)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for i := range newBots {
				newBots[i].ID = next_bot_ID
				next_bot_ID++
				newBots = append(newBots, newBots[i])
				// fmt.Printf("The member is now: %+v\n", newMembers[i])
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newBots)
		}
	}

}
