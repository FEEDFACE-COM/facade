//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
)

var TagDefaults TagConfig = TagConfig{
	Shader: nil,
	Duration: 1.0,
	Shuffle: false,
	Slot: 10,
}

func (config *TagConfig) Desc() string {
	ret := "tags["
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}

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

func (config *TagConfig) AddFlags(flagset *flag.FlagSet) {
	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset)
	}
	flagset.Float64Var(&config.Duration, "life", TagDefaults.Duration, "tag lifetime")
	flagset.Uint64Var(&config.Slot, "slot", TagDefaults.Slot, "tag slot count")
	flagset.BoolVar(&config.Shuffle, "shuffle", TagDefaults.Shuffle, "shuffle slot?")
	flagset.StringVar(&config.Fill, "fill", TagDefaults.Fill, "tag fill pattern")
}

func (config *TagConfig) VisitFlags(flagset *flag.FlagSet) bool {
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
	if shader := config.GetShader(); shader != nil {
		if shader.VisitFlags(flagset) {
			ret = true
		}
	}

	return ret
}

func (config *TagConfig) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	ret += ShaderDefaults.Help()
	tmp := flag.NewFlagSet("set", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"life", "slot", "shuffle", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			fun(flg)
		}
	}
	//tmp.VisitAll(fun)
	return ret
}
