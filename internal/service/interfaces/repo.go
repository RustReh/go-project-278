package interfaces

import (
	"context"

	"github.com/RustReh/go-project-278/internal/service/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (domain.Link, error)
	GetByShortName(ctx context.Context, shortName string) (domain.Link, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]domain.Link, error)
	Create(ctx context.Context, vo domain.LinkShortenedVO) (domain.Link, error)
	Update(ctx context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error)
	Delete(ctx context.Context, id int64) (int, error)

	CreateVisit(ctx context.Context, vo domain.LinkVisitVO) (domain.LinkVisit, error)
	CountVisits(ctx context.Context) (int64, error)
	ListVisits(ctx context.Context, offset, limit int) ([]domain.LinkVisit, error)
}
