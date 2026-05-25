package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/service/domain"
	"github.com/RustReh/go-project-278/internal/service/interfaces"
)

const redirectStatus = http.StatusFound

type VisitService struct {
	repo interfaces.Repository
}

func NewVisitService(repo interfaces.Repository) *VisitService {
	return &VisitService{repo: repo}
}

// VisitsPage — страница списка посещений.
type VisitsPage struct {
	Visits []domain.LinkVisit
	Total  int64
	Start  int
	End    int
}

func (s *VisitService) Redirect(
	ctx context.Context,
	code, ip, userAgent, referer string,
) (originalURL string, status int, err error) {
	link, err := s.repo.GetByShortName(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", 0, apperr.NotFound("Link not found")
		}
		return "", 0, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while get link by short name",
			map[string]any{"short_name": code},
			err,
		)
	}

	_, err = s.repo.CreateVisit(ctx, domain.LinkVisitVO{
		LinkId:    link.Id,
		Ip:        ip,
		UserAgent: userAgent,
		Referer:   referer,
		Status:    redirectStatus,
	})
	if err != nil {
		return "", 0, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while create link visit",
			map[string]any{"link_id": link.Id},
			err,
		)
	}

	return link.OriginalUrl, redirectStatus, nil
}

func (s *VisitService) ListVisits(ctx context.Context, start, end int) (VisitsPage, error) {
	limit := end - start
	if limit < 0 {
		return VisitsPage{}, apperr.Validation("invalid range", map[string]any{"start": start, "end": end})
	}

	total, err := s.repo.CountVisits(ctx)
	if err != nil {
		return VisitsPage{}, apperr.WithPayload(
			apperr.CodeInternal,
			"Error while count link visits",
			nil,
			err,
		)
	}

	var visits []domain.LinkVisit
	if start < int(total) && limit > 0 {
		visits, err = s.repo.ListVisits(ctx, start, limit)
		if err != nil {
			return VisitsPage{}, apperr.WithPayload(
				apperr.CodeInternal,
				"Error while list link visits",
				map[string]any{"start": start, "end": end},
				err,
			)
		}
	} else {
		visits = []domain.LinkVisit{}
	}

	return VisitsPage{
		Visits: visits,
		Total:  total,
		Start:  start,
		End:    end,
	}, nil
}
