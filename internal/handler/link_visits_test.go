package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

func TestLinkVisits_GetAll_WithRangeHeader(t *testing.T) {
	r, repo := setupTestRouter(t)
	ctx := t.Context()

	linkSvc := service.NewLinkService(repo, handlerBaseURL)
	link, _ := linkSvc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/a",
		ShortName:   "a",
	})

	visitSvc := service.NewVisitService(repo)
	for range 3 {
		if _, _, err := visitSvc.Redirect(ctx, "a", "10.0.0.1", "ua", ""); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)
	req.Header.Set("Range", "[0,2]")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d, body: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Range"); got != "link_visits 0-2/3" {
		t.Fatalf("Content-Range: got %q", got)
	}

	var resp []schemas.LinkVisitResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("len: got %d, want 2", len(resp))
	}
	if resp[0].LinkID != link.Id || resp[0].Status != 302 {
		t.Fatalf("first visit: %+v", resp[0])
	}
}

func TestLinkVisits_GetAll_WithQueryRange(t *testing.T) {
	r, repo := setupTestRouter(t)
	visitSvc := service.NewVisitService(repo)

	linkSvc := service.NewLinkService(repo, handlerBaseURL)
	_, _ = linkSvc.CreateLink(t.Context(), domain.LinkVO{
		OriginalUrl: "https://example.com/x",
		ShortName:   "x",
	})
	_, _, _ = visitSvc.Redirect(t.Context(), "x", "1.1.1.1", "ua", "")

	req := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=[0,10]", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestLinkVisits_MissingRange_400(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d", rec.Code)
	}
}
