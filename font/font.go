
package font

import (
    "fmt"
    log "../log"
)


type Font struct {
    face string
    //texture data etc    
}


func NewFont() *Font {
    ret := &Font{}
    ret.Configure( NewConfig() )
    return ret
}


func (font *Font) Configure(config *Config) {
    log.Debug("configure font: %s",config.Describe())  
    font.face = config.Face
}

func (font *Font) Describe() string {
    return fmt.Sprintf("font[%s]",font.face)
}