package main

import (
	"golang.org/x/text/encoding/charmap"
	"log"
	"strings"
)

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
	/*config := &serial.Config{
		Name:        "COM3",
		Baud:        9600,
		ReadTimeout: 2,
		Size:        8,
	}*/
}
