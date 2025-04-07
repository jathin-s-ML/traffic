package main

import (
	"fmt"
	"log"
	"os"

	"traffic-generator/config"
	"traffic-generator/generator"
)

func main() {
	// Clear log file on restart
	err := os.WriteFile("log.txt", []byte{}, 0644)
	if err != nil {
		fmt.Println("Error clearing log file:", err)
		return
	}

	// Load Configuration
	cfg, err := config.ReadConfig() // ✅ Rename local variable to `cfg`
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	fmt.Println("Starting Traffic Generator...")
	generator.Simulator(cfg.APICount, cfg.APIRate, cfg.CollectorURL) // ✅ Use `generator.Simulator`
	fmt.Println("Traffic Generator finished.")
}
