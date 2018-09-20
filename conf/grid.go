

package conf

import (
    "flag"    
    "fmt"
)


type GridConfig struct {
    Width uint
    Height uint
}


type PageDirection string
const (
    PageUp   PageDirection = "up"
    PageDown PageDirection = "down"
)


func NewGridConfig() *GridConfig {
    return &GridConfig{Width: 16, Height: 8}
}


func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Width,"w",config.Width,"grid width")
    flags.UintVar(&config.Height,"h",config.Height,"grid height")

}

func (config *GridConfig) Desc() string { return fmt.Sprintf("grid[%dx%d]",config.Width,config.Height) }






