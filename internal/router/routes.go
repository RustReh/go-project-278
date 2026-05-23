package router

import (
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, links *handler.LinksHandler) {
	r.GET("/ping", handler.Ping)

	api := r.Group("/api")
	linksGroup := api.Group("/links")
	{
		linksGroup.GET("", links.GetAll)
		linksGroup.GET("/:id", links.GetByID)
		linksGroup.DELETE("/:id", links.Delete)
		linksGroup.POST("", links.Create)
		linksGroup.PUT("/:id", links.Update)
	}
}
