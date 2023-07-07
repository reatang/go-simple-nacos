package config

//go:generate gonacos_config --config=SomeConfig --group=13 --dataid=config --codec=yaml
type SomeConfig struct {
	TestConfig  string `yaml:"TestConfig"`
	TestConfig2 string `yaml:"TestConfig2"`
	TestConfig3 string `yaml:"TestConfig3"`
	TestConfig4 string `yaml:"TestConfig4"`
	TestConfig5 int64  `yaml:"TestConfig5"`
	TestConfig6 int64  `yaml:"TestConfig6"`
}
