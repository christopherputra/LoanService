package functions

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = "8UYAJSGFU124G"

func GenerateToken(customerXid string) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"data": map[string]string{
			"customer_xid": customerXid,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
