package api

import "github.com/gin-gonic/gin"

func Render(code int, message any, context *gin.Context) {
	if message == nil {
		context.HTML(code, "index.html", nil)
		return
	}
	context.HTML(code, "index.html", gin.H{
		"message": message,
	})
}

func Abort(code int, message string, context *gin.Context) {
	context.AbortWithStatusJSON(code, gin.H{"message": message})
}

func GetPath(file string, context *gin.Context) string {
	scheme := "http"
	if context.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + context.Request.Host + "/files/" + file
}
