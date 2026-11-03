package main

import "fmt"

func printPort(port Port) {
	fmt.Printf("%s:%s - %s, %d\n", port.Address, port.Number, port.Process, port.PID)
}
