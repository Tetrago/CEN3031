//go:build debug
// +build debug

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tetrago/motmot/api/internal/auth"
	"github.com/tetrago/motmot/api/internal/globals"
	"github.com/tetrago/motmot/api/internal/user"
)

func TestRegistrationAndLogin(t *testing.T) {
	router := setupRouter()
	globals.Database = setupDatabase()
	defer globals.Database.Close()

	ident, _ := user.MakeIdentifier()

	body, _ := json.Marshal(user.RegisterRequest{
		DisplayName: "Name",
		Email:       ident,
		Password:    "password",
	})

	req := httptest.NewRequest("POST", "/api/v1/user/register", bytes.NewReader(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	body, _ = json.Marshal(auth.LoginRequest{
		Email:    ident,
		Password: "password",
	})

	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestCourseGroup(t *testing.T) {
	router := setupRouter()
	globals.Database = setupDatabase()
	defer globals.Database.Close()

	req := httptest.NewRequest("GET", "/api/v1/course/group/COP/4600", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	id := w.Body.String()

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	assert.Equal(t, id, w.Body.String())

	req = httptest.NewRequest("GET", "/api/v1/course/group/COP/1234", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestBio(t *testing.T) {
	router := setupRouter()
	globals.Database = setupDatabase()
	defer globals.Database.Close()

	ident, _ := user.MakeIdentifier()

	body, _ := json.Marshal(user.RegisterRequest{
		DisplayName: "Name",
		Email:       ident,
		Password:    "password",
	})

	req := httptest.NewRequest("POST", "/api/v1/user/register", bytes.NewReader(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(user.BioRequest{
		Bio: "New bio",
	})

	req = httptest.NewRequest("POST", "/api/v1/user/bio", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	body, _ = json.Marshal(auth.LoginRequest{
		Email:    ident,
		Password: "password",
	})

	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, cookie.Name, "token")

	body, _ = json.Marshal(user.BioRequest{
		Bio: "New bio",
	})

	req = httptest.NewRequest("POST", "/api/v1/user/bio", bytes.NewReader(body))
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", cookie.Value))

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestJoin(t *testing.T) {
	router := setupRouter()
	globals.Database = setupDatabase()
	defer globals.Database.Close()

	req := httptest.NewRequest("GET", "/api/v1/course/group/COP/4600", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	groupId, _ := strconv.Atoi(w.Body.String())

	ident, _ := user.MakeIdentifier()

	body, _ := json.Marshal(user.RegisterRequest{
		DisplayName: "Name",
		Email:       ident,
		Password:    "password",
	})

	req = httptest.NewRequest("POST", "/api/v1/user/register", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(user.JoinRequest{GroupID: int64(groupId)})
	req = httptest.NewRequest("POST", "/api/v1/user/join", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	body, _ = json.Marshal(auth.LoginRequest{
		Email:    ident,
		Password: "password",
	})

	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, cookie.Name, "token")

	body, _ = json.Marshal(user.JoinRequest{GroupID: int64(groupId)})
	req = httptest.NewRequest("POST", "/api/v1/user/join", bytes.NewReader(body))
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", cookie.Value))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
