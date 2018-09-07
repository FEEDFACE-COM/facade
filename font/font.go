
package font

import (
    "fmt"
//    "errors"
    "io/ioutil"
    "image"
    "image/draw"
    xfont "golang.org/x/image/font"
    "golang.org/x/image/math/fixed"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    log "../log"
)


const fontsize = 144.0
const rowheight = 1.5
const dpi = 144.0

const GlyphCols  = 0x10
const GlyphRows  = 0x10
const GlyphCount = GlyphCols * GlyphRows


type Font struct {
    face string
    font *truetype.Font  
}


type GlyphMap struct {
    Texture *image.RGBA
    Height int
    Width int
    Widths [GlyphCols][GlyphRows]int
}


type TextTexture struct {
    Texture *image.RGBA
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


func (font *Font) findSizes() ([GlyphCols][GlyphRows]int, fixed.Point26_6) {
    var widths [GlyphCols][GlyphRows]int    
    var max fixed.Point26_6 = freetype.Pt(0,0)


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
            s := fmt.Sprintf("%c",rune(c))
            if c < 0x020 || ( c>= 0x7f && c < 0xa0) {
                s = " "
            }
            dim,err := ctx.DrawString(s, freetype.Pt(0,0))
            if err != nil {
                log.Error("could not draw glyph 0x%02x for sizing: %s",c,err)
                continue
            }
            widths[x][y] = dim.X.Ceil()
            log.Debug("%3d,%3x\t0x%02x <%s> %v ~ %d",x,y,c,s,dim,widths[x][y])

            if dim.X > max.X { max.X = dim.X }
            c += 1
            
        }
    }

    max.Y = ctx.PointToFixed( rowheight * fontsize )
    return widths,max

}


func (font *Font) RenderGlyphMap() (*GlyphMap, error) {

    var ret *GlyphMap = &GlyphMap{}
    var max fixed.Point26_6
    
    ret.Widths,max = font.findSizes()
    ret.Width  = max.X.Ceil()
    ret.Height = max.Y.Ceil()



    log.Debug("max is %v",max)

    ret.Texture = image.NewRGBA( image.Rect(0,0,GlyphCols*ret.Width,GlyphRows*ret.Height) )

    var color = image.White
//    var Background = image.Transparent
    var Background = image.Black

    draw.Draw( ret.Texture, ret.Texture.Bounds(), Background, image.ZP, draw.Src)


    ctx := freetype.NewContext()
    ctx.SetFont(font.font)
    ctx.SetDPI( dpi )
    ctx.SetFontSize( fontsize )
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc( color )
    ctx.SetDst(ret.Texture)
    ctx.SetClip(ret.Texture.Bounds())
    

    
    c := 0x0
    for y:=0; y<GlyphRows; y++ {
        for x:=0; x<GlyphCols; x++ {
            s := fmt.Sprintf("%c",rune(c))
            if c < 0x20 || ( c >= 0x7f && c < 0xa0 ) {
                s = " "
            }

            pos := freetype.Pt( x*ret.Width, y*ret.Height )
            ctx.DrawString(s,pos)
            c += 0x1
        }
    }
    
    log.Debug("rendered glyphmap for %s",font.Describe())
    
    return ret,nil
    
}


func (font *Font) Describe() string {
    return fmt.Sprintf("font[%s]",font.face)
}

