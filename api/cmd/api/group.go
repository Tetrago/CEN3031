package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

type groupAllResponseItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// All godoc
// @Summary Gets groups
// @Description Gets all public groups
// @Tags group
// @Produce json
// @Success 200 {array} groupAllResponseItem
// @Failure 500
// @Router /group/all [get]
func groupAll(g *gin.Context) {
	var dest []model.Room
	stmt := SELECT(Room.ID, Room.Name, Room.Description).FROM(Room)

	if err := stmt.Query(Database, &dest); err != nil {
		fmt.Printf("[/group/all] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	g.JSON(http.StatusOK, lo.Map(dest, func(x model.Room, _ int) groupAllResponseItem {
		return groupAllResponseItem{x.ID, x.Name, x.Description}
	}))
}

type groupGetResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Get godoc
// @Summary Get group
// @Description Gets group information
// @Tags group
// @Produce json
// @Success 200 {object} groupGetResponse
// @Failure 400
// @Failure 500
// @Param id path int64 true "Group ID"
// @Router /group/get/{id} [get]
func groupGet(g *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}

	if err := g.ShouldBindUri(&uri); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest model.Room
	stmt := SELECT(Room.ID, Room.Name, Room.Description).FROM(Room).WHERE(Room.ID.EQ(Int64(uri.ID)))

	if err := stmt.Query(Database, &dest); err != nil {
		switch err {
		default:
			fmt.Printf("[/group/get] Error querying database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
			return
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
			return
		}
	}

	g.JSON(http.StatusOK, groupGetResponse{
		dest.ID,
		dest.Name,
		dest.Description,
	})
}

func GroupHandler(r *gin.RouterGroup) {
	r.GET("/all", groupAll)
	r.GET("/get/:id", groupGet)
}
