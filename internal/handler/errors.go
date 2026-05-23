package handler

import (
	"net/http"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/gin-gonic/gin"
)

func writeAppError(c *gin.Context, err error) {
	appErr, ok := apperr.AsAppError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Code:    string(apperr.CodeInternal),
			Message: "internal server error",
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

	c.JSON(status, schemas.ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Payload: appErr.Payload,
	})
}

func parseIDParam(c *gin.Context) (int64, error) {
	return parsePathInt64(c.Param("id"))
}
