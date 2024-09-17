package dto

type LoginWithEmailRequest struct {
	Email    *string
	UserName *string
	Password *string
}

type LoginReturn struct {
	AccessToken  *string
	RefreshToken *string
	SessionId    uint32
}

type LoginAspNetUserReturn struct {
	AccessToken  *string
	RefreshToken *string
	SessionId    *string
}
type LoginResponse struct {
	CommonResponse
	LoginAspNetUserReturn
}

type TokenData struct {
	UserId         string  `json:"user_id"`
	SessionId      string  `json:"session_id"`
	Email          string  `json:"email"`
	UserName       string  `json:"username"`
	FullName       string  `json:"fullname"`
	IsRefreshToken string  `json:"isrefreshtoken"`
	Exp            float64 `json:"exp"`
	Iat            float64 `json:"iat"`
}
