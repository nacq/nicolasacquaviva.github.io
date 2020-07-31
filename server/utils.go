package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid && claims["name"] == os.Getenv("JWT_NAME") {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	})
}

func CheckOrigin(r *http.Request) bool {
	if IsProduction() {
		origin := r.Header.Get("Origin")

		return origin == "https://nicolasacquaviva.com" || origin == "https://www.nicolasacquaviva.com"
	}

	return true
}

func GetIPFromRequest(r *http.Request) string {
	forwardedFor := r.Header.Get("x-forwarded-for")

	if forwardedFor != "" {
		return forwardedFor
	}

	return r.RemoteAddr
}

func IsProduction() bool {
	return os.Getenv("MODE") == "production"
}
