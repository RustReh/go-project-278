package main

import (
	"log"
	"net/http"

	"github.com/RustReh/go-project-278/internal/config"
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if cfg.SentryDsn == "" {
		log.Println("SENTRY_DSN is empty — events will NOT be sent to Sentry")
	} else if err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDsn,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
	} else {
		log.Println("Sentry initialized")
	}
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	router.GET("/ping", handler.Ping)

	s := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	log.Printf("Starting server at %s", cfg.Addr)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
