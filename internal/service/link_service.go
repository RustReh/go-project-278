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
	maxOriginalURLLen       = 2048
	maxShortNameLen         = 64
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
	shortURL := s.baseURL + strings.TrimSpace(shortName)
	if len(shortURL) > maxShortURLLen {
		return "", apperr.Validation(
			"short_url is too long",
			map[string]any{
				"max_length": maxShortURLLen,
				"short_url":  shortURL,
			},
		)
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

func validateCreateLinkVO(vo domain.LinkVO) error {
	original := strings.TrimSpace(vo.OriginalUrl)
	shortName := strings.TrimSpace(vo.ShortName)

	if original == "" {
		return apperr.Validation(
			"original_url is required",
			map[string]any{"original_url": vo.OriginalUrl},
		)
	}
	if len(original) > maxOriginalURLLen {
		return apperr.Validation(
			"original_url is too long",
			map[string]any{"max_length": maxOriginalURLLen},
		)
	}
	if shortName != "" && len(shortName) > maxShortNameLen {
		return apperr.Validation(
			"short_name is too long",
			map[string]any{"max_length": maxShortNameLen},
		)
	}
	return nil
}

func validateUpdateLinkVO(vo domain.LinkVO) error {
	original := strings.TrimSpace(vo.OriginalUrl)
	shortName := strings.TrimSpace(vo.ShortName)

	if original == "" || shortName == "" {
		return apperr.Validation(
			"original_url and short_name are required",
			map[string]any{
				"original_url": vo.OriginalUrl,
				"short_name":   vo.ShortName,
			},
		)
	}
	if len(original) > maxOriginalURLLen {
		return apperr.Validation(
			"original_url is too long",
			map[string]any{"max_length": maxOriginalURLLen},
		)
	}
	if len(shortName) > maxShortNameLen {
		return apperr.Validation(
			"short_name is too long",
			map[string]any{"max_length": maxShortNameLen},
		)
	}
	return nil
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
		return link, apperr.Conflict("Link with this short_name already exists")
	case errors.Is(err, repository.ErrInvalidInput):
		return link, apperr.Validation("Invalid link data", map[string]any{
			"short_name": shortName,
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
		return domain.Link{}, apperr.Validation("invalid link id", map[string]any{"link_id": id})
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

func (s *LinkService) GetAllLinks(ctx context.Context) ([]domain.Link, error) {
	links, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while list links",
			nil,
			err,
		)
	}
	return links, nil
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
	if err := validateCreateLinkVO(linkVO); err != nil {
		return domain.Link{}, err
	}

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
		return domain.Link{}, apperr.Validation("invalid link id", map[string]any{"link_id": id})
	}
	if err := validateUpdateLinkVO(linkVO); err != nil {
		return domain.Link{}, err
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
