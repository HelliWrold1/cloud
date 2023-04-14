package model

import (
	"github.com/zhufuyi/sponge/pkg/mysql"
)

type User struct {
	mysql.Model `gorm:"embedded"` // embed id and time

	Username string `gorm:"column:username;type:varchar(255)" json:"username"`
	Password string `gorm:"column:password;type:varchar(255)" json:"password"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
