package server

import (
	"net/http"
	"os"
)

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
