package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ time.Time

// CreateUserRequest create params
// todo fill in the binding rules https://github.com/go-playground/validator
type CreateUserRequest struct {
	Username string `json:"username" binding:""`
	Password string `json:"password" binding:""`
}

// UpdateUserByIDRequest update params
type UpdateUserByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username string `json:"username" binding:""`
	Password string `json:"password" binding:""`
}

// GetUserByIDRespond respond detail
type GetUserByIDRespond struct {
	ID string `json:"id"` // covert to string id

	Username string `json:"username"`
	//Password  string    `json:"password"` // 隐藏密码
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteUsersByIDsRequest request form ids
type DeleteUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUsersByIDsRequest request form ids
type GetUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUsersRequest request form params
type GetUsersRequest struct {
	query.Params // query parameters
}

// ListUsersRespond list data
type ListUsersRespond []struct {
	GetUserByIDRespond
}

// UpdateUserPasswordRequest update user's password
type UpdateUserPasswordRequest struct {
	Username    string `json:"username" binding:""`
	OldPassword string `json:"old_password" binding:""`
	NewPassword string `json:"new_password" binding:""`
}

// UpdateUserPasswordResponse update user's password
type UpdateUserPasswordResponse struct {
	Username string `json:"username" binding:""`
	Password string `json:"password" binding:""`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:""`
	Password string `json:"password" binding:""`
}
