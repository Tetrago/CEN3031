package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tetrago/motmot/api/internal/globals"
)

func Middleware() func(*gin.Context) {
	return func(c *gin.Context) {
		if raw, err := c.Cookie("token"); err != nil {
			c.SetCookie("token", "", -1, "/", globals.Opts.Hostname, globals.Opts.SslEnabled, true)
			c.AbortWithStatus(http.StatusUnauthorized)
		} else if token, err := ParseToken(raw); err != nil {
			c.SetCookie("token", "", -1, "/", globals.Opts.Hostname, globals.Opts.SslEnabled, true)
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.Set("token", token)
			c.Next()
		}
	}
}

func ExpectToken(c *gin.Context) *Token {
	if v, ok := c.Get("token"); !ok {
		panic("missing authentication middleware where expected")
	} else if token, ok := v.(*Token); !ok {
		panic("invalid authentication middleware; unexpected type")
	} else {
		return token
	}
}
