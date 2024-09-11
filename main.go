package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// CalculateChecksum calculates the NMEA checksum for a given sentence
func CalculateChecksum(sentence string) string {
	var checksum byte
	for i := 0; i < len(sentence); i++ {
		checksum ^= sentence[i]
	}
	return fmt.Sprintf("%02X", checksum)
}

// BuildNMEASentence builds the NMEA sentence with the checksum
func BuildNMEASentence(fields ...string) string {
	baseSentence := strings.Join(fields, ",")
	checksum := CalculateChecksum(baseSentence)
	return fmt.Sprintf("$%s*%s", baseSentence, checksum) // Don't append CRLF here, let SendNMEASentence handle it
}

// SendNMEASentence sends the NMEA sentence over the specified serial port with CRLF
func SendNMEASentence(portName string, baudRate int, sentence string) error {
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

	_, err = port.Write([]byte(fullSentence))
	if err != nil {
		return err
	}

	// Wait a little before reading the response
	time.Sleep(500 * time.Millisecond)

	// Buffer to store incoming data
	buf := make([]byte, 128)
	lineCount := 0

	// Read up to 10 lines from the port
	for lineCount < 10 {
		n, err := port.Read(buf)
		if err != nil {
			return err
		}
		if n > 0 {
			// Print the received data
			fmt.Print(string(buf[:n]))

			// Count the number of newlines to track lines
			lineCount += strings.Count(string(buf[:n]), "\n")
		}
	}
	defer port.Close()

	fmt.Println("\n\nCheck $PERDACK sentence")//別にプログラム的に探してもいい
	return nil
}

func main() {
	// Parse command-line arguments
	port := flag.String("p", "/dev/tty.usbserial-10", "Serial port to use")//TB-1 on Mac
	baudRate := flag.Int("s", 38400, "Baud rate for the serial port (default: 38400)")
	getFlag := flag.String("g", "", "API,CFG,SYS to generate NMEA sentences")
	executeFlag := flag.String("z", "", "API,CFG,SYS with parameters to send via serial port")
	baudUpdateFlag := flag.String("S", "", "Send NMEA sentence to update baud rate")
	versionFlag := flag.Bool("V", false, "Execute SYS [VERSION] command")
	helpFlag := flag.Bool("help", false, "Display help message")
	flag.Parse()

	// Display help if needed
	if *helpFlag {
		ShowHelp()
		return
	}

	// If no valid options are provided, show an error
	if *getFlag == "" && *executeFlag == "" && !*versionFlag && *baudUpdateFlag == "" {
		log.Fatal("You must specify either -g or -z or -S or -V with appropriate parameters")
		Showh()
	}

	// Send a sentence to update the baud rate if -S flag is provided
	if *baudUpdateFlag != "" {
		sentence := BuildNMEASentence("PERDCFG", "UART1", *baudUpdateFlag)
		fmt.Println("Sending baud rate update:", sentence)
		err := SendNMEASentence(*port, *baudRate, sentence)
		if err != nil {
			log.Fatalf("Failed to send baud rate update sentence: %v", err)
		}
	}

	// Execute SYS [VERSION] command if -V flag is provided
	if *versionFlag {
		nmeaSentence := BuildNMEASentence("PERDSYS", "VERSION")
		fmt.Printf("Executing SYS [VERSION]: %s\n", nmeaSentence)
		err := SendNMEASentence(*port, *baudRate, nmeaSentence)
		if err != nil {
			log.Fatalf("Failed to send SYS [VERSION] command: %v", err)
		}
		return
	}

	// Query NMEA sentences if -g flag is provided
	if *getFlag != "" {
		commands := strings.Split(*getFlag, ",")
		api := commands[0]
		sentence := BuildNMEASentence("PERDAPI", api, "QUERY")
		fmt.Println("Sending:", sentence)
		err := SendNMEASentence(*port, *baudRate, sentence)
		if err != nil {
			log.Fatalf("Failed to send NMEA sentence: %v", err)
		}
	}

	// Send NMEA sentences over serial port if -z flag is provided
	if *executeFlag != "" {
		// Split -z argument into parts (e.g., "PPS,VCLK,1,0,200,0,0")
		commands := strings.Split(*executeFlag, ",")
		api := commands[0]
		params := commands[1:]

		// Special case for CFG-NMEAOUT
		if api == "NMEAOUT" {
			sentence := BuildNMEASentence(append([]string{"PERDCFG", "NMEAOUT"}, params...)...)
			fmt.Println("Sending:", sentence)
			err := SendNMEASentence(*port, *baudRate, sentence)
			if err != nil {
				log.Fatalf("Failed to send NMEA sentence: %v", err)
			}
		} else {
			// General case for other APIs
			sentence := BuildNMEASentence(append([]string{"PERDAPI", api}, params...)...)
			fmt.Println("Sending:", sentence)
			err := SendNMEASentence(*port, *baudRate, sentence)
			if err != nil {
				log.Fatalf("Failed to send NMEA sentence: %v", err)
			}
		}
	}
}
