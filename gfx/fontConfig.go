
package gfx

import (
//    "fmt"
    "flag"
)
//
//type FontConfig map[string]interface{}
//func (config *FontConfig) Name() (string,bool) { ret,ok := (*config)["font"].(string); return ret,ok }
//func (config *FontConfig) SetName(val string) { (*config)["font"] = val }
//
//
//func (config *FontConfig) AddFlags(flags *flag.FlagSet) {
//    flags.StringVar(&config.Name,"font",config.Name,"use fontface `font`" )
//}
//
//func NewFontConfig() *FontConfig { 
//    ret := make(FontConfig)
//    ret.SetName("RobotoMono")
//    return &ret
//}
//
//func (config *FontConfig) Desc() string {
//    ret := "font["
//    if name,ok := config.Name(); ok {
//        ret += name
//    }
//    ret += "]"
//    return ret
//}
//

type FontConfig string
func (config *FontConfig) Name() (string,bool) { return string(*config),true }
func (config *FontConfig) SetName(val string) { *config = FontConfig(val) }

func (config *FontConfig) AddFlags(flags *flag.FlagSet) {
//    flags.StringVar(&config.Name,"font",config.Name,"use fontface `font`" )    
}

func NewFontConfig() *FontConfig {
    ret := FontConfig("RobotoMono")
    return &ret
}

func (config *FontConfig) Desc() string { 
    ret := "font["
    if name,ok := config.Name(); ok {
        ret += name
    }
    ret += "]"
    return ret
}

