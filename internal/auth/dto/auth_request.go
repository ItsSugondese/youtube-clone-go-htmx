package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"userType" validate:"required"`
}
