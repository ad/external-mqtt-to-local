package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HAItem struct {
	DevID       string     `json:"dev_id"`
	Gps         [2]float64 `json:"gps"`
	GpsAccuracy float64    `json:"gps_accuracy"`
	HostName    string     `json:"host_name"`
	Battery     float64    `json:"battery"`
	// LocationName string     `json:"location_name"`
}

func processHomeassistant(haItem HAItem) error {
	jsonStr, errMarshal := json.Marshal(haItem)
	if errMarshal != nil {
		return errMarshal
	}

	url := config.HomeassistantURL

	req, errNewRequest := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if errNewRequest != nil {
		return errNewRequest
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.HomeassistantToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return errDo
	}
	defer resp.Body.Close()

	fmt.Println("request body:", string(jsonStr))
	if resp.StatusCode != 200 {
		if config.Debug {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
		}

		body, _ := io.ReadAll(resp.Body)

		if config.Debug {
			fmt.Println("response Body:", string(body))
		}
	}

	return nil
}
