package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"kyri56xcaesar/discord-guild-web-app/internal/database"
	"kyri56xcaesar/discord-guild-web-app/internal/models"
)

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"dec": func(i int) int {
		return i - 1
	},
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "welcome.html"), "welcome.html", nil, nil)
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

	RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "bots.html"), "bots.html", nil, bots)
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

	RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "hof.html"), "hof.html", funcMap, members)
}

// Serve clients.html
func fetchNews(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	// RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "news.html"), "news.html", nil, nil)
}

func fetchCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v from %v", r.Method, r.URL.Path, r.UserAgent())

	data := make(map[string]string)
	data["test"] = "adeio"
	RespondWithJSON(w, http.StatusOK, data)
}

func fetchFeaturedNews(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v from %v", r.Method, r.URL.Path, r.UserAgent())

	// RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "featured.html"), "featured.html", nil, nil)
}

func fetchPoll(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	// RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "poll.html"), "poll.html", nil, nil)
}

func submitNews(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v from %v", r.Method, r.URL.Path, r.UserAgent())
}

func votePoll(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v from %v", r.Method, r.URL.Path, r.UserAgent())
}

// Require Autnentication!
func AscLoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v request on path: %v", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			RespondWithError(w, http.StatusBadRequest, "Invalid form data")
			return
		}

		user := models.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		if err := user.VerifyUser(); err != nil {
			log.Printf("Error verifying user: %v", err)
			RespondWithError(w, http.StatusBadRequest, "invalid user credentials")
			return
		}

		time.Sleep(time.Second * 2)
		// TODO: Authenticate user

		if user.Username == "diego" && user.Password == "diego" {
			RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "ascend-dashboard.html"), "ascend-dashboard.html", nil, user)
		} else {
			RespondWithError(w, http.StatusUnauthorized, "Unauthorized.")
		}

	case http.MethodGet:
		RespondWithTemplate(w, http.StatusOK, filepath.Join("web", "templates", "ascend-login.html"), "ascend-login.html", nil, nil)

	default:
		RespondWithError(w, http.StatusMethodNotAllowed, "Not allowed.")
	}
}
