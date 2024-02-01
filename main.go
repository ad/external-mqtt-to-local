package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	conf "github.com/ad/external-mqtt-to-local/config"
	"github.com/ad/external-mqtt-to-local/homeassistant"
	lstnr "github.com/ad/external-mqtt-to-local/listener"
)

var (
	config   *conf.Config
	haSender *homeassistant.HASender
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	confLoad, errInitConfig := conf.InitConfig()
	if errInitConfig != nil {
		log.Fatal(errInitConfig)
	}

	config = confLoad

	haSender = homeassistant.InitHASender(config)

	listener, errInitListener := lstnr.InitListener(config, haSender)
	if errInitListener != nil {
		log.Fatal(errInitListener)
	}

	defer listener.Disconnect()

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
}
