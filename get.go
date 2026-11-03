package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

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
