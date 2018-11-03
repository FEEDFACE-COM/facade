
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


type Config struct {
    mode   Mode

    grid   *GridConfig
    lines  *LinesConfig
    test   *TestConfig

    font   *gfx.FontConfig
    camera *gfx.CameraConfig    
    mask   *gfx.MaskConfig

    debug  bool
}




func (config *Config) Mode()   Mode              { return config.mode   }
func (config *Config) Grid()   *GridConfig       { return config.grid   }
func (config *Config) Lines()  *LinesConfig      { return config.lines  }
func (config *Config) Test()   *TestConfig       { return config.test   }
func (config *Config) Font()   *gfx.FontConfig   { return config.font   }
func (config *Config) Camera() *gfx.CameraConfig { return config.camera }
func (config *Config) Mask()   *gfx.MaskConfig   { return config.mask   }
func (config *Config) Debug()  bool              { return config.debug  }



func NewConfig(mode Mode) *Config {
    ret := &Config{mode: mode}
    if mode == GRID  { ret.grid  = NewGridConfig() }
    if mode == LINES { ret.lines = NewLinesConfig() }
    if mode == TEST  { ret.test  = NewTestConfig() }
    ret.font =   gfx.NewFontConfig()
    ret.camera = gfx.NewCameraConfig()
    ret.mask =   gfx.NewMaskConfig()
    return ret
}

func (config *Config) Clean() {
    if config.grid != nil { config.grid.Clean() }     
}


func (config *Config) FlagSet() *flag.FlagSet {
    ret := flag.NewFlagSet(string(config.mode), flag.ExitOnError)
    if config.grid   != nil { config.grid.AddFlags(ret) }
    if config.lines  != nil { config.lines.AddFlags(ret) }
    if config.test   != nil { config.test.AddFlags(ret) }
    if config.font   != nil { config.font.AddFlags(ret) }
    if config.camera != nil { config.camera.AddFlags(ret) }
    if config.mask   != nil { config.mask.AddFlags(ret) }
    ret.BoolVar(&config.debug,"D",config.debug,"Draw Debug" )
    
    return ret
}



func (config *Config) Desc() string {
    ret := fmt.Sprintf("conf[%s",string(config.mode))
    if config.grid   != nil { ret += " " + config.grid.Desc() }
    if config.lines  != nil { ret += " " + config.lines.Desc() }
    if config.test   != nil { ret += " " + config.test.Desc() }
    if config.font   != nil { ret += " " + config.font.Desc() }
    if config.camera != nil { ret += " " + config.camera.Desc() }
    if config.mask   != nil { ret += " " + config.mask.Desc() }
    if config.debug { ret += " DEBUG" }
    ret += "]"
    return ret
}


