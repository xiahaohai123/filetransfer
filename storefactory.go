package filetransfer

type StoreConfig struct {
	Type   string `yaml:"type"`
	Config Config `yaml:"config"`
}

type Config struct {
	Redis RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func getStoreConfig() StoreConfig {
	return StoreConfig{}
}
