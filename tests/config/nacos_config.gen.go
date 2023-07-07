package config

import (
	"sync/atomic"

	gonacos "github.com/reatang/go-simple-nacos"
)

var __SomeConfig = atomic.Value{}

func RegisterSomeConfig(ncc *gonacos.NacosConfigClient) {
	ncc.WatchF("13", "config", func(namespace, group, dataId, data string) {
		var c SomeConfig
		err := gonacos.DecodeYaml(data, &c)
		if err == nil {
			__SomeConfig.Store(c)
		}
	})
}

func GetSomeConfig() SomeConfig {
	return __SomeConfig.Load().(SomeConfig)
}
