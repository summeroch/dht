package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	App struct {
		Mode string `yaml:"mode"`
	}
	Elasticsearch struct {
		Addresses []string `yaml:"addresses"`
		Username  string   `yaml:"username"`
		Password  string   `yaml:"password"`
		Index     string   `yaml:"index"`
	}
	Tracker struct {
		List []string `yaml:"list"`
	}
}

func GetConfig() *Config {
	c := Config{}
	configFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}
