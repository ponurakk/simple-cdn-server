package main

import (
	"cdn-server/api"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	Token string `form:"token"`
}

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

	router.POST("/", func(ctx *gin.Context) {
		var req UploadRequest
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Token != config.Token {
			api.Render(http.StatusBadRequest, gin.H{"message": "Invalid token"}, ctx)
			return
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			api.Render(http.StatusBadRequest, gin.H{"message": "Please select a file"}, ctx)
			return
		}

		combination := api.GenerateRandomString(5)
		newFileName := combination + filepath.Ext(file.Filename)
		if err := ctx.SaveUploadedFile(file, "files/"+newFileName); err != nil {
			api.Abort(http.StatusInternalServerError, "Unable to save the file", ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please select a file"})
			return
		}

		api.Abort(http.StatusAccepted, "File has been uploaded and available at "+api.GetPath(newFileName, ctx), ctx)
		return
	})

	fmt.Println(config)

	// Run the server
	router.Run(":" + config.Port)
}
