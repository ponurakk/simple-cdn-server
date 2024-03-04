package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Token string `json:"token"`
	jwt.StandardClaims
}

type UploadRequest struct {
	Token string `form:"token"`
}

func ValidateToken(config Config, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey), nil
	})
	return token, err
}

func Login(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds UploadRequest
		if err := c.ShouldBind(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		expectedPassword, err := ReadConfig()
		if err != nil || expectedPassword.Token != creds.Token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Token: creds.Token,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.JwtKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Set token in cookie
		c.SetCookie("jwt", tokenString, int(expirationTime.Unix()), "/", "", false, true)

		c.String(http.StatusOK, `
      <form hx-post="/logout" hx-swap="outerHTML" hx-trigger="submit">
      <button type="submit">Logout</button>
      </form>
    `)
	}
}

func Logout(c *gin.Context) {
	// Remove cookie
	c.SetCookie("jwt", "", -1, "/", "", false, true)

	c.String(http.StatusOK, `
    <form hx-post="/login" hx-swap="outerHTML" hx-trigger="submit">
    <input type="text" placeholder="Token" name="token">
    <button type="submit">Login</button>
    </form>
  `)
}

func AuthMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := ValidateToken(config, tokenString)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
