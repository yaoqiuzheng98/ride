package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Amap AmapConfig `toml:"amap"`
}

type AmapConfig struct {
	Key           string `toml:"key"`
	DefaultRegion string `toml:"default_region"`
}

var (
	once    sync.Once
	_config *Config
	loadErr error
)

// Load reads config.toml from the given path (once). Subsequent calls return
// the cached config. If path is empty, "config.toml" is used.
func Load(path string) (*Config, error) {
	once.Do(func() {
		if path == "" {
			path = "config.toml"
		}
		data, err := os.ReadFile(path)
		if err != nil {
			loadErr = fmt.Errorf("config: read %s: %w", path, err)
			return
		}
		var c Config
		if err := toml.Unmarshal(data, &c); err != nil {
			loadErr = fmt.Errorf("config: parse %s: %w", path, err)
			return
		}
		_config = &c
	})
	return _config, loadErr
}

// Get returns the loaded config or panics if not loaded.
func Get() *Config {
	if _config == nil {
		panic("config not loaded, call config.Load first")
	}
	return _config
}
