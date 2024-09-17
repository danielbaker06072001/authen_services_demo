package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	dto "authen-service/DTO"
	"authen-service/appConfig/common"
	"authen-service/models"
	"authen-service/pb"
	"authen-service/repositories"
	OtherService "authen-service/services/Other"
	tokenizer "authen-service/services/TokenManage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	AUTHORIZATION_HEADER = "authorization"
	AUTHORIZATION_BEARER = "bearer"
	STATUS_UNAUTHORIZED  = "unauthorized"
	STATUS_AUTHORIZED    = "authorized"
	SUCCESS              = "success"
	SESSION_VERIFIED     = "session verified"
	SESSION_EXPIRED      = "session expired"
	NONE                 = "none"
	SEND_EMAIL_AGAIN     = "send email again"
	BAD_REQUEST          = "bad request"
	MISSING_SESSION      = "missing session"

	USER_ROLE_NORMAL_USER  = "NORMAL_USER"
	USER_ROLE_SYSTEM_ADMIN = "SYSTEM ADMINISTRATOR"
	USER_ROLE_SYSTEM_USER  = "SYSTEM USER"
	USER_ROLE_KOL_USER     = "KOL USER"
	USER_ROLE_CLIENT_USER  = "CLIENT USER"
)

// * This is fucntion handle the Register Account
// *Input : Email, Password, PhoneNumber, UserName
// *Return: Result, ErrorMessage (common return type)
func (s *Server) RegisterASPNetUser(ctx context.Context, req *pb.ASPNetUserRegisterRequest) (*pb.ASPNetUserRegisterResponse, error) {
	// *Initialize database information
	db := s.db_authen
	redisCache := s.redisCache
	resp := &pb.ASPNetUserRegisterResponse{}

	// Set default mssage for result and ErrorMessage
	// TODO: Later change this to cosnt or enum (not fix string)
	result := "register successful"
	errorMessage := "none"

	// *Initialize repositories and services for ASPNetUser
	aspNetUserRepo := repositories.NewAspNetUserRepository(db)
	aspNetUserService := OtherService.NewAspNetUserService(redisCache, aspNetUserRepo)

	// *Set up the new ASPNetUser object
	newAspUser := &models.AspNetUser{
		UserName:          &req.UserName,
		PhoneNumber:       &req.PhoneNumber,
		Email:             &req.Email,
		PasswordHash:      &req.Password,
		AccessFailedCount: 0,
	}

	// * Get username from email if username is empty
	if req.UserName != "" {
		*newAspUser.UserName = OtherService.ExtractNameFromEmail(req.Email)
	}

	// *Register new ASPNetUser (insert into table AspNetUsers)
	_, err := aspNetUserService.RegisterASPNetUser(newAspUser)
	if err != nil {
		log.Println("RegisterASPNetUser:", err)
		result = "cannot creat new User"
		errorMessage = err.Error()
		resp.Result = result
		resp.ErrorMessage = errorMessage
		return resp, nil
	}

	resp.ErrorMessage = errorMessage
	resp.Result = result
	return resp, nil
}

// *This function log User in with Email/Username + password
// *Input: Email/Username, Password (string)
// *Return: AccessToken, RefreshToken, SessionID
func (s *Server) LoginWithEmailByGin(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Initialize the login Request object
		// Use DTO to bind the request body to the object
		var loginRequest *dto.LoginWithEmailRequest
		if err := c.ShouldBind(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Check binding
		fmt.Println("Timestamp: shouldBind ", time.Now())

		// Check if any of the field is not valid, return error
		var identifier *string
		if loginRequest.Email != nil {
			identifier = loginRequest.Email
		} else if loginRequest.UserName != nil {
			identifier = loginRequest.UserName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neither username nor email cannot empty"})
			return
		}

		// Initialize the login response object
		loginResponse := &dto.LoginResponse{
			CommonResponse:        dto.CommonResponse{},
			LoginAspNetUserReturn: dto.LoginAspNetUserReturn{},
		}

		// Initialize the redis cache and the user service
		redisCache := s.redisCache

		// Initialize the user repository
		userRepo := repositories.NewAspNetUserRepository(db)
		userService := OtherService.NewAspNetUserService(redisCache, userRepo)

		// Check binding
		fmt.Println("Timestamp: initialize user repository ", time.Now())

		// Check if the user is valid
		user, err := userService.CheckASPNetUserCredentials(identifier, loginRequest.Password)
		if err != nil {
			loginResponse.CommonResponse.Result = "Unauthorized"
			loginResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusUnauthorized, loginResponse)
			return
		}

		// *Check the user role, if user is NORMAL USER then they won't be able to login
		// if err = userService.CheckASPNetUserRoles(user.Id); err != nil {
		// 	loginResponse.CommonResponse.Result = "Unauthorized - Can't login as normal user"
		// 	loginResponse.CommonResponse.ErrorMessage = err.Error()
		// 	c.JSON(http.StatusUnauthorized, loginResponse)
		// 	return
		// }

		// Check binding
		fmt.Println("Timestamp: checkUserCredentials ", time.Now())

		// *Set up the login ASPNet user
		// *Input: AspNetUSer object
		// ? What this do: SetUpLoginAspNetUser
		//	- set up login for user using AspNetUser
		// 	- Save user session into database + redis (AspNetUserSession table + Redis_13)
		// *Return: AccessToken, RefreshToken, SessionID
		loginResult, err := userService.SetUpLoginASPNetUser(user)
		if err != nil {
			loginResponse.CommonResponse.Result = "Internal Server Error"
			loginResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusInternalServerError, loginResponse)
			return
		}

		// Check binding
		fmt.Println("Timestamp: setupLoginASp + create token", time.Now())

		// Reponse back to the user
		// *Return: AccessToken, RefreshToken, SessionID
		loginResponse.CommonResponse.Result = "Success"
		loginResponse.LoginAspNetUserReturn.AccessToken = loginResult.AccessToken
		loginResponse.LoginAspNetUserReturn.RefreshToken = loginResult.RefreshToken
		loginResponse.LoginAspNetUserReturn.SessionId = loginResult.SessionId
		c.JSON(http.StatusOK, loginResponse)
		return
	}
}

func (s *Server) RegisterASPNetUserGin(db_authen *gorm.DB, db_data *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Initialize the register request object (what input do we need from user)
		var registerASPNetUserReq *dto.RegisterASPNetUserRequest

		if err := c.ShouldBind(&registerASPNetUserReq); err != nil {
			log.Println("UserName before pass to db: ", *registerASPNetUserReq.UserName)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		// Initialize the register response object (typically hold token + error message)
		registerASPNetUserResponse := dto.RegisterASPNetUserResponse{
			CommonResponse:           dto.CommonResponse{},
			RegisterASPNetUserReturn: dto.RegisterASPNetUserReturn{},
		}

		if registerASPNetUserReq.UserName == nil {
			*registerASPNetUserReq.UserName = OtherService.ExtractNameFromEmail(*registerASPNetUserReq.Email)
		}

		// Set data from user's input to the new ASPNetUser object
		newUser := &models.AspNetUser{
			UserName:     registerASPNetUserReq.UserName,
			PhoneNumber:  registerASPNetUserReq.PhoneNumber,
			Email:        registerASPNetUserReq.Email,
			PasswordHash: registerASPNetUserReq.Password,
		}

		// Initialize redisCacahe + userRepo + userService from repositories and services
		redisCache := s.redisCache
		userRepo := repositories.NewAspNetUserRepository(db_authen)
		userService := OtherService.NewAspNetUserService(redisCache, userRepo)

		// * Create new User to the AspNetUser table
		registerReturn, err := userService.RegisterASPNetUser(newUser)
		if err != nil {
			registerASPNetUserResponse.CommonResponse.Result = "Fail to register AspNetUser"
			registerASPNetUserResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusInternalServerError, registerASPNetUserResponse)
			return
		}

		// * Binding User to the most basic role, further role binding will be done by admin on CMS
		userRepoRole := repositories.NewAspNetRolesRepository(db_authen)
		err = userRepoRole.Create(registerASPNetUserReq.UserRole, &newUser.Id)
		if err != nil {
			registerASPNetUserResponse.CommonResponse.Result = "Fail to bind role to user"
			registerASPNetUserResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusInternalServerError, registerASPNetUserResponse)
			return
		}

		// TODO : CREATE Userprofile based on userID
		userProfileRepo := repositories.NewUserProfilesRepository(db_data)
		err = userProfileRepo.Create(registerASPNetUserReq, newUser)
		if err != nil {
			registerASPNetUserResponse.CommonResponse.Result = "Fail to create user profile"
			registerASPNetUserResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusInternalServerError, registerASPNetUserResponse)
			return
		}

		registerASPNetUserResponse.CommonResponse.Result = "success"
		registerASPNetUserResponse.CommonResponse.ErrorMessage = ""
		registerASPNetUserResponse.RegisterASPNetUserReturn.AccessToken = registerReturn.AccessToken
		registerASPNetUserResponse.RegisterASPNetUserReturn.RefreshToken = registerReturn.RefreshToken
		registerASPNetUserResponse.RegisterASPNetUserReturn.SessionId = registerReturn.SessionId
		c.JSON(http.StatusOK, registerASPNetUserResponse)
	}
}

/*
* Set the session link to the Redis Cache
 */
func (s *Server) SetLinkSession(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var linkSessionRequest *dto.LinkSessionReq
		if err := c.ShouldBind(&linkSessionRequest); err != nil {
			log.Println("SetLinkSession handler:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		sessionId := common.GenerateGUID()
		linkSessionResponse := &dto.LinkSessionResponse{
			CommonResponse: dto.CommonResponse{},
		}

		if linkSessionRequest.Link == "" || linkSessionRequest.UserEmail == "" {
			linkSessionResponse.CommonResponse.Result = BAD_REQUEST
			linkSessionResponse.CommonResponse.ErrorMessage = "missing link"
			c.JSON(http.StatusBadRequest, linkSessionResponse)
			return
		}

		linkService := OtherService.NewLinkSessionService(s.redisCache)
		log.Println("link catched", linkSessionRequest.Link)
		// linkSaved := linkSessionRequest.Link + sessionId

		err := linkService.SetLinkSession(sessionId, linkSessionRequest.UserEmail)
		if err != nil {
			linkSessionResponse.CommonResponse.Result = "Fail"
			linkSessionResponse.CommonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusNoContent, linkSessionResponse.CommonResponse)
			return
		}

		linkSessionResponse.SessionId = sessionId
		c.JSON(http.StatusOK, linkSessionResponse)
	}
}

func (s *Server) GetLinkSession(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var linkSessionRequest *dto.LinkSessionReq
		if err := c.ShouldBind(&linkSessionRequest); err != nil {
			log.Println("SetLinkSession handler:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		linkSessionResponse := &dto.LinkSessionResponse{
			CommonResponse: dto.CommonResponse{},
			SessionId:      "",
		}
		if linkSessionRequest.Link == "" {
			linkSessionResponse.CommonResponse.Result = "Bad Request"
			linkSessionResponse.CommonResponse.ErrorMessage = "missing link"
			c.JSON(http.StatusBadRequest, linkSessionResponse)
			return
		}

		sessionId := c.GetHeader("Session-Id")
		if sessionId == "" {
			linkSessionResponse.CommonResponse.Result = "Bad Request"
			linkSessionResponse.CommonResponse.ErrorMessage = "missing session"
			c.JSON(http.StatusBadRequest, linkSessionResponse)
			return
		}

		linkService := OtherService.NewLinkSessionService(s.redisCache)
		res, err := linkService.GetLinkSession(linkSessionRequest.Link)
		if err != nil {
			linkSessionResponse.CommonResponse.Result = "No result"
			linkSessionResponse.CommonResponse.ErrorMessage = "not found"
			c.JSON(http.StatusOK, linkSessionResponse)
			return
		} else if res == "" {
			linkSessionResponse.CommonResponse.Result = "fail to get link"
			linkSessionResponse.CommonResponse.ErrorMessage = "wrong link"
			c.JSON(http.StatusBadRequest, linkSessionResponse)
			return
		}

		linkSessionResponse.CommonResponse.Result = "ok"
		linkSessionResponse.CommonResponse.ErrorMessage = "none"
		linkSessionResponse.SessionId = sessionId
		c.JSON(http.StatusOK, linkSessionResponse)

	}
}

// * This API will check the validicity of the token (AccessToken + RefreshToken)
func (s *Server) VerifyToken(db *gorm.DB) func(context *gin.Context) {
	return func(context *gin.Context) {
		// Initialize the response object
		var rTokenRes models.TokenVerificationResponse // Add semicolon here
		db := s.db_authen

		// Get the authorization token from the header
		authorizationHeader := context.GetHeader(AUTHORIZATION_HEADER)
		valid, err := OtherService.CheckAuthorizationHeaderStruct(authorizationHeader)

		if err != nil && !valid {
			err := errors.New("unsupported authorization type")
			rTokenRes.Result = STATUS_UNAUTHORIZED
			rTokenRes.ErrorMessage = err.Error()
			rTokenRes.Status = http.StatusUnauthorized
			context.JSON(http.StatusUnauthorized, rTokenRes)
			return
		}

		fields := strings.Fields(authorizationHeader)

		// * Call services to check the validicity of the token
		var roles *[]string
		rList, roles, err := OtherService.ParseAndCheckVadicityToken(fields, db)

		if err != nil {
			rTokenRes.Result = STATUS_UNAUTHORIZED
			rTokenRes.ErrorMessage = err.Error()
			rTokenRes.Status = http.StatusUnauthorized
			context.JSON(http.StatusUnauthorized, rTokenRes)
		}

		rTokenRes.Result = STATUS_AUTHORIZED
		rTokenRes.ErrorMessage = ""
		rTokenRes.Status = http.StatusOK
		rTokenRes.Role = *roles
		rTokenRes.Data = rList
		context.JSON(http.StatusOK, rTokenRes)
		return
	}
}

func (s *Server) RefreshToken(db *gorm.DB) func(context *gin.Context) {
	return func(context *gin.Context) {
		// Initialize the response object
		rTokenRes := &dto.LoginResponse{
			CommonResponse:        dto.CommonResponse{},
			LoginAspNetUserReturn: dto.LoginAspNetUserReturn{},
		}

		// * Check the authorization of the token
		authorizationHeader := context.GetHeader(AUTHORIZATION_HEADER)
		valid, err := OtherService.CheckAuthorizationHeaderStruct(authorizationHeader)

		if err != nil && !valid {
			err := errors.New("unsupported authorization type")
			rTokenRes.CommonResponse.Result = STATUS_UNAUTHORIZED
			rTokenRes.CommonResponse.ErrorMessage = err.Error()
			context.JSON(http.StatusUnauthorized, rTokenRes)
			return
		}

		fields := strings.Fields(authorizationHeader)

		// * Call services to check the validicity of the token
		var data *dto.TokenData // Declare the variable "data"
		data, _, err = OtherService.ParseAndCheckVadicityToken(fields, db)
		if err != nil {
			rTokenRes.CommonResponse.Result = STATUS_UNAUTHORIZED
			rTokenRes.CommonResponse.ErrorMessage = err.Error()
			context.JSON(http.StatusUnauthorized, rTokenRes)
		}

		// * Check if this is a refresh token
		isRefreshToken, _ := strconv.ParseBool(data.IsRefreshToken)
		if !isRefreshToken {
			rTokenRes.CommonResponse.Result = "Unauthorized - Token is invalid"
			rTokenRes.CommonResponse.ErrorMessage = "Invalid token type"
			context.JSON(http.StatusUnauthorized, rTokenRes)
			return
		}

		// * Create AspNetUserSession that saved in Redis + Database
		aspNetUserSesison := models.AspNetUserSession{
			UserProfileId: data.UserId,
			UserName:      data.UserName,
			ExpiredAt:     time.Now().Add(common.ACCESSTOKEN_DURATION),
		}

		userRepo := repositories.NewAspNetUserRepository(db)
		userSessionId, err := userRepo.SaveAspNetUserSession(&aspNetUserSesison)
		if err != nil {
			rTokenRes.CommonResponse.Result = "Fail to save user session"
			rTokenRes.CommonResponse.ErrorMessage = err.Error()
			context.JSON(http.StatusInternalServerError, rTokenRes)
			return
		}

		aspNetUserSesison.SessionId = userSessionId
		_, errSaveSession := s.redisCache.SetAspNetUserSession(aspNetUserSesison)
		// *
		if errSaveSession != nil {
			rTokenRes.CommonResponse.Result = "Fail to save user session on Redis"
			rTokenRes.CommonResponse.ErrorMessage = errSaveSession.Error()
			context.JSON(http.StatusInternalServerError, rTokenRes)
			return
		}

		// * Token is authorized
		// * Create new user Claims to refresh/regenerate new accessToken and refreshToken
		userClaims := tokenizer.ASPNetUserClaim{
			UserId:    data.UserId,
			SessionId: data.SessionId,
			Email:     data.Email,
			UserName:  data.UserName,
		}

		accessToken, refreshToken, err := OtherService.GenerateAspNetUserToken(userClaims)
		if err != nil {
			rTokenRes.CommonResponse.Result = "Fail to generate token"
			rTokenRes.CommonResponse.ErrorMessage = err.Error()
			context.JSON(http.StatusInternalServerError, rTokenRes)
			return
		}

		// * If token is generated successfully, return the new token
		rTokenRes.LoginAspNetUserReturn.AccessToken = &accessToken
		rTokenRes.LoginAspNetUserReturn.RefreshToken = &refreshToken
		rTokenRes.LoginAspNetUserReturn.SessionId = &userSessionId
		rTokenRes.CommonResponse.Result = SUCCESS
		context.JSON(http.StatusOK, rTokenRes)
		return
	}
}

/*
* Check for the validicty of the link session, if not valid then return 401
 */
func (s *Server) CheckChangePasswordLinkSession(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var linkSessionRequest dto.LinkSessionReq
		if err := c.ShouldBindJSON(&linkSessionRequest); err != nil {
			log.Println("checkLinkSession handler:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		resetPasswordLink := &dto.ResetPasswordLinkResponse{
			CommonResponse: dto.CommonResponse{},
			Link:           "",
		}

		sessionId := c.GetHeader("sessionId")
		if sessionId == "" {
			resetPasswordLink.CommonResponse.Result = BAD_REQUEST
			resetPasswordLink.CommonResponse.ErrorMessage = MISSING_SESSION
			c.JSON(http.StatusBadRequest, resetPasswordLink)
			return
		}

		linkService := OtherService.NewLinkSessionService(s.redisCache)
		// ? This take linkSession to check for sessionID in Redis ???
		res, err := linkService.GetLinkSession(linkSessionRequest.Link)
		if err != nil || res == "" {
			resetPasswordLink.CommonResponse.Result = SEND_EMAIL_AGAIN
			resetPasswordLink.CommonResponse.ErrorMessage = SESSION_EXPIRED
			resetPasswordLink.Link = "https://www.weallnet.com/vi/quen-mat-khau"
			c.JSON(http.StatusOK, resetPasswordLink)
			return
		}
		resetPasswordLink.CommonResponse.Result = SESSION_VERIFIED
		resetPasswordLink.CommonResponse.ErrorMessage = NONE
		resetPasswordLink.Link = "https://www.weallnet.com/vi/dat-lai-mat-khau"
		c.JSON(http.StatusOK, resetPasswordLink)
	}
}

func (s *Server) ChangePasswordLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var changePasswordRequest dto.ChangePasswordRequest
		if err := c.ShouldBind(&changePasswordRequest); err != nil {
			log.Println("checkLinkSession handler:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		commonResponse := &dto.CommonResponse{
			Result:       "",
			ErrorMessage: "",
		}

		if changePasswordRequest.UserEmail == "" || changePasswordRequest.NewPassword == "" {
			commonResponse.Result = "no content"
			commonResponse.ErrorMessage = BAD_REQUEST
			c.JSON(http.StatusBadRequest, commonResponse)
			return
		}

		sessionId := c.Param("sessionId")
		if sessionId == "" {
			commonResponse.Result = "no content"
			commonResponse.ErrorMessage = BAD_REQUEST
			c.JSON(http.StatusBadRequest, commonResponse)
			return
		}

		linkService := OtherService.NewLinkSessionService(s.redisCache)
		getSession, err := linkService.GetLinkSession(sessionId)
		log.Println("email: ", getSession)
		log.Println("User Email: ", changePasswordRequest.UserEmail)
		if getSession != changePasswordRequest.UserEmail || err != nil {
			commonResponse.Result = "Unauthorized"
			commonResponse.ErrorMessage = "access denied"
			c.JSON(http.StatusUnauthorized, commonResponse)
			return
		}

		userRepo := repositories.NewAspNetUserRepository(db)
		userService := OtherService.NewAspNetUserService(s.redisCache, userRepo)
		err = userService.ChangeASPNetUserPassword(changePasswordRequest.UserEmail, changePasswordRequest.NewPassword)
		if err != nil {
			commonResponse.Result = "Internal Server Error"
			commonResponse.ErrorMessage = err.Error()
			c.JSON(http.StatusInternalServerError, commonResponse)
			return
		}
		commonResponse.Result = "ok"
		commonResponse.ErrorMessage = "success"
		c.JSON(http.StatusOK, commonResponse)
	}
}
