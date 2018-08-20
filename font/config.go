
package font

import (
    "fmt"
    "flag"
)


type Config struct {
    Face string
}

func (config *Config) AddFlags(flags *flag.FlagSet) {
    flags.StringVar(&config.Face,"font",config.Face,"use fontface `font`" )
}

func NewConfig() *Config { return &Config{Face: "Monaco"} }


func (config *Config) Describe() string { return fmt.Sprintf("font[%s]",config.Face) }

