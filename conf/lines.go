package conf

import (
    "flag"    
    "fmt"
)


type LinesConfig struct {
    LineCount uint
}




func NewLinesConfig() *LinesConfig {
    return &LinesConfig{LineCount: 1}
}


func (config *LinesConfig) AddFlags(flags *flag.FlagSet) {
//    flags.UintVar(&config.Height,"h",config.Height,"grid height")
    flags.UintVar(&config.LineCount,"l",config.LineCount,"line count")

}

func (config *LinesConfig) Desc() string { return fmt.Sprintf("lines[%d]",config.LineCount) }


