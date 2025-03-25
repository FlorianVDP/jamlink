package security

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"os"
	"time"

	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SecurityService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
	GenerateJWT(id *uuid.UUID, email *string, duration time.Duration, tokenType string) (string, error)
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
	GetJWTInfo(tokenString string) (uuid.UUID, error)
	GenerateSecureRandomString(n int) (string, error)
}

type securityService struct{}

func NewSecurityService() SecurityService {
	return &securityService{}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (s *securityService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", ErrPasswordHashing
	}

	return string(bytes), nil
}

func (s *securityService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func (s *securityService) GenerateJWT(id *uuid.UUID, email *string, duration time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
		"type": tokenType,
	}

	if id != nil {
		claims["id"] = id.String()
	}
	if email != nil {
		claims["email"] = *email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", ErrJWTGeneration
	}

	return tokenString, nil
}

func (s *securityService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWTSigningMethod
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrCannotExtractClaims
	}

	return claims, nil
}

func (s *securityService) ParseVerificationJWT(tokenStr string) (string, error) {
	claims, err := s.ValidateJWT(tokenStr)
	if err != nil {
		return "", err
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid token: missing email")
	}
	return email, nil
}

func (s *securityService) GetJWTInfo(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	rawID, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidUserID
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}

	return id, nil
}

func (s *securityService) GenerateSecureRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", ErrSecureRandomGeneration
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
