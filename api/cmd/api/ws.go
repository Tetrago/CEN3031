package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gorilla/websocket"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(user int64, group int64, conn *websocket.Conn) {
	quit := make(chan bool)
	sendMessage := make(chan string)

	go func() {
		for !<-quit {
			t, p, err := conn.ReadMessage()
			if err != nil || t == websocket.CloseMessage {
				quit <- true
			}

			sendMessage <- string(p)
		}
	}()

	go func() {
		for !<-quit {
			message := <-sendMessage

			stmt := RoomMessage.INSERT(
				RoomMessage.UserID,
				RoomMessage.RoomID,
				RoomMessage.Contents,
				RoomMessage.Iat,
			).MODEL(model.RoomMessage{
				UserID:   user,
				RoomID:   group,
				Contents: message,
				Iat:      time.Now().Unix(),
			})

			if _, err := stmt.Exec(Database); err != nil {
				quit <- true
			}
		}
	}()
}

// WebSocket godoc
// @Summary Opens a WebSocket
// @Description Opens a WebSocket for a user on a group
// @Tags ws
// @Consume json
// @Success 202
// @Failure 400
// @Failure 502
// @Param token query string true "Authentication token"
// @Param group query int64 true "Active group ID"
// @Router /ws [get]
func ws(g *gin.Context) {
	var request struct {
		Token   string `json:"token"`
		GroupID int64  `json:"group"`
	}

	if err := g.BindQuery(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	token, err := ProcessToken(request.Token)
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var user model.UserAccount
	if err := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.Ident))).Query(Database, &user); err != nil {
		switch err {
		default:
			fmt.Printf("[/ws] Failed to query database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}
	}

	conn, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		g.Status(http.StatusBadGateway)
		return
	}

	go websocketHandler(user.ID, request.GroupID, conn)
	g.Status(http.StatusAccepted)
}

func WSHandler(r *gin.RouterGroup) {
	r.GET("/", ws)
}
