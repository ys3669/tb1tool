//roop.go
package main
import (
	"github.com/tarm/serial"
)

func roop(portName string, baudRate int, sentence string) error {
	config := &serial.Config{
		Name: portName,
		Baud: baudRate,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		return err
	}

	// Append CRLF to the sentence
	fullSentence := sentence + "\r\n"
	for true {
		_, err = port.Write([]byte(fullSentence))
		if err != nil {
			return err
		}
	}
	defer port.Close()
	return nil
}
