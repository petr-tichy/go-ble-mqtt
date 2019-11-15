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
)

var errNotBlueMaestro = errors.New("not an BlueMaestro tag")

// BlueMaestro is a clone of Messages Type
type BlueMaestro Messages

func (b BlueMaestro) temperature(d []byte) {
	b["temperature"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 1)
}

func (b BlueMaestro) humidity(d []byte) {
	b["humidity"] = fmtDecimal(int(binary.BigEndian.Uint16(d)), 1)
}

func (b BlueMaestro) dewPoint(d []byte) {
	b["dew_point"] = fmtDecimal(int(int16(binary.BigEndian.Uint16(d))), 1)
}

func (b BlueMaestro) batteryLevel(d []byte) {
	b["battery_level"] = fmtDecimal(int(d[0]), 0)
}

func NewBlueMaestro(a *gatt.Advertisement) (BlueMaestro, error) {
	if a.Company != "Blue Maestro Limited" {
		return nil, errNotBlueMaestro
	}

	data := a.ManufacturerData

	if len(data) != 16 || data[2] != 0x17 {
		//log.Printf("BlueMaestro not found or short data: %v, %v\n", len(data), data[2])
		return nil, errNotBlueMaestro
	}

	r := make(BlueMaestro)
	r.batteryLevel(data[3:4])
	r.temperature(data[8:10])
	r.humidity(data[10:12])
	r.dewPoint(data[12:14])

	return r, nil
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
