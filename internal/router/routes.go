package router

import (
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	links *handler.LinksHandler,
	visits *handler.LinkVisitsHandler,
	redirect *handler.RedirectHandler,
) {
	r.TrustedPlatform = gin.PlatformCloudflare

	r.GET("/ping", handler.Ping)
	r.GET("/r/:code", redirect.Redirect)

	api := r.Group("/api")
	linksGroup := api.Group("/links")
	{
		linksGroup.GET("", links.GetAll)
		linksGroup.GET("/:id", links.GetByID)
		linksGroup.DELETE("/:id", links.Delete)
		linksGroup.POST("", links.Create)
		linksGroup.PUT("/:id", links.Update)
	}
	api.GET("/link_visits", visits.GetAll)
}
