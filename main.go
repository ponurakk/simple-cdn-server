package main

import (
	"cdn-server/api"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(config api.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := `<form hx-post="/logout" hx-swap="outerHTML" hx-trigger="submit">
      <button type="submit">Logout</button>
      </form>`

		tokenString, err := ctx.Cookie("jwt")
		if err != nil {
			data = `<form hx-post="/login" hx-swap="outerHTML" hx-trigger="submit">
      <input type="text" placeholder="Token" name="token">
      <button type="submit">Login</button>
      </form>`
		}

		token, err := api.ValidateToken(config, tokenString)

		if err != nil || !token.Valid {
			data = `<form hx-post="/login" hx-swap="outerHTML" hx-trigger="submit">
      <input type="text" placeholder="Token" name="token">
      <button type="submit">Login</button>
      </form>`
		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{"form": template.HTML(data)})
	}
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

	router.GET("/", index(config))

	router.POST("/login", api.Login(config))

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(api.AuthMiddleware(config))
	{
		authenticated.POST("/logout", api.Logout)
		authenticated.POST("/", api.FileSend)
		authenticated.GET("/list", api.FilesList)
		authenticated.POST("/delete/:file", api.FileDelete)
	}

	router.Run(":" + config.Port)
}
