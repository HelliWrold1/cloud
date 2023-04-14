package model

import (
	"fmt"
)

func migration() {
	// 设置存储引擎为InnoDB
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		// 自动迁移，建表
		AutoMigrate(
			&User{},
			&Frame{},
			&Downlink{},
		)
	if err != nil {
		fmt.Println("migrate err", err)
	}
}
