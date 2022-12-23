/*
* 读取conf目录下的config.yaml配置文件，并返回一个map用于查询
* 指定配置文件的路径方法还没有想好具体怎么实现，因为计划是写一个web管理页面，在管理页面去指定配置文件路径
* 本来是想用hash表查找的，但是不会写，干脆就map查找吧，可能会比数组查找效率高一点点
*
* Read the config.yaml configuration file in the conf directory and return a map for querying
* Originally, I wanted to use the hash table to look up, but I can't write it, so just map it
* Originally, I wanted to use a hash table to look, but I can't write, just map lookup, which may be a little more efficient than array lookup
 */

package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

// YamlConfig 总体配置的结构体
type YamlConfig struct {
	Server []struct {
		Domain   string `yaml:"domain"`
		Upstream bool   `yaml:"upstream"`
		Urls     []struct {
			Url    string `yaml:"url"`
			Weight int    `yaml:"weight"`
		} `yaml:"urls"`
		Default bool `yaml:"default"`
	} `yaml:"server"`
	Intercept []string `yaml:"intercept"`
}

// ServerConfig 用来生成map的server结构体
type ServerConfig struct {
	Upstream bool `yaml:"upstream"`
	Urls     []struct {
		Url    string `yaml:"url"`
		Weight int    `yaml:"weight"`
	} `yaml:"urls"`
	Default bool `yaml:"default"`
}

func ReadConfig() (map[string]ServerConfig, []string) {
	// 读取配置
	data, err := os.ReadFile("./conf/config.yaml")
	if err != nil {
		panic(err)
	}

	// 解析yaml配置文件
	var configData YamlConfig
	// 定义map
	serverData := make(map[string]ServerConfig)
	err = yaml.Unmarshal(data, &configData)
	pluginList := configData.Intercept
	if err != nil {
		panic(err)
	}
	// map赋值
	for _, server := range configData.Server {
		serverData[server.Domain] = ServerConfig{
			Upstream: server.Upstream,
			Urls:     server.Urls,
			Default:  server.Default,
		}
	}
	return serverData, pluginList
}
