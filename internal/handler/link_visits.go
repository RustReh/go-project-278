package handler

import (
	"net/http"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkVisitsHandler struct {
	service *service.VisitService
}

func NewLinkVisitsHandler(visitService *service.VisitService) *LinkVisitsHandler {
	return &LinkVisitsHandler{service: visitService}
}

// GetAll — GET /api/link_visits
func (h *LinkVisitsHandler) GetAll(c *gin.Context) {
	start, end, err := parseListRange(c.Query("range"), c.GetHeader("Range"))
	if err != nil {
		writeAppError(c, err)
		return
	}

	page, err := h.service.ListVisits(c.Request.Context(), start, end)
	if err != nil {
		writeAppError(c, err)
		return
	}

	resp := make([]schemas.LinkVisitResponse, 0, len(page.Visits))
	for _, v := range page.Visits {
		resp = append(resp, toLinkVisitResponse(v))
	}

	c.Header("Content-Range", contentRangeHeader("link_visits", page.Start, page.End, page.Total))
	c.JSON(http.StatusOK, resp)
}
