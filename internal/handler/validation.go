package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func writeInvalidRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, schemas.InvalidRequestResponse{
		Error: "invalid request",
	})
}

func writeValidationErrors(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, schemas.ValidationErrorsResponse{
		Errors: fields,
	})
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		if isInvalidJSON(err) {
			writeInvalidRequest(c)
			return false
		}
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			writeValidationErrors(c, formatValidationErrors(verr))
			return false
		}
		writeInvalidRequest(c)
		return false
	}
	return true
}

func isInvalidJSON(err error) bool {
	var syntax *json.SyntaxError
	var unmarshal *json.UnmarshalTypeError
	return errors.As(err, &syntax) ||
		errors.As(err, &unmarshal) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF)
}

func formatValidationErrors(verr validator.ValidationErrors) map[string]string {
	out := make(map[string]string, len(verr))
	for _, fe := range verr {
		out[validationFieldKey(fe)] = fe.Error()
	}
	return out
}

func validationFieldKey(fe validator.FieldError) string {
	switch fe.Field() {
	case "OriginalURL":
		return "original_url"
	case "ShortName":
		return "short_name"
	default:
		return toSnakeCase(fe.Field())
	}
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
