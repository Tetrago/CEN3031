package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/util"
)

type loginRequest struct {
	Ident    string `json:"ident"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login user
// @Description Loging to user and authenticate with the backend
// @Tags auth
// @Produce json
// @Consume json
// @Success 200 {object} loginResponse
// @Failure 400
// @Failure 500
// @Param request body loginRequest true "User login information"
// @Router /v1/auth/login [post]
func Login(g *gin.Context) {
	var request loginRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest model.UserAccount
	stmt := SELECT(UserAccount.Identifier, UserAccount.Hash).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(request.Ident)))

	if err := stmt.Query(Database, &dest); err != nil {
		switch err {
		default:
			fmt.Printf("[/v1/auth/login] Error querying database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	if util.Hash(request.Password) != dest.Hash {
		g.Status(http.StatusBadRequest)
		return
	}

	if str, err := util.MakeToken(util.TokenContents{Ident: request.Ident}); err != nil {
		fmt.Printf("[/v1/auth/login] Error making token: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.JSON(http.StatusOK, loginResponse{str})
	}
}

type renewStruct struct {
	Token string `json:"token"`
}

// Renew godoc
// @Summary Renew token
// @Description Renews token, preventing timeouts
// @Tags auth
// @Produce json
// @Consume json
// @Success 200 {object} renewStruct
// @Failure 400
// @Failure 500
// @Param request body renewStruct true "Token to renew"
// @Router /v1/auth/renew [post]
func Renew(g *gin.Context) {
	var request renewStruct
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	if contents, err := util.ProcessToken(request.Token); err != nil {
		g.Status(http.StatusBadRequest)
	} else {
		if str, err := util.MakeToken(*contents); err != nil {
			fmt.Printf("[/v1/auth/renew] Error making token: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		} else {
			g.JSON(http.StatusOK, renewStruct{str})
		}
	}
}

func AuthHandler(r *gin.RouterGroup) {
	r.POST("/login", Login)
	r.POST("/renew", Renew)
}
