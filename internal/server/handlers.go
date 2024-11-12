package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"kyri56xcaesar/discord_bots_app/internal/database"
	"kyri56xcaesar/discord_bots_app/internal/models"

	"github.com/gorilla/mux"
)

const (
	typeMember string = "members"
	typeBot    string = "bots"
	typeLine   string = "lines"
)

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
		var err error
		var newMembers []models.Member

		// need to buffer it first
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to buffer body. %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to read body")
			return
		}

		err = json.Unmarshal(body, &newMembers)
		if err == nil { // Its an array given
			_, err := dbh.InsertMultipleMembers(newMembers)
			if err != nil {
				log.Printf("Error inserting members: %v", err)
				RespondWithError(w, http.StatusInternalServerError, "Failed to insert members")
				return
			}

			RespondWithJSON(w, http.StatusCreated, "Created multiple members")
			return
		}

		// Check if its a single member
		newMember := models.Member{}
		err = json.Unmarshal(body, &newMember)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		log.Print("Its a single member!")

		_, err = dbh.InsertMember(newMember)
		if err != nil {
			// log.Printf("Error inserting member %v...: %v", newMember, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Failed to insert member")
			return
		}

		RespondWithJSON(w, http.StatusCreated, "Created member.")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func GMultipleData(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	dataT := strings.SplitN(r.URL.String(), "/", 4)[2]

	log.Printf("DataT: %s", dataT)

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		RespondWithError(w, http.StatusBadRequest, "Must provide identifiers")
		return
	}
	identifiers := strings.Split(idsParam, ",")

	switch dataT {
	case typeMember:
		members, err := dbh.GetMultipleMembersByIdentifiers(identifiers)
		if err != nil {
			log.Printf("Error getting multiple %v. : %v", dataT, err)
			RespondWithError(w, http.StatusInternalServerError, "Error getting data")
			return
		}
		RespondWithJSON(w, http.StatusOK, members)
	case typeBot:
		bots, err := dbh.GetMultipleBotsByIdentifiers(identifiers)
		if err != nil {
			log.Printf("Error getting multiple %v. : %v", dataT, err)
			RespondWithError(w, http.StatusInternalServerError, "Error getting data")
			return
		}
		RespondWithJSON(w, http.StatusOK, bots)
	case typeLine:
		lines, err := dbh.GetMultipleLinesByIdentifiers(identifiers)
		if err != nil {
			log.Printf("Error getting multiple %v. : %v", dataT, err)
			RespondWithError(w, http.StatusInternalServerError, "Error getting data")
			return
		}
		RespondWithJSON(w, http.StatusOK, lines)
	default:
		// Creazy to be here
		RespondWithError(w, http.StatusInternalServerError, "Wierd error")
	}
}

func UDMultipleData(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	dataT := strings.SplitN(strings.SplitN(r.URL.String(), "/", 4)[3], "?", 2)[0]

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing 'ids' query parameters")
		return
	}
	identifiers := strings.Split(idsParam, ",")

	var (
		err     error
		message string
	)

	switch dataT {
	case typeMember:
		switch r.Method {
		case http.MethodDelete:
			message, err = dbh.DeleteMultipleMembersByIdentifiers(identifiers)
		case http.MethodPut:
		default:

		}
	case typeBot:
		switch r.Method {
		case http.MethodDelete:
			message, err = dbh.DeleteMultipleBotsByIdentifiers(identifiers)
		case http.MethodPut:
		default:

		}
	case typeLine:
		switch r.Method {
		case http.MethodDelete:
			message, err = dbh.DeleteMultipleLinesByIdentifiers(identifiers)
		case http.MethodPut:
		default:

		}
	default:
		// Impossible to reach here!
		log.Print("Shoulnd't be here!")
		RespondWithError(w, http.StatusInternalServerError, "Invalid")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, message)
		return
	}

	RespondWithJSON(w, http.StatusOK, message)
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
	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func DataIndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dataT := strings.SplitN(r.URL.String(), "/", 4)[3]

	switch dataT {
	case typeMember:

		RespondWithJSON(w, http.StatusOK, keysSliceFromMap(database.AllowedMemberCols))
	case typeBot:
		RespondWithJSON(w, http.StatusOK, keysSliceFromMap(database.AllowedBotCols))
	case typeLine:
		RespondWithJSON(w, http.StatusOK, keysSliceFromMap(database.AllowedLineCols))
	default:
		// Shouldnt reach here...
		RespondWithError(w, http.StatusInternalServerError, "Wierd Error.")
	}
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	dataT := strings.SplitN(r.URL.String(), "/", 5)[3]

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	if identifier == "" {
		log.Print("No identifier")
		RespondWithError(w, http.StatusBadRequest, "Must provide an Identifier")
		return
	}

	switch dataT {
	case typeMember:
		data, err := dbh.GetMemberIdentifiers(identifier)
		if err != nil {
			log.Printf("Failed to get member data by identifier %v", identifier)
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data")
			return
		}

		RespondWithJSON(w, http.StatusOK, data)

	case typeBot:
		data, err := dbh.GetBotIdentifiers(identifier)
		if err != nil {
			log.Printf("Failed to get bot data by identifier %v", identifier)
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data")
			return
		}

		RespondWithJSON(w, http.StatusOK, data)

	case typeLine:
		data, err := dbh.GetLineIdentifiers(identifier)
		if err != nil {
			log.Printf("Failed to get line data by identifier %v", identifier)
			RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data")
			return
		}

		RespondWithJSON(w, http.StatusOK, data)

	default:
		// impossible to reach here
		log.Print("You've made impossible not possible! Nice!")
		RespondWithError(w, http.StatusInternalServerError, "Nice!")
	}
}

func GuildHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	RespondWithJSON(w, http.StatusFound, "{'guilds'}")
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
		var err error

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to buffer body. %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to read body")
			return
		}

		var newBots []models.Bot
		// Test if multiple
		err = json.Unmarshal(body, &newBots)
		if err == nil {
			// Its multiple bots
			if _, err := dbh.InsertMultipleBots(newBots); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to insert bots")
				return
			}
			RespondWithJSON(w, http.StatusCreated, "Created bots")
			return
		}

		var newBot models.Bot
		err = json.Unmarshal(body, &newBot)
		if err != nil {
			log.Printf("Invalid JSON format. %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON Format")
			return
		}
		if _, err := dbh.InsertBot(&newBot); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to insert bot")
			return
		}

		RespondWithJSON(w, http.StatusCreated, "Created bot")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func LinesHandler(w http.ResponseWriter, r *http.Request) {
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
		var err error
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read body buffer. %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to read body")
			return
		}

		var newLines []models.Line

		err = json.Unmarshal(body, &newLines)
		if err == nil {
			// Multiple lines indeed
			if _, err := dbh.InsertMultipleLines(newLines); err != nil {
				log.Printf("Error inserting multiple lines in the DB... " + err.Error())
				RespondWithError(w, http.StatusBadRequest, "Failed to insert lines")
				return
			}

			RespondWithJSON(w, http.StatusCreated, "Created lines")
			return

		}

		var newLine models.Line
		err = json.Unmarshal(body, &newLine)
		if err != nil {
			log.Printf("Invalid JSON format %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		if _, err := dbh.InsertLine(&newLine); err != nil {
			log.Printf("Error inserting line in the DB... %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to insert line")
			return
		}

		RespondWithJSON(w, http.StatusCreated, "Created line")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")

	}
}

func LineHandler(w http.ResponseWriter, r *http.Request) {
	dbh := database.GetConnector(DBName)

	vars := mux.Vars(r)
	identifier := vars["identifier"]
	log.Printf("%v request on path: %v with identifier %v", r.Method, r.URL.Path, identifier)

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

		if _, err := dbh.UpdateLineByIdentifier(updatedLine, identifier); err != nil {
			log.Printf("Error updating line %v... : %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not update line")
			return
		}

		RespondWithJSON(w, http.StatusOK, updatedLine)

	case http.MethodDelete:
		if _, err := dbh.DeleteLineByIdentifier(identifier); err != nil {
			log.Printf("Error deleting line :%v ...: %v", identifier, err.Error())
			RespondWithError(w, http.StatusBadRequest, "Could not delete the line")
			return
		}
		RespondWithJSON(w, http.StatusCreated, "Delete success")

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v (Metrics)", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)

	mtype := "all"
	split := strings.SplitN(r.URL.Path, "/", 4)

	last := split[len(split)-1]
	if last != "" {
		mtype = last
	}
	log.Printf("Split: %v and last part is: %v", split, last)

	metrics, err := dbh.Metrics(mtype)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data metrics")
		return
	}

	RespondWithJSON(w, http.StatusOK, metrics)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v (Not Found)", r.Method, r.URL.Path)
	RespondWithError(w, http.StatusNotFound, "Not Found")
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v (Not Allowed)", r.Method, r.URL.Path)
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
