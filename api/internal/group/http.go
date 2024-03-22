package group

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/internal/globals"
)

type AllResponseItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// All godoc
// @Summary Gets groups
// @Description Gets all public groups
// @Tags group
// @Produce json
// @Success 200 {array} AllResponseItem
// @Failure 500
// @Router /group/all [get]
func All(c *gin.Context) {
	var dest []model.Room
	stmt := SELECT(Room.ID, Room.Name, Room.Description).FROM(Room)

	if err := stmt.Query(globals.Database, &dest); err != nil {
		fmt.Printf("[/group/all] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, lo.Map(dest, func(x model.Room, _ int) AllResponseItem {
		return AllResponseItem{x.ID, x.Name, x.Description}
	}))
}

type GetResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Get godoc
// @Summary Get group
// @Description Gets group information
// @Tags group
// @Produce json
// @Success 200 {object} GetResponse
// @Failure 400
// @Failure 500
// @Param id path int64 true "Group ID"
// @Router /group/get/{id} [get]
func Get(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var dest model.Room
	stmt := SELECT(Room.ID, Room.Name, Room.Description).FROM(Room).WHERE(Room.ID.EQ(Int64(uri.ID)))

	if err := stmt.Query(globals.Database, &dest); err == qrm.ErrNoRows {
		c.Status(http.StatusBadRequest)
	} else if err != nil {
		fmt.Printf("[/group/get] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, GetResponse{
			dest.ID,
			dest.Name,
			dest.Description,
		})
	}
}

type HistoryResponseItem struct {
	ID         int64  `json:"message_id"`
	Identifier string `json:"user_ident"`
	Contents   string `json:"contents"`
	IssuedAt   int64  `json:"iat"`
}

// History godoc
// @Summary Gets group messages
// @Description Gets message history from a group in descending order
// @Tags group
// @Produce json
// @Success 200 {array} HistoryResponseItem
// @Failure 500
// @Param id     path  int64 true "Group ID"
// @Param limit  query int64 true "Max number of messages to retreive (<= 20)"
// @Param before query int64 true "UTC time cutoff; searches in reverse from this point (inclusive)"
// @Router /group/history/{id} [get]
func History(c *gin.Context) {
	var uri struct {
		ID int64 `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var request struct {
		Limit  int64 `form:"limit" binding:"required"`
		Before int64 `form:"before" binding:"required"`
	}

	if err := c.BindQuery(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if request.Limit > 20 {
		request.Limit = 20
	}

	stmt := SELECT(
		RoomMessage.ID, RoomMessage.Contents, RoomMessage.Iat,
		UserAccount.Identifier,
	).FROM(
		RoomMessage.INNER_JOIN(UserAccount, RoomMessage.UserID.EQ(UserAccount.ID)),
	).WHERE(
		RoomMessage.RoomID.EQ(Int64(uri.ID)).AND(RoomMessage.Iat.LT_EQ(Int64(request.Before))),
	).ORDER_BY(RoomMessage.Iat.DESC()).LIMIT(request.Limit)

	var dest []struct {
		model.RoomMessage

		User model.UserAccount
	}

	if err := stmt.Query(globals.Database, &dest); err != nil && err != qrm.ErrNoRows {
		fmt.Printf("[/group/history] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, lo.Map(dest, func(x struct {
		model.RoomMessage
		User model.UserAccount
	}, _ int) HistoryResponseItem {
		return HistoryResponseItem{
			x.ID,
			x.User.Identifier,
			x.Contents,
			x.Iat,
		}
	}))
}

type PopularResponseItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Popular godoc
// @Summary Gets popular groups
// @Description Gets the most popular groups by member count
// @Tags group
// @Produce json
// @Success 200 {array} PopularResponseItem
// @Failure 500
// @Param count path int64 true "Count of groups to return"
// @Router /group/popular/{count} [get]
func Popular(c *gin.Context) {
	var uri struct {
		Count int64 `uri:"count" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var dest []model.Room
	stmt := SELECT(
		Room.ID,
		Room.Name,
	).FROM(
		Room.
			LEFT_JOIN(UserRoom, Room.ID.EQ(UserRoom.RoomID)),
	).GROUP_BY(Room.ID).ORDER_BY(COUNT(Room.ID).DESC()).LIMIT(uri.Count)

	if err := stmt.Query(globals.Database, &dest); err != nil && err != qrm.ErrNoRows {
		fmt.Printf("[/group/popular] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, lo.Map(dest, func(x model.Room, _ int) PopularResponseItem {
		return PopularResponseItem{
			x.ID,
			x.Name,
		}
	}))
}

func HttpHandler(r *gin.RouterGroup) {
	g := r.Group("/group")
	g.GET("/all", All)
	g.GET("/get/:id", Get)
	g.GET("/history/:id", History)
	g.GET("/popular/:count", Popular)
}
