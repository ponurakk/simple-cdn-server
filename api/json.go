package api

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ReadJson(ctx *gin.Context) []FileInfo {
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

func SaveJson(ctx *gin.Context, data []FileInfo) {
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		Abort(http.StatusInternalServerError, "Error parsing JSON file", ctx)
	}

	if err := os.WriteFile("info.json", updatedJSON, 0644); err != nil {
		Abort(http.StatusInternalServerError, "Error saving JSON file", ctx)
	}
}

func AppendToJson(ctx *gin.Context, file *multipart.FileHeader, newFileName string) {
	data := ReadJson(ctx)

	newFile := FileInfo{
		FileName: newFileName,
		Original: file.Filename,
		Size:     file.Size,
	}

	data = append(data, newFile)
	SaveJson(ctx, data)
}

func RemoveFromJson(ctx *gin.Context, fileName string) {
	data := ReadJson(ctx)

	var indexToRemove int = -1
	for i, data := range data {
		if data.FileName == fileName {
			indexToRemove = i
			break
		}
	}

	if indexToRemove != -1 {
		data = append(data[:indexToRemove], data[indexToRemove+1:]...)
	} else {
		Abort(http.StatusInternalServerError, "Object not found in the JSON array.", ctx)
		return
	}

	SaveJson(ctx, data)
}
