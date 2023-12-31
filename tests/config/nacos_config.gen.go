// Code generated by "gonacos_config --config=SomeConfig --dataid=config --codec=yaml"; DO NOT EDIT.
package config

import "fmt"
import gonacos "github.com/reatang/go-simple-nacos"
import "sync/atomic"

var __SomeConfig = atomic.Value{}

func RegisterSomeConfig(ncc *gonacos.NacosConfigClient, defConf *SomeConfig) {
	if defConf != nil {
		__SomeConfig.Store(*defConf)
	}

	ncc.WatchF("", "config", func(namespace, group, dataId, data string) {
		var c SomeConfig
		err := gonacos.DecodeYaml(data, &c)
		if err != nil {
			fmt.Println("[gonacos]dataid:config,error:", err.Error())
			return
		}
		__SomeConfig.Store(c)
	})
}

func GetSomeConfig() SomeConfig {
	return __SomeConfig.Load().(SomeConfig)
}
