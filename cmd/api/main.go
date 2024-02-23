package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type Application struct {
	config Config
	logger *log.Logger
	DB     *gorm.DB
}


func init() {
	LoadEnv()
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENVIRONMENT"), "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DATABASE_DSN"), "Postgres dsn")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	DB, err := db.DB()
	if err != nil {
		logger.Fatal(err)
	}
	defer DB.Close()

	logger.Printf("database connection pool established")

	app := &Application{
		config: cfg,
		logger: logger,
		DB:     db,
	}

	app.migrations()

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
