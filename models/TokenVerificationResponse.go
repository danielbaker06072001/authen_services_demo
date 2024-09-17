package models

import dto "authen-service/DTO"

type ValidateCommonResult struct {
	IsValid        bool     `json:"isValid"`
	ValidateFields []string `json:"validateFields"`
}

type TokenVerificationResponse struct {
	Result         string               `json:"result"`
	ErrorMessage   string               `json:"errorMessage"`
	ValidateResult ValidateCommonResult `json:"validateResult"`
	Role           []string             `json:"role"`
	Status         int                  `json:"status"`
	Data           *dto.TokenData       `json:"data"`
}
