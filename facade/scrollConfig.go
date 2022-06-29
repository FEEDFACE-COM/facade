package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"fmt"

	//	"fmt"
	"strings"
)

var ScrollDefaults = ScrollConfig{
	Repeat:    true,
	CharCount: 32,
}

func (config *ScrollConfig) Desc() string {
	ret := "scroll["

	rok := config.GetSetRepeat()
	if rok {
		if !config.GetRepeat() {
			ret += "!"
		}
		ret += "r "
	}

	if config.GetSetCharCount() {
		ret += fmt.Sprintf("#%d ", config.GetCharCount())
	}

	if config.GetSetFill() {
		ret += "f:" + config.GetFill() + " "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *ScrollConfig) AddFlags(flagset *flag.FlagSet) {
	patterns := "title,index,alpha,clear"
	flagset.BoolVar(&config.Repeat, "repeat", ScrollDefaults.Repeat, "repeat last?")
	flagset.Uint64Var(&config.CharCount, "c", ScrollDefaults.CharCount, "char count")
	flagset.StringVar(&config.Fill, "fill", ScrollDefaults.Fill, "fill pattern ("+patterns+")")
}

func (config *ScrollConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "repeat":
			config.SetRepeat = true
			ret = true
		case "c":
			config.SetCharCount = true
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
	for _, s := range []string{"repeat", "c", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
