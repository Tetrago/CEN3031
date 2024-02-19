package main

import (
    "github.com/gin-gonic/gin"
    swaggerfiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    _ "github.com/tetrago/motmot/api/docs"
    "github.com/tetrago/motmot/api/v1"
)

func main() {
    r := gin.Default()

    v1.Setup(r.Group("/api/v1"))

    r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
    r.Run(":8080")
}
