
package facade

import (
    "fmt"
    "flag"
    gfx "../gfx"
)

var DEFAULT_MODE Mode = GRID
var DEFAULT_DIRECTORY = "~/src/gfx/facade"



type Mode string
const (
    GRID  Mode = "grid"
    LINES Mode = "lines"
    WORD  Mode = "word"
    CHAR  Mode = "char"   
    TEST  Mode = "test" 
)

var Modes = []Mode{GRID,LINES,TEST}





type RawText string


type Config map[string]interface{}


func (config *Config) Encode() string {

    return ""   
    
}



func (config *Config) Mode()   (Mode             ,bool ) { tmp,ok := (*config)["mode"]; ret,ok2 := tmp.(Mode); return ret,ok&&ok2 }
func (config *Config) Grid()   (*GridConfig      ,bool ) { tmp,ok := (*config)["grid"]; ret,ok2 := tmp.(*GridConfig); return ret,ok&&ok2 }
func (config *Config) Lines()  (*LinesConfig     ,bool ) { tmp,ok := (*config)["lines"]; ret,ok2 := tmp.(*LinesConfig); return ret,ok&&ok2 }
func (config *Config) Test()   (*TestConfig      ,bool ) { tmp,ok := (*config)["test"]; ret,ok2 := tmp.(*TestConfig); return ret,ok&&ok2 }
func (config *Config) Font()   (*gfx.FontConfig  ,bool ) { tmp,ok := (*config)["font"]; ret,ok2 := tmp.(*gfx.FontConfig); return ret,ok&&ok2 }
func (config *Config) Camera() (*gfx.CameraConfig,bool ) { tmp,ok := (*config)["camera"]; ret,ok2 := tmp.(*gfx.CameraConfig); return ret,ok&&ok2 }
func (config *Config) Mask()   (*gfx.MaskConfig  ,bool ) { tmp,ok := (*config)["mask"]; ret,ok2 := tmp.(*gfx.MaskConfig); return ret,ok&&ok2 }
func (config *Config) Debug()  ( bool            ,bool ) { tmp,ok := (*config)["debug"]; ret,ok2 := tmp.(bool); return ret,ok&&ok2 }





func NewConfig(mode Mode) *Config {
    ret :=  Config( make( map[string]interface{} ) )
    ret["mode"] = mode
    switch mode {
        case GRID:  ret["grid"] = NewGridConfig()
        case LINES: ret["lines"] = NewLinesConfig() 
        case TEST:  ret["test"] = NewTestConfig() 
    }
    ret["font"] =   gfx.NewFontConfig()
    ret["camera"] = gfx.NewCameraConfig()
    ret["mask"] =   gfx.NewMaskConfig()
    return &ret
}

func (config *Config) Clean() {
    if tmp,ok := config.Grid(); ok { tmp.Clean() }
}


func (config *Config) FlagSet() *flag.FlagSet {
    mode,ok := config.Mode(); if !ok { return nil }
    ret := flag.NewFlagSet(string(mode), flag.ExitOnError)
//    if config.grid   != nil { config.grid.AddFlags(ret) }
//    if config.lines  != nil { config.lines.AddFlags(ret) }
//    if config.test   != nil { config.test.AddFlags(ret) }
//    if config.font   != nil { config.font.AddFlags(ret) }
//    if config.camera != nil { config.camera.AddFlags(ret) }
//    if config.mask   != nil { config.mask.AddFlags(ret) }
//    ret.BoolVar(&config.debug,"D",config.debug,"Draw Debug" )
    
    return ret
}



func (config *Config) Desc() string {
    mode,_ := config.Mode();
    ret := fmt.Sprintf("conf[%s",string(mode))
    if tmp,ok := config.Grid()   ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Lines()  ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Test()   ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Font()   ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Camera() ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Mask()   ; ok  { ret += " " + tmp.Desc() }
    if tmp,ok := config.Debug(); ok&&tmp { ret += " DEBUG" }
    ret += "]"
    return ret
}


