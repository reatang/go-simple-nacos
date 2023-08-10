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

// NacosConfigClient 配置中心客户端
type NacosConfigClient struct {
	// 默认分组，一般被用作环境区分，如：dev、test、prod
	defaultGroup string

	config *NacosConf
	client config_client.IConfigClient
}

// NewNacosConfigClient 初始化函数
func NewNacosConfigClient(conf NacosConf, defaultGroup string) (*NacosConfigClient, error) {
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
		defaultGroup: defaultGroup,
		config:       &conf,
		client:       client,
	}, nil
}

// WatchF 监听配置变化
func (ncc *NacosConfigClient) WatchF(group, dataid string, onChange func(namespace, group, dataId, data string)) {
	if group == "" {
		group = ncc.defaultGroup
	}

	// 1、获取配置
	params := vo.ConfigParam{
		DataId: dataid,
		Group:  group,
	}
	content, err := ncc.client.GetConfig(params)
	if err != nil {
		panic(err)
	}

	// 初始化一次
	clientConfig, _ := ncc.client.(*config_client.ConfigClient).GetClientConfig()
	onChange(clientConfig.NamespaceId, group, dataid, content)

	// 2、监听配置
	watchParams := params
	watchParams.OnChange = onChange
	err = ncc.client.ListenConfig(watchParams)
	if err != nil {
		panic(err)
	}
}
