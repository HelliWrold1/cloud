// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// downlink http service level error code
// each resource name corresponds to a unique number (http type), the number range is 1~100, if there is the same number, trigger panic
var (
	downlinkNO       = 59
	downlinkName     = "downlink"
	downlinkBaseCode = errcode.HCode(downlinkNO)

	ErrCreateDownlink = errcode.NewError(downlinkBaseCode+1, "failed to create "+downlinkName)
	ErrDeleteDownlink = errcode.NewError(downlinkBaseCode+2, "failed to delete "+downlinkName)
	ErrUpdateDownlink = errcode.NewError(downlinkBaseCode+3, "failed to update "+downlinkName)
	ErrGetDownlink    = errcode.NewError(downlinkBaseCode+4, "failed to get "+downlinkName+" details")
	ErrListDownlink   = errcode.NewError(downlinkBaseCode+5, "failed to get list of "+downlinkName)
	// for each error code added, add +1 to the previous error code
)
