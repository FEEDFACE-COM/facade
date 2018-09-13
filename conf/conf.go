
package conf

import (
    "fmt"
    "flag"
    grid "../grid"
)


type Config struct {
    Mode Mode
    Grid *grid.Config
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

var DEFAULT Mode = GRID


var DIRECTORY = "/Users/folkert/src/gfx/facade/asset/"

func NewConfig(mode Mode) *Config {
    ret := &Config{Mode: mode}
    switch mode {
        case GRID:
            ret.Grid = grid.NewConfig()
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


