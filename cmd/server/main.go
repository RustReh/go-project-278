package main

import (
	"log"
	"net/http"

	"github.com/RustReh/go-project-278/internal/config"
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

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
