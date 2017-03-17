package main

import (
	"fmt"
	"time"
)

var (
	work       = time.Minute * 25
	shortBreak = time.Second * 10
	longBreak  = time.Second * 30
)

func main() {

	// hh, mm := timeconv(4)
	// fmt.Printf("%02d:%02d", hh, mm)
	fmt.Println(work + work)

}

// func timeconv(pomo int) (int, int) {
// 	// min := pomo * (work / (time.Minute * 1))
// 	// hh := min / 60
// 	// mm := min % 60
// 	// return hh, mm
// }
