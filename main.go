package main

import (
	"fmt"
	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

const serverTopic = "state/servers/spacezero/status"
const serialPort = "/dev/serial0"

const noRSSI = -1 << 16

type mqttConfig struct {
	client MQTT.Client
	qos    int
}

// Messages type that the plugins should return
type Messages map[string]string

var plugins = []interface{}{
	NewRuuvi,
	NewRuuviV2,
	NewXiaomi,
	NewBlueMaestro,
	NewDZDSmart,
}

func (mqtt *mqttConfig) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	// log.Printf("Got packet %#v\n", a)
	id := p.ID()
	for _, f := range plugins {
		if m, err := f.(func(*gatt.Advertisement) (Messages, error))(a); err == nil {
			mqtt.publish(m, id, rssi)
			break
		}
	}

	/*
		if a.CompanyID == 76 {
			log.Printf("P: %#v, A: %#v", p.ID(), a)
			mqtt.publish(Messages{"state": "present"}, id, rssi)
		}
	*/
}

func (mqtt *mqttConfig) publish(m map[string]string, id string, rssi int) {
	var prefix = "state/sensors"
	prefix = prefix + "/" + strings.ToLower(id) + "/"
	for k, v := range m {
		mqtt.client.Publish(prefix+k, byte(mqtt.qos), false, v)
	}
	if rssi != noRSSI {
		mqtt.client.Publish(prefix+"rssi", byte(mqtt.qos), false, fmt.Sprintf("%d", rssi))
	}
}

func (mqtt *mqttConfig) initMQTT() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://192.168.8.10:1883")
	// opts.SetClientID(*id)
	opts.SetWill(serverTopic, "died", 0, true)
	opts.SetOrderMatters(false)
	mqtt.client = MQTT.NewClient(opts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("publisher started")
	mqtt.client.Publish(serverTopic, 0, true, nil)
}

func onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning for broadcasts...")
		d.Scan([]gatt.UUID{}, true)
		return
	default:
		d.StopScanning()
	}
}

func handleStop(mqtt *mqttConfig) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	ticker := time.NewTicker(10 * time.Second)

	run := true
	for run {
		select {
		case <-signals:
			ticker.Stop()
			mqtt.client.Publish(serverTopic, 0, true, "stopped")
			mqtt.client.Disconnect(10_000)
			run = false
		case <-ticker.C:
			mqtt.client.Publish(serverTopic, 0, false, "alive")
		}
	}
}

func main() {
	mqtt := mqttConfig{nil, 0}
	mqtt.initMQTT()
	clientOptions := option.DefaultClientOptions
	device, err := gatt.NewDevice(clientOptions...)
	if err != nil {
		log.Fatalf("failed to open device, err: %s\n", err)
		return
	}
	setScanParam(device)
	device.Handle(gatt.PeripheralDiscovered(mqtt.onPeripheralDiscovered))
	if device.Init(onStateChanged) != nil {
		log.Fatalf("device init failed: %s", err)
		return
	}
	go startSerialHandler(&mqtt)
	handleStop(&mqtt)
}
