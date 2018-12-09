
package gfx

import(
//	"fmt"
)

type Config map[string]interface{}


func (config *Config) ApplyConfig(cfg *Config) {
	for key,val := range(*cfg) {
		(*config)[key] = val	
	}	
}