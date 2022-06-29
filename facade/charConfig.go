package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"fmt"

	//	"fmt"
	"strings"
)

var CharDefaults = CharConfig{
	Repeat:    true,
	CharCount: 32,
}

func (config *CharConfig) Desc() string {
	ret := "chars["

	if config.GetSetCharCount() {
		ret += fmt.Sprintf("#%d ", config.GetCharCount())
	}

	rok := config.GetSetRepeat()
	if rok {
		if !config.GetRepeat() {
			ret += "!"
		}
		ret += "⟳ "
	}

	if config.GetSetFill() {
		ret += "f:" + config.GetFill() + " "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *CharConfig) AddFlags(flagset *flag.FlagSet) {
	patterns := "title,index,alpha,clear"
	flagset.Uint64Var(&config.CharCount, "c", CharDefaults.CharCount, "char count")
	flagset.BoolVar(&config.Repeat, "repeat", CharDefaults.Repeat, "repeat last?")
	flagset.StringVar(&config.Fill, "fill", CharDefaults.Fill, "fill pattern ("+patterns+")")
}

func (config *CharConfig) VisitFlags(flagset *flag.FlagSet) bool {
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

func (config *CharConfig) Help() string {
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
