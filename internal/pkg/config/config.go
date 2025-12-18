package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBName       string `yaml:"db_name"`
	DBUsername   string `yaml:"db_username"`
	Port         string `yaml:"port"`
	JWTKey       string `yaml:"jwt_key"`
	RedisHost    string `yaml:"redis_host"`
	RedisDB      int    `yaml:"redis_db"`
	RedisExpires int    `yaml:"redis_expires"`
	SenderEmail  string `yaml:"sender_email"`
	AppPassword  string `yaml:"app_password"`
	DBPassword   string `yaml:"db_password"`
	DBHost       string `yaml:"db_host"`
}

func GetConfig() *Config {
	cfg := Config{}
	yamlFile, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &cfg
}
