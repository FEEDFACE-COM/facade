package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"fmt"

	//	"fmt"
	"strings"
)

var CharDefaults = CharConfig{
	CharCount: 32,
	Speed:     1.0,
	Repeat:    true,
}

func (config *CharConfig) Desc() string {
	ret := "chars["

	if config.GetSetCharCount() {
		ret += fmt.Sprintf("#%d ", config.GetCharCount())
	}

	if config.GetSetSpeed() {
		ret += fmt.Sprintf("s%.1f ", config.GetSpeed())
	}

	rok := config.GetSetRepeat()
	if rok {
		if !config.GetRepeat() {
			ret += "!"
		}
		ret += "⟳ "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *CharConfig) FillPatterns() []string {
	return []string{"title", "index", "alpha", "clear"}
}

func (config *CharConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Uint64Var(&config.CharCount, "c", CharDefaults.CharCount, "char count")
	flagset.Float64Var(&config.Speed, "speed", CharDefaults.Speed, "scroll speed")
	flagset.BoolVar(&config.Repeat, "repeat", CharDefaults.Repeat, "repeat last?")
}

func (config *CharConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "c":
			config.SetCharCount = true
			ret = true
		case "speed":
			config.SetSpeed = true
			ret = true
		case "repeat":
			config.SetRepeat = true
			ret = true
		}
	})

	return ret
}

func (config *CharConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("char", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"c", "speed", "repeat"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
