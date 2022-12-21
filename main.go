package main

import (
	"gateway/intercept"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), intercept.Reserve, intercept.Logger, intercept.RemoveServer)
	err := r.Run(":80")
	if err != nil {
		panic("服务启动失败！请检查端口占用")
	}
}
