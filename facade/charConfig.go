

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
//	"fmt"
	"strings"
)
    
var CharDefaults CharConfig = CharConfig{
	Shader: nil,
	Scroll: nil,
}

func (config *CharConfig) Desc() string {
	ret := "chars["
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}

	if scroll := config.GetScroll(); scroll != nil {
		ret += scroll.Desc() + " "
	}


	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *CharConfig) AddFlags(flagset *flag.FlagSet) {
	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset, Mode_CHARS)
	}
	if config.GetScroll() != nil {
		config.GetScroll().AddFlags(flagset)
	}
}


func (config *CharConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	if shader := config.GetShader(); shader != nil {
		if shader.VisitFlags(flagset) {
			ret = true
		}
	}
	if scroll := config.GetScroll(); scroll != nil {
		if scroll.VisitFlags(flagset) {
			ret = true
		}
	}

	return ret
}

func (config *CharConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("char", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}
