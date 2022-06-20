package ce102

import (
	"github.com/tarm/serial"
	"log"
	"time"
)

type Ce102 struct {
	port *serial.Port
	buff [14]int
	crc8 int
	i    int
}

func (c *Ce102) Connect(conf *serial.Config) error {
	port, err := serial.OpenPort(conf)
	if err != nil {
		return err
	}

	c.port = port
	return nil
}

func (c *Ce102) BadCommand() {
	c.startCePacket()
	c.sendByteToCE(0)
	c.endCEPacket()
}

func (c *Ce102) ReadSerialNumber(addrD uint16) {
	c.sendCommandToCE(int(addrD), ReadSerialNumber_Command)
}

func (c *Ce102) Ping(addrD uint16) {
	c.sendCommandToCE(int(addrD), Ping_Command)
}

func (c *Ce102) ReadTariffSum(addrD uint16) {
	c.sendCommandToCE(int(addrD), ReadTariffSum_Command)
}

func (c *Ce102) sendCommandToCE(addrD, command int) {
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
	c.sendByteToRS485(END_CH)
}
func (c *Ce102) startCePacket() {
	c.sendByteToRS485(END_CH)
}
func (c *Ce102) sendByteToRS485(outByte uint16) {
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
