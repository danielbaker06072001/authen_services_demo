package tokenizer

import (
	"authen-service/appConfig/common"
	config "authen-service/appConfig/config"
	"encoding/base64"
	"errors"

	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Tokenizer interface {
	NewAccessToken(claims UserClaims) (string, error)
	ParseAccessToken(accessToken string) (*UserClaims, error)
	NewRefreshToken(claims UserClaims) (string, error)
	ParseRefreshToken(refreshToken string) *UserClaims
	VerifyToken(Token string) (*UserClaims, error)
	NewAspUserAccessToken(claims ASPNetUserClaim) (string, error)
	NewAspNetUserRefreshToken(claims ASPNetUserClaim) (string, error)
	ParseAspNetUserAccessToken(token string) (*ASPNetUserClaim, error)
	VerifyAspNetUserToken(token string) (*ASPNetUserClaim, error)
}

type tokenizerProvider struct {
	secretKey string
}

func NewTokenizerProvider(secretKey config.AppConfig) *tokenizerProvider {
	return &tokenizerProvider{secretKey: secretKey.JWTKey}
}

func (p *tokenizerProvider) NewAccessToken(claims UserClaims) (string, error) {
	secretKeyBytes, _ := base64.StdEncoding.DecodeString(p.secretKey)
	newClaims := &UserClaims{
		UserId:         claims.UserId,
		Email:          claims.Email,
		SessionId:      claims.SessionId,
		DisplayName:    claims.DisplayName,
		IsRefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(common.ACCESSTOKEN_DURATION).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	genToken, err := accessToken.SignedString(secretKeyBytes)
	if err != nil {
		log.Println("fail here 2:", err)
	}
	return genToken, nil
}

func (p *tokenizerProvider) ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	return parsedAccessToken.Claims.(*UserClaims), nil
}

func (p *tokenizerProvider) NewRefreshToken(claims UserClaims) (string, error) {
	secretKeyBytes, _ := base64.StdEncoding.DecodeString(p.secretKey)
	newClaims := &UserClaims{
		UserId:         claims.UserId,
		Email:          claims.Email,
		DisplayName:    claims.DisplayName,
		IsRefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return refreshToken.SignedString(secretKeyBytes)
}

func (p *tokenizerProvider) ParseRefreshToken(refreshToken string) *UserClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	return parsedRefreshToken.Claims.(*UserClaims)
}

func (p *tokenizerProvider) VerifyToken(token string) (*UserClaims, error) {

	claims, err := p.ParseAccessToken(token)
	if err != nil {
		// Handle errors returned by ParseAccessToken
		return nil, fmt.Errorf("failed to parse access token: %v", err)
	}
	if claims == nil {
		return nil, errors.New("invalid access token")
	}
	// Check if the claims are valid
	if err := claims.Valid(); err != nil {
		// Handle invalid claims
		return nil, fmt.Errorf("invalid access token: %v", err)
	}

	return claims, nil
}

func (p *tokenizerProvider) NewAspUserAccessToken(claims ASPNetUserClaim) (string, error) {
	secretKeyBytes, _ := base64.StdEncoding.DecodeString("mysecretkey")
	newClaims := &ASPNetUserClaim{
		UserId:         claims.UserId,
		SessionId:      claims.SessionId,
		UserName:       claims.UserName,
		Email:          claims.Email,
		IsRefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(common.ACCESSTOKEN_DURATION).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims)
	genToken, err := accessToken.SignedString(secretKeyBytes)
	if err != nil {
		log.Println("failed at generating access token: ", err)
	}

	return genToken, nil
}

func (p *tokenizerProvider) NewAspNetUserRefreshToken(claims ASPNetUserClaim) (string, error) {
	secretKeyBytes, _ := base64.StdEncoding.DecodeString(p.secretKey)
	newClaims := &ASPNetUserClaim{
		UserId:         claims.UserId,
		SessionId:      claims.SessionId,
		UserName:       claims.UserName,
		Email:          claims.Email,
		IsRefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims)
	genToken, err := accessToken.SignedString(secretKeyBytes)
	if err != nil {
		log.Println("failed at generating access token: ", err)
	}

	return genToken, nil
}
func (p *tokenizerProvider) ParseAspNetUserAccessToken(token string) (*ASPNetUserClaim, error) {
	parsedToken, _ := jwt.ParseWithClaims(token, &ASPNetUserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	return parsedToken.Claims.(*ASPNetUserClaim), nil
}

func (p *tokenizerProvider) VerifyAspNetUserToken(token string) (*ASPNetUserClaim, error) {
	claims, err := p.ParseAspNetUserAccessToken(token)
	if err != nil {
		// Handle errors returned by ParseAccessToken
		return nil, fmt.Errorf("failed to parse access token: %v", err)
	}
	if claims == nil {
		return nil, errors.New("invalid access token")
	}
	// Check if the claims are valid
	if err := claims.Valid(); err != nil {
		// Handle invalid claims
		return nil, fmt.Errorf("invalid access token: %v", err)
	}

	return claims, nil
}
