package jwtAuth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var secret = []byte("karpliak vasyl")

func CreateToken(email string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTDecode(tokenString string) (string, string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("не вдалося розпарсити токен: %w", err)
	}

	emailRaw, okEmail := claims["email"]
	roleRaw, okRole := claims["role"]

	if !okEmail || !okRole {
		return "", "", fmt.Errorf("відсутні обов'язкові поля в токені")
	}

	email, okEmailCast := emailRaw.(string)
	role, okRoleCast := roleRaw.(string)

	if !okEmailCast || !okRoleCast {
		return "", "", fmt.Errorf("не вдалося привести типи до string")
	}

	return email, role, nil
}

func GetJWTDataFromCookie(c *gin.Context) (string, string, error) {
	tokenString, err := c.Request.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	email, role, err := JWTDecode(tokenString.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
	}

	return email, role, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
