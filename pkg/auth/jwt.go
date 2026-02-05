package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: []byte(secret)}
}

func (s *JWTService) GenerateAccessToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// срок действия нашего токена с(issued) - до(expires)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
			}
			return s.secretKey, nil
		})

	if err != nil {
		return nil, err
	}

	// мы проводим сначала проверку на принадлежность нашего токена к типу *Claims, и если ок
	// только потом мы проверяем внутренне его валидность(по внутр. правилам)
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetRefreshTokenTTL возвращает время жизни refresh токена
func (s *JWTService) GetRefreshTokenTTL() time.Duration {
	return refreshTokenTTL
}
