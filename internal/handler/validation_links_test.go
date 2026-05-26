package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RustReh/go-project-278/internal/schemas"
)

func TestLinks_Create_InvalidJSON_400(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d", rec.Code)
	}
	var resp schemas.InvalidRequestResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error != "invalid request" {
		t.Fatalf("got %#v", resp)
	}
}

func TestLinks_Create_InvalidURL_422(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/links",
		strings.NewReader(`{"original_url":"not-a-url","short_name":"abc"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp schemas.ValidationErrorsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	msg, ok := resp.Errors["original_url"]
	if !ok || !strings.Contains(msg, "url") {
		t.Fatalf("errors: %#v", resp.Errors)
	}
}

func TestLinks_Create_ShortNameTooShort_422(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/links",
		strings.NewReader(`{"original_url":"https://example.com/x","short_name":"ab"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: got %d", rec.Code)
	}
}

func TestLinks_Update_MissingShortName_422(t *testing.T) {
	r, _ := setupTestRouter(t)

	post := httptest.NewRequest(http.MethodPost, "/api/links",
		strings.NewReader(`{"original_url":"https://example.com/x","short_name":"xxx"}`))
	post.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	r.ServeHTTP(postRec, post)

	put := httptest.NewRequest(http.MethodPut, "/api/links/1",
		strings.NewReader(`{"original_url":"https://example.com/y"}`))
	put.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	r.ServeHTTP(putRec, put)

	if putRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: got %d, body: %s", putRec.Code, putRec.Body.String())
	}
}
