package server

import (
	"log"
	"net/http"

	"kyri56xcaesar/discord_bots_app/internal/config"
	"kyri56xcaesar/discord_bots_app/internal/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	Router    *mux.Router
	Conf_path string //sqlite .db filepath
}

func NewServer(conf string) *Server {
	server := &Server{
		Router:    mux.NewRouter(),
		Conf_path: conf,
	}

	server.routes()

	return server
}

func (s *Server) routes() {
	// Root handler for health check
	s.Router.HandleFunc("/", RootHandler)

	// Subrouter for /guild
	guildRouter := s.Router.PathPrefix("/guild").Subrouter()
	guildRouter.HandleFunc("/{type}", MembersHandler).Methods("GET", "POST")

	membersRouter := guildRouter.PathPrefix("/member").Subrouter()
	membersRouter.HandleFunc("/", RootMemberHandler).Methods("GET", "POST")
	membersRouter.HandleFunc("/{identifier}", MemberHandler).Methods("GET", "UPDATE", "DELETE")

}

func (s *Server) Start() {

	config, err := config.LoadConfig(s.Conf_path)

	log.Printf("Config file: %v", config)

	if err != nil {
		log.Fatalf("Required environment variables are missing: %+ v", err)
	}

	if err = database.InitDB(config.DBfile); err != nil {
		log.Fatalf("Error during db initialization: %v", err)
	}

	// Enable CORS for all routes
	o := handlers.AllowedOrigins([]string{"*"}) // Allow all origins
	h := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Content-Type",
		"Accept",
		"Authorization",
	})
	m := handlers.AllowedMethods([]string{
		"GET",
		"POST",
		"OPTIONS",
	})

	// Start the server
	log.Println("Server listening on port " + config.HTTPPort)
	if err := http.ListenAndServe(":"+config.HTTPPort, handlers.CORS(o, h, m)(s.Router)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
