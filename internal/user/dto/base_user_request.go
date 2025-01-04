package dto

import "github.com/google/uuid"

type UserRequest struct {
	ID       uuid.UUID `json:"id"`
	UserType string    `json:"userType" binding:"required,validUserType"`
	FullName *string   `json:"fullName"`
	Email    string    `json:"email" binding:"required,email"`
	Password *string   `json:"password"`
	FileId   uuid.UUID `json:"fileId"`
}
