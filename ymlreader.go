package filetransfer

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type YamlContent struct {
	Store StoreConfig `yaml:"store"`
}

func NewYmlContent(path string) YamlContent {
	if path == "" {
		path = "./config.yaml"
	}
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("problem read yaml file: %v", err)
	}
	var yamlContent YamlContent
	err = yaml.Unmarshal(fileContent, &yamlContent)
	if err != nil {
		log.Fatalf("problem unmarshal config.ymal: %v", err)
	}
	return yamlContent
}
