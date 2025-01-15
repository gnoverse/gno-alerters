package config

import (
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	ChainID string `toml:"chain_id"`

	RPC struct {
		Endpoint string `toml:"endpoint"`
	} `toml:"rpc"`

	Slack struct {
		Token     string `toml:"token"`
		ChannelID string `toml:"channel_id"`
	} `toml:"slack"`

	Alerts struct {
		ConsecutiveMissed []int `toml:"consecutive_missed"`
		stalledPeriod     time.Duration
		StalledPeriod     string `toml:"stalled_period"`
	}

	Validators map[string]string `toml:"validators"`
}

func (cfg Config) GetStalledPeriod() time.Duration {
	return cfg.Alerts.stalledPeriod
}

func (cfg Config) HasConsecutiveMissed(missed int) bool {
	for _, v := range cfg.Alerts.ConsecutiveMissed {
		if v == missed {
			return true
		}
	}
	return false
}

func (cfg Config) GetValidatorMoniker(addr string) string {
	for k, v := range cfg.Validators {
		if v == addr {
			return k
		}
	}
	return addr
}

func ParseConfig(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := Config{}

	err = toml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	config.Alerts.stalledPeriod, err = time.ParseDuration(config.Alerts.StalledPeriod)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
