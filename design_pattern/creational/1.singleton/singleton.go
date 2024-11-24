package main

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

const dsn = "root:123456@tcp(127.0.0.1:3306)/db1?parseTime=True&loc=Local"

func GetDBInstance() *gorm.DB {
	var err error
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(fmt.Errorf("connect MySQL failed, err is %w", err))
		}
	})

	return db
}

func main() {
	s := GetDBInstance()
	sqlDB, err := s.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	fmt.Println("连接数据库成功")
}
