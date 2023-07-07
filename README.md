# go-simple-nacos
更有golang风格的使用 nacos-sdk-go 客户端

https://github.com/nacos-group/nacos-sdk-go 太难用了，简化使用方案

## 安装nacos

    略....

## 使用

1、设置映射配置的结构体，并编写生成命令
```go
// file: config/nacos_config.go
package config

//go:generate gonacos_config --config=SomeConfig --dataid=config --codec=yaml
type SomeConfig struct {
	TestConfig  string `yaml:"TestConfig"`
}
```

2、跳转到 `config` 目录下，执行命令，生成配置代码
```
> go generate
```

3、变量注册监听到nacos
```go

// 初始化nacos，配置参数请看 nacos-sdk-go的文档
conf := gonacos.NacosConf{
    ...
}

ncc, err := gonacos.NewNacosConfigClinet(conf, "DEFAULT_GROUP")
if err != nil {
    panic(err)
}

// 注册
config.RegisterSomeConfig(ncc)


// 可以在业务中线程安全的使用了
sc := config.GetSomeConfig()
fmt.Println(sc.TestConfig)
```



吐槽：
nacos-sdk-go，完全是用java风格写的golang代码啊！
