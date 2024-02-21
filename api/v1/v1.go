package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var Database *sql.DB

func Handler(r *gin.RouterGroup) {
	AuthHandler(r.Group("/auth"))
	GroupHandler(r.Group("/group"))
	UserHandler(r.Group("/user"))
}
