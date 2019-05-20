package facade

import (
    "flag"    
    "fmt"
    "strings"
    "../gfx"
)

const (
	testWidth    = "width"
	testHeight   = "height"
	testBuffer   = "buffer"
	testBufLen   = "buflen"
	testSpeed    = "bufspeed"
)

type BufferName string
const (
    TERMBUFFER BufferName = "term"
    LINEBUFFER BufferName = "line"
)


type TestConfig gfx.Config



func NewTestConfig() *TestConfig {
    return &TestConfig{}
}


func (config *TestConfig) Width()     (uint,bool)       { ret,ok := (*config)[testWidth   ].(float64); return uint(ret), ok }
func (config *TestConfig) Height()    (uint,bool)       { ret,ok := (*config)[testHeight  ].(float64); return uint(ret), ok }
func (config *TestConfig) Buffer()    (BufferName,bool) { ret,ok := (*config)[testBuffer  ].(string); return BufferName(ret), ok }
func (config *TestConfig) BufLen()    (uint,bool)       { ret,ok := (*config)[testBufLen  ].(float64); return uint(ret), ok }
func (config *TestConfig) Speed()     (float64,bool)    { ret,ok := (*config)[testSpeed   ].(float64); return      ret , ok }
func (config *TestConfig) SetWidth(   val uint)      { (*config)[testWidth]    = float64(val) }
func (config *TestConfig) SetHeight(  val uint)      { (*config)[testHeight]   = float64(val) }
func (config *TestConfig) SetBuffer( val BufferName) { (*config)[testBuffer]   =  string(val) }  
func (config *TestConfig) SetBufLen( val uint)       { (*config)[testBufLen]   =  float64(val) }  
func (config *TestConfig) SetSpeed(   val float64) { (*config)[testSpeed]    =         val  }



func (config *TestConfig) ApplyConfig(cfg *TestConfig) {
	if tmp,ok := cfg.Width(); ok { config.SetWidth(tmp) }	
	if tmp,ok := cfg.Height(); ok { config.SetHeight(tmp) }	
	if tmp,ok := cfg.Speed(); ok { config.SetSpeed(tmp) }	
}

func (config *TestConfig) AddFlags(flags *flag.FlagSet) {
}

func (config *TestConfig) Desc() string { 
    ret := "test["
    {
    	w,wok := config.Width(); 
	    h,hok := config.Height();
    	l,lok := config.BufLen()
	    if wok { ret += fmt.Sprintf("%d",w) }
    	if wok || hok || lok { ret += "x" }
	    if hok { ret += fmt.Sprintf("%d",h) }
	    if lok { ret += fmt.Sprintf("+%d",l) }
	    if wok || hok || lok { ret += " " }
    }

	{
        buf,bok := config.Buffer()
        if bok { ret += string(buf) + " " }	
    }
    
    {
        spd,sok := config.Speed()
        if sok { ret += fmt.Sprintf("%.1f ",spd) }   
    }

    ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}



type TestState struct {
    Width uint
    Height uint   
    Buffer BufferName
    BufLen uint
    Speed float64
}


var TestDefaults = TestState{
    Width:       25,
    Height:       8,
    Buffer: TERMBUFFER,
    BufLen:       2,
    Speed:      1.0,
}


func (state *TestState) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&state.Width,"w",state.Width,"test width")
    flags.UintVar(&state.Height,"h",state.Height,"test height")
    flags.StringVar( (*string)(&state.Buffer), "mode",string(state.Buffer),"test buffer mode")
    flags.UintVar(&state.BufLen,"buffer",state.BufLen,"test buffer length")
    flags.Float64Var(&state.Speed,"speed",state.Speed,"scroll speed")
}

func (state *TestState) CheckFlags(flags *flag.FlagSet) (*TestConfig,bool) {
	ok := false
	ret := make(TestConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "w" { ok = true; ret.SetWidth( state.Width ) }
		if f.Name == "h" { ok = true; ret.SetHeight( state.Height ) }
		if f.Name == "mode" { ok = true; ret.SetBuffer( state.Buffer) }
		if f.Name == "speed" { ok = true; ret.SetSpeed( state.Speed ) }
		if f.Name == "buffer" { ok = true; ret.SetBufLen( state.BufLen ) }
	})
	return &ret,ok
}


func (state *TestState) Desc() string { return state.Config().Desc() }


func (state *TestState) Config() *TestConfig {
	ret := make(TestConfig)
	ret.SetWidth(state.Width)
	ret.SetHeight(state.Height)
	ret.SetBuffer(state.Buffer)
	ret.SetBufLen(state.BufLen)
	ret.SetSpeed(state.Speed)
	return &ret	
}


func (state *TestState) ApplyConfig(config *TestConfig) bool {
	changed := false
	if tmp,ok := config.Width();    ok { if state.Width    != tmp { changed = true }; state.Width = tmp }
	if tmp,ok := config.Height();   ok { if state.Height   != tmp { changed = true }; state.Height = tmp }
	if tmp,ok := config.Buffer();   ok { if state.Buffer   != tmp { changed = true }; state.Buffer = tmp }
	if tmp,ok := config.BufLen();   ok { if state.BufLen   != tmp { changed = true }; state.BufLen = tmp }
	if tmp,ok := config.Speed();    ok { if state.Speed    != tmp { changed = true }; state.Speed = tmp }
	return changed
}


