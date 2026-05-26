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

	switch appErr.Code {
	case apperr.CodeValidation:
		fields := apperr.AsFieldErrors(appErr.Payload)
		if len(fields) == 0 {
			fields = map[string]string{"_": appErr.Message}
		}
		writeValidationErrors(c, fields)
		return
	case apperr.CodeConflict:
		writeValidationErrors(c, map[string]string{
			"short_name": "short name already in use",
		})
		return
	case apperr.CodeNotFound:
		c.JSON(http.StatusNotFound, schemas.InvalidRequestResponse{
			Error: appErr.Message,
		})
		return
	}

	log.Printf("internal error [%s]: %s: %v", appErr.Code, appErr.Message, appErr.Err)
	captureSentry(c, appErr)
	c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
		Code:    string(apperr.CodeInternal),
		Message: appErr.Message,
		Payload: apperr.PayloadWithDetail(appErr.Payload, appErr.Err),
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
