package intercept

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"regexp"
)

type Plugin struct {
	Method   []string `yaml:"method"`
	Position []string `yaml:"position"`
	Regexp   string   `yaml:"regexp"`
	Count    int      `yaml:"count"`
}

func UsePlugin(pluginList []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var plugin Plugin
		for _, plugins := range pluginList {
			file := "./lib/" + plugins
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println(string(plugins) + "插件加载失败！")
				return
			}
			err = yaml.Unmarshal(data, &plugin)
			if err != nil {
				fmt.Println(string(plugins) + "插件加载失败！")
				return
			} else {
				reg, err := regexp.Compile(plugin.Regexp)
				if err != nil {
					fmt.Println("插件正则有误！")
					return
				}
				methods := make(map[string]bool, len(plugin.Method))
				for _, s := range plugin.Method {
					methods[s] = true
				}
				_, ok := methods[context.Request.Method]
				// 如果请求方式匹配到了
				if ok {
					for _, position := range plugin.Position {
						if position == "head" {
							checkHead(context, reg, plugin.Count)
						} else if position == "body" {
							checkBody(context, reg, plugin.Count)
						} else if position == "url" {
							checkUrl(context, reg, plugin.Count)
						}
					}
				}
			}
		}
	}
}

func checkHead(context *gin.Context, reg *regexp.Regexp, count int) {

}
func checkBody(context *gin.Context, reg *regexp.Regexp, count int) {
	buf, _ := io.ReadAll(context.Request.Body)
	matches := reg.FindAllString(string(buf), -1)
	if len(matches) >= count {
		context.AbortWithStatusJSON(403, gin.H{
			"警告！": "检测异常攻击行为！已记录IP！",
		})
		log.Println("检测到异常攻击行为！")
	}
	context.Request.Body = io.NopCloser(bytes.NewReader(buf))
}
func checkUrl(context *gin.Context, reg *regexp.Regexp, count int) {

}
