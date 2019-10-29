package main

import (
	"encoding/hex"
	"fmt"
	"github.com/bettercap/gatt"
	"testing"
)

func TestNewRuuvi(t *testing.T) {
	cases := []string{
		"03291A1ECE1EFC18F94202CA0B53",
		"03FF7F63FFFF7FFF7FFF7FFFFFFF",
		"0300FF6300008000800080000000",
		"0300FF6300008001800180010000",
		"0300FF630000FFFFFFFFFFFF0000",
	}

	for _, tt := range cases {
		packet, _ := hex.DecodeString("9904" + tt)

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 0x0499

		ruuvi, err := NewRuuvi(adv)
		if err != nil {
			t.Errorf("%q", err)
		}
		fmt.Printf("%#v: %+v\n", packet, ruuvi)
	}
}
