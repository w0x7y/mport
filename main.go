package main

/*
  only work on linux for now
  support for windows and macos will come in the future
  (when i will want to do it)
*/

import (
	"fmt"
	"os"
	"strings"
)

type Port struct {
	Address string
	Number  string
	Process string
	PID     int
}

func main() {
	args := os.Args[1:]
	filterPorts := make(map[string]bool)

	mode := ""
	if len(args) > 0 && (args[0] == "--info" || args[0] == "-i" || args[0] == "--kill" || args[0] == "-k") {
		mode = args[0]
		args = args[1:]

		if len(args) == 0 {
			fmt.Printf("Usage: mport %s <port>\n", mode)
			os.Exit(1)
		}
	}

	for _, arg := range args {
		filterPorts[arg] = true
	}

	ports, err := getPorts(filterPorts)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if mode != "" && len(ports) == 0 {
		fmt.Println("Nothing is listening on that port.")
		os.Exit(1)
	}

	// A process listening on IPv4 and IPv6 shows up twice; handle each PID once.
	seen := make(map[int]bool)

	for _, port := range ports {
		printPort(port)

		if mode == "" || seen[port.PID] {
			continue
		}
		seen[port.PID] = true

		// PID 0 means ss could not see the owner; killing it would signal our own process group.
		if port.PID == 0 {
			fmt.Println("  Owner not visible, try running with sudo.")
			continue
		}

		switch mode {
		case "--info", "-i":
			cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", port.PID))
			if err == nil {
				fmt.Println("  Command:", strings.TrimSpace(strings.ReplaceAll(string(cmdline), "\x00", " ")))
			}
		case "--kill", "-k":
			killProcess(port.PID)
		}
	}
}
