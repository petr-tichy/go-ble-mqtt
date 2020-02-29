package main

import (
	"bufio"
	"log"
	"strconv"
	"time"

	"github.com/VividCortex/ewma"
	"github.com/tarm/serial"
)

func startSerialHandler(mqtt *mqttConfig) {
	id := "sonar"
	lt := time.Now()
	c := &serial.Config{Name: serialPort, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(s)

	movingAverage := ewma.NewMovingAverage()

	for {
		buf, err := reader.ReadString('\r')
		if err != nil {
			continue
		}
		if buf[0] != 'R' {
			continue
		}
		i, err := strconv.Atoi(buf[1 : len(buf)-1])
		if err != nil {
			continue
		}
		movingAverage.Add(float64(i))
		if time.Since(lt) > 30*time.Second {
			m := strconv.Itoa(int(movingAverage.Value() + 0.5))
			mqtt.publish(Messages{"range": m}, id, noRSSI)
			lt = time.Now()
		}
	}
}
