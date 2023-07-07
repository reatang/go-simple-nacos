package gonacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ClientConfig = constant.ClientConfig
type ServerConfig = constant.ServerConfig

type NacosConf struct {
	// 客户端配置
	Client ClientConfig

	// 服务端配置
	Servers []ServerConfig
}

// Nacos 配置中心客户端
type NacosConfigClient struct {
	config *NacosConf
	client config_client.IConfigClient
}

// 初始化函数
func NewNacosConfigClinet(conf NacosConf) (*NacosConfigClient, error) {
	// 创建客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &conf.Client,
			ServerConfigs: conf.Servers,
		},
	)
	if err != nil {
		return nil, err
	}

	return &NacosConfigClient{
		config: &conf,
		client: client,
	}, nil
}

// WatchF 监听配置变化
func (ncc *NacosConfigClient) WatchF(groupid, dataid string, onChange func(namespace, group, dataId, data string)) {
	// 1、获取配置
	params := vo.ConfigParam{
		DataId: dataid,
		Group:  groupid,
	}
	content, err := ncc.client.GetConfig(params)
	if err != nil {
		panic(err)
	}

	// 初始化一次
	clientConfig, _ := ncc.client.(*config_client.ConfigClient).GetClientConfig()
	onChange(clientConfig.NamespaceId, groupid, dataid, content)

	// 2、监听配置
	watchParams := params
	watchParams.OnChange = onChange
	err = ncc.client.ListenConfig(watchParams)
	if err != nil {
		panic(err)
	}
}
