package filetransfer

import "log"

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

func CreateStoreByConfig() DataStore {
	storeConfig := getStoreConfig()
	if storeConfig.Type == "redis" {
		dataStore, err := handleRedisConfig(storeConfig.Config.Redis)
		if err != nil {
			log.Printf("problem create redis store: %v", err)
		} else {
			return dataStore
		}
	}
	return handleStoreConfig()
}

func handleRedisConfig(config RedisConfig) (DataStore, error) {
	store, err := NewRedisStore(config.Address, config.Password, config.DB)
	if err != nil {
		return nil, err
	} else {
		log.Printf("[info] success to create redis store")
		return store, nil
	}
}

func handleStoreConfig() DataStore {
	log.Printf("[info] success to create memory store")
	return NewMemoryStore()
}

func getStoreConfig() *StoreConfig {
	content, err := NewYamlContent("")
	if err != nil {
		log.Printf("[error]problem get yaml content: %v \n", err)
		return &StoreConfig{}
	}
	return &content.Store
}
