package conf

import (
    "flag"    
    "fmt"
)


type LineConfig struct {
    LineCount uint
}




func NewLineConfig() *LineConfig {
    return &LineConfig{LineCount: 1}
}


func (config *LineConfig) AddFlags(flags *flag.FlagSet) {
//    flags.UintVar(&config.Height,"h",config.Height,"grid height")
    flags.UintVar(&config.LineCount,"l",config.LineCount,"line count")

}

func (config *LineConfig) Desc() string { return fmt.Sprintf("lines[%d]",config.LineCount) }


