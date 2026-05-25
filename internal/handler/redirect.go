package handler

import (
	"github.com/RustReh/go-project-278/internal/service"
	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	service *service.VisitService
}

func NewRedirectHandler(visitService *service.VisitService) *RedirectHandler {
	return &RedirectHandler{service: visitService}
}

// Redirect — GET /r/:code
func (h *RedirectHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	url, status, err := h.service.Redirect(
		c.Request.Context(),
		code,
		c.ClientIP(),
		c.Request.UserAgent(),
		c.Request.Referer(),
	)
	if err != nil {
		writeAppError(c, err)
		return
	}

	c.Redirect(status, url)
}
