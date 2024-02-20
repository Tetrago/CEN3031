package v1

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

type UserUri struct {
	ID int64 `uri:"id" binding:"required"`
}

var Database *sql.DB

// Groups godoc
// @Summary Fetchs user groups
// @Description Fetchs groups a user belongs to
// @Tags user
// @Produce json
// @Success 200 {array} GroupModel
// @Param id path int true "User ID"
// @Router /v1/user/{id} [get]
func User(g *gin.Context) {
	var user UserUri
	if err := g.ShouldBindUri(&user); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest []struct {
		model.UserAccount
	}

	if err := SELECT(UserAccount.DisplayName).FROM(UserAccount.Table).WHERE(UserAccount.ID.EQ(Int64(user.ID))).Query(Database, &dest); err != nil {
		g.Status(http.StatusInternalServerError)
		return
	}

	if len(dest) != 1 {
		g.Status(http.StatusBadRequest)
		return
	}

	g.JSON(http.StatusOK, UserModel{
		dest[0].ID,
		dest[0].DisplayName,
		[]GroupModel{},
	})
}

func UserHandler(r *gin.RouterGroup) {
	r.GET("/", User)
}
