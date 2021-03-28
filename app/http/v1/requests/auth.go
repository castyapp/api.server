package requests

type CreateUserRequest struct {
	Fullname             string `validate:"required" form:"username"`
	Password             string `validate:"required" form:"password"`
	Username             string `validate:"required" form:"username"`
	Email                string `validate:"required,email" form:"email"`
	PasswordConfirmation string `validate:"required" form:"password_confirmation"`
}

type CreateAuthTokenRequest struct {
	User string `validate:"required" form:"user"`
	Pass string `validate:"required" form:"pass"`
}
