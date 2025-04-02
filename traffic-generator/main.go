package main

import (
	"fmt"
	"os"
)

func main() {
	// Clear log file on restart
	err := os.WriteFile("log.txt", []byte{}, 0644)
	if err != nil {
		fmt.Println("Error clearing log file:", err)
		return
	}

	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println("Starting Traffic Generator...")
	Simulator(config.APICount, config.APIRate, config.CollectorURL)
	fmt.Println("Traffic Generator finished.")
}
