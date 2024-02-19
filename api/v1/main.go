package v1

import (
    "net/http"

    "github.com/gin-gonic/gin"
)
// @BasePath /api/v1

// Ping godoc
// @Summary ping
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func Ping(g *gin.Context) {
    g.JSON(http.StatusOK, "pong")
}

func Setup(r *gin.RouterGroup) {
    r.GET("/ping", Ping)
}
