package main

import (
	"bufio"
	"fmt"
	"log"
	"net/netip"
	"os"
	"regexp"
	"sort"
)

// This should just be done with `cut -d ' ' -f 1 logs.txt | sort | uniq -c | sort -rn | head`

// Taken from the-power-of-go-tools/chapter07/memory/memory.go - see it for comments
// Regex from https://www.regular-expressions.info/ip.html
var ipAddressRegex = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

type frequency struct {
	address netip.Addr
	count   int
}

func main() {
	path := "logs.txt"
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open %q: %s", path, err)
	}
	defer file.Close()

	// Using netip.Addr instead of string so it handles IPs better.
	// e.g. "::1" and "0:0:0:0:0:0:0:1" are the same IPv6 address but different strings
	ipAddresses := make(map[netip.Addr]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		match := ipAddressRegex.FindString(line)
		if match == "" {
			continue
		}

		ip, err := netip.ParseAddr(match)
		if err != nil {
			log.Fatal(err)
		}

		ipAddresses[ip]++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sortedIpAddresses := make([]frequency, 0, len(ipAddresses))
	for ip, count := range ipAddresses {
		sortedIpAddresses = append(sortedIpAddresses, frequency{ip, count})
	}
	sort.Slice(sortedIpAddresses, func(i, j int) bool {
		return sortedIpAddresses[i].count > sortedIpAddresses[j].count
	})

	topN := 10
	fmt.Printf("%-16s%s\n", "Address", "Requests")
	for i, ip := range sortedIpAddresses {
		if i >= topN {
			break
		}
		fmt.Printf("%-16s%d\n", ip.address, ip.count)
	}
}
