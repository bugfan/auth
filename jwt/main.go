package main

import (
	"auth/jwt/action"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/http/login", action.Login)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
