// Code generated by "gonacos_config --embed=GlobalConfig --config=other:other_config.OtherConfig --dataid=other_config --codec=yaml"; DO NOT EDIT.
package config

import "fmt"
import gonacos "github.com/reatang/go-simple-nacos"
import "github.com/reatang/go-simple-nacos/tests/config/other_config"
import "sync"

var __GlobalConfigotherMutex = sync.RWMutex{}

func RegisterEmbedGlobalConfigother(ncc *gonacos.NacosConfigClient, c *other_config.OtherConfig) {
	ncc.WatchF("", "other_config", func(namespace, group, dataId, data string) {
		__GlobalConfigotherMutex.Lock()
		defer __GlobalConfigotherMutex.Unlock()

		err := gonacos.DecodeYaml(data, &c)
		if err != nil {
			fmt.Println("[gonacos]dataid:other_config,error:", err.Error())
		}
	})
}

func (c *GlobalConfig) Getother() other_config.OtherConfig {
	__GlobalConfigotherMutex.RLock()
	defer __GlobalConfigotherMutex.RUnlock()

	return *c.other
}
