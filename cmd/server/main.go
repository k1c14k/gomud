package main

import (
	"flag"
	"goMud/internal/config"
	"goMud/internal/net"
)

func main() {
	configPath := flag.String("config", "config/gomud.yaml", "Path to the config file")
	flag.Parse()
	config.SetConfigPath(*configPath)
	s := net.NewServer()
	s.Start()
}
