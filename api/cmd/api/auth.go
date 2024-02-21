package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

type authLoginRequest struct {
	Ident    string `json:"ident"`
	Password string `json:"password"`
}

type authLoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login user
// @Description Log in to user and authenticate with the backend
// @Tags auth
// @Produce json
// @Consume json
// @Success 200 {object} authLoginResponse
// @Failure 400
// @Failure 500
// @Param request body authLoginRequest true "User login information"
// @Router /auth/login [post]
func authLogin(g *gin.Context) {
	var request authLoginRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest model.UserAccount
	stmt := SELECT(UserAccount.Identifier, UserAccount.Hash).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(request.Ident)))

	if err := stmt.Query(Database, &dest); err != nil {
		switch err {
		default:
			fmt.Printf("[/auth/login] Error querying database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	if Hash(request.Password) != dest.Hash {
		g.Status(http.StatusBadRequest)
		return
	}

	if str, err := MakeToken(TokenContents{Ident: request.Ident}); err != nil {
		fmt.Printf("[/auth/login] Error making token: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.JSON(http.StatusOK, authLoginResponse{str})
	}
}

type authRenewRequest struct {
	Token string `json:"token"`
}

type authRenewResponse struct {
	Token string `json:"token"`
}

// Renew godoc
// @Summary Renew token
// @Description Renews token, preventing timeouts
// @Tags auth
// @Produce json
// @Consume json
// @Success 200 {object} authRenewResponse
// @Failure 400
// @Failure 500
// @Param request body authRenewRequest true "Token to renew"
// @Router /auth/renew [post]
func authRenew(g *gin.Context) {
	var request authRenewRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	if contents, err := ProcessToken(request.Token); err != nil {
		g.Status(http.StatusBadRequest)
	} else {
		if str, err := MakeToken(*contents); err != nil {
			fmt.Printf("[/auth/renew] Error making token: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		} else {
			g.JSON(http.StatusOK, authRenewResponse{str})
		}
	}
}

func AuthHandler(r *gin.RouterGroup) {
	r.POST("/login", authLogin)
	r.POST("/renew", authRenew)
}
