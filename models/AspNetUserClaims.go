package models

type AspNetUserClaims struct {
	Id         int64  `json:"id" gorm:"primaryKey; column:Id"`
	UserId     string `json:"userId" gorm:"column:UserId"`
	ClaimType  string `json:"claimType" gorm:"column:ClaimType"`
	ClaimValue string `json:"claimValue" gorm:"column:ClaimValue"`
}

func (AspNetUserClaims) TableName() string {
	return "\"AspNetUserClaims\""
}
