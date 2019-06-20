//
package facade

import(
	"strings"
    "flag"    
)


var DEFAULT_DIRECTORY = "~/src/gfx/facade"

var DEFAULT_MODE Mode = Mode_GRID


var Defaults = Config{
    Mode: DEFAULT_MODE,
    Debug: false,
}



func (config *Config) Desc() string {
    ret := "config["
    if config.GetSetMode() { ret += "mode[" + strings.ToLower( config.GetMode().String()) + "] " }

    if font:=config.GetFont(); font!=nil { ret += font.Desc() + " " }
    if cam:=config.GetCamera(); cam!=nil { ret += cam.Desc() + " " }
    if mask:=config.GetMask(); mask!=nil { ret += mask.Desc() + " " }

    if grid := config.GetGrid(); grid!=nil { ret += grid.Desc() + " " }
    if config.GetSetDebug() { if config.GetDebug() { ret += "DEBUG " } else { ret += "nobug " } }

    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
}

func (config *Config) AddFlags(flagset *flag.FlagSet) {
    flagset.BoolVar(&config.Debug,"D",Defaults.Debug,"debug draw")
    if grid:=config.GetGrid() ; grid!=nil { grid.AddFlags(flagset) }
    if font:=config.GetFont() ; font!=nil { font.AddFlags(flagset) }
    if  cam:=config.GetCamera(); cam!=nil { cam.AddFlags(flagset) }
    if mask:=config.GetMask() ; mask!=nil { mask.AddFlags(flagset) }
}

func (config *Config) VisitFlags(flagset *flag.FlagSet) {

	flagset.Visit( func(flg *flag.Flag) {
        switch flg.Name {
            case "D":     { config.SetDebug = true;     }
	   }
    })
    
    if grid:=config.GetGrid() ; grid!=nil {
        if ! grid.VisitFlags(flagset) { config.Grid = nil } // no flags used
    }
    
    if font:=config.GetFont() ; font!=nil { 
        if ! font.VisitFlags(flagset) { config.Font = nil } // no flags used
    }

    if cam:=config.GetCamera() ; cam!=nil {
        if ! cam.VisitFlags(flagset) { config.Camera = nil } // no flags used
    }
    
    if mask:=config.GetMask() ; mask!=nil { 
        if ! mask.VisitFlags(flagset) { config.Mask = nil } // no flags used
    }

}
