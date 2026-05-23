package handler

import (
	"net/http"

	"github.com/RustReh/go-project-278/internal/apperr"
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

func (h *LinksHandler) GetAll(c *gin.Context) {
	links, err := h.service.GetAllLinks(c.Request.Context())
	if err != nil {
		writeAppError(c, err)
		return
	}

	resp := make([]schemas.LinkResponse, 0, len(links))
	for _, link := range links {
		resp = append(resp, toLinkResponse(link))
	}
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
	var req schemas.CreateUpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAppError(c, apperr.Validation("invalid JSON", map[string]string{"method": "create_link"}))
		return
	}

	link, err := h.service.CreateLink(c.Request.Context(), linkToVO(req))
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

	var req schemas.CreateUpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAppError(c, apperr.Validation("invalid JSON", map[string]string{"method": "update_link"}))
		return
	}

	link, err := h.service.UpdateLink(c.Request.Context(), id, linkToVO(req))
	if err != nil {
		writeAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, toLinkResponse(link))
}
