

package facade

import (
    "flag"    
    "fmt"
    "strings"
)


type GridConfig struct {
    width uint
    height uint
    downward bool
    scroll bool
    speed  float64
    
    vert string
    frag string
    
    
    fill string
    
}


func (config *GridConfig) Width()    uint    { return config.width    }
func (config *GridConfig) Height()   uint    { return config.height   }
func (config *GridConfig) Downward() bool    { return config.downward }
func (config *GridConfig) Scroll()   bool    { return config.scroll   }
func (config *GridConfig) Speed()    float64 { return config.speed    }
func (config *GridConfig) Vert()     string  { return config.vert     }
func (config *GridConfig) Frag()     string  { return config.frag     }
func (config *GridConfig) Fill()     string  { return config.fill     }


func (config *GridConfig) SetWidth(w uint) { config.width = w }

func NewGridConfig() *GridConfig {
    return &GridConfig{
        width: 0,
        height: 8,
        scroll: true, 
        speed: 0.4,
        vert: "null",
        frag: "null",
        fill: "",
    }
}


func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(   &config.width,   "w",   config.width,   "grid width")
    flags.UintVar(   &config.height,  "h",   config.height,  "grid height")
    flags.BoolVar(   &config.downward,"d",   config.downward,"downward")
    flags.BoolVar(   &config.scroll,  "s",   config.scroll,  "scroll")
    flags.Float64Var(&config.speed,   "S",   config.speed,   "scroll speed")
    flags.StringVar( &config.vert,    "V",   config.vert,    "vertex shader")
    flags.StringVar( &config.frag,    "F",   config.frag,    "fragment shader")
    flags.StringVar( &config.fill,    "fill",config.fill,    "fill pattern")
}

func (config *GridConfig) Desc() string { 
    tmp := "↑"
    if config.downward { tmp = "↓" }

    tmp2 := ""
    if config.scroll { tmp2 = fmt.Sprintf("%.2f",config.speed)  }


    return fmt.Sprintf("grid[%dx%d %s%s %s,%s]",config.width,config.height,tmp,tmp2,config.vert,config.frag) 
}




func (config *GridConfig) Clean() {
    config.vert = strings.Replace(config.vert,"/","",-1)
    config.frag = strings.Replace(config.frag,"/","",-1)
}

