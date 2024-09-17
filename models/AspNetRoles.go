package models

type AspNetRoles struct {
	ID               string  `json:"id" gorm:"primaryKey; column:Id"`
	Name             *string `json:"name" gorm:"column:Name"`
	NormalizedName   *string `json:"normalizedName" gorm:"column:NormalizedName"`
	ConcurrencyStamp *string `json:"concurrencyStamp" gorm:"column:ConcurrencyStamp"`
}

func (AspNetRoles) TableName() string {
	return "\"AspNetRoles\""
}
