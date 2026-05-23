package testutil

import (
	"context"

	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

// MemRepo — in-memory реализация Repository для тестов (один поток, без блокировок).
type MemRepo struct {
	links map[int64]domain.Link
	seq   int64
}

func NewMemRepo() *MemRepo {
	return &MemRepo{links: make(map[int64]domain.Link)}
}

func (m *MemRepo) GetByID(ctx context.Context, id int64) (domain.Link, error) {
	link, ok := m.links[id]
	if !ok {
		return domain.Link{}, repository.ErrNotFound
	}
	return link, nil
}

func (m *MemRepo) GetAll(ctx context.Context) ([]domain.Link, error) {
	out := make([]domain.Link, 0, len(m.links))
	for _, link := range m.links {
		out = append(out, link)
	}
	return out, nil
}

func (m *MemRepo) Create(ctx context.Context, vo domain.LinkShortenedVO) (domain.Link, error) {
	for _, link := range m.links {
		if link.ShortName == vo.ShortName {
			return domain.Link{}, repository.ErrConflict
		}
	}

	m.seq++
	link := domain.Link{
		Id:          m.seq,
		OriginalUrl: vo.OriginalUrl,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortUrl,
	}
	m.links[link.Id] = link
	return link, nil
}

func (m *MemRepo) Update(ctx context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error) {
	if _, ok := m.links[id]; !ok {
		return domain.Link{}, repository.ErrNotFound
	}

	for _, link := range m.links {
		if link.Id != id && link.ShortName == vo.ShortName {
			return domain.Link{}, repository.ErrConflict
		}
	}

	link := domain.Link{
		Id:          id,
		OriginalUrl: vo.OriginalUrl,
		ShortName:   vo.ShortName,
		ShortUrl:    vo.ShortUrl,
	}
	m.links[id] = link
	return link, nil
}

func (m *MemRepo) Delete(ctx context.Context, id int64) (int, error) {
	if _, ok := m.links[id]; !ok {
		return 0, repository.ErrNotFound
	}
	delete(m.links, id)
	return 1, nil
}
