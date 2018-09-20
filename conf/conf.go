
package conf

import (
    "fmt"
    "flag"
)


type Config struct {
    Mode Mode
    Grid *GridConfig
    Line *LineConfig
    Font *FontConfig    
}





type Text string

type Mode string
const (
    GRID  Mode = "grid"
    LINE  Mode = "line"
    WORD  Mode = "word"
    CHAR  Mode = "char"    
)

var Modes = []Mode{GRID,LINE}

var DEFAULT_MODE Mode = LINE


var DIRECTORY = "/home/folkert/src/gfx/facade/asset/"

func NewConfig(mode Mode) *Config {
    ret := &Config{Mode: mode}
    switch mode {
        case GRID:
            ret.Grid = NewGridConfig()
        case LINE:
            ret.Line = NewLineConfig()
    }
    ret.Font = NewFontConfig()
    return ret
}


func (config *Config) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(config.Mode), flag.ExitOnError)
    
    if config.Grid != nil {
        config.Grid.AddFlags(ret)    
    }
    if config.Line != nil {
        config.Line.AddFlags(ret)
    }
    
    if config.Font != nil {
        config.Font.AddFlags(ret)    
    }
    return ret
}


    
    

func (config *Config) Desc() string {
    ret := fmt.Sprintf("conf[%s",string(config.Mode))
    if config.Grid != nil {
        ret += " " + config.Grid.Desc()
    }
    if config.Line != nil {
        ret += " " + config.Line.Desc()
    }
    if config.Font != nil {
        ret += " " + config.Font.Desc()
    }
    ret += "] [dir " + DIRECTORY + "]"
    return ret
}


