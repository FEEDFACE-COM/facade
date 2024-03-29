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
	Repeat:    false,
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

func (config *CharConfig) AddFlags(flagset *flag.FlagSet, basicOptions bool) {
	flagset.Uint64Var(&config.CharCount, "w", CharDefaults.CharCount, "width: chars in line")
	flagset.Float64Var(&config.Speed, "speed", CharDefaults.Speed, "scroll speed")
	if !basicOptions {
		flagset.BoolVar(&config.Repeat, "repeat", CharDefaults.Repeat, "repeat last?")
	}
}

func (config *CharConfig) VisitFlags(flagset *flag.FlagSet, basicOptions bool) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "w":
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

func (config *CharConfig) Help(basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("char", flag.ExitOnError)
	config.AddFlags(tmp, basicOptions)
	for _, s := range []string{"w", "speed", "repeat"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
