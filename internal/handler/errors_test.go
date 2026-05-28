package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/gin-gonic/gin"
)

func TestWriteAppError_DoesNotLeakInternalDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, rec := newTestContext()
	dbErr := errors.Join(repository.ErrInternal, errors.New(`relation "links" does not exist`))
	writeAppError(c, apperr.WithPayload(
		apperr.CodeInternal,
		"repository error",
		map[string]any{"short_name": "abc"},
		dbErr,
	))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["payload"] != nil {
		t.Fatalf("expected no payload, got %#v", body["payload"])
	}
	if body["message"] != "internal server error" {
		t.Fatalf("message: %#v", body["message"])
	}
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, rec
}
