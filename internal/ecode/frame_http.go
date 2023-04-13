// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// frame http service level error code
// each resource name corresponds to a unique number (http type), the number range is 1~100, if there is the same number, trigger panic
var (
	frameNO       = 57
	frameName     = "frame"
	frameBaseCode = errcode.HCode(frameNO)

	ErrCreateFrame = errcode.NewError(frameBaseCode+1, "failed to create "+frameName)
	ErrDeleteFrame = errcode.NewError(frameBaseCode+2, "failed to delete "+frameName)
	ErrUpdateFrame = errcode.NewError(frameBaseCode+3, "failed to update "+frameName)
	ErrGetFrame    = errcode.NewError(frameBaseCode+4, "failed to get "+frameName+" details")
	ErrListFrame   = errcode.NewError(frameBaseCode+5, "failed to get list of "+frameName)
	// for each error code added, add +1 to the previous error code
)
