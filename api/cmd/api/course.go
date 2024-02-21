package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/samber/lo"

	"github.com/tetrago/motmot/api/.gen/motmot/public/model"
	. "github.com/tetrago/motmot/api/.gen/motmot/public/table"
)

type course struct {
	Label string `json:"code"`
	Name  string `json:"title"`
}

func queryDepartmentCourses(dep string) ([]course, error) {
	url := fmt.Sprintf("https://catalog.ufl.edu/course-search/api/?page=fose&route=search&subject=%s", dep)
	request := fmt.Sprintf(`{"other":{"srcdb":""},"criteria":[{"field":"subject","value":"%s"}]}`, dep)

	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(request)))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var dest struct {
		Courses []course `json:"results"`
	}

	if err := json.Unmarshal(body, &dest); err != nil {
		return nil, err
	}

	return dest.Courses, nil
}

type courseDeparmentResponseItem struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

// Department godoc
// @Summary Get courses in department
// @Description Queries for all UF courses in a three-letter department prefix
// @Tags course
// @Produce json
// @Sucess 200 {array} courseDepartmentResponseItem
// @Failure 400
// @Failure 503
// @Param dep path string true "Three-letter department prefix"
// @Router /course/department/{dep} [get]
func courseDepartment(g *gin.Context) {
	var uri struct {
		Department string `uri:"dep" binding:"required"`
	}

	if err := g.ShouldBindUri(&uri); err != nil || len(uri.Department) != 3 {
		g.Status(http.StatusBadRequest)
		return
	}

	courses, err := queryDepartmentCourses(uri.Department)
	if err != nil {
		g.Status(http.StatusServiceUnavailable)
		return
	}

	g.JSON(http.StatusOK, lo.Map(courses, func(x course, _ int) courseDeparmentResponseItem {
		return courseDeparmentResponseItem(x)
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
// @Param code path string true "Four-digit course code"
// @Router /course/group/{dep}/{code} [get]
func courseGroup(g *gin.Context) {
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

	switch err := stmt.Query(Database, &dest); err {
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

	course, ok := lo.Find(courses, func(x course) bool { return x.Label == label })
	if !ok {
		g.Status(http.StatusBadRequest)
		return
	}

	ins := Room.INSERT(Room.Name, Room.Description).MODEL(model.Room{
		Name:        label,
		Description: course.Name,
	}).RETURNING(Room.ID)

	if err := ins.Query(Database, &dest); err != nil {
		fmt.Printf("[/course/group] Error querying database: %s\n", err.Error())
		g.Status(http.StatusInternalServerError)
	} else {
		g.JSON(http.StatusOK, dest.ID)
	}
}

func CourseHandler(r *gin.RouterGroup) {
	r.GET("/department/:dep", courseDepartment)
	r.GET("/group/:dep/:code", courseGroup)
}
