package tokens

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenService interface {
	GenerateAccessToken(claims JWTCustomClaims) (string, error)
	ValidateToken(tokenString string) (*JWTCustomClaims, error)
}

type tokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) TokenService {
	return &tokenService{secretKey}
}

func (t *tokenService) GenerateAccessToken(claims JWTCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encoded, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func (t *tokenService) ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
