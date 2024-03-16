package api

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	FileName string `json:"filename"`
	Original string `json:"original"`
	Size     int64  `json:"size"`
}

func formatBytes(bytes int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
	)

	size := float64(bytes)
	unit := ""

	switch {
	case bytes >= MB:
		unit = "MB"
		size /= MB
	case bytes >= KB:
		unit = "KB"
		size /= KB
	default:
		unit = "B"
	}

	return fmt.Sprintf("%.2f%s", size, unit)
}

func FileSend(ctx *gin.Context) {
	var req UploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		Abort(http.StatusBadRequest, "Please select a file", ctx)
		return
	}

	combination := GenerateRandomString(5)
	newFileName := combination + filepath.Ext(file.Filename)
	if err := ctx.SaveUploadedFile(file, "files/"+newFileName); err != nil {
		Abort(http.StatusInternalServerError, "Unable to save the file", ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please select a file"})
		return
	}

	AppendToJson(ctx, file, newFileName)

	Abort(http.StatusAccepted, "File has been uploaded and available at "+GetUrlPath(newFileName, ctx), ctx)
	return
}

func FilesList(ctx *gin.Context) {
	json := ReadJson(ctx)

	var templateList string
	for _, file := range json {
		noDot := strings.Replace(file.FileName, ".", "", -1)

		html := fmt.Sprintf(`
      <tr id="%s">
        <td>%s</td>
        <td>%s</td>
        <td>%s</td>
        <td><button hx-post="/delete/%s" hx-target="#%s">Delete</button></td>
      </tr>
      `, noDot, file.FileName, file.Original, formatBytes(file.Size), file.FileName, noDot)
		templateList += html
	}

	ctx.HTML(http.StatusOK, "files.html", gin.H{"files": template.HTML(templateList)})
}

func FileDelete(ctx *gin.Context) {
	file := ctx.Param("file")

	RemoveFromJson(ctx, file)

	err := os.Remove("files/" + file)
	if err != nil {
		Abort(http.StatusInternalServerError, "Error saving file", ctx)
		return
	}
}
