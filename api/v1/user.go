package v1

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

type UserUri struct {
	Identifier string `uri:"ident" binding:"required"`
}

var Database *sql.DB

// Groups godoc
// @Summary Fetchs user groups
// @Description Fetchs groups a user belongs to
// @Tags user
// @Produce json
// @Success 200 {array} GroupModel
// @Param ident path string true "User identifier"
// @Router /v1/user/{ident} [get]
func User(g *gin.Context) {
	var uri UserUri
	if err := g.ShouldBindUri(&uri); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest []struct {
		model.UserAccount

		Rooms []model.Room
	}

	stmt := SELECT(
		UserAccount.AllColumns,
		Room.AllColumns,
	).FROM(
		UserAccount.
			LEFT_JOIN(UserRoom, UserAccount.ID.EQ(UserRoom.UserID)).
			LEFT_JOIN(Room, UserRoom.RoomID.EQ(Room.ID)),
	).WHERE(
		UserAccount.Identifier.EQ(String(uri.Identifier)),
	)

	if err := stmt.Query(Database, &dest); err != nil {
		switch err {
		default:
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	fmt.Println(dest)

	g.JSON(http.StatusOK, UserModel{
		dest[0].Identifier,
		dest[0].DisplayName,
		lo.Map(dest[0].Rooms, func(x model.Room, _ int) GroupModel {
			return GroupModel{x.ID, x.Name}
		}),
	})
}

func UserHandler(r *gin.RouterGroup) {
	r.GET("/", User)
}
