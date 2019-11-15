package main

import (
	"errors"
	"github.com/bettercap/gatt"
)

type DZD Messages

var errNotDZDSmart = errors.New("not DZD Smart Boiler")

var uuidMain = gatt.UUID16(0x2b98)
var uuidSec = gatt.UUID16(0x180a)

func NewDZDSmart(a *gatt.Advertisement) (DZD, error) {
	n := 0
	for i := range a.Services {
		if a.Services[i].Equal(uuidMain) || a.Services[i].Equal(uuidSec) {
			n++
		}
	}

	if n != 2 {
		return nil, errNotDZDSmart
	}

	return DZD{"state": "present"}, nil
}
