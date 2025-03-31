package middleware

import (
	"context"
	"itv/monorepo/api_gateway/configs"
	"itv/monorepo/api_gateway/constants"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// JWT ni tekshirish funksiyasi
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(constants.AuthorizationHeader)
		if authHeader == "" {
			http.Error(w, "Unauthorized: No Token", http.StatusUnauthorized)
			return
		}

		// "Bearer TOKEN" dan faqat TOKEN qismini olish
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		conf := configs.Config()
		// Tokenni tekshirish
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JWTSecretKey), nil // Bu joyda JWT maxfiy kalitingiz bo‘lishi kerak
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		// Tokenning claims (payload) qismini olish
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid Token Claims", http.StatusUnauthorized)
			return
		}

		// `user_id` ni olish
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized: User ID Not Found", http.StatusUnauthorized)
			return
		}

		fmt.Println("Authenticated user ID:", userID)

		// Contextga `userID` ni qo‘shamiz
		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
