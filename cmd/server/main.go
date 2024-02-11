package main

import (
	"goMud/internal/net"
)

func main() {
	s := net.NewServer()
	s.Start()
}
