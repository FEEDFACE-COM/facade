

package facade

import (
    "flag"    
    "fmt"
)


type GridConfig struct {
    Width uint
    Height uint
    Downward bool
    Scroll bool
    Speed  float64
}


func NewGridConfig() *GridConfig {
    return &GridConfig{
        Width: 0,
        Height: 8,
        Scroll: true, 
        Speed: 0.4,
    }
}


func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Width,"w",config.Width,"grid width")
    flags.UintVar(&config.Height,"h",config.Height,"grid height")
    flags.BoolVar(&config.Downward,"d",config.Downward,"downward")
    flags.BoolVar(&config.Scroll,"s",config.Scroll,"scroll")
    flags.Float64Var(&config.Speed,"S",config.Speed,"scroll speed")

}

func (config *GridConfig) Desc() string { 
    tmp := "↑"
    if config.Downward { tmp = "↓" }

    tmp2 := ""
    if config.Scroll { tmp2 = fmt.Sprintf("%.2f",config.Speed)  }

    return fmt.Sprintf("grid[%dx%d %s%s]",config.Width,config.Height,tmp,tmp2) 
}






