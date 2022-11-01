package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/swaggo/http-swagger"

	_ "github.com/pavel-fokin/images-storage/docs"
)

type Config struct {
	Port            string `env:"PORT" envDefault:"8080"`
	HandlerTimeout  int    `env:"IMAGES_SERVER_HANDLER_TIMEOUT" envDefault:"30"`
	ReadTimeout     int    `env:"IMAGES_SERVER_READ_TIMEOUT" envDefault:"30"`
	WriteTimeout    int    `env:"IMAGES_SERVER_WRITE_TIMEOUT" envDefault:"30"`
	ShutdownTimeout int    `env:"IMAGES_SERVER_SHUTDOWN_TIMEOUT" envDefault:"30"`
}

type Server struct {
	config Config
	server *http.Server
	router chi.Router
}

func New(config Config) *Server {
	logger := httplog.NewLogger("images-storage", httplog.Options{
		Concise: true,
		JSON:    true,
	})

	router := chi.NewRouter()

	router.Use(middleware.Timeout(
		time.Duration(config.HandlerTimeout) * time.Second),
	)
	router.Use(httplog.RequestLogger(logger))
	// router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))

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
