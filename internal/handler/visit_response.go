package handler

import (
	"time"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

func toLinkVisitResponse(v domain.LinkVisit) schemas.LinkVisitResponse {
	return schemas.LinkVisitResponse{
		ID:        v.Id,
		LinkID:    v.LinkId,
		CreatedAt: v.CreatedAt.UTC().Format(time.RFC3339),
		IP:        v.Ip,
		UserAgent: v.UserAgent,
		Status:    v.Status,
	}
}
