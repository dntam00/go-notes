package main

import (
	"fmt"
	"strconv"
	"strings"
)

func calculateBytes(addr1, addr2 string) (int64, error) {
	// Remove "0x" prefix if present
	if len(addr1) > 2 && addr1[:2] == "0x" {
		addr1 = addr1[2:]
	}
	if len(addr2) > 2 && addr2[:2] == "0x" {
		addr2 = addr2[2:]
	}

	// Parse hex addresses to integers
	val1, err := strconv.ParseInt(addr1, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hex address 1: %v", err)
	}

	val2, err := strconv.ParseInt(addr2, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hex address 2: %v", err)
	}

	// Calculate absolute difference
	diff := val2 - val1
	if diff < 0 {
		diff = -diff
	}

	return diff, nil
}

// cat /proc/1/maps | head -n 200 | awk '{split($1, addr, "-"); printf "%s,%s\n", addr[1], addr[2]}'

func main() {
	input := `00400000,00401000
			  00600000,00601000`

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			fmt.Printf("Invalid line format: %s\n", line)
			continue
		}

		addr1 := strings.TrimSpace(parts[0])
		addr2 := strings.TrimSpace(parts[1])

		bytes, err := calculateBytes(addr1, addr2)

		if err != nil {
			fmt.Printf("Error processing %s,%s: %v\n", addr1, addr2, err)
			continue
		}

		if bytes < (1024 * 1024) {
			continue
		}

		fmt.Printf("%s,%s: %d bytes, (%d KiB), (%d MB)\n", addr1, addr2, bytes, bytes/1024, bytes/1024/1024)
	}
}
