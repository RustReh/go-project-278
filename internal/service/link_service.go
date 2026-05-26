package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/service/domain"
	"github.com/RustReh/go-project-278/internal/service/interfaces"
)

const (
	maxShortURLLen          = 512
	generatedShortNameLen   = 12
	generateShortNameMaxTry = 10
)

const shortNameAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type LinkService struct {
	repo    interfaces.Repository
	baseURL string
}

func NewLinkService(repo interfaces.Repository, baseURL string) *LinkService {
	return &LinkService{
		repo:    repo,
		baseURL: normalizeBaseURL(baseURL),
	}
}

func normalizeBaseURL(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base
}

func (s *LinkService) buildShortURL(shortName string) (string, error) {
	shortURL := s.baseURL + "r/" + strings.TrimSpace(shortName)
	if len(shortURL) > maxShortURLLen {
		return "", apperr.ValidationFields(map[string]string{
			"short_url": "short_url is too long",
		})
	}
	return shortURL, nil
}

func (s *LinkService) toShortenedVO(linkVO domain.LinkVO) (domain.LinkShortenedVO, error) {
	shortURL, err := s.buildShortURL(linkVO.ShortName)
	if err != nil {
		return domain.LinkShortenedVO{}, err
	}
	return domain.LinkShortenedVO{
		OriginalUrl: strings.TrimSpace(linkVO.OriginalUrl),
		ShortName:   strings.TrimSpace(linkVO.ShortName),
		ShortUrl:    shortURL,
	}, nil
}

func generateShortName() (string, error) {
	b := make([]byte, generatedShortNameLen)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(shortNameAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = shortNameAlphabet[n.Int64()]
	}
	return string(b), nil
}

func mapCreateUpdateErr(err error, link domain.Link, shortName string) (domain.Link, error) {
	switch {
	case errors.Is(err, repository.ErrConflict):
		return link, apperr.ValidationFields(map[string]string{
			"short_name": "short name already in use",
		})
	case errors.Is(err, repository.ErrInvalidInput):
		return link, apperr.ValidationFields(map[string]string{
			"short_name": "invalid link data",
		})
	default:
		return link, apperr.WithPayload(
			apperr.CodeInternal,
			"repository error",
			map[string]any{"short_name": shortName},
			err,
		)
	}
}

func (s *LinkService) GetLinkByID(ctx context.Context, id int64) (domain.Link, error) {
	if id <= 0 {
		return domain.Link{}, apperr.ValidationFields(map[string]string{"id": "invalid link id"})
	}

	link, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return link, apperr.NotFound("Link not found")
		}
		return link, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while get link",
			map[string]any{"link_id": id},
			err,
		)
	}

	return link, nil
}

// LinksPage — страница списка ссылок для пагинации.
type LinksPage struct {
	Links []domain.Link
	Total int64
	Start int
	End   int
}

func (s *LinkService) ListLinks(ctx context.Context, start, end int) (LinksPage, error) {
	limit := end - start
	if limit < 0 {
		return LinksPage{}, apperr.ValidationFields(map[string]string{"range": "invalid range"})
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return LinksPage{}, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while count links",
			nil,
			err,
		)
	}

	var links []domain.Link
	if start < int(total) && limit > 0 {
		links, err = s.repo.List(ctx, start, limit)
		if err != nil {
			return LinksPage{}, apperr.WithPayload(
				apperr.CodeInternal,
				"Error while list links",
				map[string]any{"start": start, "end": end},
				err,
			)
		}
	} else {
		links = []domain.Link{}
	}

	return LinksPage{
		Links: links,
		Total: total,
		Start: start,
		End:   end,
	}, nil
}

func (s *LinkService) createWithGeneratedShortName(ctx context.Context, linkVO domain.LinkVO) (domain.Link, error) {
	for range generateShortNameMaxTry {
		shortName, err := generateShortName()
		if err != nil {
			return domain.Link{}, apperr.Internal("failed to generate short_name", err)
		}

		vo := linkVO
		vo.ShortName = shortName
		shortened, err := s.toShortenedVO(vo)
		if err != nil {
			return domain.Link{}, err
		}

		link, err := s.repo.Create(ctx, shortened)
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				continue
			}
			return mapCreateUpdateErr(err, link, shortName)
		}
		return link, nil
	}

	return domain.Link{}, apperr.Internal(
		"failed to generate unique short_name",
		fmt.Errorf("exhausted %d attempts", generateShortNameMaxTry),
	)
}

func (s *LinkService) CreateLink(ctx context.Context, linkVO domain.LinkVO) (domain.Link, error) {
	linkVO.OriginalUrl = strings.TrimSpace(linkVO.OriginalUrl)
	linkVO.ShortName = strings.TrimSpace(linkVO.ShortName)

	if linkVO.ShortName == "" {
		return s.createWithGeneratedShortName(ctx, linkVO)
	}

	shortened, err := s.toShortenedVO(linkVO)
	if err != nil {
		return domain.Link{}, err
	}

	link, err := s.repo.Create(ctx, shortened)
	if err != nil {
		return mapCreateUpdateErr(err, link, linkVO.ShortName)
	}

	return link, nil
}

func (s *LinkService) UpdateLink(ctx context.Context, id int64, linkVO domain.LinkVO) (domain.Link, error) {
	if id <= 0 {
		return domain.Link{}, apperr.ValidationFields(map[string]string{"id": "invalid link id"})
	}
	linkVO.OriginalUrl = strings.TrimSpace(linkVO.OriginalUrl)
	linkVO.ShortName = strings.TrimSpace(linkVO.ShortName)

	shortened, err := s.toShortenedVO(linkVO)
	if err != nil {
		return domain.Link{}, err
	}

	link, err := s.repo.Update(ctx, id, shortened)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return link, apperr.NotFound("Link not found")
		}
		return mapCreateUpdateErr(err, link, linkVO.ShortName)
	}

	return link, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid link id", map[string]any{"link_id": id})
	}

	rows, err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperr.NotFound("Link not found")
		}
		return apperr.WithPayload(
			apperr.CodeInternal,
			"Error while delete link",
			map[string]any{"link_id": id},
			err,
		)
	}
	if rows == 0 {
		return apperr.NotFound("Link not found")
	}

	return nil
}
