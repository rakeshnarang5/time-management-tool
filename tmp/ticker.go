package main

import (
	"fmt"
	"time"
)

func main() {
	d := time.Minute * 30
	//fmt.Println(d)
	for d > time.Second*0 {
		fmt.Println(d)
		d -= time.Second * 1
		time.Sleep(time.Second * 1)
	}

}
