package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; ; i++ {
		if i%2 == 1 {
			fmt.Printf("%d %s\n", -1, "connection refused")
		} else {
			fmt.Printf("%d\n", 0)
		}
		time.Sleep(3 * time.Second)
	}
}
