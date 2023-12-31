# go-simple-nacos
更有golang风格的使用 nacos-sdk-go 客户端

https://github.com/nacos-group/nacos-sdk-go 太难用了，简化使用方案

## 安装nacos

    略....

## 安装辅助生成工具

```
> go install github.com/reatang/go-simple-nacos/cmd/gonacos_config
```

## 使用

### 配置内容

```yaml
TestConfig: "aaabbbccc"
```

### 独立结构体配置映射

将nacos中的配置直接映射到该结构体的变量中，并支持动态更新

配置，并编写生成命令
```go
// file: config/config.go
package config

//go:generate gonacos_config --config=SomeConfig --dataid=config --codec=yaml
type SomeConfig struct {
	TestConfig  string `yaml:"TestConfig"`
}
```

### 嵌入结构体配置映射

将nacos中的配置映射到主结构体的某个成员中，并支持动态更新

配置，并编写生成命令

```go
// file: config/config.go
package config

type SomeConfig struct {
    TestConfig  string `yaml:"TestConfig"`
}

// 匿名写法
//go:generate gonacos_config --embed=GlobalConfig --config=SomeConfig --dataid=config --codec=yaml
type GlobalConfig struct {
    SomeConfig
}

// 有变量名的写法
//go:generate gonacos_config --embed=GlobalConfig --config=some:SomeConfig --dataid=config --codec=yaml
type GlobalConfig struct {
    some SomeConfig
}

```

### 生成代码并使用

跳转到 `config` 目录下，执行命令，生成配置代码
```
> go generate
```

变量注册监听到nacos
```go
package main

import (
	"fmt"
	gonacos "github.com/reatang/go-simple-nacos"
)

var globalConfig = config.GlobalConfig{}

func main()  {
	// 初始化nacos，配置参数请看 nacos-sdk-go的文档
	conf := gonacos.NacosConf{
		Client: gonacos.ClientConfig{
			...
			NamespaceId: "", // 重要的参数
		},
		Servers: []gonacos.ServerConfig{
			{
				...
			}
		}
	}

	// DEFAULT_GROUP 可以被认为是环境，如：dev、test、prod
	ncc, err := gonacos.NewNacosConfigClient(conf, "DEFAULT_GROUP")
	if err != nil {
		panic(err)
	}

	// 标准配置注册，第二个参数可以传配置的默认值
	config.RegisterSomeConfig(ncc, nil)

	// 嵌入配置注册
	config.RegisterEmbedGlobalConfigSomeConfig(ncc, &globalConfig.SomeConfig)


	// 可以在业务中线程安全的使用了
	sc := config.GetSomeConfig()
	fmt.Println(sc.TestConfig)

	sc = globalConfig.GetSomeConfig()
	fmt.Println(sc.TestConfig)
}
```

## 关于性能

这两种使用方式的性能直接和 `sync.RWMutex` 和 `atomic.Value` 的性能有关。

|      | 独立结构体        | 嵌入式结构体       |
|------|--------------|--------------|
| 使用特性 | atomic.Value | sync.RWMutex |
| 性能   | 522787660    | 205038752    |  


---

吐槽：

nacos-sdk-go，完全是用java风格写的golang代码啊！
