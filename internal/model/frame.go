package model

import (
	"github.com/zhufuyi/sponge/pkg/mysql"
)

type Frame struct {
	mysql.Model `gorm:"embedded"` // embed id and time

	Frame      string `gorm:"column:frame;type:json;NOT NULL" json:"frame"`
	DevAddr    string `gorm:"column:dev_addr;type:varchar(10)" json:"dev_addr"`
	DataType   int    `gorm:"column:data_type;type:int(11);NOT NULL" json:"data_type"`
	GatewayMac string `gorm:"column:gateway_mac;type:varchar(20)" json:"gateway_mac"`
}

// TableName table name
func (m *Frame) TableName() string {
	return "frame"
}
