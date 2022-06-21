package client

import (
	"errors"
	"github.com/tarm/serial"
	"log"
	"time"
)

type Ce102 struct {
	deviceName string
	port       *serial.Port
	conf       *serial.Config

	isObserveInformProcess *bool
	stopObserve            chan bool

	buff [14]int
	crc8 int
	i    int
}

func (c *Ce102) Init(deviceName string, conf *serial.Config) {
	c.deviceName = deviceName
	c.conf = conf
	c.stopObserve = make(chan bool)
	c.isObserveInformProcess = new(bool)
}

func (c *Ce102) IsObserveInformProcess() bool {
	return *c.isObserveInformProcess
}

func (c *Ce102) StopObserveInform() error {
	log.Println("StopObserveInform device:", c.deviceName)
	if !c.IsObserveInformProcess() {
		err := errors.New("device " + c.deviceName + " not observing")
		log.Println(err)
		return err
	}

	select {
	case c.stopObserve <- true:
	default:
		err := errors.New("StopObserveInform cant stop thread")
		log.Println(err)
		return err
	}

	return nil
}

func (c Ce102) StartObserveInform(save func() error, duration time.Duration) error {
	log.Println("StartObserveInform device:", c.deviceName)
	if *c.isObserveInformProcess {
		err := errors.New("already observe")
		log.Println(err)
		return err
	}

	tr := true
	c.isObserveInformProcess = &tr

	for {
		select {
		case <-time.After(duration):
			if *c.isObserveInformProcess {
				if err := save(); err != nil {
					log.Println(err)
				}
			}
		case <-c.stopObserve:
			fl := false
			c.isObserveInformProcess = &fl
			return nil
		}
	}
}

func (c *Ce102) GetDeviceName() string {
	return c.deviceName
}

func (c *Ce102) Disconnect() error {
	log.Println("Disconnect", c.deviceName, "from port", c.conf.Name)
	if err := c.port.Close(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Ce102) Connect() error {
	log.Println("Connect port:", c.conf.Name)
	port, err := serial.OpenPort(c.conf)
	if err != nil {
		return err
	}

	c.port = port
	return nil
}

func (c *Ce102) BadCommand() {
	log.Println("BadCommand")
	c.startCePacket()
	c.sendByteToCE(0)
	c.endCEPacket()
}

func (c *Ce102) ReadSerialNumber(addrD uint16) []byte {
	log.Println("ReadSerialNumber addrD", addrD)
	c.sendCommandToCE(int(addrD), ReadSerialNumber_Command)
	return []byte{}
}

func (c *Ce102) Ping(addrD uint16) []byte {
	log.Println("Ping addrD", addrD)
	c.sendCommandToCE(int(addrD), Ping_Command)
	return []byte{}
}

func (c *Ce102) ReadTariffSum(addrD uint16) []byte {
	log.Println("ReadTariffSum addrD", addrD)
	c.sendCommandToCE(int(addrD), ReadTariffSum_Command)
	return []byte{}
}

func (c *Ce102) sendCommandToCE(addrD, command int) {
	log.Println("sendCommandToCE addrD", addrD)
	c.startCePacket()

	c.crc8 = 0
	c.sendByteToCE(OPT_CH)
	c.crc8 = crc8tab[c.crc8^OPT_CH]

	addrDH := addrD >> 8
	addrDL := addrD & 0xff
	c.sendByteToCE(uint16(addrDL))
	c.crc8 = crc8tab[c.crc8^addrDL]
	c.sendByteToCE(uint16(addrDH))
	c.crc8 = crc8tab[c.crc8^addrDH]

	c.sendByteToCE(0)
	c.crc8 = crc8tab[c.crc8^0]
	c.sendByteToCE(0)
	c.crc8 = crc8tab[c.crc8^0]

	passwd := [4]int{0x0, 0x0, 0x0, 0x0}
	for i := 0; i < 4; i++ {
		c.sendByteToCE(uint16(passwd[i]))
		c.crc8 = crc8tab[c.crc8^passwd[i]]
	}

	messageLength := 0
	if command == ReadTariffSum_Command {
		messageLength = 1
	}

	serv := DIRECT_REQ_CH + CLASS_ACCESS_CH + messageLength
	c.sendByteToCE(uint16(serv))
	c.crc8 = crc8tab[c.crc8^serv]

	commandH := (command >> 8)
	commandL := (command & 0xff)
	c.sendByteToCE(uint16(commandH))
	c.crc8 = crc8tab[c.crc8^commandH]
	c.sendByteToCE(uint16(commandL))
	c.crc8 = crc8tab[c.crc8^commandL]

	if command == ReadTariffSum_Command {
		c.sendByteToCE(0)
		c.crc8 = crc8tab[c.crc8^0]
	}

	c.sendByteToCE(uint16(c.crc8))
	c.endCEPacket()
}
func (c *Ce102) sendByteToCE(outByte uint16) {
	log.Println("sendByteToCE")
	if outByte == END_CH {
		c.sendByteToRS485(END_REPL_1_CH)
		c.sendByteToRS485(END_REPL_2_CH)
	} else if outByte == ESC_CH {
		c.sendByteToRS485(ESC_REPL_1_CH)
		c.sendByteToRS485(ESC_REPL_1_CH)
	} else {
		c.sendByteToRS485(outByte)
	}
}
func (c *Ce102) endCEPacket() {
	log.Println("endCEPacket")
	c.sendByteToRS485(END_CH)
}
func (c *Ce102) startCePacket() {
	log.Println("startCePacket")
	c.sendByteToRS485(END_CH)
}
func (c *Ce102) sendByteToRS485(outByte uint16) {
	log.Println("sendByteToRS485")
	_, errWrite := c.port.Write([]byte{byte(outByte)})
	if errWrite != nil {
		log.Fatal(errWrite)
	}
	err := c.port.Flush()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
}
