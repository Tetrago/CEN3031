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

type userGetResponseGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type userGetResponse struct {
	Identifier  string                 `json:"ident"`
	DisplayName string                 `json:"display_name"`
	Groups      []userGetResponseGroup `json:"groups"`
}

// User godoc
// @Summary Fetch user
// @Description Fetches user information and groups
// @Tags user
// @Produce json
// @Success 200 {object} userGetResponse
// @Failure 400
// @Failure 500
// @Param ident path string true "User identifier"
// @Router /user/get/{ident} [get]
func userGet(g *gin.Context) {
	var uri struct {
		Identifier string `uri:"ident" binding:"required"`
	}

	if err := g.ShouldBindUri(&uri); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest struct {
		model.UserAccount

		Rooms []model.Room
	}

	stmt := SELECT(
		UserAccount.Identifier, UserAccount.DisplayName,
		Room.Name,
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
			fmt.Printf("[/user/get] Error querying database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	g.JSON(http.StatusOK, userGetResponse{
		dest.Identifier,
		dest.DisplayName,
		lo.Map(dest.Rooms, func(x model.Room, _ int) userGetResponseGroup {
			return userGetResponseGroup{x.ID, x.Name}
		}),
	})
}

func makeIdentifier() (string, error) {
	var dest struct {
		model.UserAccount
	}

generate:
	ident, err := GenerateBase64(16)
	if err != nil {
		return "", err
	}

	stmt := SELECT(UserAccount.Identifier).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(ident))).LIMIT(1)

	switch err := stmt.Query(Database, &dest); err {
	default:
		return "", err
	case nil:
		goto generate
	case qrm.ErrNoRows:
		return ident, nil
	}
}

type userRegisterRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type userRegisterResponse struct {
	Identifier string `json:"ident"`
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user given the provided arguments
// @Tags user
// @Produce json
// @Consume json
// @Success 200 {object} userRegisterResponse
// @Failure 400
// @Failure 500
// @Param request body userRegisterRequest true "User registration information"
// @Router /user/Register [post]
func userRegister(g *gin.Context) {
	var request userRegisterRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	ident, err := makeIdentifier()
	if err != nil {
		fmt.Printf("[/user/register] Error generating identifier: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	var dest model.UserAccount

	stmt := UserAccount.INSERT(UserAccount.Identifier, UserAccount.DisplayName, UserAccount.Hash, UserAccount.Email).
		MODEL(model.UserAccount{
			Identifier:  ident,
			DisplayName: request.DisplayName,
			Hash:        Hash(request.Password),
			Email:       request.Email,
		}).
		RETURNING(UserAccount.AllColumns)

	if err := stmt.Query(Database, &dest); err != nil {
		fmt.Printf("[/user/register] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	}

	g.JSON(http.StatusOK, userRegisterResponse{dest.Identifier})
}

func UserHandler(r *gin.RouterGroup) {
	r.GET("/user/:ident", userGet)
	r.POST("/register", userRegister)
}
