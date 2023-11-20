package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xyproto/randomstring"
)

var ErrInvalidToken = fmt.Errorf("invalid token")

type (
	Manager struct {
		secretKey []byte
		mu        *sync.RWMutex
	}

	JWTClaims struct {
		jwt.StandardClaims
		UserID int64 `json:"userID"`
	}
)

func NewManager(secretKey string) *Manager {
	if secretKey == "" {
		secretKey = randomstring.CookieFriendlyString(16)
	}
	return &Manager{
		mu:        &sync.RWMutex{},
		secretKey: []byte(secretKey),
	}
}

func (m *Manager) GenerateJWT(userID int64) (string, error) {
	// Создание пользовательских клеймов
	claims := JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "otus_user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписание токена с использованием секретного ключа
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) CheckToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
