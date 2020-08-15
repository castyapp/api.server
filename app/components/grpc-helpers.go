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
		return 200, nil, true
	case codes.NotFound:
		code, response = respond.Default.NotFound()
		return
	case codes.PermissionDenied:
		code, response = respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Permission Denied!")
		return
	case codes.Unauthenticated:
		code, response = respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!")
		return
	default:

		if s, sok := status.FromError(err); sok {
			code, response = respond.Default.SetStatusCode(http.StatusBadRequest).
				SetStatusText("failed").
				RespondWithMessage(s.Message())
			return
		}

		code, response = respond.Default.SetStatusCode(500).
			SetStatusText("failed").
			RespondWithMessage("Internal server error!")
		return
	}

}