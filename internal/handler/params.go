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
		return 0, apperr.ValidationFields(map[string]string{"id": "invalid id"})
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

func linkVOFromCreate(p createLinkPayload) domain.LinkVO {
	return domain.LinkVO{
		OriginalUrl: p.OriginalURL,
		ShortName:   p.ShortName,
	}
}

func linkVOFromUpdate(p updateLinkPayload) domain.LinkVO {
	return domain.LinkVO{
		OriginalUrl: p.OriginalURL,
		ShortName:   p.ShortName,
	}
}
