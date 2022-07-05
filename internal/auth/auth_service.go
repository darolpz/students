package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var ErrSignToken = errors.New("couldn't signing token")

type IAuthService interface {
	GenerateJWT(email, role string) (string, error)
	GenerateHashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type authService struct {
	secretkey string
}

func NewAuthService(secret string) IAuthService {
	return authService{secretkey: secret}
}

func (s authService) GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(s.secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrSignToken, err)
	}
	return tokenString, nil
}

func (s authService) GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s authService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
