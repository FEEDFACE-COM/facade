
package font

import (
    log "../log"
)


type Font struct {
    Config
    //texture data etc    
}


func NewFont() *Font {
    return &Font{Config: *NewConfig()}    
}


func (font *Font) Configure(config *Config) {
    log.Debug("configure font: %s",config.Describe())  
    font.Config = *config  
}