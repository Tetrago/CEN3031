package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/internal/crypt"
	"github.com/tetrago/motmot/api/internal/globals"
)

type LoginRequest struct {
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
// @Param request body LoginRequest true "User login information"
// @Router /auth/login [post]
func Login(g *gin.Context) {
	var request LoginRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var dest model.UserAccount
	stmt := SELECT(UserAccount.Identifier, UserAccount.Hash).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(request.Email)))

	if err := stmt.Query(globals.Database, &dest); err == qrm.ErrNoRows {
		g.Status(http.StatusBadRequest)
		return
	} else if err != nil {
		fmt.Printf("[/auth/login] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	if crypt.Hash(request.Password) != dest.Hash {
		g.Status(http.StatusBadRequest)
		return
	}

	if raw, err := NewToken().SetUserIdentifier(dest.Identifier).Serialize(); err != nil {
		fmt.Printf("[/auth/login] Error making token: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.SetCookie("token", raw, 86400, "/", globals.Opts.Hostname, globals.Opts.SslEnabled, true)
		g.Status(http.StatusOK)
	}
}

func HttpHandler(r *gin.RouterGroup) {
	g := r.Group("/auth")
	g.POST("/login", Login)
}
