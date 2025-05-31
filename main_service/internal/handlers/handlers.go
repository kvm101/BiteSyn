package handlers

import (
	"fmt"
	"log"
	"net/http"
	"restaurant_reviews/database"
	"restaurant_reviews/internal"
	"restaurant_reviews/internal/jwtAuth"
	"restaurant_reviews/internal/nlp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterHandler(c *gin.Context) {
	var user internal.User
	err := c.BindJSON(&user)
	if err != nil {
		return
	}

	insertResult, err := database.RegisterUser(user.Email, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"InsertID": insertResult})
}

func LoginHandler(c *gin.Context) {
	var user internal.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err = database.GetUser(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := jwtAuth.CreateToken(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func GetUserHandler(c *gin.Context) {
	email, role, err := jwtAuth.GetJWTDataFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": email, "role": role})
}

func FeedBackHandler(c *gin.Context) {
	var review internal.Review

	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if review.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review text is required"})
		return
	}
	if review.Rating < 0 || review.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 0 and 5"})
		return
	}

	nlpReview := nlp.CheckMessage(review.Text)
	if !nlpReview.Status {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "NLP service error",
			"details": nlpReview.TextReview,
		})
		return
	}

	review.CreatedAt = time.Now().UTC()

	if review.ID == "" {
		review.ID = primitive.NewObjectID().Hex()
	}

	mixRating := (review.Rating * 0.7) + (nlpReview.Rating * 0.3)

	if mixRating < 0 {
		mixRating = 0
	} else if mixRating > 5 {
		mixRating = 5
	}

	result, err := database.CreateFeedBack(review.ID, review.UserID, review.RestaurantID, review.Text, mixRating)
	if err != nil {
		log.Printf("Error creating feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"review":     result.ID,
		"rating":     result.Rating,
		"text":       result.Text,
		"created_at": result.CreatedAt,
	})
}

func DeleteUserHandler(c *gin.Context) {
	_, role, err := jwtAuth.GetJWTDataFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid ObjectID:", err)
		return
	}

	numId, _ := strconv.Atoi(objID.Hex())

	if role == "admin" {
		_, err := database.DeleteUser(numId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Deleted user": id})
		return

	}
}
