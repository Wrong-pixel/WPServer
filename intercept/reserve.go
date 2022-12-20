/*
* 反向代理的逻辑，没啥技术含量，顺便实现了下负载均衡
* 本来是想用协程和连接池增强并发性能的，但是不会写
* 目前在1000请求/秒的并发下丢包率在0.3%
*
* The logic of the reverse proxy, there is no technical content, and the load balancing is achieved by the way
* I wanted to use coroutines and connection pooling to enhance concurrency performance, but I couldn't write it
* At present, the packet loss rate is 0.3% at 1000 requests/sec
 */

package intercept

import (
	"fmt"
	"gateway/conf"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http/httputil"
	"net/url"
)

func Reserve(c *gin.Context) {
	host := c.Request.Host
	configData := conf.ReadConfig()
	config, ok := configData[host]
	// 如果请求域名不在配置文件中，有可能是通过IP或者localhost进行访问的
	if !ok {
		c.String(200, "Welcome to the Wrong-Pixel server！")
		return
	}

	var target *url.URL
	if !config.Upstream || len(config.Urls) == 1 {
		target, _ = url.Parse(config.Urls[0].Url)
	} else {
		// 计算权重值总和
		totalWeight := 0
		for _, urlConfig := range config.Urls {
			totalWeight += urlConfig.Weight
		}

		// 生成一个1到权重值的随机数
		i := rand.Intn(totalWeight) + 1

		// 以随机数作为下标选择反向代理的服务器
		for _, urlConfig := range config.Urls {
			if i <= urlConfig.Weight {
				target, _ = url.Parse(urlConfig.Url)
				break
			}
			i -= urlConfig.Weight
		}
	}

	// 创建一个反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 更新请求头使得可以使用SSL重定向
	c.Request.URL.Host = target.Host
	c.Request.URL.Scheme = target.Scheme
	c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))
	c.Request.Host = target.Host

	// 反向代理请求
	proxy.ServeHTTP(c.Writer, c.Request)
	logMessage := fmt.Sprintf("%s requested %s and was forwarded to %s", c.ClientIP(), host, target.String())
	log.Printf(logMessage)
}
