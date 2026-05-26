package main

import (
	"log"
	"time"

	"github.com/RustReh/go-project-278/internal/app"
	"github.com/RustReh/go-project-278/internal/config"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	application, err := app.Setup(cfg)
	if err != nil {
		log.Fatalf("App setup: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("close application: %v", err)
		}
		sentry.Flush(2 * time.Second)
	}()

	log.Printf("Starting server at %s", cfg.Addr)
	if err := application.Server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
