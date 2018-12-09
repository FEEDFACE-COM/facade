
package facade

import (
//    "fmt"
	"strings"
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

type State struct {
	Mode    Mode
	Grid   *GridState
	Font   *gfx.FontState
	Camera *gfx.CameraState
	Mask   *gfx.MaskState
	Debug  bool	
}

var Defaults = State{
	Mode:  DEFAULT_MODE,
    Debug:         false,
}


const (
	facadeMode   = "mode"
	facadeGrid   = "grid"
	facadeFont   = "font"
	facadeCamera = "camera"	
	facadeMask   = "mask"
	facadeDebug  = "debug"
)

type Config gfx.Config

func (config *Config) Mode()   (Mode,bool)              { ret,ok := (*config)[facadeMode].(string);                   return             Mode(ret),ok }
func (config *Config) Grid()   (GridConfig,bool)        { ret,ok := (*config)[facadeGrid].(map[string]interface{});   return       GridConfig(ret),ok }
func (config *Config) Font()   (gfx.FontConfig,bool)    { ret,ok := (*config)[facadeFont].(map[string]interface{});   return   gfx.FontConfig(ret),ok }
func (config *Config) Camera() (gfx.CameraConfig,bool)  { ret,ok := (*config)[facadeCamera].(map[string]interface{}); return gfx.CameraConfig(ret),ok }
func (config *Config) Mask()   (gfx.MaskConfig,bool)    { ret,ok := (*config)[facadeMask].(map[string]interface{});   return   gfx.MaskConfig(ret),ok }
func (config *Config) Debug()  (bool,bool)              { ret,ok := (*config)[facadeDebug].(bool);                    return                  ret, ok }

func (config *Config) SetMode(val Mode)               { (*config)[facadeMode] = string(val) }
func (config *Config) SetGrid(val GridConfig)         { (*config)[facadeGrid] = map[string]interface{}(val) }
func (config *Config) SetFont(val gfx.FontConfig)     { (*config)[facadeFont] = map[string]interface{}(val) }
func (config *Config) SetCamera(val gfx.CameraConfig) { (*config)[facadeCamera] = map[string]interface{}(val) }
func (config *Config) SetMask(val gfx.MaskConfig)     { (*config)[facadeMask] = map[string]interface{}(val) }
func (config *Config) SetDebug(val bool)              { (*config)[facadeDebug] = val }




func NewState(mode Mode) *State {
	ret := Defaults
	switch mode {
		case GRID, LINES, TEST:
			ret.Mode = mode
	}
	switch ret.Mode {
		case GRID: ret.Grid = &GridState{}
	}
	ret.Font = &gfx.FontState{}
	ret.Camera = &gfx.CameraState{}
	ret.Mask = &gfx.MaskState{}
	return &ret	
}



//func (delta *Delta) Clean() {
////    if delta.Grid() != nil { delta.Grid().Clean() }     //???????
//}


func (state *State) AddFlags(flags *flag.FlagSet) {
    if state.Grid != nil { state.Grid.AddFlags(flags) }
    if state.Font != nil { state.Font.AddFlags(flags) }
    if state.Camera != nil { state.Camera.AddFlags(flags) }
    if state.Mask != nil { state.Mask.AddFlags(flags) }
    flags.BoolVar(&state.Debug,"D",state.Debug,"Draw Debug" )
}	


func (state *State) CheckFlags(flags *flag.FlagSet) *Config {	
	ret := make(Config)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "D" { ret.SetDebug( state.Debug ) }
	})
	if state.Font != nil   { 
		fontConfig,ok := state.Font.CheckFlags(flags) 
		if ok { ret.SetFont( *fontConfig ) }
	}
	if state.Camera != nil { 
		cameraConfig,ok := state.Camera.CheckFlags(flags) 
		if ok { ret.SetCamera( *cameraConfig ) }
	}
	if state.Mask != nil   { 
		maskConfig,ok := state.Mask.CheckFlags(flags) 
		if ok { ret.SetMask( *maskConfig ) }
	}
	if state.Grid != nil { 
		gridConfig,ok := state.Grid.CheckFlags(flags)  
		if ok { ret.SetGrid( *gridConfig ) }
	}
	return &ret
}


func (config *Config) Desc() string {
    ret := "conf["
    if mode,ok := config.Mode(); ok { ret += string(mode) + " " }
    if grid,ok := config.Grid(); ok { ret += grid.Desc() + " " }
    if font,ok := config.Font(); ok { ret += font.Desc() + " "  }
    if camera,ok := config.Camera(); ok  { ret += camera.Desc() + " " }
    if mask,ok := config.Mask(); ok  { ret += mask.Desc() + " " }
    if debug,ok := config.Debug(); ok {
	    if debug { ret += "DEBUG " } else { ret += "nobug" }
	}
    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
}


func (state *State) Desc() string { return state.Config().Desc() }



func (state *State) Config() *Config {
	ret := make(Config)
	ret.SetMode(state.Mode)	
	// etc...
	ret.SetDebug(state.Debug)
	return &ret
}

func (state *State) ApplyConfig(config *Config) {
	if tmp,ok := config.Mode();  ok { state.Mode = tmp }
	// etc...
	if tmp,ok := config.Debug(); ok { state.Debug = tmp }	
}


