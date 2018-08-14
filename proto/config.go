
package proto

import (
    "flag"
)


type Mode string
const (
    PAGER  Mode = "pager"
    CLOUD  Mode = "cloud"
    SCROLL Mode = "scroll"    
)
var Modes = []Mode{PAGER,CLOUD,SCROLL}

type Config struct {
    Mode Mode
    Pager *Pager
    Font  *Font
}




func NewConfig(mode Mode) *Config {
    ret := &Config{Mode:mode}
    switch mode {
        case "pager":
            ret.Pager = NewPager()
    }
    ret.Font = NewFont()
    return ret
}


func (config *Config) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(config.Mode), flag.ExitOnError)
    if config.Pager != nil {
        config.Pager.AddFlags(ret)
    }
    if config.Font != nil {
        config.Font.AddFlags(ret)
    }
    return ret
}


func (config *Config) Desc() string {
    ret := "conf[" + string(config.Mode)
    ret += ","
    if config.Pager != nil { 
        ret = ret + config.Pager.Desc() + ","
    }
    if config.Font != nil { 
        ret = ret + config.Font.Desc() + ","
    }
    ret += "]"
    return ret
}



