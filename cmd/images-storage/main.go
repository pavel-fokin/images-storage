package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"pavel-fokin/images-storage/internal/imagesstorage"
	"pavel-fokin/images-storage/internal/server"
	"pavel-fokin/images-storage/internal/storage"
)

type Config struct {
	Server  server.Config
	Storage storage.Config
}

func ReadConfig() *Config {
	envfile := "local.env"
	if os.Getenv("IMAGES_STORAGE_ENV_FILE") != "" {
		envfile = os.Getenv("IMAGES_STORAGE_ENV_FILE")
	}

	err := godotenv.Load(envfile)
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}

func main() {
	config := ReadConfig()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	storage := storage.New(config.Storage)

	imagesstorage := imagesstorage.New(storage)

	httpServer := server.New(config.Server)
	httpServer.SetupImagesAPIRoutes(imagesstorage)

	go httpServer.Start()

	<-sig

	httpServer.Shutdown()
}
