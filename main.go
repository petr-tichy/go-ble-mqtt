package main

import (
	"fmt"
	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
	"github.com/bettercap/gatt/linux/cmd"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const serverTopic = "state/servers/spacezero/status"

type mqttConfig struct {
	client MQTT.Client
	qos    int
}

// Messages type that the plugins should return
type Messages map[string]string

func (mqtt *mqttConfig) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	id := p.ID()
	// log.Printf("Got packet %#v\n", a)
	if m, err := NewRuuvi(a); err == nil {
		mqtt.publish(m, id, rssi)
	} else if m, err := NewXiaomi(a); err == nil {
		mqtt.publish(m, id, rssi)
	} else if m, err := NewBlueMaestro(a); err == nil {
		mqtt.publish(m, id, rssi)
	} else if m, err := NewDZDSmart(a); err == nil {
		mqtt.publish(m, id, rssi)
	}
	/*
		else if a.CompanyID == 76 {
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
	mqtt.client.Publish(prefix+"rssi", byte(mqtt.qos), false, strconv.Itoa(rssi))
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

func main() {
	mqtt := mqttConfig{nil, 0}
	mqtt.initMQTT()
	clientOptions := option.DefaultClientOptions
	device, err := gatt.NewDevice(clientOptions...)
	if err != nil {
		log.Fatalf("failed to open device, err: %s\n", err)
		return
	}
	scanParam := gatt.LnxSetScanParams(&cmd.LESetScanParameters{
		LEScanType:           0x00, // [0x00]: passive, 0x01: active
		LEScanInterval:       16,   // [0x10]: 0.625ms * 16
		LEScanWindow:         16,   // [0x10]: 0.625ms * 16
		OwnAddressType:       0x00, // [0x00]: public, 0x01: random
		ScanningFilterPolicy: 0x00, // [0x00]: accept all, 0x01: ignore non-white-listed.
	})
	err = device.Option(scanParam)
	if err != nil {
		log.Fatalf("setting device optins failed: %s", err)
		return
	}

	device.Handle(gatt.PeripheralDiscovered(mqtt.onPeripheralDiscovered))
	err = device.Init(onStateChanged)
	if err != nil {
		log.Fatalf("device init failed: %s", err)
		return
	}

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
