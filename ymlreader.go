package filetransfer

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"runtime"
)

type YamlContent struct {
	Store StoreConfig `yaml:"store"`
}

func NewYamlContent(path string) (*YamlContent, error) {
	if path == "" {
		path = getDefaultConfigPath()
	}
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var yamlContent YamlContent
	err = yaml.Unmarshal(fileContent, &yamlContent)
	if err != nil {
		return nil, err
	}
	return &yamlContent, nil
}

func getDefaultConfigPath() string {
	os := runtime.GOOS
	if os == "linux" {
		return "/etc/filetransfer/config.yml"
	} else {
		return "./config.yml"
	}
}
