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
}

func (config *TagConfig) Desc() string {
	ret := "tags["
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *TagConfig) AddFlags(flagset *flag.FlagSet) {
	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset)
	}
}

func (config *TagConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
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
	for _, s := range []string{} {
		if flg := tmp.Lookup(s); flg != nil {
			fun(flg)
		}
	}
	//tmp.VisitAll(fun)
	return ret
}
