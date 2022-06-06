//
//
package facade

//
//import (
//	"FEEDFACE.COM/facade/gfx"
//	"flag"
//	"strings"
//)
//
//var TagDefaults TagConfig = TagConfig{
//	Shader: nil,
//	Set:    nil,
//}
//
//func (config *TagConfig) Desc() string {
//	ret := "tags["
//	if shader := config.GetShader(); shader != nil {
//		ret += shader.Desc() + " "
//	}
//
//	if set := config.GetSet(); set != nil {
//		ret += set.Desc() + " "
//	}
//
//	ret = strings.TrimRight(ret, " ")
//	ret += "]"
//	return ret
//}
//
//func (config *TagConfig) AddFlags(flagset *flag.FlagSet) {
//	if config.GetShader() != nil {
//		config.GetShader().AddFlags(flagset, Mode_TAGS)
//	}
//	if config.GetSet() != nil {
//		config.GetSet().AddFlags(flagset)
//	}
//}
//
//func (config *TagConfig) VisitFlags(flagset *flag.FlagSet) bool {
//	ret := false
//	if shader := config.GetShader(); shader != nil {
//		if shader.VisitFlags(flagset) {
//			ret = true
//		}
//	}
//	if set := config.GetSet(); set != nil {
//		if set.VisitFlags(flagset) {
//			ret = true
//		}
//	}
//
//	return ret
//}
//
//func (config *TagConfig) Help() string {
//	ret := SetDefaults.Help()
//	tmp := flag.NewFlagSet("tag", flag.ExitOnError)
//	config.AddFlags(tmp)
//	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
//	return ret
//}
