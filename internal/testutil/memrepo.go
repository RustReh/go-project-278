package testutil

import (
	"context"
	"sort"

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

func (m *MemRepo) sortedLinks() []domain.Link {
	out := make([]domain.Link, 0, len(m.links))
	for _, link := range m.links {
		out = append(out, link)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })
	return out
}

func (m *MemRepo) GetByID(_ context.Context, id int64) (domain.Link, error) {
	link, ok := m.links[id]
	if !ok {
		return domain.Link{}, repository.ErrNotFound
	}
	return link, nil
}

func (m *MemRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.links)), nil
}

func (m *MemRepo) List(_ context.Context, offset, limit int) ([]domain.Link, error) {
	all := m.sortedLinks()
	if offset >= len(all) || limit <= 0 {
		return []domain.Link{}, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (m *MemRepo) Create(_ context.Context, vo domain.LinkShortenedVO) (domain.Link, error) {
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

func (m *MemRepo) Update(_ context.Context, id int64, vo domain.LinkShortenedVO) (domain.Link, error) {
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

func (m *MemRepo) Delete(_ context.Context, id int64) (int, error) {
	if _, ok := m.links[id]; !ok {
		return 0, repository.ErrNotFound
	}
	delete(m.links, id)
	return 1, nil
}
