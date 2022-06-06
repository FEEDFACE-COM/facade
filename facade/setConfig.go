//
//
package facade

//
import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"fmt"
	"strings"
)

var SetDefaults SetConfig = SetConfig{
	Slots:     8,
	Lifetime: 0.0,
	Watermark: 0.8,
	Shuffle:  false,
	Aging:    false,
}

func (config *SetConfig) Desc() string {
	ret := "set["

	sok := config.GetSetSlots()
	if sok {
		ret += fmt.Sprintf("#%d ", config.GetSlots())
	}

	lok := config.GetSetLifetime()
	if lok {
		ret += fmt.Sprintf("%.1fl ", config.GetLifetime())
	}

	wok := config.GetSetWatermark()
	if wok {
		ret += fmt.Sprintf("%.1fm ", config.GetWatermark())
	}

	uok := config.GetShuffle()
	if uok {
		ret += "⧢ "
	}

	aok := config.GetAging()
	if aok {
		ret += "a "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *SetConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Float64Var(&config.Lifetime, "life", SetDefaults.Lifetime, "word lifetime")
	flagset.Float64Var(&config.Watermark, "mark", SetDefaults.Watermark, "clear watermark")
	flagset.Uint64Var(&config.Slots, "n", SetDefaults.Slots, "word count")
	flagset.BoolVar(&config.Shuffle, "shuffle", SetDefaults.Shuffle, "shuffle words?")
	flagset.BoolVar(&config.Aging, "aging", SetDefaults.Aging, "age words?")
	flagset.StringVar(&config.Fill, "fill", SetDefaults.Fill, "fill pattern")
}

func (config *SetConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "n":
			config.SetSlots = true
			ret = true
		case "life":
			config.SetLifetime = true
			ret = true
		case "mark":
			config.SetWatermark = true
			ret = true
		case "shuffle":
			config.SetShuffle = true
			ret = true
		case "aging":
			config.SetAging = true
			ret = true
		case "fill":
			config.SetFill = true
			ret = true
		}
	})

	return ret
}

func (config *SetConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("set", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"n", "life", "mark", "shuffle", "aging", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
