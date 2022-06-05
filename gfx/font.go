//go:build (linux && arm) || DARWIN_GUI
// +build RENDERER

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	xfont "golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"strings"
)

const DEBUG_FONT = false

var fonts = map[string]*Font{}

const GlyphMapCols = 0x20
const GlyphMapRows = 0x08
const GlyphMapCount = GlyphMapCols * GlyphMapRows

const FontScratchSize = 8192

type Font struct {
	name string

	font    *truetype.Font
	context *freetype.Context
	scratch *image.RGBA

	max  struct{ w, h int }
	size [GlyphMapCols][GlyphMapRows]Size

	glyphMap *image.RGBA
}

func (font *Font) Size(col, row int) Size {
	if col < GlyphMapCols && row < GlyphMapRows {
		return font.size[col][row]
	}
	return Size{0, 0}
}

func (font *Font) Ratio() float32 { return float32(font.max.w) / float32(font.max.h) }

func (font *Font) MaxSize() Size {
	return Size{W: float32(font.max.w), H: float32(font.max.h)}
}

func NewFont(name string, scratch *image.RGBA) *Font {
	return &Font{name: name, scratch: scratch}
}

func (font *Font) GetName() string {
	if font == nil {
		return ""
	}
	return font.name
}

func (font *Font) Desc() string {
	//    tw,th := font.MaxSize().W, font.MaxSize().H
	//    mw, mh := font.MaxSize().W * GlyphMapCols, font.MaxSize().H * GlyphMapRows
	//    gw, gh := GlyphMapCols, GlyphMapRows
	//    return fmt.Sprintf("font[ %dx%d %s %.0fx%.0f %.0fx%.0f]",gw,gh,font.config.Name,tw,th,mw,mh)

	ret := "font["
	ret += font.name
	if font.max.w > 0 && font.max.h > 0 {
		ret += fmt.Sprintf(" %.2f", font.Ratio())
	}
	ret += "]"
	return ret
}

func (font *Font) loadData(data []byte) error {
	var err error

	if DEBUG_FONTSERVICE {
		log.Debug("%s parse data", font.Desc())
	}
	font.font, err = freetype.ParseFont(data)
	if err != nil {
		return log.NewError("fail parse: %s", err)
	}

	font.context = freetype.NewContext()
	font.context.SetFont(font.font)
	//    font.scratch = image.NewRGBA( image.Rect(0,0,ScratchSize,ScratchSize) )

	if DEBUG_FONTSERVICE {
		log.Debug("%s find sizes", font.Desc())
	}
	font.size, font.max = font.findSizes()

	return nil

}

//func (font *Font) Init() {
//
//    if font.context != nil {
//        log.Error("skip init %s",font.Desc())
//        return
//    }
//    font.context = freetype.NewContext()
//    font.context.SetFont(font.font)
//
//    font.Size, font.max = font.findSizes()
//    if DEBUG_FONTSERVICE { log.Debug("init %s",font.Desc()) }
//}

//func (font *Font) Close() { // do me do me
//    font.scratch = nil // bad idea better keep around
//    font.context = nil //bad idea better keep around
//    font.font = nil // bad idea better keep around
//    font.config.Name = "-" + font.config.Name
//    log.Debug("close %s",font.Desc())
//}

func (font *Font) RenderMap(debug bool) (*image.RGBA, error) {

	if font.glyphMap != nil && !debug {
		return font.glyphMap, nil
	}

	width := font.max.w
	height := font.max.h

	imageWidth := GlyphMapCols * width
	imageHeight := GlyphMapRows * height

	back := BackgroundColor
	if debug {
		back = DebugColor
	}

	ret := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(ret, ret.Bounds(), back, image.ZP, draw.Src)

	ctx := font.context
	ctx.SetDPI(dpi)
	ctx.SetFontSize(pointSize)
	ctx.SetHinting(xfont.HintingNone)
	ctx.SetSrc(ForegroundColor)
	ctx.SetDst(ret)
	ctx.SetClip(ret.Bounds())

	var c byte = 0x00
	for y := 0; debug && y < GlyphMapRows; y++ {

		for x := 0; x < GlyphMapCols; x++ {

			pos := freetype.Pt(x*width, y*height+height)

			min := []int{
				pos.X.Floor(),
				pos.Y.Floor() - height,
			}
			max := []int{
				min[0] + int(font.size[x][y].W),
				min[1] + height,
			}

			back := image.NewUniform(color.RGBA{0x40, 0x40, 0x40, 0xff})
			fore := image.NewUniform(color.RGBA{0xC0, 0xC0, 0xC0, 0xff})

			bounds := image.Rectangle{Min: image.Point{min[0], min[1]}, Max: image.Point{max[0], max[1]}}
			if (y%2 == 0 && c%2 == 0) || (y%2 != 0 && c%2 != 0) {
				draw.Draw(ret, bounds, fore, image.ZP, draw.Src)
			} else {
				draw.Draw(ret, bounds, back, image.ZP, draw.Src)
			}

			c += 0x1
		}
	}

	c = 0x00
	for y := 0; y < GlyphMapRows; y++ {

		for x := 0; x < GlyphMapCols; x++ {

			str := font.stringForByte(c)

			magic_offset := height / 5. // this should really come from the font geometrics
			pos := freetype.Pt(x*width, y*height+height-magic_offset)

			ctx.SetSrc(ForegroundColor)
			if debug && ((y%2 == 0 && c%2 == 0) || (y%2 != 0 && c%2 != 0)) {
				r, g, b, _ := ForegroundColor.RGBA()
				back := color.RGBA{255 - uint8(r), 255 - uint8(g), 255 - uint8(b), 255}
				ctx.SetSrc(image.NewUniform(back))
			}

			ctx.DrawString(str, pos)
			c += 0x1
		}
	}

	ctx.SetDst(nil)
	ctx.SetClip(image.Rect(0, 0, 0, 0))

	if !debug {
		font.glyphMap = ret
	}

	if DEBUG_FONTSERVICE {
		log.Debug("%s rendered glyphmap: %dx%d glyphs as %dx%d img", font.Desc(), GlyphMapCols, GlyphMapRows, imageWidth, imageHeight)
	}

	return ret, nil

}

func (font *Font) stringForByte(b byte) string {
	if b < 0x20 || (b >= 0x7f && b < 0xa0) {
		return " "
	}
	switch strings.ToLower(font.name) {

	case "monaco":
		switch b {
		case 0xA6, 0xAD, 0xB2, 0xB3, 0xB7, 0xB9, 0xBC, 0xBD, 0xBE, 0xD0, 0xD7, 0xDD, 0xDE, 0xF0, 0xFD, 0xFE:
			if DEBUG_FONT {
				log.Debug("%s special-case monaco char 0x%02x '%c'", font.Desc(), b, rune(b))
			}
			return " "
		default:
		}

	case "ocraext":
		if b == 0xB7 {
			if DEBUG_FONT {
				log.Debug("%s special-case ocraext char 0x%02x '%c'", font.Desc(), b, rune(b))
			}
			return " "
		}

	case "robotomono":
		if b == 0xA0 {
			if DEBUG_FONT {
				log.Debug("%s special-case robotomono char 0x%02x '%c'", font.Desc(), b, rune(b))
			}
			return " "
		}
	}
	return fmt.Sprintf("%c", rune(b))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (font *Font) RenderText(text string, DEBUG bool) (*image.RGBA, error) {

	const (
		pointSize  = 48.0
		dpi        = 144.0
		rowSpacing = 1.25
	)

	background := BackgroundColor
	if DEBUG {
		background = DebugColor
	}

	ctx := font.context
	ctx.SetDPI(dpi)
	ctx.SetFontSize(pointSize)
	ctx.SetHinting(xfont.HintingFull)
	ctx.SetSrc(ForegroundColor)
	ctx.SetDst(font.scratch)
	ctx.SetClip(font.scratch.Bounds())

	dim, _ := ctx.DrawString(text, freetype.Pt(0, 0))

	imageWidth := dim.X.Ceil() + dim.X.Ceil()/16
	imageHeight := ctx.PointToFixed(rowSpacing * pointSize).Ceil()

	ret := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	draw.Draw(ret, ret.Bounds(), background, image.ZP, draw.Src)

	ctx.SetDst(ret)
	ctx.SetClip(ret.Bounds())
	ctx.DrawString(text, freetype.Pt(imageWidth/32, 4*imageHeight/5))

	ctx.SetDst(nil)
	ctx.SetClip(image.Rect(0, 0, 0, 0))

	if DEBUG_FONTSERVICE {
		log.Debug("%s rendered '%s': %dx%d glyphs as %dx%d img", font.Desc(), text[0:min(len(text), 8)], GlyphMapCols, GlyphMapRows, imageWidth, imageHeight)
	}

	return ret, nil
}

const (
	pointSize  = 28.0
	dpi        = 144.0
	rowSpacing = 1.25
)

func (font *Font) findSizes() ([GlyphMapCols][GlyphMapRows]Size, struct{ w, h int }) {
	var size [GlyphMapCols][GlyphMapRows]Size
	var max struct{ w, h int }

	ctx := font.context
	ctx.SetDPI(dpi)
	ctx.SetFontSize(pointSize)
	ctx.SetHinting(xfont.HintingNone)
	ctx.SetSrc(image.White)
	ctx.SetDst(font.scratch)
	ctx.SetClip(image.Rect(0, 0, 1024, 1024))

	max.h = ctx.PointToFixed(rowSpacing * pointSize).Ceil()

	var c byte = 0x00
	for y := 0; y < GlyphMapRows; y++ {

		for x := 0; x < GlyphMapCols; x++ {

			str := font.stringForByte(c)
			dim, err := ctx.DrawString(str, freetype.Pt(0, 0))
			if err != nil {
				log.Error("%s fail draw glyph 0x%02x for sizing: %s", font.Desc(), c, err)
				continue
			}
			size[x][y].W = float32(dim.X.Ceil())
			size[x][y].H = float32(ctx.PointToFixed(rowSpacing * pointSize).Ceil())

			if DEBUG_FONT {
				log.Debug("%s glyph '%s' 0x%02x is %.1fx%.1f", font.Desc(), str, c, size[x][y].W, size[x][y].H)
			}
			//            scale := fixed.I( 1 )
			//            idx := font.font.Index( rune(c) )
			//            hmetric := font.font.HMetric(scale,idx)
			//            vmetric := font.font.VMetric(scale,idx)
			//            log.Debug("glyph %c\twidth %v\theight %v",c,hmetric,vmetric)

			//            bounds := font.font.Bounds(scale)

			//            magic := pointSize * dpi / 72.0

			//            log.Debug("glyph %c at index %d has horizontal %v vertical %v\trenders %v\tmagic yields %5.2f",c,idx,hmetric,vmetric,dim.X,magic)

			//            log.Debug("font scale %v bounds %v\t\tglyph %c\trenders %v\thorizontal %v\tvertical %v",scale,bounds,c,dim.X,hmetric,vmetric)

			if dim.X.Ceil() > max.w {
				max.w = dim.X.Ceil()
			}
			c += 1

		}
	}
	ctx.SetDst(nil)
	ctx.SetClip(image.Rect(0, 0, 0, 0))
	if DEBUG_FONTSERVICE {
		log.Debug("%s found max size %.1dx%.1d", font.Desc(), max.w, max.h)
	}
	return size, max

}
