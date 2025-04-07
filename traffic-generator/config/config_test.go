package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigParser_ValidConfig(t *testing.T) {
	rawConfig := map[string]string{
		"NO_OF_API":     "10",
		"API_RATE":      "5/s",
		"COLLECTOR_URL": "http://traffic-stats-col:8080/collect",
	}

	config, err := ConfigParser(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, 10, config.APICount)
	assert.Equal(t, 200*time.Millisecond, config.APIRate) // 5 requests per second = 200ms interval
	assert.Equal(t, "http://traffic-stats-col:8080/collect", config.CollectorURL)
}

func TestConfigParser_InvalidAPICount(t *testing.T) {
	rawConfig := map[string]string{
		"NO_OF_API":     "-1", // Invalid value
		"API_RATE":      "5/s",
		"COLLECTOR_URL": "http://traffic-stats-col:8080/collect",
	}

	config, err := ConfigParser(rawConfig)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "invalid NO_OF_API value")
}

func TestConfigParser_InvalidAPIRateFormat(t *testing.T) {
	rawConfig := map[string]string{
		"NO_OF_API":     "10",
		"API_RATE":      "1000x/m", // Invalid format
		"COLLECTOR_URL": "http://traffic-stats-col:8080/collect",
	}

	config, err := ConfigParser(rawConfig)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "invalid API_RATE format")
}

func TestConfigParser_EmptyCollectorURL(t *testing.T) {
	rawConfig := map[string]string{
		"NO_OF_API":     "10",
		"API_RATE":      "5/s",
		"COLLECTOR_URL": "", // Missing URL
	}

	config, err := ConfigParser(rawConfig)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "COLLECTOR_URL not set")
}

// Helper function to create a temporary config file
func createTempConfigFile(content string) (string, error) {
	file, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestReadConfig_ValidFile(t *testing.T) {
	mockConfig := `
NO_OF_API: "10"
API_RATE: "5/s"
COLLECTOR_URL: "http://traffic-stats-col:8080/collect"
`
	// Create a temporary config file
	tempFile, err := createTempConfigFile(mockConfig)
	assert.NoError(t, err)
	defer os.Remove(tempFile) // Cleanup after test

	// Rename temp file to "config.yaml" temporarily
	oldPath := "config.yaml"
	os.Rename(tempFile, oldPath)
	defer os.Rename(oldPath, tempFile) // Restore the file after test

	// Run the test
	config, err := ReadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 10, config.APICount)
	assert.Equal(t, "http://traffic-stats-col:8080/collect", config.CollectorURL)
}

func TestReadConfig_FileNotFound(t *testing.T) {
	// _ = os.Remove("config.yaml") // Ensure file does not exist

	config, err := ReadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}
