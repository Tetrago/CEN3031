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
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login user
// @Description Log in to user and authenticate with the backend
// @Tags auth
// @Produce json
// @Consume json
// @Success 200
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
	stmt := SELECT(UserAccount.Identifier, UserAccount.Hash).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(request.Email)))

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

	if str, err := MakeToken(TokenContents{Identifier: dest.Identifier}); err != nil {
		fmt.Printf("[/auth/login] Error making token: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.SetCookie("token", str, 86400, "/", Routing.Hostname, false, true)
		g.Status(http.StatusOK)
	}
}

func AuthHandler(r *gin.RouterGroup) {
	r.POST("/login", authLogin)
}
