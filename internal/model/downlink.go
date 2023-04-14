package model

import (
	"github.com/zhufuyi/sponge/pkg/mysql"
)

type Downlink struct {
	mysql.Model `gorm:"embedded"` // embed id and time

	DownLink string `gorm:"column:down_link;type:json" json:"down_link"`
	DevAddr  string `gorm:"column:dev_addr;type:varchar(10)" json:"dev_addr"`
}

// TableName table name
func (m *Downlink) TableName() string {
	return "downlink"
}
