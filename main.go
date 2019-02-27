package main

import (
	"fmt"
	"log"
	"mas/consensus"
)

func main() {
	consensus.RaftInit()
	log.Println("Server created")
	fmt.Scanln()
}
