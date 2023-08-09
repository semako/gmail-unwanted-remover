package config

import "time"

type Config struct {
	CredentialsFilePath string        `yaml:"credentialsFilePath"`
	TokenFilePath       string        `yaml:"tokenFilePath"`
	CheckInterval       time.Duration `yaml:"checkInterval"`
	Matcher             Matcher       `yaml:"matcher"`
}

type Matcher struct {
	AllowedWords []string `yaml:"allowedWords"`
	StopList     StopList `yaml:"stopList"`
}

type StopList struct {
	Words   []string `yaml:"words"`
	Domains []string `yaml:"domains"`
	Emails  []string `yaml:"emails"`
}
