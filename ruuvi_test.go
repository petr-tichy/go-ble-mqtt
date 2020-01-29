package main

import (
	"encoding/hex"
	"github.com/bettercap/gatt"
	"reflect"
	"testing"
)

func TestNewRuuvi(t *testing.T) {
	cases := []struct {
		packet  string
		message Messages
	}{
		{"03291A1ECE1EFC18F94202CA0B53",
			Messages{"battery": "2.899", "dx": "-1.000", "dy": "-1.726", "dz": "0.714", "humidity": "20.5", "pressure": "1027.66", "temperature": "26.30"}},
		{"03FF7F63FFFF7FFF7FFF7FFFFFFF",
			Messages{"battery": "65.535", "dx": "32.767", "dy": "32.767", "dz": "32.767", "humidity": "127.5", "pressure": "1155.35", "temperature": "127.99"}},
		{"0300FF6300008000800080000000",
			Messages{"battery": "0.000", "dx": "-32.768", "dy": "-32.768", "dz": "-32.768", "humidity": "0.0", "pressure": "500.00", "temperature": "-127.99"}},
		{"0300FF6300008001800180010000",
			Messages{"battery": "0.000", "dx": "-32.767", "dy": "-32.767", "dz": "-32.767", "humidity": "0.0", "pressure": "500.00", "temperature": "-127.99"}},
		{"0300FF630000FFFFFFFFFFFF0000",
			Messages{"battery": "0.000", "dx": "-0.001", "dy": "-0.001", "dz": "-0.001", "humidity": "0.0", "pressure": "500.00", "temperature": "-127.99"}},
	}

	for _, tt := range cases {
		packet, _ := hex.DecodeString("9904" + tt.packet)

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 0x0499

		ruuvi, err := NewRuuvi(adv)
		if err != nil {
			t.Errorf("%q", err)
		}
		if !reflect.DeepEqual(ruuvi, tt.message) {
			t.Errorf("%q != %q", ruuvi, tt.message)
		}
	}
}
