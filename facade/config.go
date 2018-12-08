
package facade

import (
//    "fmt"
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


func (config *Config) Mode()   (Mode,bool)              { ret,ok := (*config)["mode"].(string); return Mode(ret),ok }
func (config *Config) Font()   (gfx.FontConfig,bool)    { ret,ok := (*config)["font"].(string); return gfx.FontConfig(ret),ok }
func (config *Config) Debug()  (bool,bool)              { ret,ok := (*config)["debug"].(bool); return ret,ok }
func (config *Config) Grid()   (GridDelta,bool)        { ret,ok := (*config)["grid"].(map[string]interface{}); return GridDelta(ret),ok }
func (config *Config) Camera() (gfx.CameraConfig,bool)  { ret,ok := (*config)["camera"].(map[string]interface{}); return gfx.CameraConfig(ret),ok }
func (config *Config) Mask()   (gfx.MaskConfig,bool)    { ret,ok := (*config)["mask"].(map[string]interface{}); return gfx.MaskConfig(ret),ok }
 

func (config *Config) SetMode(val Mode)               { (*config)["mode"] = string(val) }
func (config *Config) SetFont(val gfx.FontConfig)     { (*config)["font"] = string(val) }
func (config *Config) SetDebug(val bool)              { (*config)["debug"] = val }
func (config *Config) SetGrid(val GridDelta)         { (*config)["grid"] = map[string]interface{}(val) }
func (config *Config) SetCamera(val gfx.CameraConfig) { (*config)["camera"] = map[string]interface{}(val) }
func (config *Config) SetMask(val gfx.MaskConfig)     { (*config)["mask"] = map[string]interface{}(val) }


var defaults = struct {
    debug bool        
}{
    true,
}





func NewConfig(mode Mode) *Config {
    ret := make(Config)
    ret.SetMode(mode)
    switch mode {
        case GRID: ret.SetGrid( *NewGridDelta() )    
    }
    ret.SetFont( *gfx.NewFontConfig() )
    ret.SetCamera( *gfx.NewCameraConfig() )
    ret.SetMask( *gfx.NewMaskConfig() )
    ret.SetDebug(defaults.debug)
    return &ret
}

func (config *Config) Clean() {
//    if config.Grid != nil { config.Grid.Clean() }     
}


func (config *Config) FlagSet() *flag.FlagSet {
    mode,_ := config.Mode()
    ret := flag.NewFlagSet(string(mode), flag.ExitOnError)
//    if grid,ok := config.Grid(); ok     { grid.AddFlags(ret) }
//    if lines,ok := config.Lines(); ok   { lines.AddFlags(ret) }
//    if test,ok := config.Test(); ok     { test.AddFlags(ret) }
    if font,ok := config.Font(); ok     { font.AddFlags(ret) }
    if camera,ok := config.Camera(); ok { camera.AddFlags(ret) }
    if mask,ok := config.Mask(); ok     { mask.AddFlags(ret) }
    
    
    ret.Bool("D",defaults.debug,"Draw Debug" )
    return ret
}



func (config *Config) Desc() string {
    ret := "conf["
    if mode,ok := config.Mode(); ok { ret += string(mode) }
    if font,ok := config.Font(); ok { ret += " " + font.Desc()   }
    if camera,ok := config.Camera(); ok  { ret += " " + camera.Desc() }
    if mask,ok := config.Mask(); ok  { ret += " " + mask.Desc() }
    if grid,ok := config.Grid(); ok { ret += " " + grid.Desc() }
    if debug,ok := config.Debug(); ok && debug { ret += " DEBUG" }
    ret += "]"
    return ret
}


