package main

import (
	"encoding/hex"
	"github.com/bettercap/gatt"
	"reflect"
	"testing"
)

// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md
func TestNewRuuviV2(t *testing.T) {
	cases := []struct {
		packet  string
		message Messages
	}{
		{ // sample
			"0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F",
			Messages{"battery": "2.977", "dx": "0.004", "dy": "-0.004", "dz": "1.036", "humidity": "53.4900", "movement_count": "66", "pressure": "100.044", "temperature": "24.300", "tx_power": "4"},
		}, { // max
			"057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F",
			Messages{"battery": "3.646", "dx": "32.767", "dy": "32.767", "dz": "32.767", "humidity": "163.8350", "movement_count": "254", "pressure": "115.534", "temperature": "163.835", "tx_power": "20"},
		}, { // min
			"058001000000008001800180010000000000CBB8334C884F",
			Messages{"battery": "1.600", "dx": "-32.767", "dy": "-32.767", "dz": "-32.767", "humidity": "0.0000", "movement_count": "0", "pressure": "50.000", "temperature": "-163.835", "tx_power": "-40"},
		}, { // undef
			"058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF",
			Messages{"battery": "3.647", "dx": "-32.768", "dy": "-32.768", "dz": "-32.768", "humidity": "163.8375", "movement_count": "255", "pressure": "115.535", "temperature": "-163.840", "tx_power": "22"},
		},
	}

	for _, tt := range cases {
		packet, err := hex.DecodeString("9904" + tt.packet)
		if err != nil {
			t.Error("Parsing err")
		}

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 0x0499

		ruuvi, err := NewRuuviV2(adv)
		if err != nil {
			t.Errorf("%q", err)
		}

		if !reflect.DeepEqual(ruuvi, tt.message) {
			t.Errorf("%q != %q", ruuvi, tt.message)
		}
	}
}
