package infrastructure

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_ACCESS_TOKEN_SECRET = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
var JWT_REFRESH_TOKEN_SECRET = os.Getenv("JWT_REFRESH_TOKEN_SECRET")


func GenerateAccessToken(userID,role string)(string ,error){
	claims:=jwt.MapClaims{
		"user_id":userID,
		"role":role,
		"exp":time.Now().Add(time.Minute * 15).Unix(),
	}


	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_ACCESS_TOKEN_SECRET)
}

func GenerateRefreshToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JWT_REFRESH_TOKEN_SECRET)
}

func VerifyRefreshToken(tokenStr string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return JWT_REFRESH_TOKEN_SECRET, nil
    })
    if err != nil || !token.Valid {
        return "", err
    }

    claims := token.Claims.(jwt.MapClaims)
    return claims["user_id"].(string), nil
}