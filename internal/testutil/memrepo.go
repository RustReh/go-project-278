package testutil

import (
	"context"
	"sort"
	"time"

	"github.com/RustReh/go-project-278/internal/repository"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

// MemRepo — in-memory реализация Repository для тестов (один поток, без блокировок).
type MemRepo struct {
	links    map[int64]domain.Link
	visits   map[int64]domain.LinkVisit
	seq      int64
	visitSeq int64
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		links:  make(map[int64]domain.Link),
		visits: make(map[int64]domain.LinkVisit),
	}
}

func (m *MemRepo) sortedLinks() []domain.Link {
	out := make([]domain.Link, 0, len(m.links))
	for _, link := range m.links {
		out = append(out, link)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })
	return out
}

func (m *MemRepo) sortedVisits() []domain.LinkVisit {
	out := make([]domain.LinkVisit, 0, len(m.visits))
	for _, v := range m.visits {
		out = append(out, v)
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

func (m *MemRepo) GetByShortName(_ context.Context, shortName string) (domain.Link, error) {
	for _, link := range m.links {
		if link.ShortName == shortName {
			return link, nil
		}
	}
	return domain.Link{}, repository.ErrNotFound
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

func (m *MemRepo) CreateVisit(_ context.Context, vo domain.LinkVisitVO) (domain.LinkVisit, error) {
	m.visitSeq++
	visit := domain.LinkVisit{
		Id:        m.visitSeq,
		LinkId:    vo.LinkId,
		Ip:        vo.Ip,
		UserAgent: vo.UserAgent,
		Referer:   vo.Referer,
		Status:    vo.Status,
		CreatedAt: time.Now().UTC(),
	}
	m.visits[visit.Id] = visit
	return visit, nil
}

func (m *MemRepo) CountVisits(_ context.Context) (int64, error) {
	return int64(len(m.visits)), nil
}

func (m *MemRepo) ListVisits(_ context.Context, offset, limit int) ([]domain.LinkVisit, error) {
	all := m.sortedVisits()
	if offset >= len(all) || limit <= 0 {
		return []domain.LinkVisit{}, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}
