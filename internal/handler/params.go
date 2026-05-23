package handler

import (
	"strconv"

	"github.com/RustReh/go-project-278/internal/apperr"
	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/RustReh/go-project-278/internal/service/domain"
)

func parsePathInt64(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, apperr.Validation("invalid id", map[string]any{"id": raw})
	}
	return id, nil
}

func toLinkResponse(link domain.Link) schemas.LinkResponse {
	return schemas.LinkResponse{
		ID:          link.Id,
		OriginalURL: link.OriginalUrl,
		ShortName:   link.ShortName,
		ShortURL:    link.ShortUrl,
	}
}

func linkToVO(l schemas.CreateUpdateLinkRequest) domain.LinkVO {
	return domain.LinkVO{
		ShortName:   l.ShortName,
		OriginalUrl: l.OriginalURL,
	}
}
