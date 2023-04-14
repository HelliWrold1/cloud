package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ time.Time

// CreateDownlinkRequest create params
// todo fill in the binding rules https://github.com/go-playground/validator
type CreateDownlinkRequest struct {
	DownLink string `json:"down_link" binding:""`
	DevAddr  string `json:"dev_addr" binding:""`
}

// UpdateDownlinkByIDRequest update params
type UpdateDownlinkByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	DownLink string `json:"down_link" binding:""`
	DevAddr  string `json:"dev_addr" binding:""`
}

// GetDownlinkByIDRespond respond detail
type GetDownlinkByIDRespond struct {
	ID string `json:"id"` // covert to string id

	DownLink  string    `json:"down_link"`
	DevAddr   string    `json:"dev_addr"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteDownlinksByIDsRequest request form ids
type DeleteDownlinksByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetDownlinksByIDsRequest request form ids
type GetDownlinksByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetDownlinksRequest request form params
type GetDownlinksRequest struct {
	query.Params // query parameters
}

// ListDownlinksRespond list data
type ListDownlinksRespond []struct {
	GetDownlinkByIDRespond
}
