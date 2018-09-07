
package font

import (
    "fmt"
//    "errors"
    "io/ioutil"
    "image"
    "image/draw"
    xfont "golang.org/x/image/font"
//    "golang.org/x/image/math/fixed"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    log "../log"
)


const fontsize = 144.0
const rowheight = 1.5
const dpi = 144.0
const maxwidth = 8192

var foreground = image.White
var background = image.Black
//const background = image.Transparent

const GlyphCols  = 0x10
const GlyphRows  = 0x10
const GlyphCount = GlyphCols * GlyphRows


type Font struct {
    face string
    font *truetype.Font  
}


type GlyphTexture struct {
    Texture *image.RGBA
    Height int
    Width int
    Widths [GlyphCols][GlyphRows]int
}


type TextTexture struct {
    Texture *image.RGBA
    Width int
    Height int
    Text string
}


func NewFont() *Font {
    return &Font{}
}


func (font *Font) Configure(config *Config, directory string) {
    log.Debug("configure font: %s",config.Describe())
    err := font.loadFont(directory+config.Face)
    if err != nil {
        log.Error("fail to config %s: %s",config.Describe(),err)
        return
    }
    font.face = config.Face
    
}


func (font *Font) loadFont(fontfile string) error {
    var data []byte 
    var err error
    for _,ext := range []string{ ".ttc", ".ttf" } {
        data, err = ioutil.ReadFile(fontfile + ext )
        if err == nil {
            log.Debug("load font file %s",fontfile+ext)
            break   
        }
    }
    if err != nil {
        log.Error("fail to read font %s: %s",fontfile,err)
        return err
    }
    font.font,err = freetype.ParseFont(data)
    if err != nil {
        log.Error("fail to parse font %s: %s",fontfile,err)
        return err
    }
    return nil
}


func (font *Font) findSizes() ([GlyphCols][GlyphRows]int, struct{w int; h int}) {
    var widths [GlyphCols][GlyphRows]int    
    var max = struct {w int; h int} { 0, 0 }
//    var max fixed.Point26_6 = freetype.Pt(0,0)


    tmp := image.NewRGBA( image.Rect(0,0,1024,1024) )

    ctx := freetype.NewContext()
    ctx.SetFont(font.font)
    ctx.SetDPI(dpi)
    ctx.SetFontSize(fontsize)
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc(image.White)
    ctx.SetDst(tmp)
    ctx.SetClip(tmp.Bounds())
    
    
    c := 0x00
    for y:=0; y<GlyphRows; y++ {
        for x:=0; x<GlyphCols; x++ {
            str := fmt.Sprintf("%c",rune(c))
            if c < 0x020 || ( c>= 0x7f && c < 0xa0) {
                str = " "
            }
            dim,err := ctx.DrawString(str, freetype.Pt(0,0))
            if err != nil {
                log.Error("could not draw glyph 0x%02x for sizing: %s",c,err)
                continue
            }
            widths[x][y] = dim.X.Ceil()
//            log.Debug("%d,%d\t0x%02x <%s> %v ~ %d",x,y,c,str,dim,widths[x][y])

            if dim.X.Ceil() > max.w { max.w = dim.X.Ceil() }
            c += 1
            
        }
    }

    max.h = ctx.PointToFixed( rowheight * fontsize ).Ceil()
    return widths,max

}


func (font *Font) RenderGlyphTexture() (*GlyphTexture, error) {

    var ret *GlyphTexture = &GlyphTexture{}
    var max struct{w int; h int}
    
    ret.Widths,max = font.findSizes()
    ret.Width  = max.w
    ret.Height = max.h



    log.Debug("max is %v",max)

    ret.Texture = image.NewRGBA( image.Rect(0,0,GlyphCols*ret.Width,GlyphRows*ret.Height) )


    draw.Draw( ret.Texture, ret.Texture.Bounds(), background, image.ZP, draw.Src)


    ctx := freetype.NewContext()
    ctx.SetFont(font.font)
    ctx.SetDPI( dpi )
    ctx.SetFontSize( fontsize )
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc( foreground )
    ctx.SetDst(ret.Texture)
    ctx.SetClip(ret.Texture.Bounds())
    

    
    c := 0x0
    for y:=0; y<GlyphRows; y++ {
        for x:=0; x<GlyphCols; x++ {
            str := fmt.Sprintf("%c",rune(c))
            if c < 0x20 || ( c >= 0x7f && c < 0xa0 ) {
                str = "X"
            }

            pos := freetype.Pt( x*ret.Width, y*ret.Height + ret.Height)
//            draw.Draw(ret.Texture,image.Rect(x*ret.Width,y*ret.Height,(x+1)*ret.Width,(y+1)*ret.Height) )
            ctx.DrawString(str,pos)
            c += 0x1
        }
    }
    
    log.Debug("rendered glyphs for %s",font.Describe())
    
    return ret,nil
    
}


func (font *Font) RenderTextTexture(text string) (*TextTexture, error) {

    var ret *TextTexture = &TextTexture{Text: text}
    
    tmp := image.NewRGBA( image.Rect(0,0,maxwidth,maxwidth) )
    ctx := freetype.NewContext()
    ctx.SetFont(font.font)
    ctx.SetDPI( dpi )
    ctx.SetFontSize( fontsize )
    ctx.SetHinting( xfont.HintingFull )
    ctx.SetSrc( foreground )
    ctx.SetDst(tmp)
    ctx.SetClip(tmp.Bounds())
    
    str := ret.Text
    
    dim,_ := ctx.DrawString(str,freetype.Pt(0,0) )
    
    ret.Width = dim.X.Ceil()
    ret.Height = ctx.PointToFixed( rowheight * fontsize ).Ceil()

    log.Debug("got dimensions %dx%d for text '%s'",ret.Width,ret.Height,str)
    
    ret.Texture = image.NewRGBA( image.Rect(0,0,ret.Width,ret.Height) )
    draw.Draw( ret.Texture, ret.Texture.Bounds(), background, image.ZP, draw.Src)
    
    
    ctx.SetDst(ret.Texture)
    ctx.SetClip(ret.Texture.Bounds())
    ctx.DrawString(str,freetype.Pt(0, 2*ret.Height/3))
    
    log.Debug("rendered text '%s' for %s",ret.Text,font.Describe())
    return ret,nil
        
}




func (font *Font) Describe() string {
    return fmt.Sprintf("font[%s]",font.face)
}

