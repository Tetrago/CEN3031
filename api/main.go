package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/tetrago/motmot/api/docs"
	v1 "github.com/tetrago/motmot/api/v1"
)

// @BasePath /api

func setupDatabase() *sql.DB {
	var connectString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.DatabaseHostname,
		opts.DatabasePort,
		opts.DatabaseUsername,
		opts.DatabasePassword,
		opts.DatabaseName,
	)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}

	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	v1.Handler(r.Group("/api/v1"))

	if Debug {
		r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return r
}

var opts struct {
	DatabaseHostname string `short:"r" long:"hostname" description:"Hostname of PostgreSQL database" required:"true"`
	DatabasePort     int    `short:"c" long:"connection" description:"Port used for database" default:"5432"`
	DatabaseName     string `short:"d" long:"database" description:"Name of database" required:"true"`
	DatabaseUsername string `short:"u" long:"user" description:"Name of database user" required:"true"`
	DatabasePassword string `short:"p" long:"password" description:"Password of database user" required:"true"`
	Port             int    `short:"s" long:"server" description:"Port API is served from" default:"8080"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	db := setupDatabase()
	defer db.Close()

	v1.Database = db

	r := setupRouter()
	r.Run(fmt.Sprintf(":%d", opts.Port))
}
