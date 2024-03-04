package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

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

func saveJson(ctx *gin.Context, file *multipart.FileHeader, newFileName string) {
	jsonData, err := os.ReadFile("info.json")
	if err != nil {
		Abort(http.StatusInternalServerError, "Error reading JSON file", ctx)
	}

	var dataArray []FileInfo

	if err := json.Unmarshal(jsonData, &dataArray); err != nil {
		Abort(http.StatusInternalServerError, "Error parsing JSON file", ctx)
	}

	newFile := FileInfo{
		FileName: newFileName,
		Original: file.Filename,
		Size:     file.Size,
	}

	dataArray = append(dataArray, newFile)

	updatedJSON, err := json.MarshalIndent(dataArray, "", "  ")
	if err != nil {
		Abort(http.StatusInternalServerError, "Error parsing JSON file", ctx)
	}

	if err := os.WriteFile("info.json", updatedJSON, 0644); err != nil {
		Abort(http.StatusInternalServerError, "Error saving JSON file", ctx)
	}
}

func readJson(ctx *gin.Context) []FileInfo {
	jsonData, err := os.ReadFile("info.json")
	if err != nil {
		Abort(http.StatusInternalServerError, "Error reading JSON file", ctx)
	}

	var data []FileInfo
	if err := json.Unmarshal(jsonData, &data); err != nil {
		Abort(http.StatusInternalServerError, "Error parsing json", ctx)
	}
	return data
}

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

	saveJson(ctx, file, newFileName)

	Abort(http.StatusAccepted, "File has been uploaded and available at "+GetPath(newFileName, ctx), ctx)
	return
}

func FilesList(ctx *gin.Context) {
	json := readJson(ctx)

	var templateList string
	for _, file := range json {
		html := fmt.Sprintf(`
      <tr>
        <td>%s</td>
        <td>%s</td>
        <td>%s</td>
        <td><button>Delete</button></td>
      </tr>
      `, file.FileName, file.Original, formatBytes(file.Size))
		templateList += html
	}

	ctx.HTML(http.StatusOK, "files.html", gin.H{"files": template.HTML(templateList)})
}
