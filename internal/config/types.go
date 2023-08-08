package config

import "time"

type Config struct {
	CheckInterval time.Duration `yaml:"checkInterval"`
	StopList      StopList      `yaml:"stopList"`
}

type StopList struct {
	Words   []string `yaml:"words"`
	Domains []string `yaml:"domains"`
	Emails  []string `yaml:"emails"`
}
