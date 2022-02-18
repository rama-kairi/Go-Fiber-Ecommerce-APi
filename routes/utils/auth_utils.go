package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rama-kairi/fiber-api/config"
	"github.com/rama-kairi/fiber-api/database"
	"github.com/rama-kairi/fiber-api/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthUtils is a struct that contains all the functions that are used to authenticate the user.

type TokenType struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// HashPassword - hashes the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash - checks if the password matches the hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAccessToken - generates the access token.
func GenerateToken(userID uint, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    tokenType,
		"admin":   false,
		"iss":     "fiber-api",
		"iat":     time.Now().Unix(),
	}
	if tokenType == "access" {
		claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.GetConfig().Jwt.AccessExpireMin)).Unix()
	} else if tokenType == "refresh" {
		claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.GetConfig().Jwt.RefreshExpireMin)).Unix()
	} else {
		return "Please pass access or refresh in tokenType", nil
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(config.GetConfig().Jwt.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken - Validates the token.
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.GetConfig().Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, err
}

// DecodeToken - decodes the token.
func DecodeToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, isOk := token.Claims.(jwt.MapClaims)
	if _tokentype := claims["type"]; _tokentype != tokenType {
		return nil, fmt.Errorf("invalid token type: %s", tokenType)
	}

	if isOk && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// IsAuthenticated - checks if the user is authenticated.
func IsAuthenticated(username string, password string) (bool, models.User) {
	db := database.Database.Db

	var user models.User
	db.Where("email = ?", username).First(&user)
	if user.ID == 0 {
		return false, user
	}
	return CheckPasswordHash(password, user.Password), user
}

// PasswordValidator - validates the password.
func PasswordValidator(password string) (bool, string) {
	switch {
	case len(password) < 8:
		return false, "Password must be at least 8 characters long"
	case len(password) > 128:
		return false, "Password must be less than 128 characters long"
	case !regexp.MustCompile(`[A-Z]+`).MatchString(password):
		return false, "Password must contain at least one uppercase letter"
	case !regexp.MustCompile(`[a-z]+`).MatchString(password):
		return false, "Password must contain at least one lowercase letter"
	case !regexp.MustCompile(`\d+`).MatchString(password):
		return false, "Password must contain at least one number"
	case !regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`).MatchString(password):
		return false, "Password must contain at least one special character"
	default:
		return true, ""
	}
}
