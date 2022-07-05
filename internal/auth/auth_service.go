package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSignToken = errors.New("couldn't signing token")

	ErrParseToken   = errors.New("couldn't parse token")
	ErrInvalidToken = errors.New("token is not valid")
)

type JWTClaim struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type IAuthService interface {
	GenerateJWT(email, role string) (string, error)
	GenerateHashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	CheckToken(signedToken string) (*jwt.Token, error)
}

type authService struct {
	secretkey string
}

func NewAuthService(secret string) IAuthService {
	return authService{secretkey: secret}
}

func (s authService) GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(s.secretkey)
	claims := &JWTClaim{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

func (s authService) CheckToken(signedToken string) (*jwt.Token, error) {
	var mySigningKey = []byte(s.secretkey)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return t, fmt.Errorf("there was an error in parsing")
			}
			return []byte(mySigningKey), nil
		})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}
