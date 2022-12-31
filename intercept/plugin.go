package intercept

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Plugin struct {
	Method   []string `yaml:"method"`
	Position []string `yaml:"position"`
	Regexp   string   `yaml:"regexp"`
	Count    int      `yaml:"count"`
}

// UsePlugin 是用于漏洞检测的中间件，同时负责攻击log的输出
func UsePlugin(pluginList []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var plugin Plugin
		for _, pluginName := range pluginList {
			file := "./plugin/" + pluginName
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("没有找到" + pluginName + "，请检查插件路径")
				return
			}
			err = yaml.Unmarshal(data, &plugin)
			if err != nil {
				fmt.Println(pluginName + "加载失败！" + err.Error())
				return
			}
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
						checkHead(context, reg, plugin.Count, pluginName)
					} else if position == "body" {
						checkBody(context, reg, plugin.Count, pluginName)
					} else if position == "url" {
						checkUrl(context, reg, plugin.Count, pluginName)
					}
				}
			}
		}
	}
}

// checkHead 请求头的检查
func checkHead(context *gin.Context, reg *regexp.Regexp, count int, pluginName string) {
	buf := context.Request.Header
	for _, value := range buf {
		for _, data := range value {
			matches := reg.FindAllString(strings.ToLower(data), -1)
			if len(matches) >= count {
				showAttackLog(context, pluginName)
				context.HTML(403, "403.html", gin.H{
					"time":   time.Now(),
					"Client": context.Request.Header.Get("User-Agent"),
					"Ip":     context.ClientIP(),
				})
				context.Abort()
				return
			}
		}
	}
}

// checkBody 请求体的内容的检查
func checkBody(context *gin.Context, reg *regexp.Regexp, count int, pluginName string) {
	buf, _ := io.ReadAll(context.Request.Body)
	matches := reg.FindAllString(strings.ToLower(string(buf)), -1)
	if len(matches) >= count {
		showAttackLog(context, pluginName)
		context.HTML(403, "403.html", gin.H{
			"time":   time.Now(),
			"Client": context.Request.Header.Get("User-Agent"),
			"Ip":     context.ClientIP(),
		})
		context.Abort()
		return
	}
	context.Request.Body = io.NopCloser(bytes.NewReader(buf))
}

// checkUrl 请求url的检查
func checkUrl(context *gin.Context, reg *regexp.Regexp, count int, pluginName string) {
	buf := context.Request.RequestURI
	buf, _ = url.QueryUnescape(buf)
	matches := reg.FindAllString(strings.ToLower(buf), -1)
	if len(matches) >= count {
		showAttackLog(context, pluginName)
		context.HTML(403, "403.html", gin.H{
			"time":   time.Now(),
			"Client": context.Request.Header.Get("User-Agent"),
			"Ip":     context.ClientIP(),
		})
		context.Abort()
		return
	}
}

// showAttackLog 打印攻击日志
func showAttackLog(c *gin.Context, pluginName string) {
	logFile, err := os.OpenFile("./log/deny.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)
	start := time.Now()
	timestamp := time.Now()
	latency := timestamp.Sub(start)

	if latency > time.Minute {
		latency = latency.Truncate(time.Second)
	}
	var denyLog = fmt.Sprintf("[WPServer] %v |%s %s %s| %3v | %s %s => %s %s |%s %s%s %s %#v",
		time.Now().Format("2006/01/02 - 15:04:05"),
		red, c.Request.Method, reset,
		latency,
		proxy, c.ClientIP(),
		c.Request.Host, reset,
		red, "检测到攻击行为！匹配模板：", pluginName, reset,
		c.Request.URL.Path,
	)
	fmt.Println(denyLog)
	var writeLog = fmt.Sprintf("|%s| %3v | %s => %s |%s%s %#v",
		c.Request.Method,
		latency,
		c.ClientIP(),
		c.Request.Host,
		"检测到攻击行为！匹配模板：", pluginName,
		c.Request.URL.Path,
	)
	log.Println(writeLog)
}
