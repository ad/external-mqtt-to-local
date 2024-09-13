package listener

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	conf "github.com/ad/external-mqtt-to-local/config"
	"github.com/ad/external-mqtt-to-local/homeassistant"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Listener struct {
	config   *conf.Config
	haSender *homeassistant.HASender
	Client   mqtt.Client
}

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

func InitListener(config *conf.Config, haSender *homeassistant.HASender) (*Listener, error) {
	listener := &Listener{
		config:   config,
		haSender: haSender,
	}

	opts := mqtt.NewClientOptions()
	opts.SetAutoReconnect(true)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.BrokerURL, config.BrokerPort))
	opts.SetClientID("external-mqtt-to-local")
	opts.SetUsername(config.BrokerUsername)
	opts.SetPassword(config.BrokerPassword)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	listener.Client = client

	go listener.sub(client)

	go func() {
		for range time.Tick(time.Second * 30) {
			pub(config, client)
		}
	}()

	return listener, nil
}

func (l *Listener) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	if l.config.Debug {
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

	haItem := &homeassistant.HAItem{
		DevID: fmt.Sprintf("findmy_%s", strings.Replace(l.config.DeviceID, "-", "", -1)),
		Gps: [2]float64{
			data.Lat,
			data.Lon,
		},
		GpsAccuracy: float64(data.Acc),
		HostName:    l.config.DeviceID,
		Battery:     float64(data.Batt),
	}

	errProcessHomeassistant := l.haSender.ProcessHomeassistant(haItem)
	if errProcessHomeassistant != nil {
		fmt.Printf("errProcessHomeassistant: %s for %q", errProcessHomeassistant.Error(), msg.Payload())
	}
}

func (l *Listener) sub(client mqtt.Client) {
	topic := l.config.BrokerTopic
	token := client.Subscribe(topic, 1, l.messagePubHandler)
	token.Wait()
	if l.config.Debug {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

func (l *Listener) Disconnect() {
	if l.Client != nil {
		l.Client.Disconnect(250)
	}
}

func pub(config *conf.Config, client mqtt.Client) {
	data := Data{
		Type:      "ping",
		CreatedAt: time.Now().Unix(),
	}

	result, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return
	}

	client.Publish(config.BrokerTopic, 1, false, result).Wait()
}
