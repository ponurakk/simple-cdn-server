package api

import "github.com/gin-gonic/gin"

func Abort(code int, message string, context *gin.Context) {
	context.AbortWithStatusJSON(code, gin.H{"message": message})
}

func GetUrlPath(file string, context *gin.Context) string {
	scheme := "http"
	if context.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + context.Request.Host + "/files/" + file
}
