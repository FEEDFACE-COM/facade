
package font

import (
    "fmt"
    "io/ioutil"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    log "../log"
)


type Font struct {
    face string
    font *truetype.Font  
}


func NewFont() *Font {
    ret := &Font{}
//    ret.Configure( NewConfig() )
    return ret
}


func (font *Font) Configure(config *Config, directory string) {
    log.Debug("configure font: %s",config.Describe())
    err := font.loadFont(directory+config.Face+".ttc")
    if err != nil {
        log.Error("fail to config font %s: %s",config.Describe(),err)
        return
    }
    font.face = config.Face
    err = font.renderGridMap()
    if err != nil {
        log.Error("fail to render grid map %s: %s",font.face, err)    
    }
    
}


func (font *Font) loadFont(fontfile string) error {
    data, err := ioutil.ReadFile(fontfile  )
    if err != nil {
        return err
    }
    font.font,err = freetype.ParseFont(data)
    if err != nil {
        return err
    }
    log.Debug("load font file %s",fontfile)
    return nil
}

func (font *Font) renderGridMap() error {
    
    
    
    
    
    
    return nil
}


func (font *Font) Describe() string {
    return fmt.Sprintf("font[%s]",font.face)
}