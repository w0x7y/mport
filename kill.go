package main

import (
	"fmt"
	"os"
	"strings"
)

func killProcess(pid int) {
	fmt.Printf("Are you sure you want to kill PID %d? (y/n): ", pid)

	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(answer) == "y" {
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = process.Kill()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Process killed.")
	} else {
		fmt.Println("Cancelled.")
	}
}
