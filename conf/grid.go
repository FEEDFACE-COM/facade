

package conf

import (
    "flag"    
    "fmt"
)


type GridConfig struct {
    Width uint
    Height uint
    Downward bool
}


func NewGridConfig() *GridConfig {
    return &GridConfig{Width: 0, Height: 2}
}


func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Width,"w",config.Width,"grid width")
    flags.UintVar(&config.Height,"h",config.Height,"grid height")
    flags.BoolVar(&config.Downward,"d",config.Downward,"downward")

}

func (config *GridConfig) Desc() string { 
    tmp := "↑"
    if config.Downward {
        tmp = "↓"
    }
    return fmt.Sprintf("grid[%dx%d %s]",config.Width,config.Height,tmp) 
}






