package main

import (
	"context"
	"log"
	"os/signal"
	_ "subscriptions_rest/docs"
	"subscriptions_rest/internal/config"
	"subscriptions_rest/internal/repository/postgres"
	"subscriptions_rest/internal/server/http"
	"sync"
	"syscall"
)

//	@title			Subscriptions Service
//	@version		1.0
//	@description	This is a service for checking subscriptions.

//	@license.name	Subscriptions Service 1.0

// @host		localhost
// @BasePath	/
//
//go:generate swag init -g ./cmd/subscriptions/main.go --pd --parseInternal --output ../../docs --dir ../..
func main() {
	// TODO
	// check empty user id
	cfg, err := config.New("../../config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	repository, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	httpServer := http.New(cfg, repository)

	go func() {
		defer wg.Done()
		err := httpServer.Start(cfg.Port)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	<-ctx.Done()
	err = httpServer.Stop()
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
