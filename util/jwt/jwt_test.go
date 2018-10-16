package jwt

import (
	"encoding/json"
	"log"
	"testing"
)

func TestJwt(t *testing.T) {
	m := make(map[string]interface{})
	m["name"] = "test"
	m["age"] = 25
	bs, _ := json.Marshal(m)

	token, err := GetJWT(string(bs))
	log.Println("GetJWT:", err, token)
	log.Println()
	body, err := VerifyJWT(token)
	log.Println("VerifyJWT:", err, string(body))
}
