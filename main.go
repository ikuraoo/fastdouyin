package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/service"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
