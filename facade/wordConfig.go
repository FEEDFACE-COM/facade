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
	MaxLength: 20,
	Lifetime:  8.0,
	Watermark: 0.8,
	Shuffle:   false,
	Aging:     false,
	Unique:    false,
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

	if config.GetSetUnique() {
		if !config.GetUnique() {
			ret += "!"
		}
		ret += "û "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *WordConfig) FillPatterns() []string {
	return []string{"title", "index", "alpha", "clear"}
}

func (config *WordConfig) AddFlags(flagset *flag.FlagSet, basicOptions bool) {
	flagset.Uint64Var(&config.Slots, "n", WordDefaults.Slots, "number of word slots")
	flagset.Float64Var(&config.Lifetime, "life", WordDefaults.Lifetime, "word lifetime")
	flagset.Float64Var(&config.Watermark, "mark", WordDefaults.Watermark, "buffer fill mark")
	flagset.BoolVar(&config.Shuffle, "shuffle", WordDefaults.Shuffle, "shuffle words?")
	if !basicOptions {
		flagset.Uint64Var(&config.MaxLength, "m", WordDefaults.MaxLength, "word max length")
		flagset.BoolVar(&config.Aging, "age", WordDefaults.Aging, "decaying words?")
		flagset.BoolVar(&config.Unique, "uniq", WordDefaults.Unique, "only unique words?")
	}
}

func (config *WordConfig) VisitFlags(flagset *flag.FlagSet, basicOptions bool) bool {
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
		case "age":
			config.SetAging = true
			ret = true
		case "uniq":
			config.SetUnique = true
			ret = true

		}
	})
	return ret
}

func (config *WordConfig) Help(basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("words", flag.ExitOnError)
	config.AddFlags(tmp,basicOptions)
	for _, s := range []string{"n", "m", "life", "mark", "shuffle", "uniq", "age"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
