package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	port int
	env string
}

type Application struct {
	config Config
	logger *log.Logger
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := &Application{
		config: cfg,
		logger: logger,
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr: addr,
		Handler: app.routes(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

}