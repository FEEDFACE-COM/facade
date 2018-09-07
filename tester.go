
package main

import (
    "fmt"
    "os"
    "bufio"
    "image/png"
    "image"
    "errors"
    log "./log"
    conf "./conf"
    font "./font"
    "golang.org/x/image/math/fixed"
)

type Tester struct {font *font.Font}
func NewTester() *Tester { return &Tester{} }


func (tester *Tester) Configure(config *conf.Config) {
    
    tester.font = font.NewFont()
    tester.font.Configure(config.Font,conf.DIRECTORY)
    log.Info("got font %s",tester.font.Describe())
    
    
}

func mask(in fixed.Int26_6) int { return (0xffffffff & int(in)) }
func str(s string, in fixed.Int26_6) string { return fmt.Sprintf("%s\t0x%08x\t%v\t %d ≤ %d ≤ %d",s,mask(in),in.String(),in.Floor(),in.Round(),in.Ceil()) }

func (tester *Tester) testFixed(config *conf.Config) error {
    log.Info("testing fixed...")    
    
    one_one_fourth := fixed.Int26_6( 1<<6 + 1<<4 )
    one := fixed.Int26_6(1 << 6)
    two := fixed.Int26_6(2 << 6)
    three := fixed.Int26_6(3 << 6)
    four := fixed.Int26_6(4 << 6 )

    half := fixed.Int26_6( 0x20 )
    fourth := fixed.Int26_6( 0x10 )

    four_halves := half.Mul(four)
    four_fourths := fourth.Mul(four)
    fourth_fours := four.Mul(fourth)
    five := one_one_fourth.Mul(four)
    
    log.Info(str("¼",fourth))
    log.Info(str("½",half))
    log.Info(str("1",one))
    log.Info(str("¼x4",four_fourths))
    log.Info(str("4x¼",fourth_fours))
    log.Info(str("1¼",one_one_fourth))
    log.Info(str("2",two))
    log.Info(str("4x½",four_halves))
    log.Info(str("3",three))
    log.Info(str("4",four))
    log.Info(str("4x1¼",five))
    
    return nil
}

func (tester *Tester) testCharMap() error {
    
    var charmap *font.GlyphTexture
    var err error
    
    charmap, err = tester.font.RenderGlyphTexture()
    if err != nil {
        log.PANIC("fail to render glyphmap for %s: %s",tester.font.Describe(),err)
    }
    saveIMG(charmap.Texture, "charmap")
    
    return nil    
}


func (tester *Tester) testTextTex(str string) error {
    
    
    var texttex *font.TextTexture
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


