package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	switch c.Query("sentry_test") {
	case "error":
		err := errors.New("sentry test: manual error from /ping")
		sentry.CaptureException(err)
		sentry.Flush(2 * time.Second)
		c.String(http.StatusInternalServerError, "error sent to sentry")
		return
	case "panic":
		log.Println("sentry_test=panic: raising panic (check server log and Sentry Issues)")
		panic("sentry test: panic from /ping")
	}

	response := schemas.PingResponse{
		Message: "pong",
		Status:  http.StatusOK,
	}
	c.String(response.Status, response.Message)
}
