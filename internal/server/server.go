package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"kyri56xcaesar/discord_bots_app/internal/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	Router   *mux.Router
	ConfPath string //sqlite .db filepath
}

func NewServer(conf string) *Server {
	server := &Server{
		Router:   mux.NewRouter(),
		ConfPath: conf,
	}

	server.routes()

	return server
}

func (s *Server) routes() {

	s.Router.StrictSlash(true)

	// Root handler for health check
	s.Router.HandleFunc("/", RootHandler)
	s.Router.HandleFunc("/healthz", HealthCheck)
	s.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./cmd/api/web/assets"))))

	// Subrouter for /guild
	guildRouter := s.Router.PathPrefix("/guild").Subrouter()
	guildRouter.HandleFunc("/", GuildHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/members", MembersHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/bots", BotsHandler).Methods("GET", "POST")
	guildRouter.HandleFunc("/lines", RootLineHandler).Methods("GET", "POST")

	membersRouter := guildRouter.PathPrefix("/member").Subrouter()
	membersRouter.HandleFunc("/", RootMemberHandler).Methods("GET", "POST")
	membersRouter.HandleFunc("/{identifier}/", MemberHandler).Methods("GET", "PUT", "DELETE")

	botsRouter := guildRouter.PathPrefix("/bot").Subrouter()
	botsRouter.HandleFunc("/", RootBotHandler).Methods("GET", "POST")
	botsRouter.HandleFunc("/{identifier:[0-9]+}/", BotHandler).Methods("GET", "POST", "PUT", "DELETE")

	lineRouter := guildRouter.PathPrefix("/line").Subrouter()
	lineRouter.HandleFunc("/", RootLineHandler).Methods("GET", "POST")
	lineRouter.HandleFunc("/{identifier:[0-9]+}/", LineHandler).Methods("GET", "PUT", "DELETE")

	s.Router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	s.Router.MethodNotAllowedHandler = http.HandlerFunc(notAllowedHandler)
}

func (s *Server) Start() {

	config, err := loadConfig(s.ConfPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Printf("Config file: %+v", config)

	curpath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	scriptPath := fmt.Sprintf("%s%s", curpath, database.InitSQLScriptPath)

	if err = database.InitDB(config.DBfile, scriptPath); err != nil {
		log.Fatalf("Error during db initialization: %v", err)
	}

	//Enable CORS for all routes
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins(config.AllowedOrigins),
		handlers.AllowedHeaders(config.AllowedHeaders),
		handlers.AllowedMethods(config.AllowedMethods),
	)

	srv := &http.Server{
		Handler: corsOptions(s.Router),
		Addr:    ":" + config.HTTPPort,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server listening on port " + config.HTTPPort)
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

					// Throttle restarts to prevent excessive operations
					if time.Since(lastRestart) > 5*time.Second {
						log.Println("Restarting server...")
						s.restartServer(srv)
						lastRestart = time.Now()
					} else {
						log.Println("Server restart throttled")
					}
				}()

			case syscall.SIGUSR2:
				go func() {
					log.Println("Entering script execution prompt...")
					s.runSQLScript(config.DBfile)
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

	config, err := loadConfig(s.ConfPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Printf("Config file: %+v", config)

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
	var scriptPath string
	fmt.Print("Enter the path to the SQL initialization script: ")
	fmt.Scanln(&scriptPath)

	// Check if the file exists before proceeding
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Printf("SQL script not found at path: %s", scriptPath)
		return
	}

	if err := database.InitDB(dbpath, scriptPath); err != nil {
		log.Printf("Error running SQL init script: %v", err)
		return
	}
	log.Println("SQL init script executed successfully")
}
