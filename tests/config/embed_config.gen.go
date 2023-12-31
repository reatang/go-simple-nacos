// Code generated by "gonacos_config --embed=GlobalConfig --config=SomeConfig --dataid=config --codec=yaml"; DO NOT EDIT.
package config

import "fmt"
import gonacos "github.com/reatang/go-simple-nacos"
import "sync"

var __GlobalConfigSomeConfigMutex = sync.RWMutex{}

func RegisterEmbedGlobalConfigSomeConfig(ncc *gonacos.NacosConfigClient, c *SomeConfig) {
	ncc.WatchF("", "config", func(namespace, group, dataId, data string) {
		__GlobalConfigSomeConfigMutex.Lock()
		defer __GlobalConfigSomeConfigMutex.Unlock()

		err := gonacos.DecodeYaml(data, &c)
		if err != nil {
			fmt.Println("[gonacos]dataid:config,error:", err.Error())
		}
	})
}

func (c *GlobalConfig) GetSomeConfig() SomeConfig {
	__GlobalConfigSomeConfigMutex.RLock()
	defer __GlobalConfigSomeConfigMutex.RUnlock()

	return c.SomeConfig
}
