package dto

// type RegisterASPNetUser struct{}

type RegisterASPNetUserRequest struct {
	UserName    *string
	PhoneNumber *string
	FullName    *string
	Email       *string
	Password    *string
	UserRole    *string
}

type RegisterASPNetUserReturn struct {
	AccessToken  *string
	RefreshToken *string
	SessionId    *string
}
type RegisterASPNetUserResponse struct {
	CommonResponse
	RegisterASPNetUserReturn
}
