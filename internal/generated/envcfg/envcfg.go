package envcfg

/* DO NOT EDIT. This file generated due to '@envconfig' annotation */

import (
	"fmt"

	a "github.com/freekup/mini-wallet/internal/app/infra"
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", LoadApp)
	typapp.Provide("", LoadCacheCfg)
	typapp.Provide("pg", LoadPgDatabaseCfg)
	typapp.Provide("", LoadEchoCfg)
	typapp.Provide("", LoadKafkaCfg)
}

// LoadApp load env to new instance of App
func LoadApp() (*a.App, error) {
	var cfg a.App
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadCacheCfg load env to new instance of CacheCfg
func LoadCacheCfg() (*a.CacheCfg, error) {
	var cfg a.CacheCfg
	prefix := "CACHE"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadPgDatabaseCfg load env to new instance of DatabaseCfg
func LoadPgDatabaseCfg() (*a.DatabaseCfg, error) {
	var cfg a.DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadEchoCfg load env to new instance of EchoCfg
func LoadEchoCfg() (*a.EchoCfg, error) {
	var cfg a.EchoCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadKafkaCfg load env to new instance of KafkaCfg
func LoadKafkaCfg() (*a.KafkaCfg, error) {
	var cfg a.KafkaCfg
	prefix := "BROKER_KAFKA"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}
