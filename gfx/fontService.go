
package gfx

import (
    "fmt"
    "strings"
    "io/ioutil"
    "encoding/base64"
    "image"
    "image/color"
    log "../log"       
    
)


const DEBUG_FONTSERVICE = true



var ForegroundColor = image.White
var BackgroundColor = image.Transparent
var DebugColor = image.NewUniform( color.RGBA{R: 255, G: 0, B: 0, A: 255} )




type FontService struct {
    
    fonts map[string]*Font


    directory string
    
    
}



func NewFontService(directory string) *FontService {
    ret := &FontService{directory: directory}
    ret.fonts = make( map[string]*Font )
    return ret
}


func (service *FontService) GetFont(name string) (*Font,error) {
    name = strings.ToLower(name)
    ret := service.fonts[ name ]
    if ret == nil {
        return nil, log.NewError("no font for name %s",name)
    }
    return ret, nil
}

func (service *FontService) LoadFont(name string) error {
    name = strings.ToLower(name)
    

    var err error
    var data []byte = []byte{}

    
    if font:=service.fonts[name];   font!=nil {
        return nil
//        return log.NewError("font %s already loaded: %s",name,font.Desc())
    }
    
    
    {

        path,err := service.getFilePathForName(name)
    
        if err == nil {
            
            if DEBUG_FONTSERVICE { log.Debug("%s read font %s from %s",service.Desc(),name,path) }
            data, err = ioutil.ReadFile(path)
            if err != nil {
                return log.NewError("fail to read font from %s: %s",path,err)
            }
            
        } else  { // no file found, try embedded

            if DEBUG_FONTSERVICE { log.Debug("%s no file for font %s: %s",service.Desc(),name,err ) }

            encoded := FontAsset[name]
            if encoded == "" {
                return log.NewError("no data for embedded font %s",name)    
            }

            if DEBUG_FONTSERVICE { log.Debug("%s decode embedded font %s",service.Desc(),name) }
            data,err = base64.StdEncoding.DecodeString( encoded )
            if err != nil {
                return log.NewError("fail to decode embedded font %s: %s",name,err)    
            }
        }    
    }

    
    if len(data) <= 0 {
        if DEBUG_FONTSERVICE { log.Debug("%s no data for font %s",service.Desc(),name) }
        return log.NewError("no data for font %s",name)
        
    }

    font := NewFont(name)
    err = font.loadData( data )
    if err != nil {
        log.Debug("%s fail load font %s data: %s",service.Desc(),name,err)
        return log.NewError("fail to load font %s data",name)
    }


    if DEBUG_FONTSERVICE { log.Debug("%s add font %s",service.Desc(),name) }
    service.fonts[name] = font


    return nil
    
    
}


func (service *FontService) getFilePathForName(fontName string) (string,error) {
    var extensions =[]string{ ".ttf", ".ttc" }
    var err error
    files, err := ioutil.ReadDir(service.directory)
    if err != nil {
        return "", log.NewError("fail list fonts: %s",err)
    }
    for _, f := range files {
        for _, ext := range extensions {
            if strings.ToLower(f.Name()) == strings.ToLower(fontName+ext) {
                if DEBUG_FONTSERVICE { log.Debug("%s found font file %s",service.Desc(),f.Name()) }
                return service.directory + "/" + f.Name(), nil
            }
        }
    }
    return "",log.NewError("no files for font %s",fontName)
}



func (service *FontService) Desc() string {
    ret := "fontservice["
    ret += fmt.Sprintf("%d fonts",len(service.fonts))
    ret += "]"
    return ret
}


