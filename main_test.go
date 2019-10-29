package main

import (
	"fmt"
	"github.com/bettercap/gatt"
	"testing"
)

func TestParser(t *testing.T) {
	cases := [][]byte{
		{27, 0x16, 51, 1, 1, 45, 3, 232, 254, 189, 0, 203, 1, 25, 1, 188, 1, 19, 1, 162, 1, 21, 1, 172, 0, 108, 2, 196, 145}, // BlueMaestro 2nd
		{0xff, 0xff, 0xc0, 0x66, 0x64, 0x29, 0x3b, 0x11},
	}

	for _, tt := range cases {
		packet := tt
		fmt.Println(len(packet))
		m := new(MQTTConfig)
		p := new(gatt.Peripheral)

		adv := new(gatt.Advertisement)
		// adv.Unmarshall(packet)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 307

		m.onPeripheralDiscovered(*p, adv, -99)
	}

}
