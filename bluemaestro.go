/*
&{LocalName:TH M Flags:BR/EDR Not Supported
CompanyID:307
Company:Blue Maestro Limited
ManufacturerData:[51 1 1 45 3 232 254 189 0 203 1 25 1 188 1 19 1 162 1 21 1 172 0 108 2 196 145]
ServiceData:[] Services:[] OverflowService:[] TxPowerLevel:0 Connectable:true SolicitedService:[]
Raw:[51 1 1 45 3 232 254 189 0 203 1 25 1 188 1 19 1 162 1 21 1 172 0 108 2 196 145]}
*/

package main

import (
	"encoding/binary"
	"errors"
	"github.com/bettercap/gatt"
	"log"
)

var errNotBlueMaestro = errors.New("Not an BlueMaestro tag")

type BlueMaestro Messages

func (x BlueMaestro) setTemperature(d []byte) {
	x["temperature"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 10)
}

func (x BlueMaestro) setHumidity(d []byte) {
	x["humidity"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 10)
}

func (x BlueMaestro) setDewPoint(d []byte) {
	x["dew_point"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 10)
}

func (x BlueMaestro) setBatteryLevel(d []byte) {
	x["battery_level"] = fmtDecimal(int(d[0]), 1)
}

/*
        data = {}
        raw_data = packet.retrieve('Payload for mfg_specific_data')
        if raw_data:
            pckt = raw_data[0].val
            mfg_id = unpack('<H', pckt[:2])[0]
            if mfg_id == BLUEMAESTRO:
                data['version'] = unpack('<B', pckt[2:3])[0]
                data['batt_lvl'] = unpack('<B', pckt[3:4])[0]
                data['logging'] = unpack('>H', pckt[4:6])[0]
                data['interval'] = unpack('>H', pckt[6:8])[0]
                data['temperature'] = unpack('>h', pckt[8:10])[0]/10
                data['humidity'] = unpack('>h', pckt[10:12])[0]/10
                data['pressure'] = unpack('>h', pckt[12:14])[0]/10
        return data


*/

func NewBlueMaestro(a *gatt.Advertisement) (BlueMaestro, error) {
	if a.CompanyID != 0x0133 {
		return nil, errNotBlueMaestro
	}
	
	data := a.ManufacturerData
	
	if len(data) != 30 || data[2] != 0x17 {
		log.Println("BlueMaestro not found or short data")
		return nil, errNotBlueMaestro
	}
	
	r := make(BlueMaestro)
	r.setBatteryLevel(data[3:4])
	r.setTemperature(data[8:10])
	r.setHumidity(data[10:12])
	r.setDewPoint(data[12:14])

	log.Print("BlueMaestro OK")
	return r, nil
}