package server

import (
	"fmt"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	Port            string `env:"PORT" envDefault:"8080"`
	ReadTimeout     int    `env:"IMAGES_SERVER_READ_TIMEOUT" envDefault:"5"`
	WriteTimeout    int    `env:"IMAGES_SERVER_WRITE_TIMEOUT" envDefault:"5"`
	ShutdownTimeout int    `env:"IMAGES_SERVER_SHUTDOWN_TIMEOUT" envDefault:"5"`
}

type Server struct {
	config Config
	server *http.Server
	router chi.Router
}

func New(config Config) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
	}

	return &Server{
		config: config,
		server: server,
		router: router,
	}
}

func (s *Server) Start() {
	fmt.Println("Starting images-storage...", s.config.Port)
	s.server.ListenAndServe()
}

func (s *Server) Shutdown() {
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(
		context.Background(), time.Duration(s.config.ShutdownTimeout)*time.Second,
	)
	defer cancelShutdownCtx()

	s.server.Shutdown(shutdownCtx)
}
