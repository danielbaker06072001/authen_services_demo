package repositories

import (
	"authen-service/models"

	"gorm.io/gorm"
)

type IAspNetRoles interface {
	Create(role *string, userId *string) error
	GetUserRoles(userId *string) ([]string, error)
}

type aspNetUserRolesRepositories struct {
	db *gorm.DB
}

func NewAspNetRolesRepository(db *gorm.DB) IAspNetRoles {
	return &aspNetUserRolesRepositories{db: db}
}

func (r *aspNetUserRolesRepositories) GetUserRoles(userId *string) ([]string, error) {
	var roles []string
	if err := r.db.Table(`AspNetUserRoles`).
		Select(`"AspNetRoles"."NormalizedName"`).
		Joins(`LEFT JOIN "AspNetRoles" ON "AspNetUserRoles"."RoleId" = "AspNetRoles"."Id"`).
		Joins(`LEFT JOIN "AspNetUsers" ON "AspNetUserRoles"."UserId" = "AspNetUsers"."Id"`).
		Where(`"AspNetUserRoles"."UserId" = ?`, *userId).
		Scan(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// * This function will check if the role already exists in "AspNetUserRoles" table
// If not already exist, create new role (not applied here, since role are predefined)
func (r *aspNetUserRolesRepositories) Create(role *string, userID *string) error {
	// Check if Role is already exist in AspNetRoles table
	// Retrive the role ID for later used in AspNetUserRoles table
	var RoleId string
	if err := r.db.Table(`AspNetRoles`).
		Select(`"Id"`).
		Where(`"NormalizedName" = ?`, *role).
		Scan(&RoleId).Error; err != nil {
		return err
	}

	// * Set up and create new AspNetUserRoles
	// * Add new user to the role
	// @param: UserId, RoleId string, get from AspNetRoles 1table and create new id
	NewAspNetUserRoles := models.AspNetUserRoles{
		UserId: *userID,
		RoleId: RoleId,
	}
	// Create new record in AspNetUserRoles table
	if err := r.db.Table(`AspNetUserRoles`).
		Create(&NewAspNetUserRoles).
		Error; err != nil {
		return err
	}

	// Add claims to users
	// * Claim act as assertions about the user, such as their role, permisions and other attributes
	// * Claims are added to the user's identity and can be used to make authorization decisions
	NewAspNetUserClaims := models.AspNetUserClaims{
		UserId:     *userID,
		ClaimType:  "role",
		ClaimValue: *role,
	}
	// Insert new record to AspNetUserClaims table
	if err := r.db.Table(`AspNetUserClaims`).
		Create(&NewAspNetUserClaims).
		Error; err != nil {
		return err
	}

	return nil
}
