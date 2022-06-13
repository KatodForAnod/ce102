package main

import "fmt"
import "log"
import "github.com/tarm/serial"

func main() {
	config := &serial.Config{
		Name: "/dev/ttyAMA0",
		Baud: 9600,
		ReadTimeout: 1,
		Size: 8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := string(buf[:n])
		fmt.Println(s)
	}
}
