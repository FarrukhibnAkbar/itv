package utils

import (
	"crypto/sha256"
	"itv/monorepo/api_gateway/configs"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT uchun maxfiy kalit
var conf = configs.Config()

type Tokens struct {
	Access  string
	Refresh string
}

// GenerateNewStaffTokens func for generate a new Access & Refresh tokens.
func GenerateNewTokens(credentials map[string]string) (*Tokens, error) {
	// Generate JWT Access token.
	accessToken, err := generateNewAccessToken(credentials)
	if err != nil {
		return nil, err
	}

	// Generate JWT Refresh token.
	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

// Token yaratish funksiyasi
func generateNewAccessToken(credentials map[string]string) (string, error) {

	// Token uchun payload (claims)
	claims := jwt.MapClaims{
		"user_id": credentials["user_id"],
		"role":    credentials["role"],
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Tokenning muddati (24 soat)
	}

	// JWTni yaratish
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tokenni imzolash (sign) va string ko‘rinishda qaytarish
	t, err := token.SignedString([]byte(conf.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("error in SignetString of generate token %v", err)
	}

	return t, err
}

func generateNewRefreshToken() (string, error) {
	// Create a new SHA256 hash.
	sha256 := sha256.New()
	conf := configs.Config()

	// Create a new now date and time string with salt.
	refresh := conf.JWTSecretKey + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	_, err := sha256.Write([]byte(refresh))
	if err != nil {
		// Return error, it refresh token generation failed.
		return "", err
	}

	// Set expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(24)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hex.EncodeToString(sha256.Sum(nil)) + "." + expireTime

	return t, nil
}

type UserClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GetUserMetadata extracts user metadata from JWT token
func GetUserMetadata(signingKey, tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, errors.New("authorization token missing")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims if valid
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return map[string]interface{}{
			"id":   claims.ID,
			"role": claims.Role,
		}, nil
	}

	return nil, errors.New("invalid token")
}