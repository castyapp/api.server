package components

import (
	"github.com/MrJoshLab/go-respond"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func ParseGrpcErrorResponse(err error) (code int, response interface{}, ok bool) {

	switch statusCode := status.Code(err); statusCode {
	case codes.OK:
		return http.StatusOK, nil, true
	case codes.NotFound:
		code, response = respond.Default.NotFound()
		return
	case codes.PermissionDenied:
		code, response = respond.Default.SetStatusCode(http.StatusForbidden).
			SetStatusText("Failed!").
			RespondWithMessage("Permission Denied!")
		return
	case codes.Unauthenticated:
		code, response = respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!")
		return
	case codes.InvalidArgument:

		if s, ok := status.FromError(err); ok {
			validationErrors := make(map[string] interface{}, 0)
			for _, validationErr := range s.Proto().Details {
				validationErrors[validationErr.TypeUrl] = validationErr.Value
			}
			code, response = respond.Default.SetStatusCode(420).
				SetStatusText("failed").
				RespondWithResult(validationErrors)
			return
		}

		code, response = respond.Default.SetStatusCode(420).
			SetStatusText("failed").
			RespondWithMessage("Invalid arguments!")
		return

	default:
		code, response = respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("failed").
			RespondWithMessage("Internal server error. Please try again later!")
		return
	}
}