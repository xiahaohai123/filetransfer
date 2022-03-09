package filetransfer_test

import (
	"os"
	"reflect"
	"summersea.top/filetransfer"
	"testing"
)

const memoryStoreType = "MemoryStore"
const redisStoreType = "redisStore"

func TestCreateStore(t *testing.T) {
	t.Run("create specified store", func(t *testing.T) {
		testCases := []struct {
			baseConfig filetransfer.YamlContent
			wantType   string
		}{
			{filetransfer.YamlContent{Store: filetransfer.StoreConfig{Type: "memory"}}, memoryStoreType},
			{filetransfer.YamlContent{Store: filetransfer.StoreConfig{Type: "mmory"}}, memoryStoreType},
			{filetransfer.YamlContent{}, memoryStoreType},
			{filetransfer.YamlContent{Store: filetransfer.StoreConfig{Type: "redis"}}, redisStoreType},
			{filetransfer.YamlContent{Store: filetransfer.StoreConfig{Type: "redis", Config: filetransfer.Config{Redis: filetransfer.RedisConfig{Address: "localhost:6379"}}}}, redisStoreType},
			{filetransfer.YamlContent{Store: filetransfer.StoreConfig{Type: "redis", Config: filetransfer.Config{Redis: filetransfer.RedisConfig{Address: "localhost:6381"}}}}, memoryStoreType},
		}

		for _, test := range testCases {
			file := createYmlConfigFile()
			marshal := getConfigBytes(t, test.baseConfig)
			_, _ = file.Write(marshal)
			_ = file.Close()

			dataStore := filetransfer.CreateStoreByConfig()
			assertStringEqual(t, reflect.ValueOf(dataStore).Elem().Type().Name(), test.wantType)
			_ = os.Remove(file.Name())
		}
	})
}
