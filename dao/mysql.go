package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitMysql() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//Db.AutoMigrate(&controller.User{})
	return
}
