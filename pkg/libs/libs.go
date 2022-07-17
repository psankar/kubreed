package libs

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ConfigEnvVar = "KUBREED_CONFIG"
)

const (
	minDuration = time.Millisecond * 10
	maxDuration = time.Minute * 2
	minAPI      = 1
	maxAPI      = 20
)

type Config struct {
	APICount       int      `json:"apiCount"`
	RPS            int      `json:"rps"`
	RemoteServices []string `json:"RemoteServices"`

	// See https://stackoverflow.com/questions/48050945/how-to-unmarshal-json-into-durations
	// to understand why we need this mess
	ResponseTime         time.Duration `json:"-"`
	ResponseTimeInternal string        `json:"responseTime"`
}

func GetConfigFromJSON(s string) (*Config, error) {
	var c Config
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}

	c.ResponseTime, err = time.ParseDuration(c.ResponseTimeInternal)
	if err != nil {
		return nil, err
	}

	if c.APICount < minAPI || c.APICount > maxAPI {
		return nil, fmt.Errorf("API Count should be between (%d, %d)",
			minAPI, maxAPI)
	}

	if c.ResponseTime < minDuration || c.ResponseTime > maxDuration {
		return nil, fmt.Errorf("c.ResponseTime should be between (%v, %v)",
			minDuration, maxDuration)
	}

	return &c, err
}
