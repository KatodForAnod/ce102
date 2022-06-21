package controller

import (
	"ce102/client"
	"ce102/memory"
	"errors"
	"log"
)

type IoTsController struct {
	ioTDevices map[string]client.Ce102
	mem        memory.Memory
}

func (c *IoTsController) Init(mem memory.Memory) {
	c.ioTDevices = make(map[string]client.Ce102)
	c.mem = mem
}

func (c *IoTsController) AddIoTsClients(devices []client.Ce102) error {
	log.Println("AddIoTsClients")
	for _, device := range devices {
		if _, isExist := c.ioTDevices[device.GetDeviceName()]; isExist {
			err := errors.New("device " + device.GetDeviceName() + " already exist")
			log.Println(err)
			return err
		}
	}

	for _, device := range devices {
		err := device.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		c.ioTDevices[device.GetDeviceName()] = device
	}
	return nil
}

func (c *IoTsController) RemoveIoTsClients(devicesName []string) error {
	log.Println("RemoveIoTsClients")
	var founded []client.Ce102
	for _, deviceName := range devicesName {
		if iot, isExist := c.ioTDevices[deviceName]; !isExist {
			err := errors.New("device " + deviceName + " not exist")
			log.Println(err)
			return err
		} else {
			founded = append(founded, iot)
		}
	}

	for _, tClient := range founded {
		if tClient.IsObserveInformProcess() {
			if err := tClient.StopObserveInform(); err != nil {
				log.Println(err)
			}
		}
		if err := tClient.Disconnect(); err != nil {
			log.Println(err)
		}
		delete(c.ioTDevices, tClient.GetDeviceName())
	}

	return nil
}

func (c *IoTsController) StopObserveIoTDevice(deviceName string) error {
	log.Println("Stop observe device:", deviceName)
	iot, isExist := c.ioTDevices[deviceName]
	if !isExist {
		err := errors.New("device not exist")
		log.Println(err)
		return err
	}

	if err := iot.StopObserveInform(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
