package config

import (
	"os"
	"testing"
)

const confBody = `{
    "proxy_server_addr":"127.0.0.1:5300"
}`

func createConfig(t *testing.T) error {
	file, err := os.CreateTemp("", configPath)
	if err != nil {
		t.Error("cant create temp conf file")
		return err
	}

	configPath = file.Name()

	_, err = file.WriteString(confBody)
	if err != nil {
		t.Error("cant write to temp conf file")
		return err
	}

	return nil
}

func deleteConfig(t *testing.T) error {
	err := os.Remove(configPath)
	if err != nil {
		t.Error("cant delete temp conf file")
		return err
	}
	return nil
}

func TestLoadConfig_Success(t *testing.T) {
	err := createConfig(t)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	err = deleteConfig(t)
	if err != nil {
		t.Error("cant delete temp conf file")
		return
	}
}

func TestLoadConfig_Fail(t *testing.T) {
	configPath = "notExist.txt"
	_, err := LoadConfig()
	if err == nil {
		t.Error("func must return error")
		return
	}
}
