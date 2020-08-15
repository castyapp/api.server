package tests

import (
	"github.com/CastyLab/api.server/app/http/v1/requests"
	"github.com/go-playground/validator/v10"
	"testing"
)

var validate = validator.New()

func TestUpdateRequestTheater(t *testing.T)  {
	req := requests.UpdateTheaterRequest{}
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		t.Errorf("Validation error, %v", validationErrors)
	}
}