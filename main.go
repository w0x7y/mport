package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

type Port struct {
	Address string
	Number  string
	Process string
	PID     int
}

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

// getPorts lists listening TCP ports, keeping only filterPorts when it is non-empty.
func getPorts(filterPorts map[string]bool) ([]Port, error) {
	output, err := exec.Command("ss", "-ltnp").Output()
	if err != nil {
		return nil, err
	}

	var ports []Port

	// Match process name and PID
	re := regexp.MustCompile(`\(\("([^"]+)",pid=(\d+)`)

	for line := range strings.SplitSeq(string(output), "\n") {
		fields := strings.Fields(line)

		// Skip header and invalid lines
		if len(fields) < 5 || fields[0] != "LISTEN" {
			continue
		}

		// Extract local address and port
		localAddress := fields[3]
		addressParts := strings.Split(localAddress, ":")

		if len(addressParts) < 2 {
			continue
		}

		address := strings.Join(addressParts[:len(addressParts)-1], ":")
		number := addressParts[len(addressParts)-1]

		// Filter by port number
		if len(filterPorts) > 0 && !filterPorts[number] {
			continue
		}

		// Extract process information
		process := "unknown"
		pid := 0

		if match := re.FindStringSubmatch(line); match != nil {
			process = match[1]
			fmt.Sscan(match[2], &pid)
		}

		ports = append(ports, Port{
			Address: address,
			Number:  number,
			Process: process,
			PID:     pid,
		})
	}

	// Named processes first
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Process != "unknown" &&
			ports[j].Process == "unknown"
	})

	return ports, nil
}

func printPort(port Port) {
	fmt.Printf("%s:%s - %s, %d\n", port.Address, port.Number, port.Process, port.PID)
}

func main() {
	args := os.Args[1:]
	filterPorts := make(map[string]bool)

	mode := ""
	if len(args) > 0 && (args[0] == "--info" || args[0] == "--kill") {
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
		case "--info":
			cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", port.PID))
			if err == nil {
				fmt.Println("  Command:", strings.TrimSpace(strings.ReplaceAll(string(cmdline), "\x00", " ")))
			}
		case "--kill":
			killProcess(port.PID)
		}
	}
}
