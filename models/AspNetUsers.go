package models

import "time"

// authen user
type AspNetUser struct {
	Id                          string    `json:"id"  gorm:"primaryKey;  column:Id"`
	AccessFailedCount           int64     `json:"accessFailedCount" gorm:"column:AccessFailedCount"`
	AppleRegisterTime           time.Time `json:"appleRegisterTime" gorm:"column:AppleRegisterTime"`
	IsDeactive                  bool      `json:"isDeactive" gorm:"column:IsDeactive"`
	TikTokRegisterTime          time.Time `json:"tikTokRegisterTime" gorm:"column:TikTokRegisterTime"`
	FacebookRegisterTime        time.Time `json:"facebookRegisterTime" gorm:"column:FacebookRegisterTime"`
	InstagramRegisterTime       time.Time `json:"instagramRegisterTime" gorm:"column:InstagramRegisterTime"`
	GoogleRegisterTime          time.Time `json:"googleRegisterTime" gorm:"column:GoogleRegisterTime"`
	DiscordRegisterTime         time.Time `json:"discordRegisterTime" gorm:"column:DiscordRegisterTime"`
	LinkedInRegisterTime        time.Time `json:"linkedInRegisterTime" gorm:"column:LinkedInRegisterTime"`
	EmailConfirmed              bool      `json:"emailConfirmed" gorm:"column:EmailConfirmed"`
	PhoneNumberConfirmed        bool      `json:"phoneNumberConfirmed" gorm:"column:PhoneNumberConfirmed"`
	TwoFactorEnabled            bool      `json:"twoFactorEnabled" gorm:"column:TwoFactorEnabled"`
	LockoutEnd                  time.Time `json:"lockoutEnd" gorm:"column:LockoutEnd"`
	TelegramRegisterTime        time.Time `json:"telegramRegisterTime" gorm:"column:TelegramRegisterTime"`
	ZaloRegisterTime            time.Time `json:"zaloRegisterTime" gorm:"column:ZaloRegisterTime"`
	LockoutEnabled              bool      `json:"lockoutEnabled" gorm:"column:LockoutEnabled"`
	IsDeactiveAccount           bool      `json:"isDeactiveAccount" gorm:"column:IsDeactiveAccount"`
	FacebookUserFullName        *string    `json:"facebookUserFullName" gorm:"column:FacebookUserFullName"`
	FacebookUserGivenName       *string    `json:"facebookUserGivenName" gorm:"column:FacebookUserGivenName"`
	FacebookUserId              *string    `json:"facebookUserId" gorm:"column:FacebookUserId"`
	FacebookUserPhoneCode       *string    `json:"facebookUserPhoneCode" gorm:"column:FacebookUserPhoneCode"`
	FacebookUserPhoneNumber     *string    `json:"facebookUserPhoneNumber" gorm:"column:FacebookUserPhoneNumber"`
	FacebookUserProfileImageUrl *string    `json:"facebookUserProfileImageUrl" gorm:"column:FacebookUserProfileImageUrl"`
	FacebookUserSurname         *string    `json:"facebookUserSurname" gorm:"column:FacebookUserSurname"`
	GoogleUserEmail             *string    `json:"googleUserEmail" gorm:"column:GoogleUserEmail"`
	GoogleUserFullName          *string    `json:"googleUserFullName" gorm:"column:GoogleUserFullName"`
	GoogleUserGivenName         *string    `json:"googleUserGivenName" gorm:"column:GoogleUserGivenName"`
	GoogleUserId                *string    `json:"googleUserId" gorm:"column:GoogleUserId"`
	GoogleUserPhoneCode         *string    `json:"googleUserPhoneCode" gorm:"column:GoogleUserPhoneCode"`
	GoogleUserPhoneNumber       *string    `json:"googleUserPhoneNumber" gorm:"column:GoogleUserPhoneNumber"`
	GoogleUserProfileImageUrl   *string    `json:"googleUserProfileImageUrl" gorm:"column:GoogleUserProfileImageUrl"`
	GoogleUserSurname           *string    `json:"googleUserSurname" gorm:"column:GoogleUserSurname"`
	OpenId                      *string    `json:"openId" gorm:"column:OpenId"`
	LinkedInUserEmail           *string    `json:"linkedInUserEmail" gorm:"column:LinkedInUserEmail"`
	LinkedInUserGivenName       *string    `json:"linkedInUserGivenName" gorm:"column:LinkedInUserGivenName"`
	LinkedInUserId              *string    `json:"linkedInUserId" gorm:"column:LinkedInUserId"`
	LinkedInUserPhoneCode       *string    `json:"linkedInUserPhoneCode" gorm:"column:LinkedInUserPhoneCode"`
	LinkedInUserPhoneNumber     *string    `json:"linkedInUserPhoneNumber" gorm:"column:LinkedInUserPhoneNumber"`
	LinkedInUserProfileImageUrl *string    `json:"linkedInUserProfileImageUrl" gorm:"column:LinkedInUserProfileImageUrl"`
	LinkedInUserSurname         *string    `json:"linkedInUserSurname" gorm:"column:LinkedInUserSurname"`
	TikTokUserEmail             *string    `json:"tikTokUserEmail" gorm:"column:TikTokUserEmail"`
	TikTokUserGivenName         *string    `json:"tikTokUserGivenName" gorm:"column:TikTokUserGivenName"`
	TikTokUserId                *string    `json:"tikTokUserId" gorm:"column:TikTokUserId"`
	TikTokUserPhoneCode         *string    `json:"tikTokUserPhoneCode" gorm:"column:TikTokUserPhoneCode"`
	TikTokUserPhoneNumber       *string    `json:"tikTokUserPhoneNumber" gorm:"column:TikTokUserPhoneNumber"`
	TikTokUserSurname           *string    `json:"tikTokUserSurname" gorm:"column:TikTokUserSurname"`
	ZaloUserEmail               *string    `json:"zaloUserEmail" gorm:"column:ZaloUserEmail"`
	ZaloUserGivenName           *string    `json:"zaloUserGivenName" gorm:"column:ZaloUserGivenName"`
	ZaloUserId                  *string    `json:"zaloUserId" gorm:"column:ZaloUserId"`
	ZaloUserPhoneCode           *string    `json:"zaloUserPhoneCode" gorm:"column:ZaloUserPhoneCode"`
	ZaloUserPhoneNumber         *string    `json:"zaloUserPhoneNumber" gorm:"column:ZaloUserPhoneNumber"`
	ZaloUserProfileImageUrl     *string    `json:"zaloUserProfileImageUrl" gorm:"column:ZaloUserProfileImageUrl"`
	ZaloUserSurname             *string    `json:"zaloUserSurname" gorm:"column:ZaloUserSurname"`
	DiscordUserEmail            *string    `json:"discordUserEmail" gorm:"column:DiscordUserEmail"`
	DiscordUserGivenName        *string    `json:"discordUserGivenName" gorm:"column:DiscordUserGivenName"`
	DiscordUserId               *string    `json:"discordUserId" gorm:"column:DiscordUserId"`
	DiscordUserPhoneCode        *string    `json:"discordUserPhoneCode" gorm:"column:DiscordUserPhoneCode"`
	DiscordUserPhoneNumber      *string    `json:"discordUserPhoneNumber" gorm:"column:DiscordUserPhoneNumber"`
	DiscordUserProfileImageUrl  *string    `json:"discordUserProfileImageUrl" gorm:"column:DiscordUserProfileImageUrl"`
	DiscordUserSurname          *string    `json:"discordUserSurname" gorm:"column:DiscordUserSurname"`
	InstagramEmail              *string    `json:"instagramEmail" gorm:"column:InstagramEmail"`
	InstagramGivenName          *string    `json:"instagramGivenName" gorm:"column:InstagramGivenName"`
	InstagramPhoneCode          *string    `json:"instagramPhoneCode" gorm:"column:InstagramPhoneCode"`
	InstagramPhoneNumber        *string    `json:"instagramPhoneNumber" gorm:"column:InstagramPhoneNumber"`
	InstagramProfileImageUrl    *string    `json:"instagramProfileImageUrl" gorm:"column:InstagramProfileImageUrl"`
	InstagramSurname            *string    `json:"instagramSurname" gorm:"column:InstagramSurname"`
	InstagramUserId             *string    `json:"instagramUserId" gorm:"column:InstagramUserId"`
	TelegramEmail               *string    `json:"telegramEmail" gorm:"column:TelegramEmail"`
	TelegramGivenName           *string    `json:"telegramGivenName" gorm:"column:TelegramGivenName"`
	TelegramPhoneCode           *string    `json:"telegramPhoneCode" gorm:"column:TelegramPhoneCode"`
	TelegramPhoneNumber         *string    `json:"telegramPhoneNumber" gorm:"column:TelegramPhoneNumber"`
	TelegramProfileImageUrl     *string    `json:"telegramProfileImageUrl" gorm:"column:TelegramProfileImageUrl"`
	TelegramSurname             *string    `json:"telegramSurname" gorm:"column:TelegramSurname"`
	TelegramUserId              *string    `json:"telegramUserId" gorm:"column:TelegramUserId"`
	TikTokUserProfileImageUrl   *string    `json:"tikTokUserProfileImageUrl" gorm:"column:TikTokUserProfileImageUrl"`
	UserName                    *string    `json:"userName" gorm:"column:UserName"`
	NormalizedUserName          *string    `json:"normalizedUserName" gorm:"column:NormalizedUserName"`
	Email                       *string    `json:"email" gorm:"column:Email"`
	NormalizedEmail             *string    `json:"normalizedEmail" gorm:"column:NormalizedEmail"`
	PasswordHash                *string    `json:"passwordHash" gorm:"column:PasswordHash"`
	SecurityStamp               *string    `json:"securityStamp" gorm:"column:SecurityStamp"`
	ConcurrencyStamp            *string    `json:"concurrencyStamp" gorm:"column:ConcurrencyStamp"`
	PhoneNumber                 *string    `json:"phoneNumber" gorm:"column:PhoneNumber"`
	AppleUserEmail              *string    `json:"appleUserEmail" gorm:"column:AppleUserEmail"`
	AppleUserGivenName          *string    `json:"appleUserGivenName" gorm:"column:AppleUserGivenName"`
	AppleUserId                 *string    `json:"appleUserId" gorm:"column:AppleUserId"`
	AppleUserPhoneCode          *string    `json:"appleUserPhoneCode" gorm:"column:AppleUserPhoneCode"`
	AppleUserPhoneNumber        *string    `json:"appleUserPhoneNumber" gorm:"column:AppleUserPhoneNumber"`
	AppleUserProfileImageUrl    *string    `json:"appleUserProfileImageUrl" gorm:"column:AppleUserProfileImageUrl"`
	AppleUserSurname            *string    `json:"appleUserSurname" gorm:"column:AppleUserSurname"`
	FacebookUserEmail           *string    `json:"facebookUserEmail" gorm:"column:FacebookUserEmail"`
}

// Model này sử dụng trong việc Update
func (AspNetUser) TableName() string {
	return "\"AspNetUsers\""
}
