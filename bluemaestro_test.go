package main

import (
	"fmt"
	"github.com/bettercap/gatt"
	"testing"
)

func TestNewBlueMaestro(t *testing.T) {
	cases := [][]byte{
		{51, 1, 1, 45, 3, 232, 254, 189, 0, 203, 1, 25, 1, 188, 1, 19, 1, 162, 1, 21, 1, 172, 0, 108, 2, 196, 145},
	}

	for _, tt := range cases {
		packet:=tt

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 307

		messages, err := NewBlueMaestro(adv)
		if err != nil {
			t.Errorf("%q", err)
		}
		fmt.Printf("%+v\n", messages)
	}
}
