package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/pavel-fokin/images-storage/internal/imagesstorage"
	"github.com/pavel-fokin/images-storage/internal/server"
	"github.com/pavel-fokin/images-storage/internal/storage"
)

type Config struct {
	Server  server.Config
	Storage storage.Config
}

func ReadConfig() *Config {
	envFile := "local.env"
	if os.Getenv("IMAGES_STORAGE_ENV_FILE") != "" {
		envFile = os.Getenv("IMAGES_STORAGE_ENV_FILE")
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}

// @Title Images Storage API
// @Version 0.0.1
// @Description Images Storage is a service that lets you store, retrieve, and cutout images.
// @BasePath /
func main() {
	config := ReadConfig()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	storage := storage.New(config.Storage)

	imagesStorage := imagesstorage.New(storage)

	httpServer := server.New(config.Server)
	httpServer.SetupImagesAPIRoutes(imagesStorage)

	go httpServer.Start()

	<-sig

	httpServer.Shutdown()
}
