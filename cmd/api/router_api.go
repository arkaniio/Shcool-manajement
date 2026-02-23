package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware"
	serviceStudent "github.com/ArkaniLoveCoding/Shcool-manajement/service/students"
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

	// Setup mux router
	router := mux.NewRouter()

	// Apply global middleware (order matters - first to last)
	// 1. Request ID middleware - adds unique ID to each request
	router.Use(middleware.RequestIDMiddleware)
	// 2. Logger middleware - logs all HTTP requests with socket hang up detection
	router.Use(middleware.LoggerResponse)

	// Subrouter for API v1
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Testing if the server is working!
	subRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "Successfully to testing the web server, now the web server is working!",
			"data": "Hello world!"
		}`))
	})

	// The router for the services
	userStore := serviceUser.NewStore(s.db)
	userService := serviceUser.NewHandlerUser(userStore)

	// Router for the register user
	subRouter.Handle(
		"/register",
		http.HandlerFunc(
			userService.Register_Bp,
		),
	).Methods("POST")

	// Router for the login user
	subRouter.Handle(
		"/login",
		http.HandlerFunc(
			userService.Login_Bp,
		),
	).Methods("POST")

	// Router for the profile user (with auth middleware)
	subRouter.Handle(
		"/profile",
		middleware.TokenIdMiddleware(
			http.HandlerFunc(userService.Profile_Bp),
		),
	).Methods("GET")

	// Router for the update user (with auth middleware)
	subRouter.Handle(
		"/users/{id}",
			http.HandlerFunc(userService.Update_Bp),
	).Methods("PATCH")

	// Router to see the file path for frontend to catch it
	subRouter.Handle(
		"/users/profile/{filename}",
		http.HandlerFunc(
			userService.Image_Bp,
		),
	).Methods("GET")

	//router for the student routes
	studentStore := serviceStudent.NewStudentStore(s.db)
	studentService := serviceStudent.NewHandlerStudent(studentStore)

	//router for register as a student
	subRouter.Handle(
		"/students/register",
			middleware.TokenIdMiddleware(http.HandlerFunc(
				studentService.RegisterAsStudent_Bp,
			)),
	).Methods("POST")

	// Create HTTP server
	s.server = &http.Server{
		Addr:         s.Addr,
		Handler:      router,
		// Timeouts to prevent slowloris attacks
		ReadTimeout:  15 * time.Second,  // 15 minutes
		WriteTimeout: 15 * time.Second,  // 15 minutes
		IdleTimeout:  60 * time.Second,       // 60 seconds
	}

	// Start listening
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

