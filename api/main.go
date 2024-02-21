package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/tetrago/motmot/api/docs"
	"github.com/tetrago/motmot/api/util"
	v1 "github.com/tetrago/motmot/api/v1"
)

// @BasePath /api

func setupDatabase() *sql.DB {
	password, err := parseSecret(Options.DatabasePassword)
	if err != nil {
		fmt.Printf("Failed to read database password: `%s`\n", err.Error())
		os.Exit(1)
	}

	var connectString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Options.DatabaseHostname,
		Options.DatabasePort,
		Options.DatabaseUsername,
		password,
		Options.DatabaseName,
	)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}

	return db
}

func setupRouter() *gin.Engine {
	if !Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	v1.Handler(r.Group("/api/v1"))

	if Debug {
		r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return r
}

var Options struct {
	DatabaseHostname string `short:"r" long:"hostname" description:"Hostname of PostgreSQL database" required:"true"`
	DatabasePort     int    `short:"c" long:"connection" description:"Port used for database" default:"5432"`
	DatabaseName     string `short:"d" long:"database" description:"Name of database" required:"true"`
	DatabaseUsername string `short:"u" long:"user" description:"Name of database user" required:"true"`
	DatabasePassword string `short:"p" long:"password" description:"Password of database user (use file:// to reference a file)" required:"true"`
	Port             int    `short:"s" long:"server" description:"Port API is served from" default:"8080"`
	Secret           string `short:"x" long:"secret" description:"API key secret (use file:// to reference a file)" required:"true"`
}

func parseSecret(key string) ([]byte, error) {
	if !strings.HasPrefix(key, "file://") {
		return []byte(key), nil
	}

	var path string
	fmt.Sscanf(key, "file://%s", &path)

	if data, err := os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("failed to read secret file")
	} else {
		return data, nil
	}
}

func main() {
	if _, err := flags.Parse(&Options); err != nil {
		os.Exit(1)
	}

	var err error
	if util.Secret, err = parseSecret(Options.Secret); err != nil {
		fmt.Printf("Error when parsing API secret: `%s`\n", err.Error())
		os.Exit(1)
	}

	db := setupDatabase()
	defer db.Close()

	v1.Database = db

	r := setupRouter()
	r.Run(fmt.Sprintf(":%d", Options.Port))
}
