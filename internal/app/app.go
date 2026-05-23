package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/RustReh/go-project-278/internal/config"
	"github.com/RustReh/go-project-278/internal/db"
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/router"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *http.Server
	sqlDB  *sql.DB
}

func (a *App) Close() error {
	if a.sqlDB != nil {
		return a.sqlDB.Close()
	}
	return nil
}

func Setup(cfg *config.AppConfig) (*App, error) {
	if err := initSentry(cfg.SentryDsn); err != nil {
		return nil, err
	}

	sqlDB, err := db.Open(context.Background(), cfg.DataBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	repo := repository.NewPostgresRepo(sqlDB)
	linkService := service.NewLinkService(repo, cfg.BaseURL)
	linksHandler := handler.NewLinksHandler(linkService)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	router.Register(engine, linksHandler)

	return &App{
		Server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      engine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		sqlDB: sqlDB,
	}, nil
}

func initSentry(dsn string) error {
	if dsn == "" {
		log.Println("SENTRY_DSN is empty — events will NOT be sent to Sentry")
		return nil
	}
	if err := sentry.Init(sentry.ClientOptions{Dsn: dsn}); err != nil {
		return fmt.Errorf("sentry init: %w", err)
	}
	log.Println("Sentry initialized")
	return nil
}
