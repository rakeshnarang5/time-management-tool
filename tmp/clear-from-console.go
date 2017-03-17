package main

import (
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		print("akljdfkljajkdfl")
		print(i)
		time.Sleep(time.Second * 1)
		print("\033[H\033[2J")
	}

}
