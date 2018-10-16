package jwt

import (
	"encoding/json"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/bugfan/to"
	jwt "github.com/dgrijalva/jwt-go"
)

type conf struct {
	Method string // 加密算法
	Key    string // 加密key
	Issuer string // 签发者
	Expire int64  // 签名有效期
}

var Conf = conf{
	Method: "HS256",
	Key:    "sahjdjsgaudsiudhuywge99",
	Issuer: "thais",
	Expire: 60 * 60,
}

func GetJWT(data string) (string, error) {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + Conf.Expire),
		Issuer:    Conf.Issuer,
		Audience:  data,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(Conf.Key))
	if err != nil {
		return "", err
	}
	return ss, nil
}
func VerifyJWT(str string) (body string, err error) {
	t, err := jwt.Parse(str, func(*jwt.Token) (interface{}, error) {
		return []byte(Conf.Key), nil
	})
	if err != nil {
		return body, err
	}
	m := make(map[string]interface{})
	bs, err := json.Marshal(t.Claims)
	if err != nil {
		return body, err
	}
	err = json.Unmarshal(bs, &m)
	if err != nil {
		return body, err
	}
	return to.String(m["aud"]), nil
}

/*
*	JWS
 */

// VerifyJWT 验证json web token
func VerifyJWT2(token string) (ret bool, err error) {
	jwtObj, err := jws.ParseJWT([]byte(token))
	if err != nil {
		return
	}
	err = jwtObj.Validate([]byte(Conf.Key), jws.GetSigningMethod(Conf.Method))
	if err == nil {
		ret = true
	}
	return
}

// GetJWT 获取json web token
func GetJWT2(data map[string]interface{}) (token string, err error) {
	payload := jws.Claims{}
	for k, v := range data {
		payload.Set(k, v)
	}
	now := time.Now()
	payload.SetIssuer(Conf.Issuer)
	payload.SetIssuedAt(now)
	payload.SetExpiration(now.Add(time.Duration(Conf.Expire) * time.Minute))
	jwtObj := jws.NewJWT(payload, jws.GetSigningMethod(Conf.Method))
	tokenBytes, err := jwtObj.Serialize([]byte(Conf.Key))
	if err != nil {
		return
	}
	token = string(tokenBytes)
	return
}
