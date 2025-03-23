package security

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SecurityService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
	GenerateJWT(id uuid.UUID) (string, error)
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
	GetJWTInfo(tokenString string) (uuid.UUID, error)
}

type securityService struct{}

func NewSecurityService() SecurityService {
	return &securityService{}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (s *securityService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *securityService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *securityService) GenerateJWT(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func (s *securityService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot extract claims")
	}

	return claims, nil
}

func (s *securityService) GetJWTInfo(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
