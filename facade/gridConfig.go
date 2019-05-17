

package facade

import (
    "flag"    
    "fmt"
    gfx "../gfx"
    "strings"
)





const (
	gridWidth    = "width"
	gridHeight   = "height"
	gridDownward = "downward"
	gridScroll   = "scroll"
	gridSpeed    = "speed"
	gridBufLen   = "buflen"
	gridTerm     = "term"
	gridVert     = "vert"
	gridFrag     = "frag"
	gridFill     = "fill"
)

type GridConfig gfx.Config



func (config *GridConfig) Width()     (uint,bool) { ret,ok := (*config)[gridWidth   ].(float64); return uint(ret), ok }
func (config *GridConfig) Height()    (uint,bool) { ret,ok := (*config)[gridHeight  ].(float64); return uint(ret), ok }
func (config *GridConfig) Downward()  (bool,bool) { ret,ok := (*config)[gridDownward].(bool);    return      ret , ok }
func (config *GridConfig) Scroll()    (bool,bool) { ret,ok := (*config)[gridScroll  ].(bool);    return      ret , ok }
func (config *GridConfig) Speed()  (float64,bool) { ret,ok := (*config)[gridSpeed   ].(float64); return      ret , ok }
func (config *GridConfig) BufLen()    (uint,bool) { ret,ok := (*config)[gridBufLen  ].(float64); return uint(ret), ok }
func (config *GridConfig) Term()      (bool,bool) { ret,ok := (*config)[gridTerm    ].(bool);    return      ret , ok }
func (config *GridConfig) Vert()    (string,bool) { ret,ok := (*config)[gridVert    ].(string);  return      ret , ok }
func (config *GridConfig) Frag()    (string,bool) { ret,ok := (*config)[gridFrag    ].(string);  return      ret , ok }
func (config *GridConfig) Fill()    (string,bool) { ret,ok := (*config)[gridFill    ].(string);  return      ret , ok }

func (config *GridConfig) SetWidth(   val uint)    { (*config)[gridWidth]    = float64(val) }
func (config *GridConfig) SetHeight(  val uint)    { (*config)[gridHeight]   = float64(val) }
func (config *GridConfig) SetDownward(val bool)    { (*config)[gridDownward] =         val  }
func (config *GridConfig) SetScroll(  val bool)    { (*config)[gridScroll]   =         val  }
func (config *GridConfig) SetSpeed(   val float64) { (*config)[gridSpeed]    =         val  }
func (config *GridConfig) SetBufLen(  val uint)    { (*config)[gridBufLen]   = float64(val) }  
func (config *GridConfig) SetTerm(    val bool)    { (*config)[gridTerm]     =         val  }
func (config *GridConfig) SetVert(    val string)  { (*config)[gridVert]     =         val  }
func (config *GridConfig) SetFrag(    val string)  { (*config)[gridFrag]     =         val  }
func (config *GridConfig) SetFill(    val string)  { (*config)[gridFill]     =         val  }


func (config *GridConfig) Desc() string { 
    ret := "grid["
    {
    	w,wok := config.Width(); 
	    h,hok := config.Height();
	    if wok { ret += fmt.Sprintf("%d",w) }
    	if wok || hok { ret += "x" }
	    if hok { ret += fmt.Sprintf("%d",h) }
	    if wok || hok { ret += " " }
	}
    
    if tmp,ok := config.Downward(); ok { if tmp { ret += "↓ " } else { ret += "↑ " } }
    {
		s,sok := config.Scroll();
		p,pok := config.Speed();
		if sok { if s { ret += "→" } else { ret += "↛" } }
		if pok { ret += fmt.Sprintf("%.1f",p) }
		if pok || sok { ret += " " }
	}
	{
		v, vok := config.Vert()
		f, fok := config.Frag()
		if vok { ret += v }
		if vok || fok { ret += "," }
		if fok { ret += f }
		if vok || fok { ret += " " }	
	}
    if tmp,ok := config.Fill(); ok { ret += tmp + " " } 
    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
}


func (config *GridConfig) ApplyConfig(cfg *GridConfig) {
	if tmp,ok := cfg.Width(); ok { config.SetWidth(tmp) }	
	if tmp,ok := cfg.Height(); ok { config.SetHeight(tmp) }	
	if tmp,ok := cfg.Downward(); ok { config.SetDownward(tmp) }	
	if tmp,ok := cfg.Scroll(); ok { config.SetScroll(tmp) }	
	if tmp,ok := cfg.Speed(); ok { config.SetSpeed(tmp) }	
	if tmp,ok := cfg.BufLen(); ok { config.SetBufLen(tmp) }	
	if tmp,ok := cfg.Term(); ok { config.SetTerm(tmp) }	
	if tmp,ok := cfg.Vert(); ok { config.SetVert(tmp) }	
	if tmp,ok := cfg.Frag(); ok { config.SetFrag(tmp) }	
	if tmp,ok := cfg.Fill(); ok { config.SetFill(tmp) }	
}





type GridState struct {
    Width uint
    Height uint
    
    Downward bool
    Scroll bool
    Speed float64
    
    BufLen uint
    Term bool
    
    Vert string
    Frag string
    Fill string
}

var GridDefaults = GridState{
    Width:        0,
    Height:       8,
    Downward: false,
    Scroll:    true,
    Speed:      0.4,
    BufLen:       2,
    Term:     false,
    Vert:     "def",
    Frag:     "def",
    Fill:        "",
}

func (state *GridState) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&state.Width,"w",state.Width,"grid width")
    flags.UintVar(&state.Height,"h",state.Height,"grid height")
    flags.BoolVar(&state.Downward,"d",state.Downward,"downward")
    flags.BoolVar(&state.Scroll,"s",state.Scroll,"scroll")
    flags.Float64Var(&state.Speed,"S",state.Speed,"scroll speed")
    flags.UintVar(&state.BufLen,"l",state.BufLen,"buffer lines")
    flags.BoolVar(&state.Term,"t",state.Term,"ansi terminal")
    flags.StringVar(&state.Vert,"V",state.Vert,"vertex shader")
    flags.StringVar(&state.Frag,"F",state.Frag,"fragment shader")
    flags.StringVar(&state.Fill,"fill",state.Fill,"fill pattern")
}


func (state *GridState) CheckFlags(flags *flag.FlagSet) (*GridConfig,bool) {
	ok := false
	ret := make(GridConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "w" { ok = true; ret.SetWidth( state.Width ) }
		if f.Name == "h" { ok = true; ret.SetHeight( state.Height ) }
		if f.Name == "d" { ok = true; ret.SetDownward( state.Downward ) }
		if f.Name == "s" { ok = true; ret.SetScroll( state.Scroll ) }
		if f.Name == "S" { ok = true; ret.SetSpeed( state.Speed ) }
		if f.Name == "l" { ok = true; ret.SetBufLen( state.BufLen ) }
		if f.Name == "t" { ok = true; ret.SetTerm( state.Term ) }
		if f.Name == "V" { ok = true; ret.SetVert( state.Vert ) }
		if f.Name == "F" { ok = true; ret.SetFrag( state.Frag ) }
		if f.Name == "fill" { ok = true; ret.SetFill( state.Fill ) }
	})
	return &ret,ok
}


func (state *GridState) Desc() string { return state.Config().Desc() }



//func (delta *GridDelta) Clean() {
//    if vert,ok := delta.Vert(); ok { delta.SetVert( strings.Replace(vert,"/","",-1) ) }
//    if frag,ok := delta.Frag(); ok { delta.SetFrag( strings.Replace(frag,"/","",-1) ) }
//}





func (state *GridState) Config() *GridConfig {
	ret := make(GridConfig)
	ret.SetWidth(state.Width)
	ret.SetHeight(state.Height)
	ret.SetDownward(state.Downward)
	ret.SetScroll(state.Scroll)
	ret.SetSpeed(state.Speed)
	ret.SetBufLen(state.BufLen)
	ret.SetTerm(state.Term)
	ret.SetVert(state.Vert)
	ret.SetFrag(state.Frag)
	ret.SetFill(state.Fill)
	return &ret	
}

func (state *GridState) ApplyConfig(config *GridConfig) bool {
	changed := false
	if tmp,ok := config.Width();    ok { if state.Width    != tmp { changed = true }; state.Width = tmp }
	if tmp,ok := config.Height();   ok { if state.Height   != tmp { changed = true }; state.Height = tmp }
	if tmp,ok := config.Downward(); ok { if state.Downward != tmp { changed = true }; state.Downward = tmp }
	if tmp,ok := config.Scroll();   ok { if state.Scroll   != tmp { changed = true }; state.Scroll = tmp }
	if tmp,ok := config.Speed();    ok { if state.Speed    != tmp { changed = true }; state.Speed = tmp }
	if tmp,ok := config.BufLen();   ok { if state.BufLen   != tmp { changed = true }; state.BufLen = tmp }
	if tmp,ok := config.Term();     ok { if state.Term     != tmp { changed = true }; state.Term = tmp }
	if tmp,ok := config.Vert();     ok { if state.Vert     != tmp { changed = true }; state.Vert = tmp }
	if tmp,ok := config.Frag();     ok { if state.Frag     != tmp { changed = true }; state.Frag = tmp }
	if tmp,ok := config.Fill();     ok { if state.Fill     != tmp { changed = true }; state.Fill = tmp }
	return changed
}







