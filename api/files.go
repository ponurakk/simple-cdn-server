package api

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func FileSend(ctx *gin.Context) {
	var req UploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		Render(http.StatusBadRequest, gin.H{"message": "Please select a file"}, ctx)
		return
	}

	combination := GenerateRandomString(5)
	newFileName := combination + filepath.Ext(file.Filename)
	if err := ctx.SaveUploadedFile(file, "files/"+newFileName); err != nil {
		Abort(http.StatusInternalServerError, "Unable to save the file", ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please select a file"})
		return
	}

	Abort(http.StatusAccepted, "File has been uploaded and available at "+GetPath(newFileName, ctx), ctx)
	return
}

func FilesList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"files": []string{"file1", "file2"}})
}
