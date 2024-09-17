package models

type AspNetUserRoles struct {
	UserId string `json:"userId" gorm:"column:UserId"`
	RoleId string `json:"roleId" gorm:"column:RoleId"`
}

func (AspNetUserRoles) TableName() string {
	return "\"AspNetUserRoles\""
}
