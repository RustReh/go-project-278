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
	if link.ShortUrl != "https://short.io/exmpl" {
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
	if link.ShortUrl != testBaseURL+link.ShortName {
		t.Fatalf("short_url: got %q", link.ShortUrl)
	}
}

func TestCreateLink_RequiresOriginalURL(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.CreateLink(context.Background(), domain.LinkVO{ShortName: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeValidation {
		t.Fatalf("got %#v", err)
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
	if !ok || appErr.Code != apperr.CodeConflict {
		t.Fatalf("got %#v", err)
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

func TestGetAllLinks(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	if _, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/1",
		ShortName:   "a",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/2",
		ShortName:   "b",
	}); err != nil {
		t.Fatal(err)
	}

	links, err := svc.GetAllLinks(ctx)
	if err != nil {
		t.Fatalf("GetAllLinks: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("len: got %d, want 2", len(links))
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
	if updated.ShortUrl != "https://short.io/new" {
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

func TestUpdateLink_RequiresShortName(t *testing.T) {
	svc, _ := newTestService()
	ctx := context.Background()

	created, _ := svc.CreateLink(ctx, domain.LinkVO{
		OriginalUrl: "https://example.com/x",
		ShortName:   "x",
	})

	_, err := svc.UpdateLink(ctx, created.Id, domain.LinkVO{
		OriginalUrl: "https://example.com/y",
	})
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != apperr.CodeValidation {
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
