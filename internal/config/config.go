package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const configFilePath = "./cmd/config.yml"

func New() (*Config, error) {
	c := &Config{}
	dat, err := c.readYaml()
	if err != nil {
		return nil, fmt.Errorf("pkg config: failed to read yaml config: %w", err)
	}
	if err = c.decode(dat); err != nil {
		return nil, fmt.Errorf("pkg config: failed to decode yaml config: %w", err)
	}
	return c, nil
}

func (c *Config) decode(dat []byte) error {
	return yaml.Unmarshal(dat, &c)
}

func (c *Config) readYaml() ([]byte, error) {
	return os.ReadFile(configFilePath)
}
