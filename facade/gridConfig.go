

package facade

import (
    "flag"    
    "fmt"
    "strings"
)


type GridConfig struct {
    Width uint
    Height uint
    Downward bool
    Scroll bool
    Speed  float64
    
    Vert string
    Frag string
    
}


func NewGridConfig() *GridConfig {
    return &GridConfig{
        Width: 0,
        Height: 8,
        Scroll: true, 
        Speed: 0.4,
        Vert: "identity",
        Frag: "identity",
    }
}


func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Width,"w",config.Width,"grid width")
    flags.UintVar(&config.Height,"h",config.Height,"grid height")
    flags.BoolVar(&config.Downward,"d",config.Downward,"downward")
    flags.BoolVar(&config.Scroll,"s",config.Scroll,"scroll")
    flags.Float64Var(&config.Speed,"S",config.Speed,"scroll speed")
    flags.StringVar(&config.Vert,"V",config.Vert,"vertex shader")
    flags.StringVar(&config.Frag,"F",config.Frag,"fragment shader")
}

func (config *GridConfig) Desc() string { 
    tmp := "↑"
    if config.Downward { tmp = "↓" }

    tmp2 := ""
    if config.Scroll { tmp2 = fmt.Sprintf("%.2f",config.Speed)  }


    return fmt.Sprintf("grid[%dx%d %s%s %s,%s]",config.Width,config.Height,tmp,tmp2,config.Vert,config.Frag) 
}




func (config *GridConfig) Clean() {
    config.Vert = strings.Replace(config.Vert,"/","",-1)
    config.Frag = strings.Replace(config.Frag,"/","",-1)
}

