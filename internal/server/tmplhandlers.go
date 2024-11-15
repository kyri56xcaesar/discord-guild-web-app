package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"kyri56xcaesar/discord_bots_app/internal/database"
)

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"dec": func(i int) int {
		return i - 1
	},
}

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

	err = tmpl.Execute(w, bots)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error executing template")
	}
}

// Serve hof.html
func HofHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	dbh := database.GetConnector(DBName)
	if dbh == nil {
		log.Print("Error on database connection")
		RespondWithError(w, http.StatusInternalServerError, "Database connection failed")
		return
	}
	members, err := dbh.GetAllMembers()
	if err != nil {
		log.Print("Failed to retrieve members")
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve members")
		return
	}

	// Sort members on msg_count
	sort.Slice(members, func(i, j int) bool {
		return members[i].Msgcount > members[j].Msgcount
	})

	opacity := "40"

	for _, v := range members {
		v.Usercolor += opacity
	}

	tmplPath := filepath.Join("cmd", "api", "web", "templates", "hof.html")
	tmpl, err := template.New("hof.html").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		log.Print("Error loading template. " + err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Error loading template")
		return
	}

	err = tmpl.Execute(w, members)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error executing template")
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
