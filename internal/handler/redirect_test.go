package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RustReh/go-project-278/internal/service"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

func TestRedirect_302_and_records_visit(t *testing.T) {
	r, repo := setupTestRouter(t)

	linkSvc := service.NewLinkService(repo, handlerBaseURL)
	_, err := linkSvc.CreateLink(t.Context(), domain.LinkVO{
		OriginalUrl: "https://example.com/target",
		ShortName:   "abc",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/r/abc", nil)
	req.Header.Set("User-Agent", "curl/8.5.0")
	req.Header.Set("Referer", "https://ref.example/")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status: got %d, want 302, body: %s", rec.Code, rec.Body.String())
	}
	loc := rec.Header().Get("Location")
	if loc != "https://example.com/target" {
		t.Fatalf("Location: got %q", loc)
	}

	visits, err := repo.ListVisits(t.Context(), 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(visits) != 1 {
		t.Fatalf("visits: got %d, want 1", len(visits))
	}
	if visits[0].Status != http.StatusFound || visits[0].LinkId != 1 {
		t.Fatalf("visit: %+v", visits[0])
	}
	if visits[0].UserAgent != "curl/8.5.0" {
		t.Fatalf("user_agent: %q", visits[0].UserAgent)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/r/missing", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d", rec.Code)
	}
}
