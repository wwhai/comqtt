package main

import (
	"fmt"
	"github.com/wind-c/comqtt/server/listeners/auth"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"

	mqtt "github.com/wind-c/comqtt/server"
	"github.com/wind-c/comqtt/server/listeners"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("TCP"))

	server := mqtt.New()
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}
	stats := listeners.NewHTTPStats("stats", ":8080")
	err = server.AddListener(stats, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println(aurora.BgMagenta("  Started!  "))

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))

}
