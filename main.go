package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	return fmt.Sprintf("$%s*%s\r\n", baseSentence, checksum)
}

// SendNMEASentence sends the NMEA sentence over the specified serial port
func SendNMEASentence(portName string, baudRate int, sentence string) error {
	config := &serial.Config{
		Name: portName,
		Baud: baudRate,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		return err
	}

	// Write the NMEA sentence
	_, err = port.Write([]byte(sentence))
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

	return nil
}

func findSerialPort(vendorID, productID string) (string, error) {
	// This function will only be implemented for Linux
	if runtime.GOOS == "linux" {
		devicesPath := "/sys/bus/usb-serial/devices/"

		devices, err := os.ReadDir(devicesPath)
		if err != nil {
			return "", err
		}

		for _, device := range devices {
			devicePath := filepath.Join(devicesPath, device.Name(), "device", "uevent")
			data, err := os.ReadFile(devicePath)
			if err != nil {
				continue
			}

			content := string(data)
			if strings.Contains(content, "PRODUCT="+vendorID+"/"+productID) {
				devPath := filepath.Join("/dev", device.Name())
				if _, err := os.Stat(devPath); err == nil {
					return devPath, nil
				}
			}
		}
		return "/dev/ttyUSB0", nil // Example output
	}
	return "", fmt.Errorf("automatic serial port detection is not supported on this OS")
}

func main() {
	// Parse command-line arguments
	port := flag.String("D", "/dev/ttyUSB0", "Serial port to use")
	baudRate := flag.Int("s", 38400, "Baud rate for the serial port (default: 38400)")
	generateFlag := flag.String("g", "", "API,CFG,SYS to generate NMEA sentences")
	executeFlag := flag.String("z", "", "API,CFG,SYS with parameters to send via serial port")
	baudUpdateFlag := flag.String("S", "", "Send NMEA sentence to update baud rate")
	versionFlag := flag.Bool("V", false, "Execute SYS [VERSION] command")
	hFlag := flag.Bool("h", false, "Display help message")
	helpFlag := flag.Bool("help", false, "Display long help message")
	flag.Parse()

	if *hFlag {
		Showh()
		return
	}

	if *helpFlag {
		ShowHelp()
		return
	}

	if *generateFlag == "" && *executeFlag == "" && *baudUpdateFlag == "" && !*versionFlag && !*helpFlag {
		log.Fatal("You must specify either -g or -z or -S or -V with appropriate parameters")
	}

	if *port == "" {
		var err error
		*port, err = findSerialPort("0403", "6015")
		if err != nil {
			log.Fatalf("Failed to find serial port with Vendor ID: %s, Product ID: %s. Error: %v", "6015", "0403", err)
		}
	}

	if *versionFlag {
		nmeaSentence := "PERDSYS,VERSION"
		formattedSentence := BuildNMEASentence(nmeaSentence)
		fmt.Printf("Executing SYS [VERSION]: %s\n", formattedSentence)
		err := SendNMEASentence(*port, *baudRate, formattedSentence)
		if err != nil {
			log.Fatalf("Failed to send SYS [VERSION] command: %v", err)
		}
		return
	}

	// Generate NMEA sentences if -g flag is provided
	if *generateFlag != "" {
		commands := strings.Split(*generateFlag, ",")
		for _, cmd := range commands {
			sentence := BuildNMEASentence("PERDAPI", cmd, "QUERY")
			fmt.Println(sentence)
		}
	}

	// Send NMEA sentences over serial port if -z flag is provided
	if *executeFlag != "" {
		commands := strings.Split(*executeFlag, ",")
		for _, cmd := range commands {
			parts := strings.Split(cmd, " ")
			api := parts[0]
			params := parts[1:]
			var sentence string
			if api == "NMEAOUT" {
				sentence = BuildNMEASentence("PERDCFG", "NMEAOUT")
			} else if api == "VERSION" {
				sentence = BuildNMEASentence("PERDSYS", "VERSION")
			} else {
				sentence = BuildNMEASentence(append([]string{"PERDAPI", api}, params...)...)
			}
			fmt.Println("Sending:", sentence)
			err := SendNMEASentence(*port, *baudRate, sentence)
			if err != nil {
				log.Fatalf("Failed to send NMEA sentence: %v", err)
			}
		}
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
}
