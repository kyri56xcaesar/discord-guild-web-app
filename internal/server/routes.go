package server

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"

	"kyri56xcaesar/discord_bots_app/internal/database"
	"kyri56xcaesar/discord_bots_app/internal/models"

	"github.com/gorilla/mux"
)

var mu sync.Mutex
var members []models.Member
var bots []models.Bot

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]byte("OK"))
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	resourceType := vars["type"]

	switch resourceType {
	case "members":
		if r.Method == "GET" {

			res, err := database.GetAllMembers()
			if err != nil {
				log.Printf("There's been an error brother...")
			}
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)

		} else if r.Method == "POST" {

			var newMembers []models.Member
			err := json.NewDecoder(r.Body).Decode(&newMembers)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			members = append(members, newMembers...)

			_, err = database.InsertMultipleMembers(newMembers)
			if err != nil {
				log.Printf("There's been an error brother...")
			}
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newMembers)

		}

	case "bots":
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(bots)
		} else if r.Method == "POST" {

			var newBots []models.Bot
			err := json.NewDecoder(r.Body).Decode(&newBots)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			newBots = append(newBots, newBots...)

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newBots)
		}
	}

}

func MemberHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	if !IsAlphanumeric(identifier) {
		log.Print("Invalid identifier input...")
		json.NewEncoder(w).Encode("Not Allowed.")

		return
	}

	if r.Method == "GET" {
		// Call ID

		res, err := database.GetMemberByIdentifier(identifier)
		if err != nil {
			log.Printf("There's been an error brother...")
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)

	} else if r.Method == "UPDATE" {

		var updatedMember models.Member
		err := json.NewDecoder(r.Body).Decode(&updatedMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = database.UpdateMemberByIdentifier(updatedMember, identifier)
		if err != nil {
			log.Printf("There's been an error brother...")
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedMember)

	} else if r.Method == "DELETE" {

		res, err := database.DeleteMemberByIdentifier(identifier)
		if err != nil {
			log.Printf("There's been an error brother...")
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)

	}

}

func RootMemberHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		res, err := database.GetAllMembers()
		if err != nil {
			log.Printf("There's been an error brother...")
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)

	} else if r.Method == "POST" {

		var newMember models.Member
		err := json.NewDecoder(r.Body).Decode(&newMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = database.InsertMember(newMember)
		if err != nil {
			log.Printf("There's been an error brother...")
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newMember)

	}
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}
