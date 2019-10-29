package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bettercap/gatt"
)

type Ruuvi Messages

var errNotRuuvi = errors.New("Not an Ruuvi tag")

func getTemp(i, f byte) string {
	integer := int8(i & 0x7f)
	if int8(i) < 0 {
		integer = -integer
	}
	return fmt.Sprintf("%d.%d", integer, f)
}

func NewRuuvi(a *gatt.Advertisement) (Ruuvi, error) {
	if !(len(a.ManufacturerData) == 16 && a.CompanyID == 0x0499 && a.ManufacturerData[2] == 0x03) {
		return nil, errNotRuuvi
	}
	data := a.ManufacturerData[2:]

	r := make(Ruuvi)
	r["humidity"] = fmtDecimal(int(uint(data[1])*5), 10)
	r["temperature"] = getTemp(data[2], data[3])
	r["pressure"] = fmtDecimal(int(binary.BigEndian.Uint16(data[4:6]))+50000, 1000)
	r["dx"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(data[6:8]))), 1000)
	r["dy"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(data[8:10]))), 1000)
	r["dz"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(data[10:12]))), 1000)
	r["battery"] = fmtDecimal(int(binary.BigEndian.Uint16(data[12:14])), 1000)

	/*
		length=sqrt(dx**2 + dy**2 + dz**2)
	*/
	// log.Print("Ruuvi OK")
	return r, nil
}
