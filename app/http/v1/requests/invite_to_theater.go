package requests

type InviteToTheaterRequest struct {
	FriendIDs  []string  `json:"friend_ids" validate:"required,gt=0,dive,required"`
}