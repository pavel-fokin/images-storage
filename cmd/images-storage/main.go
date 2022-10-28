package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"

	"pavel-fokin/images-storage/internal/imagesstorage"
	"pavel-fokin/images-storage/internal/server"
	"pavel-fokin/images-storage/internal/storage"
)

type Config struct {
	Server  server.Config
	Storage storage.Config
}

func ReadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}

func main() {

	config := ReadConfig()

	rootCtx, cancelRootCtx := context.WithCancel(context.Background())
	defer cancelRootCtx()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	storage := storage.New(rootCtx, config.Storage)

	imagesstorage := imagesstorage.New(storage)

	httpServer := server.New(config.Server)
	httpServer.SetupImagesAPIRoutes(imagesstorage)

	go httpServer.Start()

	<-sig

	httpServer.Shutdown()
}
