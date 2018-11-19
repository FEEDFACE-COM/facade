

package facade

import (
    "flag"    
    "fmt"
    "strings"
)


type GridConfig map[string]interface{}

func (config *GridConfig) Width() (int,bool) { ret,ok := (*config)["width"].(float64); return int(ret),ok }
func (config *GridConfig) Height() (int,bool) { ret,ok := (*config)["height"].(float64); return int(ret),ok }
func (config *GridConfig) Downward() (bool,bool) { ret,ok := (*config)["downward"].(bool); return ret,ok }

func (config *GridConfig) Scroll() (bool,bool) { ret,ok := (*config)["scroll"].(bool); return ret,ok }
func (config *GridConfig) Speed() (float32,bool) { ret,ok := (*config)["speed"].(float64); return float32(ret),ok }

func (config *GridConfig) Vert() (string,bool) { ret,ok := (*config)["vert"].(string); return ret,ok }
func (config *GridConfig) Frag() (string,bool) { ret,ok := (*config)["frag"].(string); return ret,ok }
func (config *GridConfig) Fill() (string,bool) { ret,ok := (*config)["fill"].(string); return ret,ok }



func (config *GridConfig) SetWidth(val int) { (*config)["width"] = float64(val) }
func (config *GridConfig) SetHeight(val int) { (*config)["height"] = float64(val) }
func (config *GridConfig) SetDownward(val bool) { (*config)["downward"] = val }

func (config *GridConfig) SetScroll(val bool) { (*config)["scroll"] = val }
func (config *GridConfig) SetSpeed(val float32) { (*config)["speed"] = float64(val) }

func (config *GridConfig) SetVert(val string) { (*config)["vert"] = val }
func (config *GridConfig) SetFrag(val string) { (*config)["frag"] = val }
func (config *GridConfig) SetFill(val string) { (*config)["fill"] = val }


func NewGridConfig() *GridConfig {
    ret := make(GridConfig)
    ret.SetWidth(0)
    ret.SetHeight(8)
    ret.SetScroll(true)
    ret.SetSpeed(0.4)
    ret.SetVert("null")
    ret.SetFrag("null")
    return &ret
}




func (config *GridConfig) AddFlags(flags *flag.FlagSet) {
//    flags.UintVar(&config.Width,"w",config.Width,"grid width")
//    flags.UintVar(&config.Height,"h",config.Height,"grid height")
//    flags.BoolVar(&config.Downward,"d",config.Downward,"downward")
//    flags.BoolVar(&config.Scroll,"s",config.Scroll,"scroll")
//    flags.Float64Var(&config.Speed,"S",config.Speed,"scroll speed")
//    flags.StringVar(&config.Vert,"V",config.Vert,"vertex shader")
//    flags.StringVar(&config.Frag,"F",config.Frag,"fragment shader")
//    flags.StringVar(&config.Fill,"fill",config.Fill,"fill pattern")
}

func (config *GridConfig) Desc() string { 
    ret := "grid["
    if tmp,ok := config.Width(); ok { ret += fmt.Sprintf("%d",tmp) }
    ret += "x"
    if tmp,ok := config.Height(); ok { ret += fmt.Sprintf("%d",tmp) }
    if tmp,ok := config.Downward(); ok { if tmp { ret += "↓" } else { ret += "↑" } }
    ret += " "
    if tmp,ok := config.Scroll(); ok { if tmp { ret += "→" } else { ret += "x" } }
    if tmp,ok := config.Speed(); ok { ret += fmt.Sprintf("%.1f",tmp) }
    ret += " "
    if tmp,ok := config.Vert(); ok { ret += tmp }
    ret += ","
    if tmp,ok := config.Frag(); ok { ret += tmp }
    if tmp,ok := config.Fill(); ok { ret += " " + tmp }
    ret += "]"
    return ret
}




func (config *GridConfig) Clean() {
    if vert,ok := config.Vert(); ok { config.SetVert( strings.Replace(vert,"/","",-1) ) }
    if frag,ok := config.Frag(); ok { config.SetFrag( strings.Replace(frag,"/","",-1) ) }
}

