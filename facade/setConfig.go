//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
	gfx "../gfx"
)

var SetDefaults SetConfig = SetConfig{
	Duration: 2.0,
	Shuffle: false,
	Slot: 8,
}

func (config *SetConfig) Desc() string {
	ret := "set["

    dok := config.GetSetDuration()
    if dok {
        ret += fmt.Sprintf("%.1f ",config.GetDuration())    
    }

    sok := config.GetSetSlot()
    if sok {
        ret += fmt.Sprintf("%d ",config.GetSlot())
    }
    
    uok := config.GetShuffle()
    if uok {
        ret += "⧢ "
    }

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *SetConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Float64Var(&config.Duration, "life", SetDefaults.Duration, "word lifetime")
	flagset.Uint64Var(&config.Slot, "slot", SetDefaults.Slot, "word count")
	flagset.BoolVar(&config.Shuffle, "shuffle", SetDefaults.Shuffle, "word shuffle?")
	flagset.StringVar(&config.Fill, "fill", SetDefaults.Fill, "word fill pattern")
}

func (config *SetConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "life":
            config.SetDuration = true
            ret = true
		case "slot":
            config.SetSlot = true
            ret = true
        case "shuffle":
            config.SetShuffle = true
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
	for _, s := range []string{"life", "slot", "shuffle", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
    
