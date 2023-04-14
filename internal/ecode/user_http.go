// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// user http service level error code
// each resource name corresponds to a unique number (http type), the number range is 1~100, if there is the same number, trigger panic
var (
	userNO       = 51
	userName     = "user"
	userBaseCode = errcode.HCode(userNO)

	ErrCreateUser = errcode.NewError(userBaseCode+1, "failed to create "+userName)
	ErrDeleteUser = errcode.NewError(userBaseCode+2, "failed to delete "+userName)
	ErrUpdateUser = errcode.NewError(userBaseCode+3, "failed to update "+userName)
	ErrGetUser    = errcode.NewError(userBaseCode+4, "failed to get "+userName+" details")
	ErrListUser   = errcode.NewError(userBaseCode+5, "failed to get list of "+userName)
	// for each error code added, add +1 to the previous error code
)
