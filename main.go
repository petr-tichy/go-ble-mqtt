package main

import (
	"fmt"
	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
)

type MQTTConfig struct {
	client MQTT.Client
	qos    int
}

type Messages map[string]string

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for iBeacon Broadcasts...")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func (mqtt *MQTTConfig) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	rssi_s := strconv.Itoa(rssi)
	id := p.ID()
	// log.Printf("Got packet %#v\n", a)
	if m, err := NewRuuvi(a); err == nil {
		mqtt.publish(m, "state/sensors", id, rssi_s)
	} else if m, err := NewXiaomi(a); err == nil {
		mqtt.publish(m, "", id, rssi_s)
	} else if m, err := NewBlueMaestro(a); err == nil {
		mqtt.publish(m, "", id, rssi_s)
	} else if !(a.CompanyID != 76 || a.LocalName != "SpaceBoiler") {
		log.Printf("Last: %#v\n", a)
	}
}

func (mqtt *MQTTConfig) publish(m map[string]string, prefix, id, rssi string) {
	if len(prefix) == 0 {
		prefix = "test/state/sensors"
	}
	prefix = prefix + "/" + strings.ToLower(id) + "/"
	for k, v := range m {
		mqtt.client.Publish(prefix+k, byte(mqtt.qos), false, v)
	}
	mqtt.client.Publish(prefix+"rssi", byte(mqtt.qos), false, rssi)
}

func (mqtt *MQTTConfig) initMQTT() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://192.168.8.10:1883")
	// opts.SetClientID(*id)
	// opts.SetUsername(*user)
	// opts.SetPassword(*password)
	// opts.SetCleanSession(*cleansess)
	opts.SetWill("state/servers/spacezero/will", "died", 0, false)
	// opts.SetConnectRetry(true)
	opts.SetOrderMatters(false)
	mqtt.client = MQTT.NewClient(opts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Publisher Started")
}

func main() {
	mqtt := MQTTConfig{nil, 0}
	mqtt.initMQTT()
	clientOptions := option.DefaultClientOptions
	// clientOptions.Set
	device, err := gatt.NewDevice(clientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(mqtt.onPeripheralDiscovered))
	_ = device.Init(onStateChanged)
	select {}
}
