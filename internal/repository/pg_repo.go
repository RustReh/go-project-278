package repository

import (
	"context"
	"database/sql"

	"github.com/RustReh/go-project-278/internal/db/sqlc"
	"github.com/RustReh/go-project-278/internal/service/domain"
	"github.com/RustReh/go-project-278/internal/service/interfaces"
)

var _ interfaces.Repository = (*PostgresRepo)(nil)

type PostgresRepo struct {
	q *sqlc.Queries
}

func toDomain(dbLink sqlc.Link) domain.Link {
	return domain.Link{
		Id:          dbLink.ID,
		OriginalUrl: dbLink.OriginalUrl,
		ShortName:   dbLink.ShortName,
		ShortUrl:    dbLink.ShortUrl,
	}
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		q: sqlc.New(db),
	}
}

func (repo *PostgresRepo) GetByID(ctx context.Context, id int64) (domain.Link, error) {
	link, err := repo.q.GetLinkByID(ctx, id)
	if err != nil {
		return domain.Link{}, MapError(err)
	}
	return toDomain(link), nil
}

func (repo *PostgresRepo) Count(ctx context.Context) (int64, error) {
	count, err := repo.q.CountLinks(ctx)
	if err != nil {
		return 0, MapError(err)
	}
	return count, nil
}

func (repo *PostgresRepo) List(ctx context.Context, offset, limit int) ([]domain.Link, error) {
	rows, err := repo.q.ListLinks(ctx, sqlc.ListLinksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, MapError(err)
	}

	links := make([]domain.Link, 0, len(rows))
	for _, row := range rows {
		links = append(links, toDomain(row))
	}
	return links, nil
}

func (repo *PostgresRepo) Create(ctx context.Context, vo domain.LinkShortenedVO) (domain.Link, error) {
	link, err := repo.q.CreateLink(ctx, sqlc.CreateLinkParams{
		OriginalUrl: vo.OriginalUrl,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortUrl,
	})
	if err != nil {
		return domain.Link{}, MapError(err)
	}
	return toDomain(link), nil
}

func (repo *PostgresRepo) Update(ctx context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error) {
	link, err := repo.q.UpdateLink(ctx, sqlc.UpdateLinkParams{
		ID:          id,
		OriginalUrl: vo.OriginalUrl,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortUrl,
	})
	if err != nil {
		return domain.Link{}, MapError(err)
	}
	return toDomain(link), nil
}

func (repo *PostgresRepo) Delete(ctx context.Context, id int64) (int, error) {
	rows, err := repo.q.DeleteLink(ctx, id)
	if err != nil {
		return 0, MapError(err)
	}
	return int(rows), nil
}
