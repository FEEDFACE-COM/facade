package conf

import (
    "flag"    
    "fmt"
)


type LinesConfig struct {
    Height uint
}




func NewLinesConfig() *LinesConfig {
    return &LinesConfig{Height: 1}
}


func (config *LinesConfig) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&config.Height,"h",config.Height,"line count")

}

func (config *LinesConfig) Desc() string { return fmt.Sprintf("lines[%d]",config.Height) }


