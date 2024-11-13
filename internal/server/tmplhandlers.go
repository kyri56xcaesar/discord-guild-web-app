package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"kyri56xcaesar/discord_bots_app/internal/database"
)

// Serve bots.html
func BotsDHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	if dbh == nil {
		RespondWithError(w, http.StatusInternalServerError, "Database connection failed")
		return
	}
	bots, err := dbh.GetAllBots()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bots")
		return
	}

	tmplPath := filepath.Join("cmd", "api", "web", "templates", "bots.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, bots)
}

// Serve hof.html
func HofHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	if dbh == nil {
		RespondWithError(w, http.StatusInternalServerError, "Database connection failed")
		return
	}
	members, err := dbh.GetAllMembers()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve members")
		return
	}

	// Sort members on msg_count
	sort.Slice(members, func(i, j int) bool {
		return members[i].MsgCount > members[j].MsgCount
	})

	// Create a slice of maps with index and member data
	rankedMembers := make([]map[string]interface{}, len(members))
	for i, member := range members {
		rankedMembers[i] = map[string]interface{}{
			"Index":  i + 1, // 1-based index
			"Member": member,
		}
	}

	tmplPath := filepath.Join("cmd", "api", "web", "templates", "hof.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, rankedMembers)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Serve clients.html
func ClientsHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("cmd", "api", "web", "templates", "clients.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
