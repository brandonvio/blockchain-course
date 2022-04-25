package globals

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var IpPattern = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		log.Printf("ERROR: an error occurred:\n %s %v\n", target, err)
		return false
	}
	return true
}

func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	m := IpPattern.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for port := startPort; port <= endPort; port += 1 {
		for ip := startIp; ip <= endIp; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}

	log.Println("neighbors:", neighbors)
	return neighbors
}

func GetHost() string {
	hostname, err1 := os.Hostname()
	if err1 != nil {
		return "127.0.0.1"
	}
	log.Println("hostname:", hostname)
	address, err2 := net.LookupHost(hostname)
	if err2 != nil {
		return "127.0.0.1"
	}
	log.Println("address:", address)
	return address[0]
}
