package configure

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() error {
	err := InitDB()
	InitRedis()
	return err
}

func InitDB() error {

	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	//dsn := "root:root@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=true"
	var err error
	Db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	return err

}