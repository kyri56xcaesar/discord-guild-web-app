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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	// RespondWithJSON(w, http.StatusOK, "Up")
	templatePath := "/cmd/api/web/templates/index.html"
	RespondWithTemplate(w, http.StatusOK, templatePath, nil)
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	if dbh == nil {
		RespondWithError(w, http.StatusInternalServerError, "Database connection failed")
		return
	}

	switch r.Method {
	case http.MethodGet:
		res, err := dbh.GetAllMembers()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve members")
			return
		}
		if res == nil {
			// If the result is nil, return 404
			RespondWithError(w, http.StatusNotFound, "Members not found")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)

	case http.MethodPost:
		var newMembers []models.Member
		if err := json.NewDecoder(r.Body).Decode(&newMembers); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if _, err := dbh.InsertMultipleMembers(newMembers); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to insert members")
			return
		}

		RespondWithJSON(w, http.StatusCreated, newMembers)

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func MemberHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	if !IsAlphanumeric(identifier) {
		RespondWithError(w, http.StatusBadRequest, "Invalid identifier")
		return
	}

	dbh := database.GetConnector(DBName)

	switch r.Method {
	case http.MethodGet:
		// Call ID
		res, err := dbh.GetMemberByIdentifier(identifier)
		if err != nil {
			// log.Printf("Error getting member by identifier %v...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusNotFound, "Member not found")
			return
		}

		if res == nil {
			// If the result is nil, return 404
			// log.Printf("Member with identifier %v not found", identifier)
			RespondWithError(w, http.StatusNotFound, "Member not found")
			return
		}

		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, res)

	case http.MethodPut:
		var updatedMember models.Member
		if err := json.NewDecoder(r.Body).Decode(&updatedMember); err != nil {
			// log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		if _, err := dbh.UpdateMemberByIdentifier(updatedMember, identifier); err != nil {
			// log.Printf("Error updating member %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusInternalServerError, "Failed to update member")
			return
		}
		// log.Printf("\nThe result of this thing is: %+v\n", res)
		RespondWithJSON(w, http.StatusOK, updatedMember)

	case http.MethodDelete:
		if _, err := dbh.DeleteMemberByIdentifier(identifier); err != nil {
			// log.Printf("Error deleting member :%v ...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Failed to delete member")
			return
		}
		RespondWithJSON(w, http.StatusCreated, "Deletion success")

	case http.MethodPost:
		var newMember models.Member
		if err := json.NewDecoder(r.Body).Decode(&newMember); err != nil {
			// log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		if _, err := dbh.InsertMember(newMember); err != nil {
			// log.Printf("Error inserting member %v...: %v", newMember, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Did not insert the member")
			return
		}

		RespondWithJSON(w, http.StatusCreated, newMember)

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func RootMemberHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	switch r.Method {
	case http.MethodGet:
		http.Redirect(w, r, "/guild/members", http.StatusMovedPermanently)

	case http.MethodPost:
		var newMember models.Member
		if err := json.NewDecoder(r.Body).Decode(&newMember); err != nil {
			// log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		if _, err := dbh.InsertMember(newMember); err != nil {
			// log.Printf("Error inserting a single member: %v... %v", newMember, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Failed to insert member")
			return
		}
		RespondWithJSON(w, http.StatusCreated, newMember)

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func GuildHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	RespondWithJSON(w, http.StatusFound, "{'guilds'}")
}

func RootBotHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)
	dbh := database.GetConnector(DBName)

	switch r.Method {
	case http.MethodGet:
		log.Print("Redirecting to /guild/bots/")
		http.Redirect(w, r, "/guild/bots", http.StatusMovedPermanently)
	case http.MethodPost:
		var newBot models.Bot
		err := json.NewDecoder(r.Body).Decode(&newBot)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		_, err = dbh.InsertBot(&newBot)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to insert multiple lines")
		} else {
			RespondWithJSON(w, http.StatusCreated, newBot)
		}
	}
}

func BotHandler(w http.ResponseWriter, r *http.Request) {

	dbh := database.GetConnector(DBName)

	vars := mux.Vars(r)
	identifier := vars["identifier"]
	log.Printf("%v request on path: %v with identifier %v", r.Method, r.URL.Path, identifier)

	switch r.Method {
	case http.MethodGet:
		// Get specific identifier
		res, err := dbh.GetBotByIdentifier(identifier)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to get bot")
			return
		}

		if res == nil {
			// If the result is nil, return 404
			RespondWithError(w, http.StatusNotFound, "Bot not found")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)

	case http.MethodPut:
		var updatedBot models.Bot
		if err := json.NewDecoder(r.Body).Decode(&updatedBot); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if _, err := dbh.UpdateBotByIdentifier(&updatedBot, identifier); err != nil {
			log.Printf("Error updating bot %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not update bot")
			return
		}
		RespondWithJSON(w, http.StatusOK, updatedBot)

	case http.MethodDelete:
		if _, err := dbh.DeleteBotByIdentifier(identifier); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to delete bot")
			return
		}
		RespondWithJSON(w, http.StatusCreated, "Deletion success")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func BotsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	switch r.Method {
	case http.MethodGet:
		res, err := dbh.GetAllBots()
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to get bots")
			return
		}
		if res == nil {
			// If the result is nil, return 404
			RespondWithError(w, http.StatusNotFound, "Bots not found")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)

	case http.MethodPost:
		var newBots []models.Bot
		if err := json.NewDecoder(r.Body).Decode(&newBots); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON formant")
			return
		}
		if _, err := dbh.InsertMultipleBots(newBots); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to insert bots")
		} else {
			RespondWithJSON(w, http.StatusCreated, newBots)
		}

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func RootLineHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	switch r.Method {
	case http.MethodGet:
		// Get everything
		res, err := dbh.GetBotLines()
		if err != nil {
			log.Printf("Error fetching lines from the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, "Failed to get lines")
			return
		}

		if res == nil {
			// If the result is nil, return 404
			log.Printf("Lines not found")
			RespondWithError(w, http.StatusNotFound, "Lines not found")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)
	case http.MethodPost:
		var newLines []models.Line
		if err := json.NewDecoder(r.Body).Decode(&newLines); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if _, err := dbh.InsertMultipleLines(newLines); err != nil {
			log.Printf("Error inserting multiple lines in the DB... " + err.Error())
			RespondWithError(w, http.StatusBadRequest, "Failed to insert lines")
		} else {
			RespondWithJSON(w, http.StatusCreated, newLines)
		}
	}

}

func LineHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	log.Printf("Identifier: %s", identifier)

	switch r.Method {
	case http.MethodGet:
		// Get specific identifier
		res, err := dbh.GetBotLineByIdentifier(identifier)
		if err != nil {
			log.Printf("Error getting line by identifier %v...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not reach line")
			return
		}
		if res == nil {
			// If the result is nil, return 404
			log.Printf("Line with identifier %v not found", identifier)
			RespondWithError(w, http.StatusNotFound, "Line not found")
			return
		}
		RespondWithJSON(w, http.StatusOK, res)

	case http.MethodPut:
		var updatedLine models.Line
		if err := json.NewDecoder(r.Body).Decode(&updatedLine); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		if _, err := dbh.UpdateLineByIndentifier(updatedLine, identifier); err != nil {
			log.Printf("Error updating line %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not update line")
			return
		}

		RespondWithJSON(w, http.StatusOK, updatedLine)

	case http.MethodDelete:
		if _, err := dbh.DeleteLineByIndentifier(identifier); err != nil {
			log.Printf("Error deleting line :%v ...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not delete the line")
			return
		}
		RespondWithJSON(w, http.StatusCreated, "Delete success")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusNotFound, "Not Found")
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusMethodNotAllowed, "Not Allowed")
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