package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ time.Time

// CreateFrameRequest create params
// todo fill in the binding rules https://github.com/go-playground/validator
type CreateFrameRequest struct {
	Frame      string `json:"frame" binding:""`
	DevAddr    string `json:"dev_addr" binding:""`
	DataType   int    `json:"data_type" binding:""`
	GatewayMac string `json:"gateway_mac" binding:""`
}

// UpdateFrameByIDRequest update params
type UpdateFrameByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Frame      string `json:"frame" binding:""`
	DevAddr    string `json:"dev_addr" binding:""`
	DataType   int    `json:"data_type" binding:""`
	GatewayMac string `json:"gateway_mac" binding:""`
}

// GetFrameByIDRespond respond detail
type GetFrameByIDRespond struct {
	ID string `json:"id"` // covert to string id

	Frame      string    `json:"frame"`
	DevAddr    string    `json:"dev_addr"`
	DataType   int       `json:"data_type"`
	GatewayMac string    `json:"gateway_mac"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DeleteFramesByIDsRequest request form ids
type DeleteFramesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetFramesByIDsRequest request form ids
type GetFramesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetFramesRequest request form params
type GetFramesRequest struct {
	query.Params // query parameters
}

// ListFramesRespond list data
type ListFramesRespond []struct {
	GetFrameByIDRespond
}
