package main

import (
	"encoding/hex"
	"github.com/bettercap/gatt"
	"reflect"
	"testing"
)

func TestNewBlueMaestro(t *testing.T) {
	cases := []struct {
		packet  string
		message Messages
	}{
		{"3301012d03e8febd00cb011901bc011301a2011501ac006c02c491",
			Messages(nil)},
	}

	for _, tt := range cases {
		packet, _ := hex.DecodeString(tt.packet)

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 307

		messages, err := NewBlueMaestro(adv)
		if err != nil {
			t.Errorf("%q", err)
		}

		/*
			fmt.Printf("%+#v", struct {
				packet   string
				messages Messages
			}{hex.EncodeToString(packet), messages})
			fmt.Printf("%+v\n", messages)
		*/

		if !reflect.DeepEqual(messages, tt.message) {
			t.Errorf("%q != %q", messages, tt.message)
		}

	}
}
