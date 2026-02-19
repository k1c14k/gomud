package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MudlibConfig struct {
	MudlibPath string `yaml:"mudlib_path"`
}

type Config struct {
	ServerConfig ServerConfig `yaml:"server_config"`
	MudlibConfig MudlibConfig `yaml:"mudlib_config"`
}

var defaultConfig = Config{
	ServerConfig: ServerConfig{
		Host: "0.0.0.0",
		Port: 2323,
	},
	MudlibConfig: MudlibConfig{
		MudlibPath: "mudlib/",
	},
}

func loadConfigFromYamlString(yamlString string) *Config {
	config := defaultConfig

	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		log.Println("Error unmarshalling yaml:", err)
	}

	return &config
}

var config *Config = nil
var configPath *string = nil

func SetConfigPath(path string) {
	configPath = &path
}

func loadConfig() *Config {
	if configPath == nil {
		config := defaultConfig
		applyOverrides(&config)
		return &config
	}

	b, err := os.ReadFile(*configPath)
	if err != nil {
		log.Println("Error reading config file:", err)
		config := defaultConfig
		applyOverrides(&config)
		return &config
	}

	config := loadConfigFromYamlString(string(b))
	applyOverrides(config)
	return config
}

func GetConfig() Config {
	if config == nil {
		config = loadConfig()
	}

	return *config
}

var (
	HostOverride   *string
	PortOverride   *int
	MudlibOverride *string
)

func SetOverrides(host *string, port *int, mudlib *string) {
	HostOverride = host
	PortOverride = port
	MudlibOverride = mudlib
}

func applyOverrides(cfg *Config) {
	if HostOverride != nil && *HostOverride != "" {
		cfg.ServerConfig.Host = *HostOverride
	}
	if PortOverride != nil && *PortOverride != 0 {
		cfg.ServerConfig.Port = *PortOverride
	}
	if MudlibOverride != nil && *MudlibOverride != "" {
		cfg.MudlibConfig.MudlibPath = *MudlibOverride
	}
}
