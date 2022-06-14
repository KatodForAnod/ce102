package main

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
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

func ReadeWindows1251File(text string) (string, error) {
	dec := charmap.Windows1251.NewDecoder()
	out, err := dec.String(text)
	if err != nil {
		log.Fatal(err)
	}

	return out, nil
}

func WriteWindows1251File(text string) (string, error) {
	cod := charmap.Windows1251.NewEncoder()
	out, err := cod.String(text)
	if err != nil {
		log.Fatal(err)
	}

	return out, nil
}

func main() {
	/*fmt.Println([]byte{0x50, 0x4f, 0x57, 0x45, 0x50, 0x28, 0x29, 0x03, 0x64})
	fmt.Println(string([]byte{0x50, 0x4f, 0x57, 0x45, 0x50, 0x28, 0x29, 0x03, 0x64}))
	fmt.Println(string([]byte{80, 79, 87, 69, 80, 40, 41, 3, 100}))


	fmt.Println(ReadeWindows1251File(string([]byte{80 ,207, 215, 197, 80, 40 ,48, 172 ,48 ,177 ,48 ,169, 141, 10, 3, 99 })))
	//fmt.Println(hexaNumberToInteger(string([]byte{207, 2})))
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
		fmt.Println(string(buf[:n]))
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
