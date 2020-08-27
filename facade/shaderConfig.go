//
//
package facade

//
import (
	"flag"
	gfx "../gfx"
	"strings"
)

var ShaderDefaults ShaderConfig = ShaderConfig{
	Vert: "def",
	Frag: "def",
}

func (config *ShaderConfig) Desc() string {
	ret := ""

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
	ret += ""
	return ret
}

func (config *ShaderConfig) AddFlags(flagset *flag.FlagSet) {

	flagset.StringVar(&config.Vert, "vert", ShaderDefaults.Vert, "vertex shader")
	flagset.StringVar(&config.Frag, "frag", ShaderDefaults.Frag, "fragment shader")

}

func (config *ShaderConfig) VisitFlags(flagset *flag.FlagSet) bool {
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

func (config *ShaderConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("shader", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll( func (f *flag.Flag) { ret += gfx.FlagHelp(f) } )
	return ret
}


