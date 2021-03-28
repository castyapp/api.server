package requests

type CreateMessageRequest struct {
	Content string `validate:"required" form:"content"`
}
