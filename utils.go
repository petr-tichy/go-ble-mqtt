package main

import (
	"fmt"
)

func fmtDecimal(i, s int) string {
	l := s + 1
	if i < 0 {
		l++
	}
	str := fmt.Sprintf("%0*d", l, i)
	ret := str[:len(str)-s]
	if s > 0 {
		ret += "." + str[len(str)-s:]
	}
	return ret
}
