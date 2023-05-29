package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"wifi-scaner-credentials/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type string `yaml:"type" env:"TYPE" env-default:"port"`
		Host string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
		Port string `yaml:"port" env:"PORT" env-default:"8000"`
	} `yaml:"listen"`
	SSL struct {
		CertPath string `yaml:"cert_path" env:"CERT_PATH" env-default:"ssl/server.crt"`
		KeyPath  string `yaml:"key_path" env:"KEY_PATH" env-default:"ssl/server.key"`
	} `yaml:"ssl"`
	API struct {
		RegistryID        string `yaml:"registry_id" env:"REGISTRY_ID" env-required:"true"`
		AuthorizedKeyPath string `yaml:"authorized_key_path" env:"AUTHORIZED_KEY_PATH" env-default:"./authorized_key.json"`
	} `yaml:"api"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read application configuration")
		instance = &Config{}

		if err := cleanenv.ReadConfig("./config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
