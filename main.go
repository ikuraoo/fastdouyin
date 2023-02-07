package main

import (
	"github.com/ikuraoo/fastdouyin/dao"
	"github.com/ikuraoo/fastdouyin/service"

	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()

	dao.InitMysql()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
