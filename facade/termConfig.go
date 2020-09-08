//
//
package facade

//
import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"strings"
)

var TermDefaults TermConfig = TermConfig{}

func (config *TermConfig) Desc() string {
	ret := "term["
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}
	if grid := config.GetGrid(); grid != nil {
		ret += grid.Desc() + " "
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *TermConfig) AddFlags(flagset *flag.FlagSet) {
	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset, Mode_TERM)
	}

	if config.GetGrid() != nil {
		config.GetGrid().AddFlags(flagset)
	}
}

func (config *TermConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		}
	})
	if shader := config.GetShader(); shader != nil {
		if shader.VisitFlags(flagset) {
			ret = true
		}
	}

	if grid := config.GetGrid(); grid != nil {
		if grid.VisitFlags(flagset) {
			ret = true
		}
	}
	return ret
}
func (config *TermConfig) Help() string {
	ret := GridDefaults.Help()
	tmp := flag.NewFlagSet("term", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}
