package requests

type SearchUserRequest struct {
	Keyword string `validate:"required" form:"keyword"`
}

type AcceptFriendRequest struct {
	RequestID string `validate:"required" form:"request_id"`
}

type UpdatePasswordRequest struct {
	Password                string `validate:"required" form:"password"`
	NewPassword             string `validate:"required" form:"new_password"`
	NewPasswordConfirmation string `validate:"required" form:"new_password_confirmation"`
}
