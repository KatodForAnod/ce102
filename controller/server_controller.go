package controller

import (
	"ce102/config"
	"ce102/logsetting"
	"ce102/memory"
	"log"
)

type Controller struct {
	mem            memory.Memory
	ioTsController IoTsController
}

func (c *Controller) Init(mem memory.Memory, controller IoTsController) {
	c.mem = mem
	c.ioTsController = controller
}

func (c *Controller) GetInformation(deviceName string) ([]byte, error) {
	log.Println("controller get information of iot device", deviceName)

	load, err := c.mem.Load(deviceName)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return load, nil
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	log.Println("controller get lastNRowsLogs")
	file, err := logsetting.OpenLastLogFile()
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	logs, err := logsetting.GetNLastLines(file, nRows)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	return logs, nil
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	log.Println("controller AddIoTDevice")

	return nil
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	log.Println("controller RmIoTDevice")

	err := c.ioTsController.RemoveIoTsClients([]string{deviceName})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	log.Println("controller stop observe device")

	if err := c.ioTsController.StopObserveIoTDevice(deviceName); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
