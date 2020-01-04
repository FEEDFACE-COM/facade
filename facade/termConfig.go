//
//
package facade

//
import (
	"flag"
	"strings"
)

var TermDefaults TermConfig = TermConfig{
	Grid: &GridConfig{},
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

	if config.GetGrid() == nil {
		config.Grid = &GridConfig{}
	}
	if grid := config.GetGrid(); grid != nil {
		grid.AddFlags(flagset)
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
