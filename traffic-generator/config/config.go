package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	APICount     int
	APIRate      time.Duration
	CollectorURL string
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var rawConfig map[string]string
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, err
	}

	return ConfigParser(rawConfig)
}

func ConfigParser(rawConfig map[string]string) (*Config, error) {
	apiCount, err := strconv.Atoi(rawConfig["NO_OF_API"])
	if err != nil || apiCount <= 0 || apiCount > 8192 {
		return nil, fmt.Errorf("invalid NO_OF_API value")
	}

	re := regexp.MustCompile(`^(\d+)/([smhSMH])$`)
	matches := re.FindStringSubmatch(rawConfig["API_RATE"])
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid API_RATE format, use '2/s', '100/m', or '3000/h'")
	}

	rateValue, _ := strconv.Atoi(matches[1])
	timeUnit := strings.ToLower(matches[2])
	var interval time.Duration
	switch timeUnit {
	case "s":
		interval = time.Second / time.Duration(rateValue)
	case "m":
		interval = time.Minute / time.Duration(rateValue)
	case "h":
		interval = time.Hour / time.Duration(rateValue)
	default:
		return nil, fmt.Errorf("unknown time unit")
	}

	if rawConfig["COLLECTOR_URL"] == "" {
		return nil, fmt.Errorf("COLLECTOR_URL not set")
	}

	return &Config{
		APICount:     apiCount,
		APIRate:      interval,
		CollectorURL: rawConfig["COLLECTOR_URL"],
	}, nil
}
