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
		ID:          dbLink.ID,
		OriginalURL: dbLink.OriginalUrl,
		ShortName:   dbLink.ShortName,
		ShortURL:    dbLink.ShortUrl,
	}
}

func visitToDomain(v sqlc.LinkVisit) domain.LinkVisit {
	return domain.LinkVisit{
		Id:        v.ID,
		LinkId:    v.LinkID,
		Ip:        v.Ip,
		UserAgent: v.UserAgent,
		Referer:   v.Referer,
		Status:    int(v.Status),
		CreatedAt: v.CreatedAt,
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

func (repo *PostgresRepo) GetByShortName(ctx context.Context, shortName string) (domain.Link, error) {
	link, err := repo.q.GetLinkByShortName(ctx, shortName)
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
		OriginalUrl: vo.OriginalURL,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortURL,
	})
	if err != nil {
		return domain.Link{}, MapError(err)
	}
	return toDomain(link), nil
}

func (repo *PostgresRepo) Update(ctx context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error) {
	link, err := repo.q.UpdateLink(ctx, sqlc.UpdateLinkParams{
		ID:          id,
		OriginalUrl: vo.OriginalURL,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortURL,
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

func (repo *PostgresRepo) CreateVisit(ctx context.Context, vo domain.LinkVisitVO) (domain.LinkVisit, error) {
	visit, err := repo.q.CreateLinkVisit(ctx, sqlc.CreateLinkVisitParams{
		LinkID:    vo.LinkId,
		Ip:        vo.Ip,
		UserAgent: vo.UserAgent,
		Referer:   vo.Referer,
		Status:    int32(vo.Status),
	})
	if err != nil {
		return domain.LinkVisit{}, MapError(err)
	}
	return visitToDomain(visit), nil
}

func (repo *PostgresRepo) CountVisits(ctx context.Context) (int64, error) {
	count, err := repo.q.CountLinkVisits(ctx)
	if err != nil {
		return 0, MapError(err)
	}
	return count, nil
}

func (repo *PostgresRepo) ListVisits(ctx context.Context, offset, limit int) ([]domain.LinkVisit, error) {
	rows, err := repo.q.ListLinkVisits(ctx, sqlc.ListLinkVisitsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, MapError(err)
	}

	visits := make([]domain.LinkVisit, 0, len(rows))
	for _, row := range rows {
		visits = append(visits, visitToDomain(row))
	}
	return visits, nil
}
