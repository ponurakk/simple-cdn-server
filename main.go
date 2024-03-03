package main

import (
	"cdn-server/api"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := api.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Static("/files", "files")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/login", api.Login(config))

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(api.AuthMiddleware(config))
	{
		authenticated.POST("/logout", api.Logout)
		authenticated.POST("/", api.FileSend)
		authenticated.GET("/list", api.FilesList)
	}

	router.Run(":" + config.Port)
}
