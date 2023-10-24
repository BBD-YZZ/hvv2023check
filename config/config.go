package config

import (
	"xyz/tools"

	"gopkg.in/ini.v1"
)

type Config struct {
	Mysql_Config    `ini:"mysql"`
	Ceye_Api_Config `ini:"api"`
}

type Mysql_Config struct {
	Address  string `ini:"address"`
	Port     string `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type Redis_Config struct {
}

type Ceye_Api_Config struct {
	Address string `ini:"address"`
	Type    string `ini:"type"`
	Token   string `ini:"token"`
	Filter  string `ini:"filter"`
}

func GetConfig() (config *Config, err error) {
	path, _ := tools.GetAbPath()

	cfg, err := ini.Load(path + "/config/config.ini")
	if err != nil {
		return &Config{}, err
	}

	config = new(Config)
	err = cfg.MapTo(config)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}
