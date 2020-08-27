//
//
package facade

//
import (
	gfx "../gfx"
	"flag"
	"strings"
)

var WordDefaults WordConfig = WordConfig{
	Shader: nil,
	Set:    nil,
}

func (config *WordConfig) Desc() string {
	ret := "words["
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}

	if set := config.GetSet(); set != nil {
		ret += set.Desc() + " "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *WordConfig) AddFlags(flagset *flag.FlagSet) {
	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset)
	}
	if config.GetSet() != nil {
		config.GetSet().AddFlags(flagset)
	}
}

func (config *WordConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	if shader := config.GetShader(); shader != nil {
		if shader.VisitFlags(flagset) {
			ret = true
		}
	}
	if set := config.GetSet(); set != nil {
		if set.VisitFlags(flagset) {
			ret = true
		}
	}

	return ret
}

func (config *WordConfig) Help() string {
	ret := SetDefaults.Help()
	tmp := flag.NewFlagSet("word", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}
