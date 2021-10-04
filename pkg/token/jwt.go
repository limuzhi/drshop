package token

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var (
	TokenExpired        = errors.New("Token is expired")
	TokenNotValidYet    = errors.New("Token not active yet")
	TokenMalformed      = errors.New("That's not even a token")
	TokenInvalid        = errors.New("Couldn't handle this token:")
	TokenPrivateKeyNill = errors.New("private key is nill")
)

// UserClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type UserClaims struct {
	Username string   `json:"username"`
	UserID   int64    `json:"uid"`
	RoleKey  []string `json:"rolekey"`
	RoleIds  []int64  `json:"roleIds"`
	jwt.StandardClaims
}

// Token jwt
type Token struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	// 缓存key
	CacheKey string
}

func New(privateKeyByte, publicKeyByte []byte) (*Token, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return &Token{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func NewPublic(publicKeyByte []byte) (*Token, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return nil, err
	}
	return &Token{
		publicKey: publicKey,
	}, nil
}

//ParseToken 解码
func (srv *Token) ParseToken(tokenStr string) (*UserClaims, error) {
	if srv.publicKey == nil {
		return nil, errors.New("private key is nill")
	}
	tokenStr = strings.ReplaceAll(tokenStr, "Bearer ", "")
	t, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.publicKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims, nil
	}
	return nil, err
}

// ParseToken 将 User 用户信息加密为 JWT 字符串
// expireTime := time.Now().Add(time.Hour * 24 * 3).Unix() 三天后过期
func (srv *Token) CreateToken(claims *UserClaims) (string, error) {
	if srv.privateKey == nil {
		return "", TokenPrivateKeyNill
	}
	claims.NotBefore = time.Now().Unix()
	claims.IssuedAt = time.Now().Unix()
	claims.Issuer = "drpshop.dev"
	claims.Subject = "AuthToken"

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tk, err := jwtToken.SignedString(srv.privateKey)
	if err != nil {
		return "", err
	}
	return tk, nil
}

// RefreshToken 更新token
func (srv *Token) RefreshToken(tokenString string) (string, error) {
	if srv.publicKey == nil {
		return "", errors.New("private key is nill")
	}
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.publicKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return srv.CreateToken(claims)
	}
	return "", TokenInvalid
}
