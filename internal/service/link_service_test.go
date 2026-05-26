package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/RustReh/go-project-278/internal/service/domain"
	"github.com/RustReh/go-project-278/internal/testutil"
)

const testBaseURL = "https://short.io/"

func newTestService() (*service.LinkService, *testutil.MemRepo) {
	repo := testutil.NewMemRepo()
	return service.NewLinkService(repo, testBaseURL), repo
}

func TestCreateLink_WithShortName(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	link, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/a",
		ShortName:   "exmpl",
	})
	if err != nil {
		t.Fatalf("CreateLink: %v", err)
	}
	if link.ShortName != "exmpl" {
		t.Fatalf("short_name: got %q, want exmpl", link.ShortName)
	}
	if link.ShortUrl != "https://short.io/r/exmpl" {
		t.Fatalf("short_url: got %q", link.ShortUrl)
	}
}

func TestCreateLink_WithoutShortName_GeneratesUnique(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	link, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/auto",
	})
	if err != nil {
		t.Fatalf("CreateLink: %v", err)
	}
	if link.ShortName == "" {
		t.Fatal("expected generated short_name")
	}
	if len(link.ShortName) != 12 {
		t.Fatalf("generated short_name length: got %d, want 12", len(link.ShortName))
	}
	if link.ShortUrl != testBaseURL+"r/"+link.ShortName {
		t.Fatalf("short_url: got %q", link.ShortUrl)
	}
}

func TestCreateLink_Conflict(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()
	vo := domain.LinkVO{OriginalUrl: "https://example.com/1", ShortName: "dup"}

	if _, err := svc.CreateLink(ctx, vo); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/2",
		ShortName:   "dup",
	})
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeValidation {
		t.Fatalf("got %#v", err)
	}
	fields := apperr.AsFieldErrors(appErr.Payload)
	if fields["short_name"] != "short name already in use" {
		t.Fatalf("fields: %#v", fields)
	}
}

func TestGetLinkByID_NotFound(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.GetLinkByID(context.Background(), 999)
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeNotFound {
		t.Fatalf("got %#v", err)
	}
}

func TestListLinks(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	for i := 1; i <= 11; i++ {
		name := "ln" + string(rune('0'+i))
		if _, err := svc.CreateLink(ctx, domain.LinkVO{
			OriginalUrl: "https://example.com/" + name,
			ShortName:   name,
		}); err != nil {
			t.Fatal(err)
		}
	}

	page, err := svc.ListLinks(ctx, 0, 10)
	if err != nil {
		t.Fatalf("ListLinks: %v", err)
	}
	if page.Total != 11 {
		t.Fatalf("total: got %d, want 11", page.Total)
	}
	if len(page.Links) != 10 {
		t.Fatalf("len: got %d, want 10", len(page.Links))
	}
	if page.Links[0].Id != 1 || page.Links[9].Id != 10 {
		t.Fatalf("ids: first=%d last=%d", page.Links[0].Id, page.Links[9].Id)
	}

	page2, err := svc.ListLinks(ctx, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Links) != 5 {
		t.Fatalf("page2 len: got %d, want 5", len(page2.Links))
	}
	if page2.Links[0].Id != 6 || page2.Links[4].Id != 10 {
		t.Fatalf("page2 ids: first=%d last=%d", page2.Links[0].Id, page2.Links[4].Id)
	}
}

func TestUpdateLink(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	created, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/old",
		ShortName:   "old",
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := svc.UpdateLink(ctx, created.Id, domain.LinkVO{
		OriginalUrl: "https://example.com/new",
		ShortName:   "new",
	})
	if err != nil {
		t.Fatalf("UpdateLink: %v", err)
	}
	if updated.OriginalUrl != "https://example.com/new" || updated.ShortName != "new" {
		t.Fatalf("got %+v", updated)
	}
	if updated.ShortUrl != "https://short.io/r/new" {
		t.Fatalf("short_url: %q", updated.ShortUrl)
	}
}

func TestUpdateLink_NotFound(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.UpdateLink(context.Background(), 42, domain.LinkVO{
		OriginalUrl: "https://example.com/x",
		ShortName:   "x",
	})
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeNotFound {
		t.Fatalf("got %#v", err)
	}
}

func TestDeleteLink(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	created, _ := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/del",
		ShortName:   "del",
	})

	if err := svc.DeleteLink(ctx, created.Id); err != nil {
		t.Fatalf("DeleteLink: %v", err)
	}
	_, err := svc.GetLinkByID(ctx, created.Id)
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeNotFound {
		t.Fatalf("got %#v", err)
	}
}

func TestDeleteLink_NotFound(t *testing.T) {
	svc, _ := newTestService()
	err := svc.DeleteLink(context.Background(), 7)
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeNotFound {
		t.Fatalf("got %#v", err)
	}
}

func TestMemRepo_ConflictFromRepository(t *testing.T) {
	repo := testutil.NewMemRepo()
	ctx := context.Background()
	vo := domain.LinkShortenedVO{
		OriginalUrl: "https://example.com",
		ShortName:   "c",
		ShortUrl:    "https://short.io/c",
	}
	if _, err := repo.Create(ctx, vo); err != nil {
		t.Fatal(err)
	}
	_, err := repo.Create(ctx, vo)
	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("got %v", err)
	}
}
