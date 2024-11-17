package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"kyri56xcaesar/discord_bots_app/internal/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	Router   *mux.Router
	Config   *EnvConfig // sqlite .db filepath
	serverID int
}

type ServerError struct{}

func (serror *ServerError) Error() string {
	return "server error"
}

// there should be a limit to the amount of servers possible
const (
	MAX_SERVERS  int    = 100
	templatePath string = "/web/templates/index.html"
	assetsPath   string = "./web/assets"
)

// Server pool
var (
	serverPool    [MAX_SERVERS]Server
	currentIndex  int = 0
	CurrentServer *Server
	DBName        string

	poolMutex sync.Mutex
)

func NewServer(conf string) (*Server, error) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if currentIndex >= MAX_SERVERS {
		log.Println("Server pool limit reached. Cannot create more servers. UPDATE MAX SERVERS")
		return nil, &ServerError{}
	}

	server := serverPool[currentIndex]

	server.serverID = currentIndex
	server.Router = mux.NewRouter()
	server.Config = loadConfig(conf)

	log.Printf("ServerID: %d\n[CFG]...Loading configurations...\n%v\n", currentIndex, server.Config.toString())

	server.routes()

	currentIndex += 1

	return &server, nil
}

// ToDo
// CRUD: Check
// Add Pagination: ToDo
// Bulk Operations: ToDo
// Search(Filteting): ToDo
// Metrics/Stats: ToDo
// Exports/Imports: ToDo
// Utilities: Check
// Documentation: ToDo
//

func (s *Server) routes() {
	s.Router.StrictSlash(true)

	// Root handler for health check
	// Templates
	s.Router.HandleFunc("/", RootHandler)
	s.Router.HandleFunc("/dbots", BotsDHandler)
	s.Router.HandleFunc("/hof", HofHandler)

	s.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))
	// Admin (Must Verify)
	s.Router.HandleFunc("/admin/healthz", HealthCheck)

	// Subrouter for /guild
	guildRouter := s.Router.PathPrefix("/guild").Subrouter()

	// CRUD
	guildRouter.HandleFunc("/", GuildHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/members", MembersHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/bots", BotsHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/lines", LinesHandler).Methods("GET", "POST")

	// Specific CRUD
	guildRouter.HandleFunc("/member/{identifier}", MemberHandler).Methods("GET", "PUT", "DELETE")
	guildRouter.HandleFunc("/bot/{identifier:[0-9]+}", BotHandler).Methods("GET", "POST", "PUT", "DELETE")
	guildRouter.HandleFunc("/line/{identifier:[0-9]+}", LineHandler).Methods("GET", "PUT", "DELETE")

	// Filtered Search
	guildRouter.HandleFunc("/search/members", GMultipleData).Methods("GET")
	guildRouter.HandleFunc("/search/bots", GMultipleData).Methods("GET")
	guildRouter.HandleFunc("/search/lines", GMultipleData).Methods("GET")

	// Filtered Update, Delete
	guildRouter.HandleFunc("/delete/members", UDMultipleData).Methods("DELETE", "PUT", "PATCH")
	guildRouter.HandleFunc("/delete/bots", UDMultipleData).Methods("DELETE", "PUT", "PATCH")
	guildRouter.HandleFunc("/delete/lines", UDMultipleData).Methods("DELETE", "PUT", "PATCH")

	// Utility endpoints and Metrics
	guildRouter.HandleFunc("/get/members/{identifier:[a-zA-Z,]+}", DataHandler).Methods("GET")
	guildRouter.HandleFunc("/get/bots/{identifier:[a-zA-Z,]+}", DataHandler).Methods("GET")
	guildRouter.HandleFunc("/get/lines/{identifier:[a-zA-Z,]+}", DataHandler).Methods("GET")

	guildRouter.HandleFunc("/data/members", DataIndexHandler).Methods("GET")
	guildRouter.HandleFunc("/data/bots", DataIndexHandler).Methods("GET")
	guildRouter.HandleFunc("/data/lines", DataIndexHandler).Methods("GET")

	guildRouter.HandleFunc("/metrics", MetricsHandler).Methods("GET")
	guildRouter.HandleFunc("/metrics/members", MetricsHandler).Methods("GET")
	guildRouter.HandleFunc("/metrics/bots", MetricsHandler).Methods("GET")
	guildRouter.HandleFunc("/metrics/lines", MetricsHandler).Methods("GET")

	s.Router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	s.Router.MethodNotAllowedHandler = http.HandlerFunc(notAllowedHandler)
}

func (s *Server) Start() {
	log.Print("Server starting...")

	// Database Setuo
	curpath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	scriptPath := fmt.Sprintf("%s%s", curpath, database.InitSQLScriptPath)
	// log.Printf("Current directory: %s", curpath)
	// log.Printf("Script path: %s", scriptPath)

	// Init database
	if err = database.InitDB(s.Config.DBfile, scriptPath); err != nil {
		log.Printf("[INIT DB]Error during db initialization: %v", err)
		// Should handle case that it is fatal...
	}
	// Set database reference
	DBName = s.Config.DBfile

	// Service Setup
	// Enable CORS for all routes
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins(s.Config.AllowedOrigins),
		handlers.AllowedHeaders(s.Config.AllowedHeaders),
		handlers.AllowedMethods(s.Config.AllowedMethods),
	)

	srv := &http.Server{
		Handler: corsOptions(s.Router),
		Addr:    ":" + s.Config.HTTPPort,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server listening on port " + s.Config.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Set up a buffered channel for signals
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, os.Interrupt, syscall.SIGUSR1, syscall.SIGUSR2)
	// Mutex and timestamp for throttling server restarts
	var restartMutex sync.Mutex
	var lastRestart time.Time

	for {
		select {
		case sigReceived := <-sig:
			switch sigReceived {
			case os.Interrupt:
				log.Println("Received interrupt signal, shutting down gracefully...")

				// Graceful shutdown
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := srv.Shutdown(ctx); err != nil {
					log.Fatalf("Server shutdown failed: %v", err)
				}
				log.Println("Server exited properly")
				return

			case syscall.SIGUSR1:
				go func() {
					restartMutex.Lock()
					defer restartMutex.Unlock()

					// Check for a new config path from an environment variable
					newConfigPath := os.Getenv("NEW_CONFIG_PATH")
					if newConfigPath != "" {
						log.Printf("Using new config path: %s", newConfigPath)
						s.Config.ConfigPath = newConfigPath // Update the path if provided
					}

					// Throttle restarts to prevent excessive operations
					if time.Since(lastRestart) > 5*time.Second {
						s.restartServer(srv)
						lastRestart = time.Now()
					} else {
						log.Println("Server restart throttled")
					}
				}()

			case syscall.SIGUSR2:
				go func() {
					log.Println("Entering script execution prompt...")
					s.runSQLScript(s.Config.DBfile)
				}()
			}
		}
	}
}

func (s *Server) restartServer(srv *http.Server) {
	// Shutdown the current server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server for restart: %v", err)
		return
	}
	log.Println("Server shut down for restart")

	config := loadConfig(s.Config.ConfigPath)

	log.Printf("Config file: %+v", config)

	time.Sleep(5 * time.Second)
	// Restart the server with the same configuration
	go func() {
		newSrv := &http.Server{
			Handler: s.Router,
			Addr:    ":" + config.HTTPPort,
		}
		log.Println("Restarting server...")
		if err := newSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to restart: %v", err)
		}
	}()
}

// runSQLInitScript allows dynamically running an SQL script
func (s *Server) runSQLScript(dbpath string) {
	// Example: Dynamically get the SQL script path from the user

	sqlDir := "/internal/database/sqlscripts/"
	wd, _ := os.Getwd()

	sqlDir = fmt.Sprintf("%s%s", wd, sqlDir)

	scriptPath := os.Getenv("SQL_SCRIPT_PATH")
	if scriptPath != "" {
		log.Printf("Using the script path: %s", scriptPath)
	} else {
		log.Printf("SQL Env Path empty... prompting...")
		fmt.Printf("Current dir: %v\n", wd)
		out, err := exec.Command("ls", "--color=auto", sqlDir).Output()
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("Avaliable scripts:\n%s\n", string(out))

		contents := strings.Split(string(out), "\n")
		for index, content := range contents {
			if strings.HasSuffix(content, ".sql") {
				fmt.Printf("%d.  %s\n", index+1, content)
			}
		}
		fmt.Print("\n\nEnter the SQL script[name/index]: ")
		fmt.Scanln(&scriptPath)
		index, err := strconv.Atoi(scriptPath)
		if err == nil {
			if index >= 1 && index <= len(contents) {
				scriptPath = contents[index-1]
			}
		}

		log.Printf("Script chosen: %s", scriptPath)

	}

	absPath := fmt.Sprintf("%s%s", sqlDir, scriptPath)
	// Check if the file exists before proceeding
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Printf("SQL script not found at path: %s", absPath)
		return
	}

	if err := database.InitDB(dbpath, absPath); err != nil {
		log.Printf("Error running SQL init script: %v", err)
		return
	}
	log.Println("SQL init script executed successfully")
}
