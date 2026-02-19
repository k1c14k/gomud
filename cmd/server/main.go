package main

import (
	"flag"
	"goMud/internal/config"
	"goMud/internal/net"
)

func main() {
	configPath := flag.String("config", "mud.yaml", "Path to the config file")
	host := flag.String("host", "", "Server host address (e.g., 0.0.0.0)")
	port := flag.Int("port", 0, "Server port (e.g., 2323)")
	mudlib := flag.String("mudlib", "", "Path to the mudlib directory")
	flag.Parse()

	config.SetOverrides(host, port, mudlib)
	config.SetConfigPath(*configPath)
	s := net.NewServer()
	s.Start()
}
