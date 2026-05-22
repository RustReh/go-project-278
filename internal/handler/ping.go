package handler

import (
	"net/http"

	"github.com/RustReh/go-project-278/internal/schemas"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	response := schemas.PingResponse{
		Message: "pong",
		Status:  http.StatusOK,
	}
	c.String(response.Status, response.Message)
}
