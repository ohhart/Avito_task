package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenService представляет сервис для работы с токенами JWT
type TokenService struct {
	jwtKey []byte
}

// NewTokenService создает новый экземпляр TokenService
func NewTokenService(jwtKey []byte) *TokenService {
	return &TokenService{
		jwtKey: jwtKey,
	}
}

// GenerateToken генерирует JWT токен
func (ts *TokenService) GenerateToken(userID int, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа

	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(ts.jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken разбирает и верифицирует JWT токен
func (ts *TokenService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return ts.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("неверный токен")
	}
}

// VerifyAdminToken проверяет, является ли токен админским
func (ts *TokenService) VerifyAdminToken(tokenString string) error {
	claims, err := ts.ParseToken(tokenString)
	if err != nil {
		return err
	}

	if claims.Role != "admin" {
		return errors.New("требуется администраторский токен")
	}

	return nil
}
