package main

import (
	"gateway/intercept"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(intercept.Reserve, intercept.RemoveServer)
	err := r.Run(":80")
	if err != nil {
		panic("服务启动失败！请检查端口占用")
	}
}
