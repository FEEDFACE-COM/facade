
// +build linux,arm

package gfx

import (
    "fmt"
//    "errors"
    "io/ioutil"
    "image"
    "image/color"
    "image/draw"
    xfont "golang.org/x/image/font"
//    "golang.org/x/image/math/fixed"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    log "../log"
    conf "../conf"
)


var fonts = map[string]*Font {}


var foreground = image.White
//var background = image.NewUniform( color.RGBA{R: 0, G: 255, B: 255, A: 255} )
var background = image.Transparent

const GlyphCols  = 0x20
const GlyphRows  = 0x08
const GlyphCount = GlyphCols * GlyphRows


const scratchSize = 8192


const DEBUG_FONT = false


func GetFont(config *conf.FontConfig, directory string) (*Font,error) {
    if fonts[config.Name] == nil {
        tmp := NewFont(config, directory)
        err := tmp.loadFont(directory+config.Name)
        if err != nil {
            fonts[config.Name] = nil
            log.Error("fail to load font %s: %s",config.Name,err)
            return nil, log.NewError("fail to load font %s: %s",config.Name,err)
        } 
        fonts[config.Name] = tmp
    } 
    return fonts[config.Name],nil
    //note, its' still leaking tho!
    
}



type Font struct {
    config conf.FontConfig

    directory string

    font *truetype.Font
    context *freetype.Context
    scratch *image.RGBA
    max struct {w, h int}
    Size [GlyphCols][GlyphRows]Size
}


func (font *Font) Ratio() float32 { return float32(font.max.w) / float32(font.max.h) }

func (font *Font) MaxSize() Size {
    return Size{W: float32(font.max.w), H: float32(font.max.h)}    
}


func NewFont(config *conf.FontConfig, directory string) *Font {
    ret := &Font{config: *config, directory: directory}
    ret.scratch = image.NewRGBA( image.Rect(0,0,scratchSize,scratchSize) )
    return ret
}


func (font *Font) Configure(config *conf.FontConfig) {
    if config == nil { return }
    if *config == font.config { return }
    
    log.Debug("config %s -> %s",font.Desc(),config.Desc())
    font.config = *config
//    }
    
}


func (font *Font) Desc() string { 
//    tw,th := font.MaxSize().W, font.MaxSize().H
//    mw, mh := font.MaxSize().W * GlyphCols, font.MaxSize().H * GlyphRows
//    gw, gh := GlyphCols, GlyphRows
//    return fmt.Sprintf("font[ %dx%d %s %.0fx%.0f %.0fx%.0f]",gw,gh,font.config.Name,tw,th,mw,mh)
    return fmt.Sprintf("font[%s %.2f]",font.config.Name,font.Ratio())
}



func (font *Font) loadFont(fontfile string) error {
    var data []byte 
    var err error
    var ext string
    for _,ext = range []string{ ".ttc", ".ttf", ".TTC", ".TTF" } {
        data, err = ioutil.ReadFile(fontfile + ext )
        if err == nil {
            break   
        }
    }
    if err != nil {
        return log.NewError("fail to read font: %s",err)
    }
    font.font,err = freetype.ParseFont(data)
    if err != nil {
        return log.NewError("fail to parse: %s",err)
    }
    log.Debug("load font file %s",fontfile+ext)
    return nil
}

func (font *Font) Init() {
    if font.context != nil { 
        log.Debug("font ready")
        return 
    }
    font.context = freetype.NewContext()
    font.context.SetFont(font.font)

    font.Size, font.max = font.findSizes()
    log.Debug("init font[%s %dx%d]",font.config.Name,font.max.w,font.max.h)
}

func (font *Font) Close() { // do me do me 
//    font.scratch = nil
//    font.context = nil
//    font.font = nil
//    font.config.Name = "+" + font.config.Name
}


func (font *Font) RenderMapRGBA() (*image.RGBA, error) {

//    const (
//        pointSize = 72.0
//        rowSpacing = 1.0
//        dpi = 144.0  
//    )



    width := font.max.w
    height := font.max.h


    imageWidth := GlyphCols*width 
    imageHeight := GlyphRows*height

    back := background
    if DEBUG_FONT {
        back = image.NewUniform( color.RGBA{R: 255, G: 0, B: 0, A: 255} )
    }
        
    ret := image.NewRGBA( image.Rect(0,0,imageWidth,imageHeight) )
    draw.Draw( ret, ret.Bounds(), back, image.ZP, draw.Src)

    ctx := font.context
    ctx.SetDPI( dpi )
    ctx.SetFontSize( pointSize )
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc( foreground )
    ctx.SetDst(ret)
    ctx.SetClip(ret.Bounds())




    var c byte = 0x00
    for y:=0; DEBUG_FONT && y<GlyphRows; y++ {
        
        for x:=0; x<GlyphCols; x++ {
            
            
            pos := freetype.Pt( x*width, y*height+height)

            min := []int{ 
                pos.X.Floor(), 
                pos.Y.Floor()-height,  
            }
            max := []int{ 
                min[0] + int(font.Size[x][y].W),  
                min[1] + height,  
            }
            
            back := image.NewUniform( color.RGBA{0x40,0x40,0x40,0xff} )
            fore := image.NewUniform( color.RGBA{0xC0,0xC0,0xC0,0xff} )

            bounds := image.Rectangle{Min: image.Point{min[0],min[1]}, Max: image.Point{max[0],max[1]} }
            if (y % 2 == 0 && c % 2 == 0 ) || (y%2 != 0 && c %2 != 0) {
                draw.Draw( ret, bounds, fore, image.ZP, draw.Src )
            } else {
                draw.Draw( ret, bounds, back, image.ZP, draw.Src )
            }
                
                
                
            c += 0x1
        }
    }
    
    
    
    c = 0x00
    for y:=0; y<GlyphRows; y++ {
        
        for x:=0; x<GlyphCols; x++ {
            
            str := font.stringForByte(c)
            
            magic_offset := height/5. // this should really come from the font geometrics
            pos := freetype.Pt( x*width, y*height + height - magic_offset)
            
            ctx.SetSrc( foreground )
            if DEBUG_FONT && (  ( y % 2 == 0 && c % 2 == 0 ) || (y%2 != 0 && c %2 != 0) ) {
                r,g,b,_ := foreground.RGBA()
                back := color.RGBA{255 - uint8(r), 255-uint8(g), 255 - uint8(b), 255}
                ctx.SetSrc( image.NewUniform( back ) )
            }
                
            ctx.DrawString(str,pos)
            c += 0x1
        }
    }
    
    
    
    ctx.SetDst(nil)
    ctx.SetClip( image.Rect(0,0,0,0) )

    if DEBUG_FONT {
        log.Debug("rendered map %s   %dx%d glyphs in %dx%d img",font.Desc(),GlyphCols,GlyphRows,imageWidth,imageHeight)
    }
    return ret,nil
    
}

func (font *Font) stringForByte(b byte) string {
        if b < 0x20 || ( b >= 0x7f && b < 0xa0 ) {
            return "X"
        }
        if font.config.Name == "OCRAEXT" && b == 0xB7 {
            log.Debug("special-case ocraext '%c'",b)
            return " "
        }
        return fmt.Sprintf("%c",rune(b))    
}


func min(a,b int) int { 
    if a < b { return a }
    return b 
}

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
    ctx.SetDst(font.scratch)
    ctx.SetClip(font.scratch.Bounds())

    dim,_ := ctx.DrawString(text,freetype.Pt(0,0) )
    
    imageWidth := dim.X.Ceil() + dim.X.Ceil()/16
    imageHeight := ctx.PointToFixed( rowSpacing * pointSize ).Ceil()


    ret := image.NewRGBA( image.Rect(0,0,imageWidth,imageHeight) )

    draw.Draw( ret, ret.Bounds(), background, image.ZP, draw.Src)

    ctx.SetDst(ret)
    ctx.SetClip(ret.Bounds())
    ctx.DrawString(text,freetype.Pt(imageWidth/32, 4*imageHeight/5))

    ctx.SetDst(nil)
    ctx.SetClip( image.Rect(0,0,0,0) )
    
    if DEBUG_FONT {
        log.Debug("rendered '%s' %s   %dx%d glyphs in %dx%d img",text[0:min(len(text),8)],font.Desc(),GlyphCols,GlyphRows,imageWidth,imageHeight)
    }
    
    return ret,nil
}

const (
    pointSize = 28.0
    dpi = 144.0  
    rowSpacing = 1.25
)







func (font *Font) findSizes() ([GlyphCols][GlyphRows]Size, struct{w,h int}) {
    var size [GlyphCols][GlyphRows]Size
    var max struct{w,h int} 

    ctx := font.context
    ctx.SetDPI(dpi)
    ctx.SetFontSize(pointSize)
    ctx.SetHinting( xfont.HintingNone )
    ctx.SetSrc(image.White)
    ctx.SetDst(font.scratch)
    ctx.SetClip(image.Rect(0,0,1024,1024))
    
    max.h = ctx.PointToFixed( rowSpacing * pointSize ).Ceil()
    
    var c byte = 0x00
    for y:=0; y<GlyphRows; y++ {
        

        for x:=0; x<GlyphCols; x++ {
            
            str := font.stringForByte(c)
            dim,err := ctx.DrawString(str, freetype.Pt(0,0))
            if err != nil {
                log.Error("could not draw glyph 0x%02x for sizing: %s",c,err)
                continue
            }
            size[x][y].W = float32(dim.X.Ceil())
            size[x][y].H = float32(ctx.PointToFixed( rowSpacing * pointSize ).Ceil())
            
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

    return size,max

}

