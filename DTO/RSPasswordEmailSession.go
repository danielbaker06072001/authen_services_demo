package dto

type LinkSessionReq struct {
	Link      string
	UserEmail string `json:"userEmail"`
}

type LinkSessionResponse struct {
	CommonResponse
	SessionId string
}

type ResetPasswordLinkResponse struct {
	CommonResponse
	Link string
}

type ChangePasswordRequest struct {
	UserEmail   string `json:"userEmail"`
	NewPassword string `json:"newPassword"`
}
