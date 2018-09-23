

package gfx

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
    conf "../conf"
)



var foreground = image.White
var background = image.Black
//var background = image.Transparent

const GlyphCols  = 0x20
const GlyphRows  = 0x08
const GlyphCount = GlyphCols * GlyphRows


const maxDimension = 8192


type Font struct {
    Name string
    directory string

    font *truetype.Font
    context *freetype.Context
    tmp *image.RGBA
    max struct {w int; h int}
    Widths [GlyphCols][GlyphRows]int
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


func NewFont(config *conf.FontConfig, directory string) *Font {
    ret := &Font{directory: directory}
    ret.tmp = image.NewRGBA( image.Rect(0,0,maxDimension,maxDimension) )
    ret.Name = config.Name
    return ret
}


func (font *Font) Configure(config *conf.FontConfig) {
    log.Debug("configure font: %s",config.Desc())
    
    //regen?
    if font.Name != config.Name || font.font == nil {
        font.Name = config.Name
    }
    
}


func (font *Font) Desc() string {
    return fmt.Sprintf("font[%s]",font.Name)
}



func (font *Font) loadFont(fontfile string) error {
    var data []byte 
    var err error
    for _,ext := range []string{ ".ttc", ".ttf", ".TTC", ".TTF" } {
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
    log.Debug("read font file %s",fontfile)
    return nil
}

func (font *Font) Init() {
    err := font.loadFont(font.directory+font.Name)
    if err != nil {
        log.Error("fail to load font %s: %s",font.Name,err)
        return
    }
    font.context = freetype.NewContext()
    font.context.SetFont(font.font)

    font.Widths, font.max = font.findSizes()
    log.Debug("init font[%s %dx%d]",font.Name,font.max.w,font.max.h)
}



func (font *Font) RenderMapRGBA() (*image.RGBA, error) {

//    const (
//        pointSize = 72.0
//        rowSpacing = 1.0
//        dpi = 144.0  
//    )


//    var max struct{w int; h int}

//    var err error
    
    
//    _,max := font.findSizes(pointSize,dpi,rowSpacing)
    width  := font.max.w
    height := font.max.h


    ret := image.NewRGBA( image.Rect(0,0,GlyphCols*width,2*GlyphRows*height) )
    draw.Draw( ret, ret.Bounds(), background, image.ZP, draw.Src)

    ctx := font.context
    ctx.SetDPI( dpi )
    ctx.SetFontSize( pointSize )
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc( foreground )
    ctx.SetDst(ret)
    ctx.SetClip(ret.Bounds())

    c := 0x00
    for y:=0; y<GlyphRows; y++ {
        
        for x:=0; x<GlyphCols; x++ {
            str := fmt.Sprintf("%c",rune(c))
            if c < 0x20 || ( c >= 0x7f && c < 0xa0 ) {
                str = " "
            }

            pos := freetype.Pt( x*width, 2*y*height+height)
            ctx.DrawString(str,pos)
            c += 0x1
        }
    }
    log.Debug("rendered rgba map with %s",font.Desc())
    ctx.SetDst(nil)
    ctx.SetClip( image.Rect(0,0,0,0) )
    return ret,nil
    
}



//func (font *Font) RenderGlyphTexture() (*GlyphTexture, error) {
//
//
//
////    log.Debug("max is %v",max)
////    log.Debug("bounds is %v",font.font.Bounds( fixed.I(1<<6)))
////    log.Debug("%d funits per em",font.font.FUnitsPerEm())
//
//    ret.Texture, err = font.RenderGlyphRGBA(ret.Width,ret.Height,pointSize,rowSpacing,dpi)
//    if err != nil {
//        log.Error("fail render map rgba: %s",err)
//        return nil,err
//    }
////    log.Debug("rendered glyphs for %s",font.Desc())
//    
//    return ret,nil
//    
//}



func (font *Font) RenderTextRGBA(text string) (*image.RGBA, error) {
    
    
    const (
        pointSize = 72.0
        dpi = 144.0  
        rowSpacing = 1.25
    )

    ctx := font.context
    ctx.SetDPI( dpi )
    ctx.SetFontSize( pointSize )
    ctx.SetHinting( xfont.HintingFull )
    ctx.SetSrc( foreground )
    ctx.SetDst(font.tmp)
    ctx.SetClip(font.tmp.Bounds())

    dim,_ := ctx.DrawString(text,freetype.Pt(0,0) )
    
    width := dim.X.Ceil() + dim.X.Ceil()/16
    height := ctx.PointToFixed( rowSpacing * pointSize ).Ceil()

    ret := image.NewRGBA( image.Rect(0,0,width,height) )

    draw.Draw( ret, ret.Bounds(), background, image.ZP, draw.Src)

    ctx.SetDst(ret)
    ctx.SetClip(ret.Bounds())
    ctx.DrawString(text,freetype.Pt(width/32, 4*height/5))
//    log.Debug("rendered '%s' with %s",text,font.Desc())

    ctx.SetDst(nil)
    ctx.SetClip( image.Rect(0,0,0,0) )
    return ret,nil
}

const (
    pointSize = 28.0
    dpi = 144.0  
    rowSpacing = 1.25
)







func (font *Font) findSizes() ([GlyphCols][GlyphRows]int, struct{w int; h int}) {
    var widths [GlyphCols][GlyphRows]int    
    var max = struct {w int; h int} { 0, 0 }
//    var max fixed.Point26_6 = freetype.Pt(0,0)


    tmp := image.NewRGBA( image.Rect(0,0,1024,1024) )

    ctx := font.context
    ctx.SetDPI(dpi)
    ctx.SetFontSize(pointSize)
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc(image.White)
    ctx.SetDst(tmp)
    ctx.SetClip(tmp.Bounds())
    
    max.h = ctx.PointToFixed( rowSpacing * pointSize ).Ceil()
    
    c := 0x00
    for y:=0; y<GlyphRows; y++ {
        

        for x:=0; x<GlyphCols; x++ {
            str := fmt.Sprintf("%c",rune(c))
            if c < 0x020 || ( c>= 0x7f && c < 0xa0) {
                str = " "
            }
            dim,err := ctx.DrawString(str, freetype.Pt(0,0))
            if err != nil {
//                log.Error("could not draw glyph 0x%02x for sizing: %s",c,err)
                continue
            }
            widths[x][y] = dim.X.Ceil()
//            log.Debug("%d,%d\t0x%02x <%s> %v ~ %d",x,y,c,str,dim,widths[x][y])




//            scale := fixed.I( 1 )
//            idx := font.font.Index( rune(c) )
//            hmetric := font.font.HMetric(scale,idx)
//            vmetric := font.font.VMetric(scale,idx)
//            log.Debug("glyph %c\twidth %v\theight %v",c,hmetric,vmetric)

//            bounds := font.font.Bounds(scale)

//            magic := pointSize * dpi / 72.0
            
//            log.Debug("glyph %c at index %d has horizontal %v vertical %v\trenders %v\tmagic yields %5.2f",c,idx,hmetric,vmetric,dim.X,magic)

//            log.Debug("font scale %v bounds %v\t\tglyph %c\trenders %v\thorizontal %v\tvertical %v",scale,bounds,c,dim.X,hmetric,vmetric)

            if dim.X.Ceil() > max.w { max.w = dim.X.Ceil() }
            c += 1
            
        }
    }

    return widths,max

}

