package filetransfer_test

import (
	"gopkg.in/yaml.v2"
	"os"
	"summersea.top/filetransfer"
	"summersea.top/filetransfer/test"
	"testing"
)

func TestNewYamlContent(t *testing.T) {
	t.Run("no config file", func(t *testing.T) {
		content, err := filetransfer.NewYamlContent("")
		testutil.AssertNil(t, content)
		testutil.AssertNotNil(t, err)
	})

	t.Run("no content file", func(t *testing.T) {
		file := createYmlConfigFile()
		_ = file.Close()
		content, err := filetransfer.NewYamlContent("")
		testutil.AssertNil(t, err)
		testutil.AssertNotNil(t, content)
		err = os.Remove(file.Name())
		if err != nil {
			t.Errorf("failed delete: %v", err)
		}
		testutil.AssertStructEquals(t, *content, filetransfer.YamlContent{})
	})

	t.Run("content file", func(t *testing.T) {
		yamlContent := filetransfer.YamlContent{Store: filetransfer.StoreConfig{
			Config: filetransfer.Config{Redis: filetransfer.RedisConfig{Address: "localhost:6379"}},
			Type:   "redis",
		}}
		marshal := getConfigBytes(t, yamlContent)
		file := createYmlConfigFile()
		_, _ = file.Write(marshal)
		_ = file.Close()

		content, err := filetransfer.NewYamlContent("")
		if err != nil {
			t.Errorf("problem get yaml content: %v", err)
		}
		_ = os.Remove(file.Name())
		testutil.AssertStructEquals(t, *content, yamlContent)
	})
}

func createYmlConfigFile() *os.File {
	file, _ := os.OpenFile("config.yml", os.O_RDWR|os.O_CREATE, 0666)
	return file
}

func getConfigBytes(t *testing.T, content filetransfer.YamlContent) []byte {
	t.Helper()
	marshal, err := yaml.Marshal(content)
	if err != nil {
		t.Fatalf("problem marshal yaml: %v", err)
	}
	return marshal
}
