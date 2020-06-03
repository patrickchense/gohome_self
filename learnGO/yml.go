package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	VerifyToken string `yaml:"verify_token"`
	AccessToken string `yaml:"access_token"`
	AppSecret   string `yaml:app_secret`
}

func parseContentFile() string {
	contentFile, err := ioutil.ReadFile("content.yml")

	if err != nil {
		log.Printf("Error opening content file: %s\n", err)
		panic(err)
	}

	er, _ := yaml.Marshal(contentFile)

	if er != nil {
		log.Printf("marshal content file: %s\n", er)
	}

	return string(er)
}

func (c *Config) readYml() *Config {
	yamlFile, err := ioutil.ReadFile("bot.config.yml")
	if err != nil {
		log.Printf("error opening yml file: %s\n", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("unmarshal config file: %s\n", err)
	}

	fmt.Printf("here is the parsed content.yml: %s\n", c)

	return c
}

func getToken() string {
	var c Config
	c.readYml()
	v, err := json.Marshal(c)
	if err != nil {
		log.Printf("err marshalling json file: %s\n", err)
	}
	return string(v)
}
