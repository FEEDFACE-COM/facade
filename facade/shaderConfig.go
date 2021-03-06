//
//
package facade

//
import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
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

func (config *ShaderConfig) AddFlags(flagset *flag.FlagSet, mode Mode) {

	frags := " (" + availableShaders(prefixForMode(mode), ".frag") + ")"
	verts := " (" + availableShaders(prefixForMode(mode), ".vert") + ")"

	flagset.StringVar(&config.Vert, "vert", ShaderDefaults.Vert, "vertex shader"+verts)
	flagset.StringVar(&config.Frag, "frag", ShaderDefaults.Frag, "fragment shader"+frags)

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

func (config *ShaderConfig) Help(mode Mode) string {
	ret := ""
	tmp := flag.NewFlagSet("shader", flag.ExitOnError)
	config.AddFlags(tmp, mode)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}

func availableShaders(prefix, suffix string) string {
	ret := ""
	for name, _ := range ShaderAsset {
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
			ret += strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
			ret += ", "
		}
	}
	ret = strings.TrimSuffix(ret, ", ")
	return ret
}

func prefixForMode(mode Mode) string {
	switch mode {
	case Mode_LINES:
		return "grid/"
	case Mode_TERM:
		return "grid/"
	case Mode_TAGS:
		return "set/"
	case Mode_WORDS:
		return "set/"
	}
	return ""
}
