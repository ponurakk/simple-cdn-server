package api

import "github.com/gin-gonic/gin"

const template = "index.html"

func Render(code int, message any, context *gin.Context) {
	if message == nil {
		context.HTML(code, template, nil)
		return
	}
	context.HTML(code, template, gin.H{
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
