package requests

import (
	"github.com/castyapp/api.server/app/models"
	"github.com/CastyLab/grpc.proto/proto"
)

type UpdateTheaterRequest struct {
	Description       string                    `validate:"required_without_all=Privacy VideoPlayerAccess"`
	Privacy           proto.PRIVACY             `validate:"required_without_all=Description VideoPlayerAccess"`
	VideoPlayerAccess proto.VIDEO_PLAYER_ACCESS `validate:"required_without_all=Description Privacy"`
}

type InviteToTheaterRequest struct {
	FriendIDs []string `json:"friend_ids" validate:"required,gt=0,dive,required"`
}

type AddSubtitlesRequest struct {
	Subtitles []models.Subtitle `json:"subtitles" validate:"required,gt=0,dive,required"`
}
