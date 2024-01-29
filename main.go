package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//	{
//		"_type": "location",
//		"acc": 2,
//		"alt": -20,
//		"batt": 100,
//		"bs": 3,
//		"conn": "m",
//		"created_at": 1706518008,
//		"lat": 0.00,
//		"lon": 0.00,
//		"m": 1,
//		"tid": "m0",
//		"tst": 1706518008,
//		"vac": 0,
//		"vel": 0
//	}

type Data struct {
	Type      string  `json:"_type"`
	Acc       int     `json:"acc"`
	Alt       int     `json:"alt"`
	Batt      int     `json:"batt"`
	Bs        int     `json:"bs"`
	Conn      string  `json:"conn"`
	CreatedAt int64   `json:"created_at"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	M         int     `json:"m"`
	Tid       string  `json:"tid"`
	Tst       int     `json:"tst"`
	Vac       int     `json:"vac"`
	Vel       int     `json:"vel"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if config.Debug {
		fmt.Printf("Received message: %q from topic: %s\n", msg.Payload(), msg.Topic())
	}

	data := &Data{}
	errUnmarshal := json.Unmarshal(msg.Payload(), data)
	if errUnmarshal != nil {
		fmt.Printf("errUnmarshal: %s for %q", errUnmarshal.Error(), msg.Payload())
	}

	// skip bad data
	if data.Type != "location" || data.Lat == 0.00 || data.Lon == 0.00 {
		return
	}

	haItem := HAItem{
		DevID: fmt.Sprintf("findmy_%s", strings.Replace(config.DeviceID, "-", "", -1)),
		Gps: [2]float64{
			data.Lat,
			data.Lon,
		},
		GpsAccuracy: float64(data.Acc),
		HostName:    config.DeviceID,
		Battery:     float64(data.Batt),
	}

	errProcessHomeassistant := processHomeassistant(haItem)
	if errProcessHomeassistant != nil {
		fmt.Printf("errProcessHomeassistant: %s for %q", errProcessHomeassistant.Error(), msg.Payload())
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	if config.Debug {
		fmt.Println("Connected")
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if config.Debug {
		fmt.Printf("Connect lost: %v\n", err)
	}
}

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	InitConfig()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.BrokerURL, config.BrokerPort))
	opts.SetClientID("external-mqtt-to-local")
	opts.SetUsername(config.BrokerUsername)
	opts.SetPassword(config.BrokerPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go sub(client)
	// go pub(client)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	if config.Debug {
		fmt.Println("awaiting signal")
	}

	<-done
	fmt.Println("exiting")

	client.Disconnect(250)
}

func sub(client mqtt.Client) {
	topic := config.BrokerTopic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	if config.Debug {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

// func pub(client mqtt.Client) {
// 	data := Data{
// 		Type:      "location",
// 		Batt:      50,
// 		CreatedAt: time.Now().Unix(),
// 		Lat:       0.00,
// 		Lon:       0.00,
// 	}

// 	result, errMarshal := json.Marshal(data)
// 	if errMarshal != nil {
// 		return
// 	}

// 	client.Publish(config.BrokerTopic, 1, false, result).Wait()
// }
