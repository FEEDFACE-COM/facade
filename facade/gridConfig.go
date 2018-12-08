

package facade

import (
    "flag"    
    "fmt"
    "strings"
)


type GridDelta map[string]interface{}

func (delta *GridDelta) Width() (uint,bool) { ret,ok := (*delta)["width"].(float64); return uint(ret),ok }
func (delta *GridDelta) Height() (uint,bool) { ret,ok := (*delta)["height"].(float64); return uint(ret),ok }
func (delta *GridDelta) Downward() (bool,bool) { ret,ok := (*delta)["downward"].(bool); return ret,ok }
func (delta *GridDelta) Scroll() (bool,bool) { ret,ok := (*delta)["scroll"].(bool); return ret,ok }
func (delta *GridDelta) Speed() (float64,bool) { ret,ok := (*delta)["speed"].(float64); return ret,ok }
func (delta *GridDelta) Vert() (string,bool) { ret,ok := (*delta)["vert"].(string); return ret,ok }
func (delta *GridDelta) Frag() (string,bool) { ret,ok := (*delta)["frag"].(string); return ret,ok }
func (delta *GridDelta) Fill() (string,bool) { ret,ok := (*delta)["fill"].(string); return ret,ok }



func (delta *GridDelta) SetWidth(val uint) { (*delta)["width"] = float64(val) }
func (delta *GridDelta) SetHeight(val uint) { (*delta)["height"] = float64(val) }
func (delta *GridDelta) SetDownward(val bool) { (*delta)["downward"] = val }
func (delta *GridDelta) SetScroll(val bool) { (*delta)["scroll"] = val }
func (delta *GridDelta) SetSpeed(val float64) { (*delta)["speed"] = val }
func (delta *GridDelta) SetVert(val string) { (*delta)["vert"] = val }
func (delta *GridDelta) SetFrag(val string) { (*delta)["frag"] = val }
func (delta *GridDelta) SetFill(val string) { (*delta)["fill"] = val }


func NewGridDelta() *GridDelta {
    ret := make(GridDelta)
    return &ret
}

type GridConfig struct {
    width uint
    height uint
    
    downward bool
    scroll bool
    speed float64
    
    vert string
    frag string
    fill string
}

var GridDefaults = GridConfig {
    0,       //width
    8,       //height
    false,   //downward
    true,    //scroll
    0.4,     //speed
    "null",  //vert
    "null",  //frag
    "",      //fill
}



func (config *GridConfig) AddFlags(flags *flag.FlagSet) {

    flags.UintVar(&config.width,"w",config.width,"grid width")
    flags.UintVar(&config.height,"h",config.height,"grid height")
    flags.BoolVar(&config.downward,"d",config.downward,"downward")
    flags.BoolVar(&config.scroll,"s",config.scroll,"scroll")
    flags.Float64Var(&config.speed,"S",config.speed,"scroll speed")
    flags.StringVar(&config.vert,"V",config.vert,"vertex shader")
    flags.StringVar(&config.frag,"F",config.frag,"fragment shader")
    flags.StringVar(&config.fill,"fill",config.fill,"fill pattern")

}

func (delta *GridDelta) Desc() string { 
    ret := "grid["
    if tmp,ok := delta.Width(); ok { ret += fmt.Sprintf("%d",tmp) }
    ret += "x"
    if tmp,ok := delta.Height(); ok { ret += fmt.Sprintf("%d",tmp) }
    if tmp,ok := delta.Downward(); ok { if tmp { ret += "↓" } else { ret += "↑" } }
    ret += " "
    if tmp,ok := delta.Scroll(); ok { if tmp { ret += "→" } else { ret += "x" } }
    if tmp,ok := delta.Speed(); ok { ret += fmt.Sprintf("%.1f",tmp) }
    ret += " "
    if tmp,ok := delta.Vert(); ok { ret += tmp }
    ret += ","
    if tmp,ok := delta.Frag(); ok { ret += tmp }
    if tmp,ok := delta.Fill(); ok { ret += " " + tmp }
    ret += "]"
    return ret
}

func (config *GridConfig) Desc() string {
	delta := config.Delta()
	return delta.Desc()	
}

func (config *GridConfig) Delta() *GridDelta {
	ret := NewGridDelta()
	ret.SetWidth(config.width)
	ret.SetHeight(config.height)
	ret.SetDownward(config.downward)
	ret.SetScroll(config.scroll)
	ret.SetSpeed(config.speed)
	ret.SetVert(config.vert)
	ret.SetFrag(config.frag)
	ret.SetFill(config.fill)
	return ret	
}

func (config *GridConfig) ApplyDelta(delta *GridDelta) {
	if tmp,ok := delta.Width(); ok { config.width = tmp }
	if tmp,ok := delta.Height(); ok { config.height = tmp }
	if tmp,ok := delta.Downward(); ok { config.downward = tmp }
	if tmp,ok := delta.Scroll(); ok { config.scroll = tmp }
	if tmp,ok := delta.Speed(); ok { config.speed = tmp }
	if tmp,ok := delta.Vert(); ok { config.vert = tmp }
	if tmp,ok := delta.Frag(); ok { config.frag = tmp }
	if tmp,ok := delta.Fill(); ok { config.fill = tmp }
}



func (delta *GridDelta) Clean() {
    if vert,ok := delta.Vert(); ok { delta.SetVert( strings.Replace(vert,"/","",-1) ) }
    if frag,ok := delta.Frag(); ok { delta.SetFrag( strings.Replace(frag,"/","",-1) ) }
}

