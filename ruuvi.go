package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bettercap/gatt"
)

type Ruuvi Messages

var errNotRuuvi = errors.New("not an Ruuvi tag")

func getTemp(i, f byte) string {
	integer := int8(i & 0x7f)
	if int8(i) < 0 {
		integer = -integer
	}
	return fmt.Sprintf("%d.%02d", integer, f)
}

func (r Ruuvi) temperature(d []byte) {
	r["temperature"] = getTemp(d[0], d[1])
}

func (r Ruuvi) humidity(d byte) {
	r["humidity"] = fmtDecimal(int(uint(d)*5), 1)
}

func (r Ruuvi) pressure(d []byte) {
	r["pressure"] = fmtDecimal(int(binary.BigEndian.Uint16(d))+50000, 2)
}

func (r Ruuvi) accelerometer(d []byte, dir string) {
	r[dir] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 3)
}

func (r Ruuvi) battery(d []byte) {
	r["battery"] = fmtDecimal(int(binary.BigEndian.Uint16(d)), 3)
}

func NewRuuvi(a *gatt.Advertisement) (Messages, error) {
	if !(len(a.ManufacturerData) == 16 && a.CompanyID == 0x0499 && a.ManufacturerData[2] == 0x03) {
		return nil, errNotRuuvi
	}
	data := a.ManufacturerData[2:]

	r := make(Ruuvi)
	r.humidity(data[1])
	r.temperature(data[2:4])
	r.pressure(data[4:6])
	r.accelerometer(data[6:8], "dx")
	r.accelerometer(data[8:10], "dy")
	r.accelerometer(data[10:12], "dz")
	r.battery(data[12:14])

	/*
		length=sqrt(dx**2 + dy**2 + dz**2)
	*/

	return Messages(r), nil
}
