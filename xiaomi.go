package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/bettercap/gatt"
	"log"
)

var errNotXiaomi = errors.New("Not an Xiaomi tag")

var XIAOMI_UUID = gatt.UUID16(0xfe95)
var XIAOMI_SIG = []byte{0x50, 0x20, 0xaa, 0x01}

type Xiaomi Messages

func (x Xiaomi) setTemperature(d []byte) {
	x["temperature"] = fmtDecimal(int(int16(binary.LittleEndian.Uint16(d))), 10)
}

func (x Xiaomi) setHumidity(d []byte) {
	x["humidity"] = fmtDecimal(int(binary.LittleEndian.Uint16(d)), 10)
}

func (x Xiaomi) setBatteryLevel(d []byte) {
	x["battery_level"] = fmtDecimal(int(d[0]), 1)
}

//noinspection GoNilness
func NewXiaomi(a *gatt.Advertisement) (Xiaomi, error) {
	var data []byte
	for _, x := range a.ServiceData {
		if x.UUID.Equal(XIAOMI_UUID) {
			data = x.Data
		}
	}

	if len(data) > 0 || len(data) <= 13 {
		log.Println("Xiaomi not found or short data")
		return nil, errNotXiaomi
	}

	if bytes.Compare(data[0:4], XIAOMI_SIG) != 0 {
		log.Printf("Xiaomi wrong signature" )
		return nil, errNotXiaomi
	}

	payloadLength := int(data[13])
	if payloadLength != len(data)-14 {
		// Data length doesn't match
		log.Printf("Xiaomi Data length doesn't match: %#v\n", a)
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
	payloadType := data[11]

	r := make(Xiaomi)
	// r["message_number"] = fmt.Sprintf("%d", data[4])

	if payloadType == 0x0d && payloadLength == 4 {
		// Temperature & Humidity
		r.setTemperature(data[14:16])
		r.setHumidity(data[16:18])
	} else if payloadType == 0x0a && payloadLength == 1 {
		// Battery level
		r.setBatteryLevel(data[14:15])
	} else if payloadType == 0x06 && payloadLength == 2 {
		// Humidity
		r.setHumidity(data[14:16])
	} else if payloadType == 0x04 && payloadLength == 2 {
		// Temperature
		r.setTemperature(data[14:16])
	} else {
		log.Printf("Xiaomi invalid: %#v\n", a)
		return nil, errNotXiaomi
	}
	log.Print("Xiaomi OK")
	return r, nil
}

/*
   if not found:
       return None

   adv_payload = found.retrieve("Adv Payload")[0]

   if not adv_payload or len(adv_payload) <= 13:
       # Packet too short
       #print('Adv Payload too short')
       return None

   val = adv_payload.val
   if not (data[0]==0x50 and data[1]==0x20 and data[2]==0xaa and data[3]==0x01):
       return None

   payload_length = data[13]
   if payload_length != len(val) - 14:
       # Data length doesn't match
       # print('Invalid payload length %s vs %s' % (payload_length, len(val) - 14))
       return None

   # result["message_number"] = data[4]

   mac_addr = aiobs.MACAddr(None)
   mac_addr.decode(data[5:11])
   mac_addr = mac_addr.val
   if mac_addr != packet.retrieve("peer")[0].val:
       # print("Packet source MAC doesn't match indicated MAC")
       return None

   state = self.devices[mac_addr]

   payload_type = data[11]
   if payload_type == 0x0d and payload_length == 4:
       # Temperature & Humidity
       state.temperature = data[14:16]
       state.humidity = data[16:18]
   elif payload_type == 0x0a and payload_length == 1:
       # Battery level
       state.battery_level = data[14:15]
   elif payload_type == 0x06 and payload_length == 2:
       # Humidity
       state.humidity = data[14:16]
   elif payload_type == 0x04 and payload_length == 2:
       # Temperature
       state.temperature = data[14:16]
   else:
       # print('Invalid packet type %s or length' % payload_type)
       # adv_payload.show()
       return None

   if state.get('changed') is not None:
       return state

   return None

*/
