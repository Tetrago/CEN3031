package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/util"
)

const userIdentifierLength = 16

// User godoc
// @Summary Fetch user
// @Description Fetches user information and groups
// @Tags user
// @Produce json
// @Success 200 {object} UserModel
// @Failure 400
// @Failure 500
// @Param ident path string true "User identifier"
// @Router /v1/user/get/{ident} [get]
func User(g *gin.Context) {
	var request struct {
		Identifier string `uri:"ident" binding:"required"`
	}

	if err := g.ShouldBindUri(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest struct {
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
		UserAccount.Identifier.EQ(String(request.Identifier)),
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

	g.JSON(http.StatusOK, UserModel{
		dest.Identifier,
		dest.DisplayName,
		lo.Map(dest.Rooms, MapToGroupModel),
	})
}

func GetIdentifier() (string, error) {
	var dest struct {
		model.UserAccount
	}

generate:
	ident, err := util.GenerateBase64(userIdentifierLength)
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

type registerRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"hash"`
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user given the provided arguments
// @Tags user
// @Produce json
// @Consume json
// @Success 200 {object} UserModel
// @Failure 400
// @Failure 500
// @Param request body registerRequest true "User registration information"
// @Router /v1/user/register [post]
func Register(g *gin.Context) {
	var request registerRequest

	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	ident, err := GetIdentifier()
	if err != nil {
		g.Status(http.StatusInternalServerError)
		return
	}

	var dest model.UserAccount

	stmt := UserAccount.INSERT(UserAccount.Identifier, UserAccount.DisplayName, UserAccount.Hash, UserAccount.Email).
		MODEL(model.UserAccount{
			Identifier:  ident,
			DisplayName: request.DisplayName,
			Hash:        util.Hash(request.Password),
			Email:       request.Email,
		}).
		RETURNING(UserAccount.AllColumns)

	if err := stmt.Query(Database, &dest); err != nil {
		g.Status(http.StatusInternalServerError)
	}

	g.JSON(http.StatusOK, UserModel{
		dest.Identifier,
		dest.DisplayName,
		[]GroupModel{},
	})
}

func UserHandler(r *gin.RouterGroup) {
	r.GET("/get/:ident", User)
	r.POST("/register", Register)
}
