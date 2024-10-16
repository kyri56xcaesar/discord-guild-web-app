package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	lineRouter.HandleFunc("/{identifier:[0-9]}/", LineHandler).Methods("GET", "PUT", "DELETE")

	s.Router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
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

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
}
