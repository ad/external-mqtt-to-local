package homeassistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	conf "github.com/ad/external-mqtt-to-local/config"
)

type HASender struct {
	config *conf.Config
}

type HAItem struct {
	DevID       string     `json:"dev_id"`
	Gps         [2]float64 `json:"gps"`
	GpsAccuracy float64    `json:"gps_accuracy"`
	HostName    string     `json:"host_name"`
	Battery     float64    `json:"battery"`
	// LocationName string     `json:"location_name"`
}

func InitHASender(config *conf.Config) *HASender {
	haSender := &HASender{
		config: config,
	}

	return haSender
}

func (hs *HASender) ProcessHomeassistant(haItem *HAItem) error {
	jsonStr, errMarshal := json.Marshal(haItem)
	if errMarshal != nil {
		return errMarshal
	}

	url := hs.config.HomeassistantURL

	req, errNewRequest := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if errNewRequest != nil {
		return errNewRequest
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hs.config.HomeassistantToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return errDo
	}
	defer resp.Body.Close()

	if hs.config.Debug {
		fmt.Println("request body:", string(jsonStr))
	}

	if resp.StatusCode != 200 {
		if hs.config.Debug {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
		}

		body, _ := io.ReadAll(resp.Body)

		if hs.config.Debug {
			fmt.Println("response Body:", string(body))
		}
	}

	return nil
}
