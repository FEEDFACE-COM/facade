

package grid

import (
    "flag"    
    "fmt"
)


type Config struct {
    Width uint
    Height uint
}


type PageDirection string
const (
    PageUp   PageDirection = "up"
    PageDown PageDirection = "down"
)


func NewConfig() *Config {
    return &Config{Width: 16, Height: 8}
}


func (config *Config) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Width,"w",config.Width,"grid width")
    flags.UintVar(&config.Height,"h",config.Height,"grid height")

}

func (config *Config) Describe() string { return fmt.Sprintf("grid[%dx%d]",config.Width,config.Height) }






