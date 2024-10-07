package server

import (
	"encoding/json"
	"log"
	"net/http"

	"kyri56xcaesar/discord_bots_app/internal/database"
	"kyri56xcaesar/discord_bots_app/internal/models"

	"github.com/gorilla/mux"
)

const DBName string = "dads.db"

// var members []models.Member // Not needed for now.. perhaps configure some memory operation-caching work
// var bots []models.Bot // Not needed for now.. perhaps configure some memory operation-caching work

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	RespondWithJSON(w, http.StatusOK, "OK")
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

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
			return
		}

		RespondWithJSON(w, http.StatusCreated, newMembers)

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed mate")
	}

}

func MemberHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	dbh := database.GetConnector(DBName)

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

	} else if r.Method == "PUT" {

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
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

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
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	w.Write([]byte("{'Guild'}"))
}

func RootBotHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)
	dbh := database.GetConnector(DBName)

	if r.Method == "GET" {
		log.Print("Redirecting to /guild/bots/")
		http.Redirect(w, r, "/guild/bots", http.StatusMovedPermanently)
	} else if r.Method == "POST" {

		var newBot models.Bot
		err := json.NewDecoder(r.Body).Decode(&newBot)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, err = dbh.InsertBot(&newBot)
		if err != nil {
			log.Printf("Error inserting multiple lines in the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			RespondWithJSON(w, http.StatusCreated, newBot)
		}
	}
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	vars := mux.Vars(r)
	identifier, idExists := vars["identifier"]
	log.Printf("Identifier: %s", identifier)

	if r.Method == "GET" {
		if idExists {
			// Get specific identifier
			res, err := dbh.GetBotByIdentifier(identifier)
			if err != nil {
				log.Printf("Error getting bot by identifier %v...: %v", identifier, err.Error())
				RespondWithError(w, http.StatusBadRequest, "Could not reach bot")
				return
			}

			RespondWithJSON(w, http.StatusOK, res)
		} else {
			// Get everything
			res, err := dbh.GetAllBots()
			if err != nil {
				log.Printf("Error fetching bots from the DB... " + err.Error())
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			RespondWithJSON(w, http.StatusOK, res)
		}

	} else if r.Method == "PUT" {
		if idExists {

			var updatedBot models.Bot
			err := json.NewDecoder(r.Body).Decode(&updatedBot)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				RespondWithError(w, http.StatusBadRequest, "BadR")
				return
			}

			_, err = dbh.UpdateBotByIdentifier(&updatedBot, identifier)
			if err != nil {
				log.Printf("Error updating bot %v... : %v", identifier, err.Error())
				RespondWithError(w, http.StatusBadRequest, "Could not update bot")
				return
			}
			// log.Printf("\nThe result of this thing is: %+v\n", res)
			RespondWithJSON(w, http.StatusOK, updatedBot)

		} else {
			RespondWithError(w, http.StatusBadRequest, "must provide identifier")
		}
	} else if r.Method == "DELETE" {
		if idExists {
			res, err := dbh.DeleteBotByIdentifier(identifier)
			if err != nil {
				log.Printf("Error deleting bot :%v ...: %v", identifier, err.Error())
				RespondWithError(w, http.StatusBadRequest, "Could not delete the bot")
				return
			}
			RespondWithJSON(w, http.StatusCreated, res)

		} else {
			RespondWithError(w, http.StatusBadRequest, "must provide identifier")
		}

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "mot allowed mate")
	}
}

func BotsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

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

func RootBotLineHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	if r.Method == "GET" {
		// Get everything
		res, err := dbh.GetBotLines()
		if err != nil {
			log.Printf("Error fetching lines from the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, res)
	} else if r.Method == "POST" {
		var newLines []models.Line
		err := json.NewDecoder(r.Body).Decode(&newLines)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = dbh.InsertMultipleLines(newLines)
		if err != nil {
			log.Printf("Error inserting multiple lines in the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			RespondWithJSON(w, http.StatusCreated, newLines)
		}
	}

}

func BotLineHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	vars := mux.Vars(r)
	identifier := vars["identifier"]
	log.Printf("Identifier: %s", identifier)

	if r.Method == "GET" {
		// Get specific identifier
		res, err := dbh.GetBotLineByIdentifier(identifier)
		if err != nil {
			log.Printf("Error getting line by identifier %v...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not reach line")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)

	} else if r.Method == "PUT" {

		var updatedLine models.Line
		err := json.NewDecoder(r.Body).Decode(&updatedLine)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "BadR")
			return
		}

		_, err = dbh.UpdateLineByIndentifier(updatedLine, identifier)
		if err != nil {
			log.Printf("Error updating line %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not update line")
			return
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, updatedLine)

		RespondWithError(w, http.StatusBadRequest, "must provide identifier")

	} else if r.Method == "DELETE" {

		res, err := dbh.DeleteLineByIndentifier(identifier)
		if err != nil {
			log.Printf("Error deleting line :%v ...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not delete the line")
			return
		}
		RespondWithJSON(w, http.StatusCreated, res)

	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "mot allowed mate")
	}

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v (Health Check)", r.Method, r.URL.Path)

	// Open the database connection
	// Convert the schema info to JSON
	schema := database.DBHealthCheck(DBName)
	response, err := json.Marshal(schema)
	if err != nil {
		http.Error(w, "Failed to encode schema info as JSON", http.StatusInternalServerError)
		return
	}

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
