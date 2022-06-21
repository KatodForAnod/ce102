package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	ProxyServerAddr string      `json:"proxy_server_addr"`
	IoTsDevices     []IotConfig `json:"iots_devices"`
}

type IotConfig struct {
	DeviceName  string        `json:"device_name"`
	Port        string        `json:"port"`
	Baud        int           `json:"baud"`
	ReadTimeout time.Duration `json:"read_timeout''"`
	Size        byte          `json:"size"`
	//	Parity
	//	StopBits
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
