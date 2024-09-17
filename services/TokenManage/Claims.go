package tokenizer

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	UserId         uint   `json:"id"`
	SessionId      uint   `json:"session_id"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayname"`
	IsRefreshToken bool   `json:"isrefreshtoken"`
	// Role           string `json:"role"`
	jwt.StandardClaims
}

type ASPNetUserClaim struct {
	UserId         string `json:"userId"`
	SessionId      string `json:"sessionId"`
	Email          string `json:"email"`
	UserName       string `json:"username"`
	FullName       string `json:"fullname"`
	IsRefreshToken bool   `json:"isrefreshtoken"`
	jwt.StandardClaims
}
