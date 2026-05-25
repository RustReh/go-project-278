package service_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/RustReh/go-project-278/internal/service/domain"
	"github.com/RustReh/go-project-278/internal/testutil"
)

func TestVisitService_Redirect(t *testing.T) {
	repo := testutil.NewMemRepo()
	linkSvc := service.NewLinkService(repo, testBaseURL)
	visitSvc := service.NewVisitService(repo)
	ctx := context.Background()

	_, err := linkSvc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/go",
		ShortName:   "go",
	})
	if err != nil {
		t.Fatal(err)
	}

	url, status, err := visitSvc.Redirect(ctx, "go", "192.168.1.1", "Mozilla/5.0", "https://from/")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.com/go" || status != http.StatusFound {
		t.Fatalf("got url=%q status=%d", url, status)
	}

	page, err := visitSvc.ListVisits(ctx, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || page.Visits[0].Ip != "192.168.1.1" {
		t.Fatalf("page: %+v", page)
	}
}

func TestVisitService_Redirect_NotFound(t *testing.T) {
	visitSvc := service.NewVisitService(testutil.NewMemRepo())
	_, _, err := visitSvc.Redirect(context.Background(), "nope", "1.1.1.1", "", "")
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeNotFound {
		t.Fatalf("got %#v", err)
	}
}
