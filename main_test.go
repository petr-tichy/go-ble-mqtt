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
		m := new(mqttConfig)
		p := new(gatt.Peripheral)

		adv := new(gatt.Advertisement)
		// adv.Unmarshall(packet)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 307

		m.onPeripheralDiscovered(*p, adv, -99)
	}

}

func TestFmtDecimal(t *testing.T) {
	cases := []struct {
		i int
		s int
		r string
	}{{10001, 3, "10.001"},
		{1001, 3, "1.001"},
		{101, 2, "1.01"},
		{11, 1, "1.1"},
		{709, 2, "7.09"},
		{-709, 2, "-7.09"},
		{-7009, 3, "-7.009"},
		{-9, 3, "-0.009"},
		{-11, 2, "-0.11"},
		{-1, 3, "-0.001"},
		{1, 3, "0.001"},
		{99, 2, "0.99"},
		{-1000001, 3, "-1000.001"},
		{0, 3, "0.000"},
		{0, 0, "0"},
		{1000, 0, "1000"},
		//{0, 0, "0"},

	}
	for _, tt := range cases {
		if c := fmtDecimal(tt.i, tt.s); c != tt.r {
			t.Errorf("fmtDecimal(%d, %d) got %s want %s", tt.i, tt.s, c, tt.r)
		}

	}
}

func TestGetTemp(t *testing.T) {
	cases := []struct {
		i, f byte
		r    string
	}{
		{0, 0x01, "0.01"},
		{0x01, 0x01, "1.01"},
		{0, 1, "0.01"},
		{0x80 | 0x01, 1, "-1.01"},
		{0x80 | 0x01, 99, "-1.99"},
		{0x80 | 0x7f, 99, "-127.99"},
		{0x00, 00, "0.00"},
		{0x80 | 0x00, 00, "0.00"},
	}
	for _, tt := range cases {
		if c := getTemp(tt.i, tt.f); c != tt.r {
			t.Errorf("getTemp(%q, %q) got %s want %s", tt.i, tt.f, c, tt.r)
		}
	}
}
