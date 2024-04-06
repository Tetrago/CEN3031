package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/tetrago/motmot/api/docs"
	"github.com/tetrago/motmot/api/internal/auth"
	"github.com/tetrago/motmot/api/internal/course"
	"github.com/tetrago/motmot/api/internal/globals"
	"github.com/tetrago/motmot/api/internal/group"
	"github.com/tetrago/motmot/api/internal/user"
	"github.com/tetrago/motmot/api/internal/ws"
)

func setupDatabase() *sql.DB {
	var connectString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		globals.Opts.DatabaseHostname,
		globals.Opts.DatabasePort,
		globals.Opts.DatabaseUsername,
		globals.Opts.DatabasePassword,
		globals.Opts.DatabaseName,
	)

	if db, err := sql.Open("postgres", connectString); err != nil {
		panic("Failed to connect to database!")
	} else {
		return db
	}
}

func getOrigin() string {
	if globals.Opts.SslEnabled {
		return fmt.Sprintf("https://%s:%d", globals.Opts.Hostname, globals.Opts.Hostport)
	} else {
		return fmt.Sprintf("http://%s:%d", globals.Opts.Hostname, globals.Opts.Hostport)
	}
}

func setupRouter() *gin.Engine {
	if !Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{getOrigin()}
	config.AllowCredentials = true

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(cors.New(config))

	g := r.Group(globals.Opts.BasePath)

	auth.HttpHandler(g)
	course.HttpHandler(g)
	group.HttpHandler(g)
	user.HttpHandler(g)
	ws.HttpHandler(g)

	if Debug {
		docs.SwaggerInfo.BasePath = globals.Opts.BasePath
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return r
}

func main() {
	globals.Database = setupDatabase()
	defer globals.Database.Close()

	r := setupRouter()
	r.Run(fmt.Sprintf(":%d", globals.Opts.Port))
}
