package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/tetrago/motmot/api/docs"
	v1 "github.com/tetrago/motmot/api/v1"
)

// @BasePath /api

func getEnv(name string, def string) string {
	value, found := os.LookupEnv(name)
	if found {
		return value
	} else {
		return def
	}
}

func setupDatabase() *sql.DB {
	var dbHostname = os.Getenv("DB_HOSTNAME")
	var dbPort = getEnv("DB_PORT", "5432")
	var dbDatabase = getEnv("DB_DATABASE", "motmot")
	var dbUsername = getEnv("DB_USERNAME", "motmot")
	var dbPassword = os.Getenv("DB_PASSWORD")

	if len(dbHostname) == 0 || len(dbPassword) == 0 {
		fmt.Println("DB_HOSTNAME or DB_PASSWORD not set!")
		os.Exit(1)
	}

	var connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHostname, dbPort, dbUsername, dbPassword, dbDatabase)

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

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

func main() {
	db := setupDatabase()
	defer db.Close()

	v1.Database = db

	r := setupRouter()
	r.Run(fmt.Sprintf(":%s", getEnv("PORT", "8080")))
}
