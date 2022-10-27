package main

import (
  "fmt"
	"os"
	"os/signal"
	"syscall"

  "github.com/caarlos0/env/v6"

  "pavel-fokin/images-storage/internal/server"
)

type Config struct {
	Server server.Config
}

func ReadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return cfg
}

func main () {

	config := ReadConfig()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	httpServer := server.New(config.Server)
	httpServer.SetupImagesAPIRoutes()

	go httpServer.Start()

	<- sig

	httpServer.Shutdown()
}
