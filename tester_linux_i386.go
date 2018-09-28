
// +build linux arm

package main

import (
    "fmt"
    "os"
    "bufio"
    "image/png"
    "image"
//    "errors"
    log "./log"
    gfx "./gfx"
    facade "./facade"
)

type Tester struct {font *gfx.Font; name string; directory string}
func NewTester(directory string) *Tester { return &Tester{directory: directory} }


func (tester *Tester) Configure(config *facade.Config) {
    tester.name = config.Font.Name
    var err error
    gfx.SetFontDirectory(tester.directory)
    tester.font,err = gfx.GetFont(config.Font)
    if err != nil {
        log.PANIC("fail loading font %s: %s",config.Font.Name,err)
    }
    tester.font.Init()
    tester.font.Configure(config.Font)
    log.Info("got font %s",tester.font.Desc())
    
    
}

func (tester *Tester) testCharMap() (*image.RGBA,error) {
    
    ret, err := tester.font.RenderMapRGBA()
    if err != nil {
        log.Error("fail to render glyphmap for %s: %s",tester.font.Desc(),err)
        return nil,err
    }
    return ret, nil
}


func (tester *Tester) testTextTex(str string) (*image.RGBA,error) {
    
    
    ret,err := tester.font.RenderTextRGBA(str)
    
    if err != nil {
        log.Error("fail to generate texture for '%s': %s",str,err)
        return nil,err
    }
    return ret, nil
}



func (tester *Tester) Test(str string) {
    test0,_ := tester.testCharMap()
    test1,_ := tester.testTextTex(str)
    SaveRGBA(test0,fmt.Sprintf("map-%s",tester.directory+tester.name+".png"))
    SaveRGBA(test1,fmt.Sprintf("text-%s-%s",tester.name,str))
}


func SaveRGBA(img *image.RGBA,outPath string)  {

    if img == nil {
        log.Error("no image to save at "+outPath)
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


