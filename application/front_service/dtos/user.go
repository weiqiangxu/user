package dtos

type UserDto struct {
	Name string `json:"name,omitempty"`
}

type GetUserListRequest struct {
	ID int32 `json:"id,omitempty" form:"id"`
}
