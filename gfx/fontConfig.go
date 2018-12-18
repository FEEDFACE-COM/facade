
package gfx

import (
//    "fmt"
    "flag"
)


//const DEFAULT_FONT = "vt323"
const DEFAULT_FONT = "Menlo"



type FontConfig Config

const(
	fontName = "font"
)


func (config *FontConfig) Name() (string,bool) { ret,ok := (*config)[fontName].(string); return ret,ok }
func (config *FontConfig) SetName(val string) { (*config)[fontName] = val }

func (config *FontConfig) Desc() string { 
    ret := "font["
    if name,ok := config.Name(); ok {
        ret += name
    }
    ret += "]"
    return ret
}



func (config *FontConfig) ApplyConfig(cfg *FontConfig) {
	if tmp,ok := cfg.Name(); ok { config.SetName( tmp ) }	
}

type FontState struct {
	Name string
}

var FontDefaults = FontState{
	Name: DEFAULT_FONT,	
}

func (state *FontState) AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&state.Name,"font",state.Name,"font face")
}

func (state *FontState) CheckFlags(flags *flag.FlagSet) (*FontConfig,bool) {
	ok := false
	ret := make(FontConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "font" { ok = true; ret.SetName( state.Name ) }
	})
	return &ret,ok
}




func (state *FontState) Desc() string { return state.Config().Desc() }

func (state *FontState) Config() *FontConfig {
	ret := make(FontConfig)
	ret.SetName(state.Name)	
	return &ret
}


func (state *FontState) ApplyConfig(config *FontConfig) bool {
	changed := false
	if tmp,ok := config.Name(); ok { if state.Name != tmp { changed = true }; state.Name = tmp }
	return changed	
}





