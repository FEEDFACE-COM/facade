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
	Vert: "def",
	Frag: "def",
}

func (config *TagConfig) Desc() string {
	ret := "tags["
	{
		vok := config.GetSetVert()
		fok := config.GetSetFrag()
		if vok {
			ret += config.GetVert()
		}
		if vok || fok {
			ret += ","
		}
		if fok {
			ret += config.GetFrag()
		}
		if vok || fok {
			ret += " "
		}
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *TagConfig) AddFlags(flagset *flag.FlagSet) {

	flagset.StringVar(&config.Vert, "vert", GridDefaults.Vert, "vertex shader")
	flagset.StringVar(&config.Frag, "frag", GridDefaults.Frag, "fragment shader")

}

func (config *TagConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "vert":
			{
				config.SetVert = true
			}
		case "frag":
			{
				config.SetFrag = true
			}
		}
	})
	return config.SetVert || config.SetFrag
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

	tmp := flag.NewFlagSet("set", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"frag", "vert"} {
		if flg := tmp.Lookup(s); flg != nil {
			fun(flg)
		}
	}
	//tmp.VisitAll(fun)
	return ret
}
