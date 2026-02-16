package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	serviceUser "github.com/ArkaniLoveCoding/School-manajement/service/users"
)

type ApiServer struct {
	Addr   string
	db     *sqlx.DB
	server *http.Server
}

func ApiServerAddr(addr string, db *sqlx.DB) *ApiServer {
	return &ApiServer{
		Addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {

	//setup mux router
	router := mux.NewRouter()

	//subrouter of all this router on this router
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// testing if the server is working!
	subRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "Successfully to testing the web server, now the web server is working!",
			"data": "Hello world!"
		}`))
	})

	// not authenticate 

	userStore := serviceUser.NewStore(s.db)
	userService := serviceUser.NewHandlerUser(userStore)

	userStores := serviceUser.NewStore(s.db)
	userServices := serviceUser.NewHandlerUserForAuthenticate(userStores)

	// Create HTTP server
	s.server = &http.Server{
		Addr:   s.Addr,
		Handler: router,
	}

	log.Printf("Server starting on %s", s.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		return errors.New(err.Error())
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *ApiServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
