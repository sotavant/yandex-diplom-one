package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sotavant/yandex-diplom-one/internal"
	"time"
)

type claims struct {
	jwt.RegisteredClaims
	UserId int64
}

const tokenExp = time.Hour * 3
const secretKey = "someSecretSuperKey"

func BuildJWTString(userId int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserId: userId,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserId(tokenString string) (int64, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		internal.Logger.Infow("error in parse token", "err", err)
		return -1, err
	}

	if !token.Valid {
		return -1, nil
	}

	return claims.UserId, nil
}
