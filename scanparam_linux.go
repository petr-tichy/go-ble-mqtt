package main

import (
	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/linux/cmd"
	"log"
)

func setScanParam(device gatt.Device) {
	scanParam := gatt.LnxSetScanParams(&cmd.LESetScanParameters{
		LEScanType:           0x00, // [0x00]: passive, 0x01: active
		LEScanInterval:       16,   // [0x10]: 0.625ms * 16
		LEScanWindow:         16,   // [0x10]: 0.625ms * 16
		OwnAddressType:       0x00, // [0x00]: public, 0x01: random
		ScanningFilterPolicy: 0x00, // [0x00]: accept all, 0x01: ignore non-white-listed.
	})
	if err := device.Option(scanParam); err != nil {
		log.Fatalf("setting device optins failed: %s", err)
		return
	}
}
