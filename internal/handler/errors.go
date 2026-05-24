package handler

import (
	"log"
	"net/http"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func writeAppError(c *gin.Context, err error) {
	appErr, ok := apperr.AsAppError(err)
	if !ok {
		log.Printf("unhandled error: %v", err)
		captureSentry(c, err)
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Code:    string(apperr.CodeInternal),
			Message: "internal server error",
			Payload: map[string]any{"detail": apperr.RootCause(err)},
		})
		return
	}

	status := http.StatusInternalServerError
	switch appErr.Code {
	case apperr.CodeNotFound:
		status = http.StatusNotFound
	case apperr.CodeValidation:
		status = http.StatusBadRequest
	case apperr.CodeConflict:
		status = http.StatusConflict
	}

	payload := appErr.Payload
	if status >= http.StatusInternalServerError {
		log.Printf("internal error [%s]: %s: %v", appErr.Code, appErr.Message, appErr.Err)
		captureSentry(c, appErr)
		payload = apperr.PayloadWithDetail(payload, appErr.Err)
	}

	c.JSON(status, schemas.ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Payload: payload,
	})
}

func captureSentry(c *gin.Context, err error) {
	if err == nil {
		return
	}
	hub := sentry.GetHubFromContext(c.Request.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.CaptureException(err)
}

func parseIDParam(c *gin.Context) (int64, error) {
	return parsePathInt64(c.Param("id"))
}
