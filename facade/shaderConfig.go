//
//
package facade

//
import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"sort"
	"strings"
)

var ShaderDefaults ShaderConfig = ShaderConfig{
	Vert: "def",
	Frag: "def",
}

func (config *ShaderConfig) Desc() string {
	ret := "shader["

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

func (config *ShaderConfig) AddFlags(flagset *flag.FlagSet, mode Mode, basicOptions bool) {

	frags := " (" + AvailableShaders(PrefixForMode(mode), ".frag") + ")"
	verts := " (" + AvailableShaders(PrefixForMode(mode), ".vert") + ")"

	if basicOptions {
		flagset.StringVar(&config.Vert, "shape", ShaderDefaults.Vert, "shape"+verts)
	} else {
		flagset.StringVar(&config.Vert, "shape", ShaderDefaults.Vert, "vertex shader"+verts)
		flagset.StringVar(&config.Frag, "color", ShaderDefaults.Frag, "fragment shader"+frags)
	}

}

func (config *ShaderConfig) VisitFlags(flagset *flag.FlagSet, mode Mode, basicOptions bool) bool {

	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "shape":
			{
				config.SetVert = true
				if basicOptions {
					frag := PrefixForMode(mode) + config.Vert + "." + string(gfx.FragType)
					if _, ok := ShaderAsset[frag]; ok {
						config.SetFrag = true
						config.Frag = config.Vert
					} else {
						config.SetFrag = true
						config.Frag = ShaderDefaults.Frag
					}
				}
			}
		case "color":
			{
				config.SetFrag = true
			}
		}
	})
	return config.SetVert || config.SetFrag
}

func (config *ShaderConfig) Help(mode Mode, basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("shader", flag.ExitOnError)
	config.AddFlags(tmp, mode, basicOptions)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}

func AvailableShaders(prefix, suffix string) string {
	ret := ""
	names := []string{}
	for k, _ := range ShaderAsset {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
			ret += strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
			ret += ", "
		}
	}
	ret = strings.TrimSuffix(ret, ", ")
	return ret
}

func PrefixForMode(mode Mode) string {
	switch mode {
	case Mode_LINES:
		return "lines/"
	case Mode_TERM:
		return "lines/"
	case Mode_WORDS:
		return "words/"
	case Mode_CHARS:
		return "chars/"
	}
	return ""
}
