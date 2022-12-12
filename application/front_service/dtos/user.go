package dtos

type UserDto struct {
	Name string `json:"name,omitempty"`
}

type GetUserListRequest struct {
	ID int32 `json:"id,omitempty" form:"id"`
}

// UserInfoRequest contains user information
// 默认验证规则 github.com/go-playground/validator/v10@v10.10.0/baked_in.go:bakedInValidators
type UserInfoRequest struct {
	FirstName string            `validate:"required,lt=2"` // 长度小于2
	LastName  string            `validate:"required"`      // 必填字段
	Age       uint8             `validate:"gte=0,lte=130"`
	Email     string            `validate:"required,email"`
	Addresses []*AddressRequest `validate:"required,dive,required"` // a person can have a home and cottage...
}

// AddressRequest houses a users address information
type AddressRequest struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

type UserRequest struct {
	FirstName string `validate:"required,lt=2"`
}
