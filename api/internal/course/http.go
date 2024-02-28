package course

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
	"github.com/tetrago/motmot/api/internal/globals"
)

type DepartmentRequest struct {
}

type DepartmentResponseItem struct {
	Identifier  string `json:"label"`
	Description string `json:"name"`
}

type DepartmentResponse []DepartmentResponseItem

// Department godoc
// @Summary Get courses in department
// @Description Queries for all UF courses in a three-letter department prefix
// @Tags course
// @Produce json
// @Sucess 200 {object} DepartmentResponse
// @Failure 400
// @Failure 503
// @Param dep path string true "Three-letter department prefix"
// @Router /course/department/{dep} [get]
func Department(c *gin.Context) {
	var uri struct {
		Department string `uri:"dep" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil || len(uri.Department) != 3 {
		c.Status(http.StatusBadRequest)
		return
	}

	courses, err := queryDepartmentCourses(uri.Department)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	c.JSON(http.StatusOK, lo.Map(courses, func(x course, _ int) DepartmentResponseItem {
		return DepartmentResponseItem(x)
	}))
}

// Group godoc
// @Summary Get group of specified course
// @Description Gets (or creates) group of specified course
// @Tags course
// @Produce json
// @Sucess 200 {int64}
// @Failure 400
// @Failure 500
// @Failure 503
// @Param dep  path string true "Three-letter department prefix"
// @Param code path string true "Four-digit (with potential postfix) course code"
// @Router /course/group/{dep}/{code} [get]
func Group(g *gin.Context) {
	var uri struct {
		Department string `uri:"dep" binding:"required"`
		Code       string `uri:"code" binding:"required"`
	}

	if err := g.ShouldBindUri(&uri); err != nil || len(uri.Department) != 3 || (len(uri.Code) != 4 && len(uri.Code) != 5) {
		g.Status(http.StatusBadRequest)
		return
	}

	label := fmt.Sprintf("%s %s", strings.ToUpper(uri.Department), strings.ToUpper(uri.Code))

	var dest model.Room
	stmt := SELECT(Room.ID).FROM(Room).WHERE(Room.Name.EQ(String(label)))

	switch err := stmt.Query(globals.Database, &dest); err {
	default:
		fmt.Printf("[/course/group] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
		return
	case nil:
		g.JSON(http.StatusOK, dest.ID)
		return
	case qrm.ErrNoRows:
	}

	courses, err := queryDepartmentCourses(uri.Department)
	if err != nil {
		g.Status(http.StatusServiceUnavailable)
		return
	}

	course, ok := lo.Find(courses, func(x course) bool { return x.Identifier == label })
	if !ok {
		g.Status(http.StatusBadRequest)
		return
	}

	ins := Room.INSERT(Room.Name, Room.Description).MODEL(model.Room{
		Name:        label,
		Description: course.Description,
	}).RETURNING(Room.ID)

	if err := ins.Query(globals.Database, &dest); err != nil {
		fmt.Printf("[/course/group] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.JSON(http.StatusOK, dest.ID)
	}
}

func HttpHandler(r *gin.RouterGroup) {
	g := r.Group("/course")
	g.GET("/department/:dep", Department)
	g.GET("/group/:dep/:code", Group)
}
