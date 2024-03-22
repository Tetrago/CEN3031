package user

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/internal/auth"
	"github.com/tetrago/motmot/api/internal/crypt"
	"github.com/tetrago/motmot/api/internal/globals"
)

type GetResponseGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetResponse struct {
	Identifier  string              `json:"ident"`
	DisplayName string              `json:"display_name"`
	Bio         *string             `json:"bio,omitempty"`
	Groups      *[]GetResponseGroup `json:"groups"`
}

// User godoc
// @Summary Fetch user
// @Description Fetches publically available user information and groups.
// @Tags user
// @Produce json
// @Success 200 {object} GetResponse
// @Failure 400
// @Failure 500
// @Param ident path string true "User identifier"
// @Router /user/get/{ident} [get]
func Get(c *gin.Context) {
	var uri struct {
		Identifier string `uri:"ident" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.Status(http.StatusBadRequest)
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

	if err := stmt.Query(globals.Database, &dest); err == qrm.ErrNoRows {
		c.Status(http.StatusBadRequest)
	} else if err != nil {
		fmt.Printf("[/user/get] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		groups := lo.Map(dest.Rooms, func(x model.Room, _ int) GetResponseGroup {
			return GetResponseGroup{x.ID, x.Name}
		})

		c.JSON(http.StatusOK, GetResponse{
			dest.Identifier,
			dest.DisplayName,
			dest.Bio,
			&groups,
		})
	}
}

type RegisterRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	Identifier string `json:"ident"`
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user given the provided arguments
// @Tags user
// @Produce json
// @Consume json
// @Success 200 {object} RegisterResponse
// @Failure 400
// @Failure 500
// @Param request body RegisterRequest true "User registration information"
// @Router /user/register [post]
func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var existing model.UserAccount
	stmt := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(request.Email)))

	if err := stmt.Query(globals.Database, &existing); err == nil {
		c.Status(http.StatusBadRequest)
		return
	} else if err != qrm.ErrNoRows {
		fmt.Printf("[/user/register] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	ident, err := makeIdentifier()
	if err != nil {
		fmt.Printf("[/user/register] Error generating identifier: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	var dest model.UserAccount
	ins := UserAccount.INSERT(UserAccount.Identifier, UserAccount.DisplayName, UserAccount.Hash, UserAccount.Email).
		MODEL(model.UserAccount{
			Identifier:  ident,
			DisplayName: request.DisplayName,
			Hash:        crypt.Hash(request.Password),
			Email:       request.Email,
		}).
		RETURNING(UserAccount.AllColumns)

	if err := ins.Query(globals.Database, &dest); err != nil {
		fmt.Printf("[/user/register] Error querying database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, RegisterResponse{dest.Identifier})
	}
}

type PostProfilePictureRequest struct {
	Image string `json:"jpeg"`
}

// Profile Picture godoc
// @Summary Upload profile picture
// @Description Uploads a new profile picture, replacing the old one
// @Tags user
// @Consume json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Param request body PostProfilePictureRequest true "New profile picture"
// @Router /user/profile_picture [post]
func PostProfilePicture(c *gin.Context) {
	var request PostProfilePictureRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	token := auth.ExpectToken(c)

	var user model.UserAccount
	stmt := UserAccount.SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.UserIdentifier())))

	if err := stmt.Query(globals.Database, &user); err != nil {
		switch err {
		default:
			fmt.Printf("[/user/profile_picture] Failed query database: %s\n", err.Error())
			c.Status(http.StatusInternalServerError)
		case qrm.ErrNoRows:
			c.Status(http.StatusBadRequest)
		}

		return
	}

	data, err := hex.DecodeString(request.Image)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	file, err := os.Create(filepath.Join(globals.Opts.ImageFolderPath, token.UserIdentifier()))
	if err != nil {
		fmt.Printf("[/user/profile_picture] Failed to create file: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := file.Truncate(0); err != nil {
		fmt.Printf("[/user/profile_picture] Failed to truncate file: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if _, err := file.Write(data); err != nil {
		fmt.Printf("[/user/profile_picture] Failed to write file: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

var Colors = []color.RGBA{
	{0xA6, 0xAD, 0xBB, 0xFF},
	{0x00, 0xB5, 0xFF, 0xFF},
	{0x00, 0xA9, 0x6E, 0xFF},
	{0xFF, 0xBE, 0x00, 0xFF},
	{0xFF, 0x58, 0x61, 0xFF},
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
func GetProfilePicture(c *gin.Context) {
	var uri struct {
		Identifier string `uri:"ident" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if data, err := os.ReadFile(filepath.Join(globals.Opts.ImageFolderPath, uri.Identifier)); err != nil {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.SetRGBA(0, 0, Colors[crypt.HashToInt(uri.Identifier)%len(Colors)])

		var b bytes.Buffer
		w := bufio.NewWriter(&b)

		if err := jpeg.Encode(w, img, nil); err != nil {
			fmt.Printf("[/usr/profile_picture] Failed to generate temporary profile picture: %s\n", err.Error())
		} else {
			c.Data(http.StatusOK, "image/jpeg", b.Bytes())
		}
	} else {
		c.Data(http.StatusOK, "image/jpeg", data)
	}
}

type BioRequest struct {
	Bio string `json:"bio"`
}

// Bio godoc
// @Summary Updates bio
// @Description Updates a user's bio
// @Tags user
// @Produce json
// @Consume json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Param request body BioRequest true "New bio"
// @Router /user/bio [post]
func Bio(c *gin.Context) {
	token := auth.ExpectToken(c)

	var request BioRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	stmt := UserAccount.UPDATE(UserAccount.Bio).MODEL(model.UserAccount{
		Bio: &request.Bio,
	}).WHERE(UserAccount.Identifier.EQ(String(token.UserIdentifier())))

	if _, err := stmt.Exec(globals.Database); err != nil {
		fmt.Printf("[/user/bio] Failed to execute query on database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}

type JoinRequest struct {
	GroupID int64 `json:"group_id"`
}

// Join godoc
// @Summary Join group
// @Description Adds a user to a group
// @Tags user
// @Consume json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Param request body JoinRequest true "Group to join"
// @Router /user/join [post]
func Join(c *gin.Context) {
	var request JoinRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	token := auth.ExpectToken(c)

	var user model.UserAccount
	stmt := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.UserIdentifier())))

	if err := stmt.Query(globals.Database, &user); err == qrm.ErrNoRows {
		c.Status(http.StatusBadRequest)
		return
	} else if err != nil {
		fmt.Printf("[/user/join] Failed query database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	var room model.Room
	stmt = Room.SELECT(Room.ID).FROM(Room).WHERE(Room.ID.EQ(Int64(request.GroupID)))

	if err := stmt.Query(globals.Database, &room); err == qrm.ErrNoRows {
		c.Status(http.StatusBadRequest)
		return
	} else if err != nil {
		fmt.Printf("[/user/join] Failed query database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	ins := UserRoom.INSERT(UserRoom.UserID, UserRoom.RoomID).MODEL(model.UserRoom{
		UserID: user.ID,
		RoomID: room.ID,
	})

	if _, err := ins.Exec(globals.Database); err != nil {
		fmt.Printf("[/user/join] Failed to execute query on database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}

type LeaveRequest struct {
	GroupID int64 `json:"group_id"`
}

// Leave godoc
// @Summary Leave group
// @Description Removes a user from a group
// @Tags user
// @Consume json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Param request body JoinRequest true "Group to leave"
// @Router /user/leave [post]
func Leave(c *gin.Context) {
	var request LeaveRequest
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	token := auth.ExpectToken(c)

	var dest model.UserAccount
	stmt := SELECT(UserAccount.ID).FROM(UserAccount).WHERE(UserAccount.Identifier.EQ(String(token.UserIdentifier())))

	if err := stmt.Query(globals.Database, &dest); err == qrm.ErrNoRows {
		fmt.Print("[/user/leave] Unable to find user from signed token")
		c.Status(http.StatusInternalServerError)
		return
	} else if err != nil {
		fmt.Printf("[/user/leave] Failed to query database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	var entry model.UserRoom
	del := UserRoom.DELETE().WHERE(UserRoom.UserID.EQ(Int64(dest.ID)).AND(UserRoom.RoomID.EQ(Int64(request.GroupID)))).RETURNING(UserRoom.AllColumns)

	if err := del.Query(globals.Database, &entry); err == qrm.ErrNoRows {
		c.Status(http.StatusBadRequest)
	} else if err != nil {
		fmt.Printf("[/user/leave] Failed to execute query on database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}

type GroupsResponseItem struct {
	ID   int64  `json:"group_id"`
	Name string `json:"name"`
}

// Leave godoc
// @Summary Get groups
// @Description Returns all groups a user belongs to
// @Tags user
// @Produe json
// @Success 200 {array} GroupsResponseItem
// @Failure 500
// @Router /user/groups [get]
func Groups(c *gin.Context) {
	token := auth.ExpectToken(c)

	var dest []model.Room
	stmt := SELECT(Room.ID, Room.Name).FROM(
		UserAccount.
			INNER_JOIN(UserRoom, UserAccount.ID.EQ(UserRoom.UserID)).
			INNER_JOIN(Room, UserRoom.RoomID.EQ(Room.ID)),
	).WHERE(UserAccount.Identifier.EQ(String(token.UserIdentifier())))

	if err := stmt.Query(globals.Database, &dest); err != nil {
		fmt.Printf("[/user/groups] Failed to query database: %s\n", err.Error())
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, lo.Map(dest, func(x model.Room, _ int) GroupsResponseItem {
			return GroupsResponseItem{x.ID, x.Name}
		}))
	}
}

func HttpHandler(r *gin.RouterGroup) {
	g := r.Group("/user")
	g.POST("/register", Register)
	g.GET("/get/:ident", Get)
	g.GET("/profile_picture/:ident", GetProfilePicture)

	g.Use(auth.Middleware())
	g.POST("/profile_picture", PostProfilePicture)
	g.POST("/bio", Bio)
	g.POST("/join", Join)
	g.POST("/leave", Leave)
	g.GET("/groups", Groups)
}
