package model

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Database(connString string) {
	db, err := gorm.Open("mysql", "root:123456@(localhost:3306)/micro_todolist?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	if gin.Mode() == "release" {
		db.LogMode(false)
	}

	// 默认不加复数
	db.SingularTable(true)

	// 设置连接池

	// 空闲
	db.DB().SetMaxIdleConns(20)
	// 打开
	db.DB().SetMaxOpenConns(100)
	// 超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	DB = db

	migration()

}
