package main

import (
	"encoding/binary"
	"errors"
	"github.com/bettercap/gatt"
)

type RuuviV2 Messages

var errNotRuuviV2 = errors.New("not an RuuviV2 tag")

func (r RuuviV2) temperature(d []byte) {
	r["temperature"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d)))*5, 3)
}

func (r RuuviV2) humidity(d []byte) {
	r["humidity"] = fmtDecimal(int(binary.BigEndian.Uint16(d))*25, 4)
}

func (r RuuviV2) pressure(d []byte) {
	r["pressure"] = fmtDecimal(int(binary.BigEndian.Uint16(d))+50000, 3)
}

func (r RuuviV2) accelerometer(d []byte, dir string) {
	r[dir] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 3)
}

func (r RuuviV2) battery(d []byte) {
	b := binary.BigEndian.Uint16(d) >> 5
	r["battery"] = fmtDecimal(int(b)+1600, 3)
}

func (r RuuviV2) txPower(d []byte) {
	b := binary.BigEndian.Uint16(d) & 0x1F
	r["tx_power"] = fmtDecimal(int(b)*2-40, 0)
}

func (r RuuviV2) movementCount(d byte) {
	r["movement_count"] = fmtDecimal(int(uint(d)), 0)
}

func NewRuuviV2(a *gatt.Advertisement) (Messages, error) {
	if !(len(a.ManufacturerData) == 26 && a.CompanyID == 0x0499 && a.ManufacturerData[2] == 0x05) {
		return nil, errNotRuuviV2
	}
	data := a.ManufacturerData[2:]

	r := make(RuuviV2)
	r.temperature(data[1:])
	r.humidity(data[3:])
	r.pressure(data[5:])
	r.accelerometer(data[7:], "dx")
	r.accelerometer(data[9:], "dy")
	r.accelerometer(data[11:], "dz")
	r.battery(data[13:])
	r.txPower(data[13:])
	r.movementCount(data[15])

	return Messages(r), nil
}
