package main

import (
	"blockchain/globals"
	"fmt"
)

func main() {
	fmt.Printf("gethost:%s", globals.GetHost())

	//neighbors := globals.FindNeighbors(
	//	"127.0.0.1",
	//	5000,
	//	0,
	//	3,
	//	5000,
	//	5004,
	//)
	//
	//log.Printf("NEIGBORS:%s\n", neighbors)

	//serverPorts := []uint16{5000, 5001, 5002, 5003}
	//
	//for _, port := range serverPorts {
	//	fmt.Println(globals.IsFoundHost("127.0.0.1", port))
	//}
}
