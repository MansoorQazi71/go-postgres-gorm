package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(c *gin.Context) {

}

// func RequireAuth(c *gin.Context) {
// 	tokenString, err := c.Cookie("token")
// 	if err != nil {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

// 	if err != nil || !token.Valid {
// 		log.Println("Invalid token:", err)
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 		// Expiry check
// 		exp := int64(claims["exp"].(float64))
// 		if time.Now().Unix() > exp {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}

// 		// Fetch user from DB
// 		var user models.User
// 		initializers.DB.First(&user, claims["sub"])

// 		if user.ID == 0 {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}

// 		// Attach user to context
// 		c.Set("user", user)
// 		c.Next()
// 	} else {
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 	}
// }

func RequireAuth(c *gin.Context) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
		return
	}

	// Extract the token part
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		log.Println("Invalid token:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Expiry check
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Fetch user from DB
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Attach user to context
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	}
}
