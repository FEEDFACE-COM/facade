package facade

import (
    "flag"    
    "fmt"
    "../gfx"
)

const (
	testWidth    = "width"
	testHeight   = "height"
)


type TestConfig gfx.Config



func NewTestConfig() *TestConfig {
    return &TestConfig{}
}


func (config *TestConfig) Width()     (uint,bool) { ret,ok := (*config)[testWidth   ].(float64); return uint(ret), ok }
func (config *TestConfig) Height()    (uint,bool) { ret,ok := (*config)[testHeight  ].(float64); return uint(ret), ok }
func (config *TestConfig) SetWidth(   val uint)    { (*config)[testWidth]    = float64(val) }
func (config *TestConfig) SetHeight(  val uint)    { (*config)[testHeight]   = float64(val) }



func (config *TestConfig) ApplyConfig(cfg *TestConfig) {
	if tmp,ok := cfg.Width(); ok { config.SetWidth(tmp) }	
	if tmp,ok := cfg.Height(); ok { config.SetHeight(tmp) }	
}

func (config *TestConfig) AddFlags(flags *flag.FlagSet) {
}

func (config *TestConfig) Desc() string { 
    ret := "test["
    {
    	w,wok := config.Width(); 
	    h,hok := config.Height();
	    if wok { ret += fmt.Sprintf("%d",w) }
    	if wok || hok { ret += "x" }
	    if hok { ret += fmt.Sprintf("%d",h) }
	    if wok || hok { ret += " " }
	}
	ret += "]"
	return ret
}



type TestState struct {
    Width uint
    Height uint   
}


var TestDefaults = TestState{
    Width:       25,
    Height:       8,
}    


func (state *TestState) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&state.Width,"w",state.Width,"test width")
    flags.UintVar(&state.Height,"h",state.Height,"test height")
}

func (state *TestState) CheckFlags(flags *flag.FlagSet) (*TestConfig,bool) {
	ok := false
	ret := make(TestConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "w" { ok = true; ret.SetWidth( state.Width ) }
		if f.Name == "h" { ok = true; ret.SetHeight( state.Height ) }
	})
	return &ret,ok
}


func (state *TestState) Desc() string { return state.Config().Desc() }


func (state *TestState) Config() *TestConfig {
	ret := make(TestConfig)
	ret.SetWidth(state.Width)
	ret.SetHeight(state.Height)
	return &ret	
}


func (state *TestState) ApplyConfig(config *TestConfig) bool {
	changed := false
	if tmp,ok := config.Width();    ok { if state.Width    != tmp { changed = true }; state.Width = tmp }
	if tmp,ok := config.Height();   ok { if state.Height   != tmp { changed = true }; state.Height = tmp }
	return changed
}


