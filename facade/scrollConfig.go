

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	//	"fmt"
	"strings"
)
    
var ScrollDefaults ScrollConfig = ScrollConfig{
	Repeat: true,
}

func (config *ScrollConfig) Desc() string {
	ret := "scroll["

	uok := config.GetSetRepeat()
	if uok {
		ret += " "
		if config.GetRepeat() {
			ret += "+"
		} else {
			ret += "-"
		}
		ret += "r"
	}

	if config.GetSetFill() {
		ret += " f" + config.GetFill()
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *ScrollConfig) AddFlags(flagset *flag.FlagSet) {
	patterns := "title,index,alpha,clear"
	flagset.BoolVar(&config.Repeat, "repeat", ScrollDefaults.Repeat, "repeat last?")
	flagset.StringVar(&config.Fill, "fill", ScrollDefaults.Fill, "fill pattern ("+patterns+")")
}


func (config *ScrollConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "repeat":
			config.SetRepeat = true
			ret = true
		case "fill":
			config.SetFill = true
			ret = true
		}
	})

	return ret
}

func (config *ScrollConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("char", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"repeat", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
