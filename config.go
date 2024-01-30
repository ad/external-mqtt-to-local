package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const ConfigFileName = "/data/options.json"

// Config ...
type Config struct {
	DeviceID string `json:"DEVICEID"`

	HomeassistantURL   string `json:"HOMEASSISTANTURL"`
	HomeassistantToken string `json:"HOMEASSISTANTTOKEN"`

	BrokerURL      string `json:"BROKERURL"`
	BrokerPort     int    `json:"BROKERPORT"`
	BrokerUsername string `json:"BROKERUSERNAME"`
	BrokerPassword string `json:"BROKERPASSWORD"`
	BrokerTopic    string `json:"BROKERTOPIC"`
	Debug          bool   `json:"DEBUG"`
}

var config = &Config{}

func InitConfig() {
	var initFromFile = false

	if _, err := os.Stat(ConfigFileName); err == nil {
		jsonFile, err := os.Open(ConfigFileName)
		if err == nil {
			byteValue, _ := io.ReadAll(jsonFile)
			if err = json.Unmarshal(byteValue, &config); err != nil {
				fmt.Printf("error on unmarshal config from file %s\n", err.Error())
			} else {
				initFromFile = true
			}
		}
	}

	if !initFromFile {
		flag.StringVar(&config.DeviceID, "DEVICEID", lookupEnvOrString("DEVICEID", config.DeviceID), "DEVICEID")

		flag.StringVar(&config.HomeassistantURL, "HOMEASSISTANTURL", lookupEnvOrString("HOMEASSISTANTURL", config.HomeassistantURL), "Homeassistant URL")
		flag.StringVar(&config.HomeassistantToken, "HomeassistantToken", lookupEnvOrString("HomeassistantToken", config.HomeassistantToken), "Homeassistant Token")

		flag.StringVar(&config.BrokerURL, "BROKERURL", lookupEnvOrString("BROKERURL", config.BrokerURL), "Broker URL")
		flag.IntVar(&config.BrokerPort, "BROKERPORT", lookupEnvOrInt("BROKERPORT", config.BrokerPort), "Broker Port")
		flag.StringVar(&config.BrokerUsername, "BROKERUSERNAME", lookupEnvOrString("BROKERUSERNAME", config.BrokerUsername), "Broker Username")
		flag.StringVar(&config.BrokerPassword, "BROKERPASSWORD", lookupEnvOrString("BROKERPASSWORD", config.BrokerPassword), "Broker Password")
		flag.StringVar(&config.BrokerTopic, "BROKERTOPIC", lookupEnvOrString("BROKERTOPIC", config.BrokerTopic), "Broker Topic")

		flag.BoolVar(&config.Debug, "DEBUG", lookupEnvOrBool("DEBUG", config.Debug), "Debug")

		flag.Parse()
	}

	if config.HomeassistantToken == "" {
		if val, ok := os.LookupEnv("SUPERVISOR_TOKEN"); ok {
			config.HomeassistantToken = val
			config.HomeassistantURL = "http://supervisor/core/api/services/device_tracker/see"
			if config.Debug {
				fmt.Printf("supervisor_token: %s\n", config.HomeassistantToken)
			}
		} else {
			if config.Debug {
				fmt.Println("supervisor token not set")
			}
		}
	}

	if config.HomeassistantURL == "" || config.HomeassistantToken == "" || config.BrokerURL == "" || config.BrokerPort == 0 {
		log.Fatal("provide config data")
	}
}

func lookupEnvOrString(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if x, err := strconv.Atoi(val); err == nil {
			return x
		}
	}
	return defaultVal
}

func lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if x, err := strconv.ParseBool(val); err == nil {
			return x
		}
	}
	return defaultVal
}
