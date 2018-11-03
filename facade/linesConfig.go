
package facade

import (
    "flag"    
    "fmt"
)


type LinesConfig struct {
    height uint
}

func (config *LinesConfig) Height()    uint    { return config.height    }



func NewLinesConfig() *LinesConfig {
    return &LinesConfig{height: 1}
}


func (config *LinesConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.height,"h",config.height,"line count")

}

func (config *LinesConfig) Desc() string { return fmt.Sprintf("lines[%d]",config.height) }


