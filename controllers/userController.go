package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dev_mansoor/go-postgres-gorm/helpers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) (string, error) {
	// Implement password hashing logic here
	return "", nil
}

func VerifyPassword(hashedPassword, password string) bool {
	// Implement password verification logic here
	return false
}

// func Register() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var user models.User
// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 			return
// 		}

// 		validationErr := validate.Struct(user)
// 		if validationErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
// 			return
// 		}

// 		// Check email
// 		var count int64
// 		err := initializers.DB.WithContext(ctx).
// 			Model(&models.User{}).
// 			Where("email = ?", user.Email).
// 			Count(&count).Error
// 		if err != nil {
// 			log.Panic(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		}
// 		if count > 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "This email already exists"})
// 			return
// 		}

// 		// Check phone
// 		err = initializers.DB.WithContext(ctx).
// 			Model(&models.User{}).
// 			Where("phone = ?", user.Phone).
// 			Count(&count).Error
// 		if err != nil {
// 			log.Panic(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		}
// 		if count > 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "This phone already exists"})
// 			return
// 		}

// 		// TODO: Save user logic here
// 		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
// 	}
// }

func Register(c *gin.Context) {

	var body struct {
		First_name string `json:"first_name" binding:"required"`
		Last_name  string `json:"last_name" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=6"`
		Phone      string `json:"phone" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUserID := uuid.New().String()

	user := models.User{
		First_name: body.First_name,
		Last_name:  body.Last_name,
		Username:   body.Username,
		Email:      body.Email,
		Password:   string(hash),
		Phone:      body.Phone,
		User_id:    newUserID,
	}
	if err := initializers.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": user.ID})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days expiry
	})

	// Sign token
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600*24*30, "", "", false, true) // 30 days expiry

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		// "user": gin.H{
		// 	"id":       user.ID,
		// 	"email":    user.Email,
		// 	"username": user.Username,
		// },
		"token": tokenString,
	})
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		// Check if the user type matches
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		result := initializers.DB.WithContext(ctx).First(&user, userId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged in successful"})
}
