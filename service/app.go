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
	"strconv"
	"sync"
)

type Member struct {
	ID       int      `json:"id"`
	User     string   `json:"user"`
	Nick     string   `json:"nick"`
	Avatar   string   `json:"avatar"`
	JoinedAt string   `json:"joined_at"`
	Status   string   `json:"status"`
	Roles    []string `json:"roles"`
	MsgCount int      `json:"msg_count"`
}

var members []Member
var nextID = 1
var mu sync.Mutex

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	certFile := os.Getenv("CERTFILE_PATH")
	keyFile := os.Getenv("KEYFILE_PATH")
	http_port := os.Getenv("HTTP_PORT")
	https_port := os.Getenv("HTTPS_PORT")
	ip := os.Getenv("IP")

	fmt.Printf("ip: %v, certFile_path: %v\nkeyFile_path: %v\nhttp_port: %v\nhttps_port: %v\n", ip, certFile, keyFile, http_port, https_port)

	if ip == "" || (https_port == "" && http_port == "") || certFile == "" || keyFile == "" {
		log.Fatalf("Required environment variables are missing")
	}

	r := mux.NewRouter()

	// Root handler for health check
	r.HandleFunc("/", healthCheckHandler)

	// Subrouter for /guild
	guildRouter := r.PathPrefix("/guild").Subrouter()
	guildRouter.HandleFunc("/members", getMembersHandler).Methods("GET")
	guildRouter.HandleFunc("/members", addMemberHandler).Methods("POST")
	guildRouter.HandleFunc("/members/{id:[0-9]+}", getMemberByIDHandler).Methods("GET")

	// Enable CORS for all routes
	corsObj := handlers.AllowedOrigins([]string{"*"}) // Allow all origins
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	// Start the server
	log.Println("Server listening on port " + http_port)
	log.Fatal(http.ListenAndServe(":"+http_port, handlers.CORS(corsObj, headersOk, methodsOk)(r)))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getMembersHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func addMemberHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("A post request has arrived!")

	var newMembers []Member
	err := json.NewDecoder(r.Body).Decode(&newMembers)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign an ID to the new member
	// Loop through the new members, assign IDs, and add them to the in-memory store
	for i := range newMembers {
		newMembers[i].ID = nextID
		nextID++
		members = append(members, newMembers[i])

		fmt.Printf("The member is now :%+v\n", newMembers[i])

	}

	// Return the newly added members with the generated IDs
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(newMembers) // Respond with the updated list of members
}

func getMemberByIDHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// Find the member with the specified ID
	for _, member := range members {
		if member.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(member)
			return
		}
	}

	http.Error(w, "Member not found", http.StatusNotFound)
}
