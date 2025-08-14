package controllers

import (
	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers/models"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {

	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	c.Bind(&body)

	post := models.Post{
		Title: body.Title,
		Body:  body.Body,
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"body":    post,
		"message": "post created successfully",
	})
}

func PostIndex(c *gin.Context) {

	var posts []models.Post
	result := initializers.DB.Find(&posts)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"body":    posts,
		"message": "posts retrieved successfully",
	})
}
func PostShow(c *gin.Context) {

	id := c.Param("id")

	var posts models.Post
	result := initializers.DB.First(&posts, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"body":    posts,
		"message": "posts retrieved successfully",
	})
}
func PostUpdate(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	post.Title = body.Title
	post.Body = body.Body

	if err := initializers.DB.Save(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(200, gin.H{
		"body":    post,
		"message": "Post updated successfully",
	})
}
func PostDelete(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	if err := initializers.DB.Delete(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(200, gin.H{
		"body":    post,
		"message": "Post deleted successfully",
	})
}
