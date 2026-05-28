package handler

import (
	"net/http"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/gin-gonic/gin"
)

type LinksHandler struct {
	service *service.LinkService
}

func NewLinksHandler(linkService *service.LinkService) *LinksHandler {
	return &LinksHandler{service: linkService}
}

// GetAll — GET /api/links?range=[start,end]
func (h *LinksHandler) GetAll(c *gin.Context) {
	start, end, err := parseListRange(c.Query("range"), c.GetHeader("Range"))
	if err != nil {
		writeAppError(c, err)
		return
	}

	page, err := h.service.ListLinks(c.Request.Context(), start, end)
	if err != nil {
		writeAppError(c, err)
		return
	}

	resp := make([]schemas.LinkResponse, 0, len(page.Links))
	for _, link := range page.Links {
		resp = append(resp, toLinkResponse(link))
	}

	c.Header("Content-Range", contentRangeHeader("links", page.Start, page.End, page.Total))
	c.JSON(http.StatusOK, resp)
}

// Delete — DELETE /api/links/:id
func (h *LinksHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeAppError(c, err)
		return
	}

	if err := h.service.DeleteLink(c.Request.Context(), id); err != nil {
		writeAppError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByID — GET /api/links/:id
func (h *LinksHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeAppError(c, err)
		return
	}

	link, err := h.service.GetLinkByID(c.Request.Context(), id)
	if err != nil {
		writeAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, toLinkResponse(link))
}

// Create — POST /api/links
func (h *LinksHandler) Create(c *gin.Context) {
	var req createLinkPayload
	if !bindJSON(c, &req) {
		return
	}

	link, err := h.service.CreateLink(c.Request.Context(), linkVOFromCreate(req))
	if err != nil {
		writeAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toLinkResponse(link))
}

// Update — PUT /api/links/:id
func (h *LinksHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		writeAppError(c, err)
		return
	}

	var req updateLinkPayload
	if !bindJSON(c, &req) {
		return
	}

	link, err := h.service.UpdateLink(c.Request.Context(), id, linkVOFromUpdate(req))
	if err != nil {
		writeAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, toLinkResponse(link))
}
