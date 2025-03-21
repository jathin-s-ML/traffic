package main

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrafficGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Traffic Generator Suite")
}

var _ = Describe("ConfigParser", func() {
	DescribeTable("parsing configuration",
		func(rawConfig map[string]string, expectedConfig *Config, expectedError string) {
			config, err := ConfigParser(rawConfig)

			if expectedError != "" {
				Expect(err).To(HaveOccurred(), "Expected an error but got nil")
				Expect(err.Error()).To(ContainSubstring(expectedError), "Unexpected error message")
			} else {
				Expect(err).NotTo(HaveOccurred(), "Expected no error but got one")
				Expect(config).To(Equal(expectedConfig), "Parsed config does not match expected values")
			}
		},

		Entry("Valid configuration",
			map[string]string{
				"NO_OF_API":    "100",
				"API_RATE":     "10/s",
				"COLLECTOR_URL": "http://localhost:8080",
			},
			&Config{
				APICount:     100,
				CollectorURL: "http://localhost:8080",
				APIRate:      100 * time.Millisecond, 
			},
			"",
		),

		Entry("Invalid NO_OF_API",
			map[string]string{
				"NO_OF_API":    "-1",
				"API_RATE":     "10/s",
				"COLLECTOR_URL": "http://localhost:8080",
			},
			nil,
			"invalid NO_OF_API value",
		),

		Entry("Invalid API_RATE",
			map[string]string{
				"NO_OF_API":    "100",
				"API_RATE":     "invalid-rate",
				"COLLECTOR_URL": "http://localhost:8080",
			},
			nil,
			"invalid API_RATE format",
		),

		Entry("Missing COLLECTOR_URL",
			map[string]string{
				"NO_OF_API": "100",
				"API_RATE":  "10/s",
			},
			nil,
			"COLLECTOR_URL not set",
		),
	)
})
