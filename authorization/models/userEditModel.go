package models

type UserEdit struct {
	First_name  string `json:"first_name" validate:"required,min=2,max=10"`
	Last_name   string `json:"last_name" validate:"required,min=2,max=20"`
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword" validate:"required,min=12"`
	NewPassword string `json:"newPassword" validate:"required,min=12"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address" validate:"required"`
}
