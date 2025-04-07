package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

// Global logger instance
var logger *zap.Logger

// ✅ Initializes a global logger
func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = ""       // Remove timestamp
	config.EncoderConfig.LevelKey = ""      // Remove log level
	config.EncoderConfig.CallerKey = ""     // Remove caller info
	config.EncoderConfig.MessageKey = "msg" // Keep only the message

	var err error
	logger, err = config.Build()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
}
