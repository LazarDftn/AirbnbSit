package models

type LoginModel struct {
	Recaptcha *string `json:"recaptcha"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}
