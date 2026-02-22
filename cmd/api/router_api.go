package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware"
	serviceUser "github.com/ArkaniLoveCoding/Shcool-manajement/service/users"
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

	// the router for the services
	userStore := serviceUser.NewStore(s.db)
	userService := serviceUser.NewHandlerUser(userStore)

	//router for the resgister user
	subRouter.Handle(
		"/register",
		http.HandlerFunc(
			userService.Register_Bp,
		),
	).Methods("POST")

	//router for the login user
	subRouter.Handle(
		"/login",
		http.HandlerFunc(
			userService.Login_Bp,
		),
	).Methods("POST")

	//router for the profile user
	subRouter.Handle(
		"/profile",
		middleware.TokenIdMiddleware(http.HandlerFunc(
			userService.Profile_Bp,
		)),
	).Methods("GET")

	//router for the update user
	subRouter.Handle(
		"/users/{id}",
		http.HandlerFunc(
			userService.Update_Bp,
		),
	).Methods("PATCH")

	//router to see the file path for frontend to catch it 
	subRouter.Handle(
		"/users/profile/{filename}",
		http.HandlerFunc(
			userService.Image_Bp,
		),
	).Methods("GET")

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
