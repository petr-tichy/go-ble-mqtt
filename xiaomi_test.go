package main

/*
import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bettercap/gatt"
)

func TestNewXiaomi(t *testing.T) {
	cases := &gatt.Advertisement{LocalName:"MJ_HT_V1", Flags:0x6, CompanyID:0x0, Company:"", ManufacturerData:[]uint8(nil), ServiceData:[]gatt.ServiceData{{UUID:gatt.UUID{[]uint8{0x95, 0xfe}}, Data:[]uint8{0x50, 0x20, 0xaa, 0x1, 0x71, 0x28, 0x9, 0xdc, 0xa8, 0x65, 0x4c, 0xd, 0x10, 0x4, 0xf8, 0x0, 0xe3, 0x1}}, gatt.ServiceData{UUID:gatt.UUID{b:[]uint8{0xff, 0xff}}, Data:[]uint8{0xc0, 0x66, 0x64, 0x29, 0x3b, 0x11}}, gatt.ServiceData{UUID:gatt.UUID{b:[]uint8{0xff, 0xff}}, Data:[]uint8{0xc0, 0x66, 0x64, 0x29, 0x3b, 0x11}}}, Services:[]gatt.UUID{gatt.UUID{b:[]uint8{0xf, 0x18}}, gatt.UUID{b:[]uint8{0xa, 0x18}}, gatt.UUID{b:[]uint8{0xf, 0x18}}, gatt.UUID{b:[]uint8{0xa, 0x18}}}, OverflowService:[]gatt.UUID(nil), TxPowerLevel:0, Connectable:true, SolicitedService:[]gatt.UUID(nil), Raw:[]uint8{0xff, 0xff, 0xc0, 0x66, 0x64, 0x29, 0x3b, 0x11}}

	for _, tt := range cases {
		packet, _ := hex.DecodeString(tt)

		adv := new(gatt.Advertisement)
		adv.ManufacturerData = packet
		adv.Raw = packet
		adv.CompanyID = 0x0499

		ruuvi, err := NewXiaomi(adv)
		if err != nil {
			t.Errorf("%q", err)
		}
		fmt.Printf("%+v\n", ruuvi)
	}
}

*/
