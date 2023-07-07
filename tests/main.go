package main

import (
	"context"
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	gonacos "github.com/reatang/go-simple-nacos"
	"github.com/reatang/go-simple-nacos/tests/config"
)

var (
	username    = flag.String("username", "", "nacos username")
	password    = flag.String("password", "", "nacos password")
	addr        = flag.String("addr", "", "nacos addr")
	namespaceId = flag.String("namespace_id", "", "nacos namespace id")
	grpcPort    = flag.Uint64("grpc_port", 0, "nacos grpc port")
)

func main() {
	flag.Parse()
	if *username == "" || *password == "" || *addr == "" || *namespaceId == "" || *grpcPort == 0 {
		fmt.Println("参数错误")
		return
	}

	conf := gonacos.NacosConf{
		Client: gonacos.ClientConfig{
			Username:            *username,
			Password:            *password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./log",
			CacheDir:            "./cache",
			LogLevel:            "debug",

			NamespaceId: *namespaceId,
		},
		Servers: []gonacos.ServerConfig{
			{
				IpAddr:   *addr,
				Port:     443,
				Scheme:   "https",
				GrpcPort: *grpcPort,
			},
		},
	}

	ncc, err := gonacos.NewNacosConfigClinet(conf)
	if err != nil {
		panic(err)
	}

	config.RegisterSomeConfig(ncc)

	tc, _ := context.WithTimeout(context.Background(), 20*time.Second)
	n := atomic.Int64{}
	for i := 0; i < 5; i++ {
		go func() {
			for {
				n.Add(1)

				aaaConfig := config.GetSomeConfig()
				if aaaConfig.TestConfig != "aaa" && aaaConfig.TestConfig != "bbb" {
					fmt.Println("值有错误")
				}
			}
		}()
	}

	<-tc.Done()
	fmt.Printf("循环次数：%d\n", n.Load())
}
