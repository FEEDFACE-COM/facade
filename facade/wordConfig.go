//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
)

var WordDefaults WordConfig = WordConfig{
	Shader: nil,
	Set: nil,
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
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	ret += SetDefaults.Help()
	tmp := flag.NewFlagSet("word", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{} {
		if flg := tmp.Lookup(s); flg != nil {
			fun(flg)
		}
	}
	//tmp.VisitAll(fun)
	return ret
}
