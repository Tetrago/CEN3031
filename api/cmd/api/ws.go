package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type message struct {
	Identifier string `json:"user_ident"`
	Contents   string `json:"contents"`
	IssuedAt   int64  `json:"iat"`
}

var channels = make(map[int64][]chan<- message)

func wsHandler(conn *websocket.Conn) {
	defer conn.Close()

	var handshake struct {
		Token   string `form:"token" binding:"required"`
		GroupID int64  `form:"group" binding:"required"`
	}

	if err := conn.ReadJSON(&handshake); err != nil {
		return
	}

	token, err := ProcessToken(handshake.Token)
	if err != nil {
		return
	}

	var user model.UserAccount
	if err := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.Ident))).Query(Database, &user); err != nil {
		if err != qrm.ErrNoRows {
			fmt.Printf("[/ws] Failed to query database: %s\n", err.Error())
		}

		return
	}

	quit := make(chan int)
	recv := make(chan message)

	if _, ok := channels[handshake.GroupID]; !ok {
		channels[handshake.GroupID] = []chan<- message{recv}
	} else {
		channels[handshake.GroupID] = append(channels[handshake.GroupID], recv)
	}

	go func() {
		for {
			select {
			case <-quit:
				return
			case message := <-recv:
				if err := conn.WriteJSON(message); err != nil {
					quit <- 0
					return
				}
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		default:
			t, p, err := conn.ReadMessage()
			if err != nil || t == websocket.CloseMessage {
				quit <- 0
				break loop
			}

			contents := string(p)
			iat := time.Now().Unix()

			go func() {
				for _, v := range channels[handshake.GroupID] {
					if v != recv {
						v <- message{
							token.Ident,
							contents,
							iat,
						}
					}
				}
			}()

			stmt := RoomMessage.INSERT(
				RoomMessage.UserID,
				RoomMessage.RoomID,
				RoomMessage.Contents,
				RoomMessage.Iat,
			).MODEL(model.RoomMessage{
				UserID:   user.ID,
				RoomID:   handshake.GroupID,
				Contents: contents,
				Iat:      iat,
			})

			if _, err := stmt.Exec(Database); err != nil {
				quit <- 0
				break loop
			}
		}
	}

	if _, index, ok := lo.FindIndexOf(channels[handshake.GroupID], func(x chan<- message) bool { return x == recv }); !ok {
		panic("Channel bus mismatch!")
	} else {
		channels[handshake.GroupID] = append(channels[handshake.GroupID][:index], channels[handshake.GroupID][index+1:]...)
	}
}

// WebSocket godoc
// @Summary Opens a WebSocket
// @Description Opens a WebSocket for a user on a group
// @Tags ws
// @Consume json
// @Failure 400
// @Failure 502
// @Router /ws [get]
func wsGet(g *gin.Context) {
	conn, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		return
	}

	go wsHandler(conn)
}
