//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
)

var TermDefaults TermConfig = TermConfig{
	Grid: nil,
}

func (config *TermConfig) Desc() string {
	ret := "term["
	if grid := config.GetGrid(); grid != nil {
		ret += grid.Desc()
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *TermConfig) AddFlags(flagset *flag.FlagSet) {

	if config.GetGrid() != nil {
		config.GetGrid().AddFlags(flagset)
	}

}

func (config *TermConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		}
	})
	setGrid := false
	if grid := config.GetGrid(); grid != nil {
		setGrid = grid.VisitFlags(flagset)
	}
	return setGrid
}
func (config *TermConfig) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	ret += GridDefaults.Help()
	tmp := flag.NewFlagSet("term", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
