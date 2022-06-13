package main

import (
	"fmt"
	"log"
	"strings"
)
import "github.com/tarm/serial"

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func main() {
	/*var hexaNumber string
	hexaNumber = "57"
	output, err := strconv.ParseInt(hexaNumberToInteger(hexaNumber), 16, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Output %d\n", output)
	fmt.Println(string(output))
	return*/

	config := &serial.Config{
		Name:        "COM3",
		Baud:        9600,
		ReadTimeout: 2,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	var n int

	buf := make([]byte, 1024)
	stream.Write([]byte{0x01, 0x52, 0x31, 0x02})
	n, err = stream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n > 0 {
		fmt.Println(buf[:n])
	}

	buf = make([]byte, 1024)
	stream.Write([]byte{0x50, 0x4f, 0x57, 0x45, 0x50, 0x28, 0x29, 0x03, 0x64})
	n, err = stream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n > 0 {
		fmt.Println(buf[:n], "-1")
	}
	n, err = stream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n > 0 {
		fmt.Println(buf[:n], "-2")
	}
	for _, b := range buf[:n] {
		fmt.Print(string(b))
	}

	buf = make([]byte, 1024)
	stream.Write([]byte{0x02, 0x50, 0x4f, 0x57, 0x45, 0x50, 0x28, 0x30, 0x2c, 0x31, 0x30, 0x30, 0x29, 0x0d, 0x0A, 0x03, 0x63})
	n, err = stream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n > 0 {
		fmt.Println(buf[:n], "-3")
	}
}
