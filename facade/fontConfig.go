package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"sort"
	"strings"
)

const DEFAULT_FONT = "monaco"

var FontDefaults FontConfig = FontConfig{
	Name: DEFAULT_FONT,
}

func (config *FontConfig) Desc() string {
	ret := "font["
	ret += config.GetName()
	ret += "]"
	return ret
}

func (config *FontConfig) AddFlags(flagset *flag.FlagSet, basicOptions bool) {
	fonts := " (" + AvailableFonts() + ")"
	flagset.StringVar(&config.Name, "font", FontDefaults.Name, "typeface"+fonts)
}

func (config *FontConfig) VisitFlags(flagset *flag.FlagSet, basicOptions bool) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "font":
			{
				config.SetName = true
			}
		}
	})
	return config.SetName
}

func (config *FontConfig) Help(basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("font", flag.ExitOnError)
	config.AddFlags(tmp,basicOptions)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}

func AvailableFonts() string {
	ret := ""
	names := []string{}
	for k, _ := range FontAsset {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		ret += name
		ret += ", "
	}
	ret = strings.TrimSuffix(ret, ", ")
	return ret
}
