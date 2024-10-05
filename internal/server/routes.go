package server

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"kyri56xcaesar/discord_bots_app/internal/database"
	"kyri56xcaesar/discord_bots_app/internal/models"

	"github.com/gorilla/mux"
)

const DBName string = "dads.db"

// var members []models.Member // Not needed for now.. perhaps configure some memory operation-caching work
// var bots []models.Bot // Not needed for now.. perhaps configure some memory operation-caching work

func RootHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "OK")
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {

	dbh := database.GetConnection(DBName)

	if r.Method == "GET" {

		res, err := dbh.GetAllMembers()
		if err != nil {
			log.Printf("Error getting fetching members from the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, res)

	} else if r.Method == "POST" {

		var newMembers []models.Member
		err := json.NewDecoder(r.Body).Decode(&newMembers)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = dbh.InsertMultipleMembers(newMembers)
		if err != nil {
			log.Printf("Error inserting multiple members in the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			RespondWithJSON(w, http.StatusCreated, newMembers)
		}

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed mate")
	}

}

func MemberHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	dbh := database.GetConnection(DBName)

	if !IsAlphanumeric(identifier) {
		RespondWithError(w, http.StatusMethodNotAllowed, "Nah")
		return
	}

	if r.Method == "GET" {
		// Call ID

		res, err := dbh.GetMemberByIdentifier(identifier)
		if err != nil {
			log.Printf("Error getting member by identifier %v...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not reach member")
			return
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, res)

	} else if r.Method == "UPDATE" {

		var updatedMember models.Member
		err := json.NewDecoder(r.Body).Decode(&updatedMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "BadR")
			return
		}

		_, err = dbh.UpdateMemberByIdentifier(updatedMember, identifier)
		if err != nil {
			log.Printf("Error updating member %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not update member")
			return
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, updatedMember)

	} else if r.Method == "DELETE" {

		res, err := dbh.DeleteMemberByIdentifier(identifier)
		if err != nil {
			log.Printf("Error deleting member :%v ...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not delete the member")
			return
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusCreated, res)

	} else if r.Method == "POST" {

		var newMember models.Member
		err := json.NewDecoder(r.Body).Decode(&newMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Could not create member")
			return
		}

		res, err := dbh.InsertMember(newMember)
		if err != nil {
			log.Printf("Error inserting member %v...: %v", newMember, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Did not insert the member")
			return
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusCreated, res)

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed mate")
	}

}

func RootMemberHandler(w http.ResponseWriter, r *http.Request) {

	dbh := database.GetConnection(DBName)

	if r.Method == "GET" {
		http.Redirect(w, r, "/guild/members", http.StatusMovedPermanently)

	} else if r.Method == "POST" {

		var newMember models.Member
		err := json.NewDecoder(r.Body).Decode(&newMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "BadR")
			return
		}

		res, err := dbh.InsertMember(newMember)
		if err != nil {
			log.Printf("Error inserting a single member: %v... %v", newMember, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not insert member")
			return
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusCreated, res)

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed mate")
	}

}

func GuildHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{'Guild'}"))
}

func RootBotHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{'Bots'}"))
}

func BotHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		http.Redirect(w, r, "/guild/bots", http.StatusMovedPermanently)
	} else if r.Method == "POST" {
		log.Printf("TEst")
	}
	w.Write([]byte("{'Bot'}"))
}

func BotsHandler(w http.ResponseWriter, r *http.Request) {

	dbh := database.GetConnection(DBName)

	if r.Method == "GET" {

		res, err := dbh.GetAllBots()
		if err != nil {
			log.Printf("Error getting fetching bots from the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, res)

	} else if r.Method == "POST" {

		var newBots []models.Bot
		err := json.NewDecoder(r.Body).Decode(&newBots)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = dbh.InsertMutipleBots(newBots)
		if err != nil {
			log.Printf("Error inserting multiple bots in the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			RespondWithJSON(w, http.StatusCreated, newBots)
		}

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed mate")
	}
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}
