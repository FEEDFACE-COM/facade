package gfx

import ()

type Config map[string]interface{}

func (config *Config) ApplyConfig(cfg *Config) {
	for key, val := range *cfg {
		(*config)[key] = val
	}
}
