package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIP(ip string) []int {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		fmt.Println("Invalid IP")
		return nil
	}
	ipArr := make([]int, 4)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || n > 255 {
			fmt.Println("Invalid IP part:", p)
			return nil
		}
		ipArr[i] = n
	}
	return ipArr
}

func parseMask(mask string) []int {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		fmt.Println("Invalid Mask")
		return nil
	}
	maskArr := make([]int, 4)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || n > 255 {
			fmt.Println("Invalid Mask part:", p)
			return nil
		}
		maskArr[i] = n
	}
	return maskArr
}

func prefixMask(maskArr []int) int {
	prefix := 0
	foundZero := false
	for _, n := range maskArr {
		for i := 7; i >= 0; i-- {
			if n&(1<<i) != 0 {
				if foundZero {
					return -1
				}
				prefix++
			} else {
				foundZero = true
			}
		}
	}
	return prefix
}

func numberOfHosts(prefix int) int {
	return 1<<(32-prefix) - 2
}

func ipToUint32(ip []int) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) []int {
	return []int{
		int((n >> 24) & 0xFF),
		int((n >> 16) & 0xFF),
		int((n >> 8) & 0xFF),
		int(n & 0xFF),
	}
}

func networkAddress(ipArr, maskArr []int) []int {
	ip := ipToUint32(ipArr)
	mask := ipToUint32(maskArr)
	return uint32ToIP(ip & mask)
}

func broadcastAddress(networkArr, maskArr []int) []int {
	net := ipToUint32(networkArr)
	mask := ipToUint32(maskArr)
	return uint32ToIP(net | (^mask))
}

func rangeAddress(networkArr, maskArr []int) (string, string) {
	net := ipToUint32(networkArr)
	mask := ipToUint32(maskArr)
	first := net + 1
	last := (net | (^mask)) - 1
	startIP := uint32ToIP(first)
	endIP := uint32ToIP(last)
	return fmt.Sprintf("%d.%d.%d.%d", startIP[0], startIP[1], startIP[2], startIP[3]),
		fmt.Sprintf("%d.%d.%d.%d", endIP[0], endIP[1], endIP[2], endIP[3])
}

func generateAllSubnets(networkArr, maskArr []int) [][]int {
	var subnets [][]int
	prefix := prefixMask(maskArr)
	if prefix == -1 {
		return subnets
	}
	subnetSize := uint32(1 << (32 - prefix))
	start := ipToUint32(networkArr)

	for n := start; n <= 0xFFFFFFFF; n += subnetSize {
		subnets = append(subnets, uint32ToIP(n))
		if n+subnetSize > start+subnetSize*16 {
			break
		}
	}
	return subnets
}

// --- Main ---
func main() {
	var ipStr, maskStr string
	fmt.Print("Enter IP address: ")
	fmt.Scanln(&ipStr)
	fmt.Print("Enter subnet mask: ")
	fmt.Scanln(&maskStr)

	ipArr := parseIP(ipStr)
	maskArr := parseMask(maskStr)
	if ipArr == nil || maskArr == nil {
		return
	}

	prefix := prefixMask(maskArr)
	fmt.Printf("\nPrefix Mask: /%d\n", prefix)
	fmt.Printf("Number of Hosts: %d\n", numberOfHosts(prefix))

	network := networkAddress(ipArr, maskArr)
	fmt.Printf("Network Address: %v\n", network)

	broadcast := broadcastAddress(network, maskArr)
	fmt.Printf("Broadcast Address: %v\n", broadcast)

	startIP, endIP := rangeAddress(network, maskArr)
	fmt.Printf("Host Range: %s - %s\n", startIP, endIP)

	fmt.Println("\nSubnets in this network:")
	subnets := generateAllSubnets(network, maskArr)
	for i, s := range subnets {
		subStart, subEnd := rangeAddress(s, maskArr)
		fmt.Printf("Subnet %d: %d.%d.%d.%d/%d, Hosts: %s - %s\n",
			i+1, s[0], s[1], s[2], s[3], prefix, subStart, subEnd)
	}
}
