package main

import (
	"gateway/conf"
	"gateway/intercept"
	"github.com/gin-gonic/gin"
)

func main() {
	var config, pluginList = conf.ReadConfig()
	r := gin.New()
	r.Use(gin.Recovery(), intercept.UsePlugin(pluginList), intercept.Reserve(config), intercept.Logger, intercept.RemoveServer)
	err := r.Run(":80")
	if err != nil {
		panic("服务启动失败！请检查端口占用")
	}
}
