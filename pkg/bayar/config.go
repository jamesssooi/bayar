package bayar

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Config describes the configuration values for Bayar.
type Config struct {
	HostAddress          string
	PortNumber           int
	GoogleSecretFilePath string
}

var cachedCfg Config

// LoadConfig loads configuration values for Bayar. Set the BAYAR_CONFIG
// environment variable to specify the location of the configuration file.
func LoadConfig() (Config, error) {
	var cfg Config

	if cachedCfg != (Config{}) {
		return cachedCfg, nil
	}

	cfgPath, hasEnv := os.LookupEnv("BAYAR_CONFIG")
	if !hasEnv {
		cwd, _ := os.Getwd()
		cfgPath = path.Join(cwd, "bayar.json")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return cfg, err
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}

	cachedCfg = cfg

	return cfg, nil
}
