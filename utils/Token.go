package utils

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// save sign key in env
var signedString = []byte("mysupersecretkey")

func TokenGenerator(s string) (string, error)  {
	setToken := jwt.New(jwt.SigningMethodHS256)
	claims := setToken.Claims.(jwt.MapClaims)

	claims["id"] = s

	token, err := setToken.SignedString(signedString)

	if err != nil {
		return "", err
	}

	return token, nil
}

func TokenDecrypt (s string) (bool, string) {
	
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(s, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(signedString),nil
	})

	if err != nil {
		fmt.Println(err)
		return false, "Token is not valid"
	}
	
	s = fmt.Sprintf("%v", claims["id"])

	return true, s
}