package interfaces

import (
	"context"

	"github.com/RustReh/go-project-278/internal/service/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (domain.Link, error)
	GetAll(ctx context.Context) ([]domain.Link, error)
	Create(ctx context.Context, vo domain.LinkShortenedVO) (domain.Link, error)
	Update(ctx context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error)
	Delete(ctx context.Context, id int64) (int, error)
}
