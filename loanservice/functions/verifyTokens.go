package functions

import (
	"github.com/golang-jwt/jwt"
)

func VerifyToken(tokenString string) (valid bool, customer_xid string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if token == nil {
		return
	}
	valid = token.Valid
	if !valid {
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	customer_xid = data["customer_xid"].(string)

	return
}
