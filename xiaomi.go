package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/bettercap/gatt"
	"log"
)

var errNotXiaomi = errors.New("not an Xiaomi tag")

var xiaomiUUUID = gatt.UUID16(0xfe95)
var xiaomiSig = []byte{0x50, 0x20, 0xaa, 0x01}

type Xiaomi Messages

func (x Xiaomi) temperature(d []byte) {
	x["temperature"] = fmtDecimal(int(int16(binary.LittleEndian.Uint16(d))), 1)
}

func (x Xiaomi) humidity(d []byte) {
	x["humidity"] = fmtDecimal(int(binary.LittleEndian.Uint16(d)), 1)
}

func (x Xiaomi) batteryLevel(d []byte) {

}

func NewXiaomi(a *gatt.Advertisement) (Xiaomi, error) {
	var data []byte
	for _, sd := range a.ServiceData {
		if sd.UUID.Equal(xiaomiUUUID) {
			data = sd.Data
		}
	}

	if len(data) == 0 || len(data) < 14 || bytes.Compare(data[0:4], xiaomiSig) != 0 {
		// log.Printf("Xiaomi not found or short data: %#v\n", a)
		return nil, errNotXiaomi
	}

	payloadLength := int(data[13])
	if len(data) != payloadLength+14 {
		// Data length doesn't match
		// log.Printf("Xiaomi Data length doesn't match: %#v\n", a)
		return nil, errNotXiaomi
	}

	/*
		mac_addr = aiobs.MACAddr(None)
		mac_addr.decode(data[5:11])
		mac_addr = mac_addr.val
		if mac_addr != packet.retrieve("peer")[0].val:
		# print("Packet source MAC doesn't match indicated MAC")
		return None
	*/

	payloadType := binary.LittleEndian.Uint16(data[11:])

	r := make(Xiaomi)

	if payloadType == 0x1004 && payloadLength == 2 {
		r.temperature(data[14:])
	} else if payloadType == 0x1006 && payloadLength == 2 {
		r.humidity(data[14:])
	} else if payloadType == 0x100a && payloadLength == 1 {
		r.batteryLevel(data[14:])
	} else if payloadType == 0x100d && payloadLength == 4 {
		r.temperature(data[14:])
		r.humidity(data[16:])
	} else {
		log.Printf("Xiaomi invalid: %#v\n", a)
		return nil, errNotXiaomi
	}
	r["message_number"] = fmtDecimal(int(data[4]), 0)
	return r, nil
}