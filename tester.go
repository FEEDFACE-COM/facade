
package main

import (
    "os"
    "bufio"
    "image/png"
    "image"
    "errors"
    log "./log"
    conf "./conf"
    gfx "./gfx"
)

type Tester struct {font *gfx.Font}
func NewTester() *Tester { return &Tester{} }


func (tester *Tester) Configure(config *conf.Config) {
    
    tester.font = gfx.NewFont()
    tester.font.Configure(config.Font,conf.DIRECTORY)
    log.Info("got font %s",tester.font.Desc())
    
    
}

func (tester *Tester) testCharMap() error {
    
    var charmap *gfx.GlyphTexture
    var err error
    
    charmap, err = tester.font.RenderGlyphTexture()
    if err != nil {
        log.PANIC("fail to render glyphmap for %s: %s",tester.font.Desc(),err)
    }
    saveIMG(charmap.Texture, "charmap")
    
    return nil    
}


func (tester *Tester) testTextTex(str string) error {
    
    
    var texttex *gfx.TextTexture
    var err error
    texttex,err = tester.font.RenderTextTexture(str)
    
    if err != nil {
        log.Error("fail to generate texture for '%s': %s",str,err)
        return errors.New("fail to generate texture")     
    }
    
    saveIMG(texttex.Texture, "texttex")
    
    return nil    
}





func saveIMG(img *image.RGBA,outname string) {

    var outPath = conf.DIRECTORY + outname + ".png"
    
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


