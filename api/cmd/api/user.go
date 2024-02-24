package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	Bio         *string                `json:"bio,omitempty"`
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
		UserAccount.Identifier, UserAccount.DisplayName, UserAccount.Bio,
		Room.ID, Room.Name,
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
		dest.Bio,
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
// @Router /user/register [post]
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

	var models model.UserAccount
	if err := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(request.Email))).Query(Database, &models); err != nil {
		if err != qrm.ErrNoRows {
			fmt.Printf("[/user/register] Error querying database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
			return
		}
	} else {
		g.Status(http.StatusBadRequest)
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
		return
	}

	g.JSON(http.StatusOK, userRegisterResponse{dest.Identifier})
}

type userBioRequest struct {
	Token string `json:"token"`
	Bio   string `json:"bio"`
}

// Bio godoc
// @Summary Updates bio
// @Description Updates a user's bio
// @Tags user
// @Produce json
// @Consume json
// @Success 200
// @Failure 400
// @Failure 500
// @Param request body userBioRequest true "User token and new bio"
// @Router /user/bio [post]
func userBio(g *gin.Context) {
	var request userBioRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	token, err := ProcessToken(request.Token)
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	stmt := UserAccount.UPDATE(UserAccount.Bio).MODEL(model.UserAccount{
		Bio: &request.Bio,
	}).WHERE(UserAccount.Identifier.EQ(String(token.Identifier)))

	if _, err := stmt.Exec(Database); err != nil {
		fmt.Printf("[/user/bio] Failed to execute query on database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.Status(http.StatusOK)
	}
}

type userJoinRequest struct {
	Token   string `json:"token"`
	GroupID int64  `json:"group_id"`
}

// Join godoc
// @Summary Join group
// @Description Adds a user to a group
// @Tags user
// @Consume json
// @Success 200
// @Failure 400
// @Failure 500
// @Param request body userJoinRequest true "User token and group to join"
// @Router /user/join [post]
func userJoin(g *gin.Context) {
	var request userJoinRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	token, err := ProcessToken(request.Token)
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var user model.UserAccount
	if err := UserAccount.SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.Identifier))).Query(Database, &user); err != nil {
		switch err {
		default:
			fmt.Printf("[/user/join] Failed query database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	var room model.Room
	if err := Room.SELECT(Room.ID).FROM(Room).WHERE(Room.ID.EQ(Int64(request.GroupID))).Query(Database, &room); err != nil {
		switch err {
		default:
			fmt.Printf("[/user/join] Failed query database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	stmt := UserRoom.INSERT(UserRoom.UserID, UserRoom.RoomID).MODEL(model.UserRoom{
		UserID: user.ID,
		RoomID: room.ID,
	})

	if _, err := stmt.Exec(Database); err != nil {
		fmt.Printf("[/user/join] Failed to execute query on database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.Status(http.StatusOK)
	}
}

type userProfilePicturePostRequest struct {
	Token string `json:"token"`
	Image string `json:"jpeg"`
}

// Profile Picture godoc
// @Summary Upload profile picture
// @Description Uploads a new profile picture, replacing the old one
// @Tags user
// @Consume json
// @Success 200
// @Failure 400
// @Failure 500
// @Param request body userProfilePicturePostRequest true "User token and profile picture"
// @Router /user/profile_picture [post]
func userProfilePicturePost(g *gin.Context) {
	var request userProfilePicturePostRequest
	if err := g.BindJSON(&request); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	token, err := ProcessToken(request.Token)
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	var user model.UserAccount
	if err := UserAccount.SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.Identifier))).Query(Database, &user); err != nil {
		switch err {
		default:
			fmt.Printf("[/user/profile_picture] Failed query database: %s\n", err.Error())
			g.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			g.Status(http.StatusBadRequest)
		}

		return
	}

	data, err := hex.DecodeString(request.Image)
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	file, err := os.Create(filepath.Join(Options.ImageStore, token.Identifier))
	if err != nil {
		fmt.Printf("[/user/profile_picture] Failed to create file: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	if err := file.Truncate(0); err != nil {
		fmt.Printf("[/user/profile_picture] Failed to truncate file: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	if _, err := file.Write(data); err != nil {
		fmt.Printf("[/user/profile_picture] Failed to write file: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	}

	g.Status(http.StatusOK)
}

// Profile Picture godoc
// @Summary Retrieves profile picture
// @Description Gets a user's profile picture from their identifier
// @Tags user
// @Produce jpeg
// @Success 200
// @Failure 400
// @Failure 500
// @Param ident path string true "User identifier"
// @Router /user/profile_picture/{ident} [get]
func userProfilePictureGet(g *gin.Context) {
	var uri struct {
		Identifier string `uri:"ident" binding:"required"`
	}

	if err := g.ShouldBindUri(&uri); err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	if data, err := os.ReadFile(filepath.Join(Options.ImageStore, uri.Identifier)); err != nil {
		g.Status(http.StatusBadRequest)
	} else {
		g.Data(http.StatusOK, "image/jpeg", data)
	}
}

func UserHandler(r *gin.RouterGroup) {
	r.GET("/get/:ident", userGet)
	r.POST("/register", userRegister)
	r.POST("/bio", userBio)
	r.POST("/join", userJoin)
	r.POST("/profile_picture", userProfilePicturePost)
	r.GET("/profile_picture/:ident", userProfilePictureGet)
}
