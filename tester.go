
package main

import (
    "os"
    "bufio"
    "image/png"
    "image"
//    "errors"
    log "./log"
    conf "./conf"
    gfx "./gfx"
)

type Tester struct {font *gfx.Font}
func NewTester() *Tester { return &Tester{} }


func (tester *Tester) Configure(config *conf.Config) {
    
    tester.font = gfx.NewFont(nil)
    tester.font.Configure(config.Font,conf.DIRECTORY)
    log.Info("got font %s",tester.font.Desc())
    
    
}

func (tester *Tester) testCharMap() (*image.RGBA,error) {
    
    var charmap *gfx.GlyphTexture
    var err error
    
    charmap, err = tester.font.RenderGlyphTexture()
    if err != nil {
        log.Error("fail to render glyphmap for %s: %s",tester.font.Desc(),err)
        return nil,err
    }
    return charmap.Texture, nil
}


func (tester *Tester) testTextTex(str string) (*image.RGBA,error) {
    
    
    var texttex *gfx.TextTexture
    var err error
    texttex,err = tester.font.RenderTextTexture(str)
    
    if err != nil {
        log.Error("fail to generate texture for '%s': %s",str,err)
        return nil,err
    }
    return texttex.Texture, nil
}





func SaveRGBA(img *image.RGBA,outname string)  {

    var outPath = conf.DIRECTORY + "/out/" + outname + ".png"
    
    if img == nil {
        log.Error("nil image not saved at "+outPath)
        return 
    }
    
    outFile, err := os.Create(outPath)
    if err != nil {
        log.PANIC("fail to create file %s: %s",outPath,err)
    }
    defer outFile.Close()
    
    writer := bufio.NewWriter(outFile)
    if err := png.Encode(writer, img); err != nil {
        log.PANIC("fail to encode image to %s: %s",err,err)
    }
    
    writer.Flush()
    log.Info("wrote image to %s",outPath)
    
}


