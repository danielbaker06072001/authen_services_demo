package OtherService

import (
	dto "authen-service/DTO"
	"authen-service/appConfig/common"
	"authen-service/infrastructure"
	"authen-service/models"
	"authen-service/repositories"
	tokenizer "authen-service/services/TokenManage"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	NORMAL_USER = "Normal User"
)

type IASPNetUserService interface {
	RegisterASPNetUser(user *models.AspNetUser) (*dto.LoginAspNetUserReturn, error)
	SetUpLoginASPNetUser(user *models.AspNetUser) (*dto.LoginAspNetUserReturn, error)
	CheckASPNetUserCredentials(identifier, password *string) (*models.AspNetUser, error)
	ChangeASPNetUserPassword(email string, newPassword string) error
	CheckASPNetUserRoles(userId string) error
}

type aspNetUserService struct {
	userCaching infrastructure.ICaching
	userRepo    repositories.IASPNetUser
}

func NewAspNetUserService(userCaching infrastructure.ICaching, userRepo repositories.IASPNetUser) IASPNetUserService {
	return &aspNetUserService{userCaching: userCaching, userRepo: userRepo}
}

func (s *aspNetUserService) RegisterASPNetUser(user *models.AspNetUser) (*dto.LoginAspNetUserReturn, error) {
	// returnValue := dto.LoginAspNetUserReturn{}
	isValidEmail := common.IsValidEmail(*user.Email)
	if !isValidEmail {
		log.Println("RegisterASPNetUser service: invalid email")
		return nil, errors.New("invalid email")
	}
	// * Insert + Create new User to AspNetUser table
	err := s.userRepo.Create(user)
	if err != nil {
		log.Println("error aspnet user service:", err)
		return nil, err
	}

	// * After registering account, set up login and generate (access token, refresh token, session id) for user
	// ? Does the coder want the flow to be: Create => Automatically logged in
	returnValue, err := s.SetUpLoginASPNetUser(user)
	if err != nil {
		log.Println("error aspnet user service:", err)
		return nil, err
	}

	return returnValue, nil
}

func (s *aspNetUserService) SetUpLoginASPNetUser(user *models.AspNetUser) (*dto.LoginAspNetUserReturn, error) {
	// Set up login for user based on AspNetUser
	returnValue := dto.LoginAspNetUserReturn{}
	aspNetUserSession := models.AspNetUserSession{
		UserProfileId: user.Id,
		UserName:      *user.UserName,
		ExpiredAt:     time.Now().Add(common.ACCESSTOKEN_DURATION),
	}

	// Save userSessionId to AspNetUserSession table
	// Later used for check user session
	// ! Saved in AspNetUserSession Database
	userSessionId, err := s.userRepo.SaveAspNetUserSession(&aspNetUserSession)
	if err != nil {
		log.Println("failed to save user session at set up steps: ", err)
		return nil, fmt.Errorf("failed to cache user session: %v", err)
	}

	// Set cache for user session
	// Logged in user will be saved to redis server and auto delete after certain amount of time
	// ! Saved in redis server
	// * Saved data type: UserSessionId_LoginSession
	_, errSaveSession := s.userCaching.SetAspNetUserSession(aspNetUserSession)
	if errSaveSession != nil {
		return nil, fmt.Errorf("failed to cache user session: %v", err)
	}

	// Generate token for user through userClaims
	userClaims := tokenizer.ASPNetUserClaim{
		UserId:    user.Id,
		SessionId: userSessionId,
		Email:     *user.Email,
		UserName:  *user.UserName,
	}

	// Generate tokens from userClaims (AspNetUserClaim) - using JWT
	accessToken, refreshToken, err := GenerateAspNetUserToken(userClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	returnValue.AccessToken = &accessToken
	returnValue.RefreshToken = &refreshToken
	returnValue.SessionId = &userSessionId
	return &returnValue, nil
}

func (s *aspNetUserService) CheckASPNetUserCredentials(identifier, password *string) (*models.AspNetUser, error) {
	data, err := s.userRepo.GetASPNetUserByPassword(identifier, password)
	if err != nil {
		log.Println("CheckASPNetUserCredentials:", err)
		return nil, errors.New("wrong password or email")
	}

	return data, nil
}

func (s *aspNetUserService) CheckASPNetUserRoles(userId string) error {
	err := s.userRepo.GetRoles(userId)
	if err != nil {
		log.Println("CheckASPNetUserRoles:", err)
		return errors.New("get roles failed")
	}
	return nil
}

func (s aspNetUserService) ChangeASPNetUserPassword(email string, newPassword string) error {
	err := s.userRepo.ChangePassword(email, newPassword)
	if err != nil {
		log.Println("ChangeASPNetUserPassword:", err)
		return errors.New("change password failed")
	}
	return nil
}
