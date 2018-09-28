

package facade

import (
    "flag"    
    "fmt"
)


type TestConfig struct {
}




func NewTestConfig() *TestConfig {
    return &TestConfig{}
}


func (config *TestConfig) AddFlags(flags *flag.FlagSet) {
}

func (config *TestConfig) Desc() string { return fmt.Sprintf("test[]") }






