package config

import "github.com/reatang/go-simple-nacos/tests/config/other_config"

//go:generate gonacos_config --config=SomeConfig --dataid=config --codec=yaml
type SomeConfig struct {
	TestConfig  string `yaml:"TestConfig"`
	TestConfig2 string `yaml:"TestConfig2"`
	TestConfig3 string `yaml:"TestConfig3"`
	TestConfig4 string `yaml:"TestConfig4"`
	TestConfig5 int64  `yaml:"TestConfig5"`
	TestConfig6 int64  `yaml:"TestConfig6"`
}

//go:generate gonacos_config --embed=GlobalConfig --config=SomeConfig --dataid=config --codec=yaml
//go:generate gonacos_config --embed=GlobalConfig --config=Other:other_config.OtherConfig --dataid=other_config --codec=yaml
type GlobalConfig struct {
	SomeConfig

	Other other_config.OtherConfig
}
