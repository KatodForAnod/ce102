package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ProxyServerAddr string      `json:"proxy_server_addr"`
	IoTsDevices     []IotConfig `json:"iots_devices"`
}

type IotConfig struct {
}

const configPath = "conf.config"

func LoadConfig() (loadedConf Config, err error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return Config{}, err
	}

	err = json.Unmarshal(data, &loadedConf)
	return
}
