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

var WordDefaults = WordConfig{
	Slots:     8,
	MaxLength: 0,
	Lifetime:  0.0,
	Watermark: 0.5,
	Shuffle:   false,
	Aging:     false,
}

func (config *WordConfig) Desc() string {
	ret := "words["

	if config.GetSetSlots() {
		ret += fmt.Sprintf("#%d ", config.GetSlots())
	}
	if config.GetSetMaxLength() {
		ret += fmt.Sprintf("≤%d", config.GetMaxLength())
	}
	if config.GetSetLifetime() {
		ret += fmt.Sprintf("%.1fl ", config.GetLifetime())
	}
	if config.GetSetWatermark() {
		ret += fmt.Sprintf("%0.1fm ", config.GetWatermark())
	}
	if config.GetSetShuffle() {
		if !config.GetShuffle() {
			ret += "!"
		}
		ret += "⧢ "
	}
	if config.GetSetAging() {
		if !config.GetAging() {
			ret += "!"
		}
		ret += "å "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *WordConfig) FillPatterns() []string {
	return []string{"title", "index", "alpha", "clear"}
}

func (config *WordConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Uint64Var(&config.Slots, "n", WordDefaults.Slots, "word count")
	flagset.Uint64Var(&config.MaxLength, "m", WordDefaults.MaxLength, "word max length")
	flagset.Float64Var(&config.Lifetime, "life", WordDefaults.Lifetime, "word lifetime")
	flagset.Float64Var(&config.Watermark, "mark", WordDefaults.Watermark, "buffer clear mark")
	flagset.BoolVar(&config.Shuffle, "shuffle", WordDefaults.Shuffle, "shuffle words?")
	flagset.BoolVar(&config.Aging, "aging", WordDefaults.Aging, "age words?")
}

func (config *WordConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "n":
			config.SetSlots = true
			ret = true
		case "m":
			config.SetMaxLength = true
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
		}
	})
	return ret
}

func (config *WordConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("words", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"n", "m", "life", "mark", "shuffle", "aging"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
