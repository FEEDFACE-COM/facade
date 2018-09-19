
package conf

import (
    "fmt"
    "flag"
)


type Config struct {
    Mode Mode
    Grid *GridConfig
    Font *FontConfig    
}





type Text string

type Mode string
const (
    GRID   Mode = "grid"
    CLOUD  Mode = "cloud"
    SCROLL Mode = "scroll"    
)

var Modes = []Mode{GRID,CLOUD,SCROLL}

var DEFAULT_MODE Mode = GRID


var DIRECTORY = "/home/folkert/src/gfx/facade/asset/"

func NewConfig(mode Mode) *Config {
    ret := &Config{Mode: mode}
    switch mode {
        case GRID:
            ret.Grid = NewGridConfig()
    }
    ret.Font = NewFontConfig()
    return ret
}


func (config *Config) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(config.Mode), flag.ExitOnError)
    
    if config.Grid != nil {
        config.Grid.AddFlags(ret)    
    }
    
    if config.Font != nil {
        config.Font.AddFlags(ret)    
    }
    return ret
}


    
    

func (config *Config) Describe() string {
    ret := fmt.Sprintf("conf[%s]",string(config.Mode))
    if config.Grid != nil {
        ret += " " + config.Grid.Describe()
    }
    if config.Font != nil {
        ret += " " + config.Font.Describe()
    }
    ret += " [dir " + DIRECTORY + "]"
    return ret
}


