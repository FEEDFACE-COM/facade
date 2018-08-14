
package conf

import (
    "flag"
)

type Text string

type Mode string
const (
    PAGER  Mode = "pager"
    CLOUD  Mode = "cloud"
    SCROLL Mode = "scroll"    
)
var Modes = []Mode{PAGER,CLOUD,SCROLL}

type Conf struct {
    Mode Mode
    Pager *Pager
    Font  *Font
}




func NewConf(mode Mode) *Conf {
    ret := &Conf{Mode:mode}
    switch mode {
        case "pager":
            ret.Pager = NewPager()
    }
    ret.Font = NewFont()
    return ret
}


func (conf *Conf) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(conf.Mode), flag.ExitOnError)
    if conf.Pager != nil {
        conf.Pager.AddFlags(ret)
    }
    if conf.Font != nil {
        conf.Font.AddFlags(ret)
    }
    return ret
}


func (conf *Conf) Desc() string {
    ret := "conf[" + string(conf.Mode)
    ret += ","
    if conf.Pager != nil { 
        ret = ret + conf.Pager.Desc() + ","
    }
    if conf.Font != nil { 
        ret = ret + conf.Font.Desc() + ","
    }
    ret += "]"
    return ret
}



