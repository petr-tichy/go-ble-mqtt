package main

import "fmt"

func fmtDecimal(i int, s int) string {
	var j int
	if i < 0 {
		j = -i % s
	} else {
		j = i % s
	}
	return fmt.Sprintf("%d.%d", i/s, j)
}
