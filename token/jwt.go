package token

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

var (
	ErrTokenInvalid = errors.New("couldn't handle this token")
)

type Claims struct {
	jwt.RegisteredClaims
	UserId int64  // 应用内部用户id
	SignId string // 单点用户id
}

type AuthToken struct {
	SigningKey []byte                 // 密钥
	Header     map[string]interface{} // 头信息
	Raw        string                 // 加密token
}

func NewAuthToken(secret []byte) *AuthToken {
	return &AuthToken{
		SigningKey: secret,
	}
}

// CreateToken 创建一个token
func (p *AuthToken) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.SigningKey)
}

// ParseFromRequest 解析token
func (p *AuthToken) ParseFromRequest(req *http.Request) (*Claims, error) {
	tokenString, err := request.OAuth2Extractor.ExtractToken(req)
	if err != nil {
		return nil, err
	}

	return p.ParseToken(tokenString)
}

func (p *AuthToken) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return p.SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	p.Header = token.Header
	p.Raw = token.Raw

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
