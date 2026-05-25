package handler_test

import (
	"testing"

	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/RustReh/go-project-278/internal/router"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/RustReh/go-project-278/internal/testutil"
	"github.com/gin-gonic/gin"
)

const handlerBaseURL = "https://short.io/"

func setupTestRouter(t *testing.T) (*gin.Engine, *testutil.MemRepo) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := testutil.NewMemRepo()
	linkSvc := service.NewLinkService(repo, handlerBaseURL)
	visitSvc := service.NewVisitService(repo)

	r := gin.New()
	router.Register(
		r,
		handler.NewLinksHandler(linkSvc),
		handler.NewLinkVisitsHandler(visitSvc),
		handler.NewRedirectHandler(visitSvc),
	)
	return r, repo
}
