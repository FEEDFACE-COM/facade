
package conf

import (
    "fmt"
    "flag"
    gfx "../gfx"
)


type Config struct {
    Mode Mode
    Grid *GridConfig
    Lines *LinesConfig
    Test *TestConfig
    Font *gfx.FontConfig
    Camera *gfx.CameraConfig    
    Mask *gfx.MaskConfig
    Debug bool
}



type RawText string


type Mode string
const (
    GRID  Mode = "grid"
    LINES Mode = "lines"
    WORD  Mode = "word"
    CHAR  Mode = "char"   
    TEST  Mode = "test" 
)

var Modes = []Mode{GRID,LINES,TEST}

var DEFAULT_MODE Mode = TEST


var DIRECTORY = "/home/folkert/src/gfx/facade/asset/"

func NewConfig(mode Mode) *Config {
    ret := &Config{Mode: mode}
    if mode == GRID  { ret.Grid  = NewGridConfig() }
    if mode == LINES { ret.Lines = NewLinesConfig() }
    if mode == TEST  { ret.Test  = NewTestConfig() }
    ret.Font =   gfx.NewFontConfig()
    ret.Camera = gfx.NewCameraConfig()
    ret.Mask =   gfx.NewMaskConfig()
//˚    ret.Debug = true
    return ret
}


func (config *Config) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(config.Mode), flag.ExitOnError)
    if config.Grid   != nil { config.Grid.AddFlags(ret) }
    if config.Lines  != nil { config.Lines.AddFlags(ret) }
    if config.Test   != nil { config.Test.AddFlags(ret) }
    if config.Font   != nil { config.Font.AddFlags(ret) }
    if config.Camera != nil { config.Camera.AddFlags(ret) }
    if config.Mask   != nil { config.Mask.AddFlags(ret) }
    ret.BoolVar(&config.Debug,"D",config.Debug,"Draw Debug" )
    return ret
}



func (config *Config) Desc() string {
    ret := fmt.Sprintf("conf[%s",string(config.Mode))
    if config.Grid   != nil { ret += " " + config.Grid.Desc() }
    if config.Lines  != nil { ret += " " + config.Lines.Desc() }
    if config.Test   != nil { ret += " " + config.Test.Desc() }
    if config.Font   != nil { ret += " " + config.Font.Desc() }
    if config.Camera != nil { ret += " " + config.Camera.Desc() }
    if config.Mask   != nil { ret += " " + config.Mask.Desc() }
    if config.Debug { ret += " DEBUG" }
    ret += "]"
    return ret
}


