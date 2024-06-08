package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ServerConfig struct {
	Address string `yaml:"address"`
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
		Address: "0.0.0.0:2323",
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
		return &defaultConfig
	}

	b, err := os.ReadFile(*configPath)
	if err != nil {
		log.Println("Error reading config file:", err)
		return &defaultConfig
	}

	return loadConfigFromYamlString(string(b))
}

func GetConfig() Config {
	if config == nil {
		config = loadConfig()
	}

	return *config
}
