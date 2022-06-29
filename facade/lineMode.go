//go:build RENDERER
// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"fmt"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

const DEBUG_LINEMODE = true

type LineMode struct {
	width, height uint

	downward bool
	buffer   uint

	vert, frag string

	lineBuffer *LineBuffer
	termBuffer *TermBuffer

	texture *gfx.Texture
	program *gfx.Program
	object  *gfx.Object
	data    []float32

	refreshChan chan bool
}

const (
	TILECOUNT  gfx.UniformName = "tileCount"
	TILESIZE   gfx.UniformName = "tileSize"
	TILEOFFSET gfx.UniformName = "tileOffset"
	CURSORPOS  gfx.UniformName = "cursorPos"
	SCROLLER   gfx.UniformName = "scroller"
)

const (
	TILECOORD gfx.AttribName = "tileCoord"
	GRIDCOORD gfx.AttribName = "gridCoord"
)

func (mode *LineMode) ScheduleRefresh() {

	select {
	case mode.refreshChan <- true:
	default:
	}

}

func (mode *LineMode) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-mode.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}

func (mode *LineMode) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	if mode.checkRefresh() {
		if DEBUG_LINEMODE {
			log.Debug("%s refresh", mode.Desc())
		}
		mode.generateData(font)
		mode.renderMap(font) // TODO: should not renderMap anytime anything changed!
	}

	gl.ActiveTexture(gl.TEXTURE0)
	mode.program.UseProgram(debug)
	mode.object.BindBuffer()

	tileCount := mgl32.Vec2{float32(mode.width), float32(mode.height)}
	mode.program.Uniform2fv(TILECOUNT, 1, &tileCount[0])

	tileSize := mgl32.Vec2{font.MaxSize().W / font.MaxSize().H, font.MaxSize().H / font.MaxSize().H}
	mode.program.Uniform2fv(TILESIZE, 1, &tileSize[0])

	tileOffset := mgl32.Vec2{-1., 0.0}
	if mode.width%2 == 0 { //even columns
		tileOffset[0] = 0.5
	}
	if mode.height%2 == 0 { //even rows
		tileOffset[1] = -0.5
	}
	if mode.downward && mode.termBuffer == nil { // mode lines
		tileOffset[1] += 1.
	}
	mode.program.Uniform2fv(TILEOFFSET, 1, &tileOffset[0])

	cursorPos := mgl32.Vec2{-1., -1.}
	if mode.termBuffer != nil { // mode term
		x, y := mode.termBuffer.GetCursor()
		cursorPos[0] = float32(x)
		cursorPos[1] = float32(y)
	}
	mode.program.Uniform2fv(CURSORPOS, 1, &cursorPos[0])

	mode.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	mode.program.Uniform1f(gfx.FONTRATIO, font.Ratio())

	clocknow := float32(gfx.Now())
	mode.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow)

	scroller := float32(0.0)
	if mode.termBuffer == nil { // lines mode
		scroller = float32(mode.lineBuffer.GetScroller())
		if mode.downward {
			scroller *= -1.
		}
	}
	mode.program.Uniform1f(SCROLLER, scroller)

	camera.Uniform(mode.program)

	mode.texture.Uniform(mode.program)

	scale := float32(1.0)
	scale = mode.autoScale(camera, font)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	mode.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])

	mode.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+2+2)*4, (0)*4)
	mode.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+2+2)*4, (3)*4)
	mode.program.VertexAttribPointer(TILECOORD, 2, (3+2+2+2)*4, (3+2)*4)
	mode.program.VertexAttribPointer(GRIDCOORD, 2, (3+2+2+2)*4, (3+2+2)*4)

	count := int32(mode.width * (mode.height + 1))
	offset := int32(0)

	if !debug {
		mode.program.SetDebug(false)
		mode.texture.BindTexture()
		gl.DrawArrays(gl.TRIANGLES, int32(offset*2*3), (count)*(2*3))
		mode.program.SetDebug(debug)
	}

	if debug {
		gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		off := offset
		// REM, use single gl.DrawElements call instead (create indice array before)
		for r := 0; r < int(mode.height+1); r++ {
			for c := 0; c < int(mode.width); c++ {
				gl.DrawArrays(gl.LINE_STRIP, int32(off*2*3), int32(1*2*3))
				off += int32(1)
			}
			//    	   gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(mode.Width*2*3) )
			//    	   off += int32(mode.Width)
		}

	}
}

func (mode *LineMode) fill(name string) []string {
	ret := []string{}

	switch name {
	case "title":
		title := []string{""}

		if mode.width >= 80 && mode.height >= 16 {
			title = strings.Split(`
 _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _      
|_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|
|   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |
`, "\n")[1:]
		} else if mode.width >= 40 && mode.height >= 8 {
			title = strings.Split(`
 _   _   _   _   _   _
|_  |_| /   |_| | \ |_
|   | | \_  | | |_/ |_
                      
       by FEEDFACE.COM
`, "\n")[1:]
		} else if mode.width >= 13 {
			title = []string{"F A C A D E", ""}
		} else if mode.width >= 8 {
			title = []string{"FACADE", ""}
		}

		gw, gh := int(mode.width), int(mode.height)
		tw, th := len(title[0]), len(title)-1
		dw, dh := ((gw - tw) / 2.), ((gh - th) / 2.)

		if DEBUG_LINEMODE {
			log.Debug("fit %dx%d title into %dx%d grid, margin %dx%d", tw, th, gw, gh, dw, dh)
		}
		for r := 0; r < gh; r++ {
			tmp := ""
			tr := r - dh
			for c := 0; c < gw; c++ {
				tc := c - dw
				if tr >= 0 && tr < th && tc >= 0 && tc < tw {
					tmp += string(title[tr][tc])
				} else {
					tmp += " "
				}
			}
			ret = append(ret, tmp)
		}

		return ret

	case "index":
		w, h := int(mode.width), int(mode.height)
		for r := 0; r < h; r++ {
			tmp := ""
			for c := 0; c < w; c++ {
				d := "."
				if c%10 == 0 {
					d = fmt.Sprintf("%d", r%10)
				}
				if r%5 == 0 {
					d = fmt.Sprintf("%d", c%10)
				}
				if c%10 == 0 && r%5 == 0 {
					d = "#"
				}
				tmp += fmt.Sprintf("%s", d)
			}
			ret = append(ret, tmp)
		}
		return ret

	case "alpha":
		w, h := int(mode.width), int(mode.height)
		alpha := strings.Repeat(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", w)
		s := 0
		for r := 0; r < h; r++ {
			ret = append(ret, alpha[s:s+w])
			s += w
		}
		return ret

	case "clear":
		w, h := int(mode.width), int(mode.height)
		for r := 0; r < h; r++ {
			s := ""
			for c := 0; c < w; c++ {
				s += " "
			}
			ret = append(ret, s)
		}
		return ret

	}

	log.Warning("%s no such fill pattern: %s", mode.Desc(), name)
	return ret

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mode *LineMode) vertices(
	tileSize gfx.Size, // dimensions of the tile
	glyphSize gfx.Size, // dimensions of the glyph
	maxGlyphSize gfx.Size, // max dimension of glyph in font
	gridCoord gfx.Coord, // row/col coord of tile
	tilePos gfx.Point, // position of tile in grid
	texOffset gfx.Point, // glyph offset in texture
) []float32 {

	col, row := float32(gridCoord.X), float32(gridCoord.Y)

	w, h := tileSize.W, tileSize.H
	x, y := tilePos.X, tilePos.Y
	ox, oy := texOffset.X, texOffset.Y

	th := 1. / float32(gfx.GlyphMapRows)
	tw := glyphSize.W / (maxGlyphSize.W * float32(gfx.GlyphMapCols))

	return []float32{
		//vertex           texcoords         tile  grid
		-w / 2, +h / 2, 0, 0. + ox, 0. + oy, x, y, col, row,
		-w / 2, -h / 2, 0, 0. + ox, th + oy, x, y, col, row,
		+w / 2, -h / 2, 0, tw + ox, th + oy, x, y, col, row,
		+w / 2, -h / 2, 0, tw + ox, th + oy, x, y, col, row,
		+w / 2, +h / 2, 0, tw + ox, 0. + oy, x, y, col, row,
		-w / 2, +h / 2, 0, 0. + ox, 0. + oy, x, y, col, row,
	}

}

func (mode *LineMode) generateData(font *gfx.Font) {
	mode.data = []float32{}

	//    if DEBUG_LINEMODE { log.Debug("%s generate vertex data %s",mode.Desc(),font.Desc()) }

	width, height := int(mode.width), int(mode.height)

	for r := 0; r <= height; r++ {
		y := -1 * (r - height/2)

		row := r
		if mode.downward && mode.termBuffer == nil { // mode lines
			row = height - r
		}

		line := Line("")
		if mode.lineBuffer != nil {
			line = mode.lineBuffer.GetLine(uint(row))
		} else if mode.termBuffer != nil {
			line = mode.termBuffer.GetLine(uint(row))
		}

		for col := 0; col < width; col++ {
			x := col - width/2 + (width % 2)
			run := rune(' ')
			if col < len(line) {
				run = line[col]
			}

			gridCoord := gfx.Coord{col, row}
			tilePos := gfx.Point{float32(x), float32(y)}
			glyphCoord := getGlyphCoord(run)
			glyphSize := font.Size(glyphCoord.X, glyphCoord.Y)
			maxGlyphSize := font.MaxSize()

			tileSize := gfx.Size{
				W: glyphSize.W / glyphSize.H,
				H: glyphSize.H / glyphSize.H,
			}

			texOffset := gfx.Point{
				X: float32(glyphCoord.X) / (gfx.GlyphMapCols),
				Y: float32(glyphCoord.Y) / (gfx.GlyphMapRows),
			}

			mode.data = append(mode.data, mode.vertices(tileSize, glyphSize, maxGlyphSize, gridCoord, tilePos, texOffset)...)

		}

	}
	mode.object.BufferData(len(mode.data)*4, mode.data)

}

func getGlyphCoord(run rune) gfx.Coord {
	if run <= 0x20 || run >= 0x80 {
		return gfx.Coord{X: 0, Y: 0}
	}
	chr := byte(run)

	cols := byte(gfx.GlyphMapCols)

	col := chr % cols
	row := chr / cols
	return gfx.Coord{
		X: int(col),
		Y: int(row),
	}
}

func (mode *LineMode) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Info("%s init", mode.Desc())

	mode.object = gfx.NewObject("lines")
	mode.object.Init()

	mode.texture = gfx.NewTexture("lines")
	mode.texture.Init()

	mode.program = programService.GetProgram("lines", "lines/")
	err := mode.program.Link(mode.vert, mode.frag)
	if err != nil {
		log.Error("%s fail link program: %s", mode.Desc(), err)
	}

	mode.renderMap(font)

	mode.ScheduleRefresh()

}

func (mode *LineMode) renderMap(font *gfx.Font) error {

	if DEBUG_LINEMODE {
		log.Debug("%s render texture map %s", mode.Desc(), font.Desc())
	}

	rgba, err := font.RenderMap(false)
	if err != nil {
		log.Error("%s fail render font map: %s", mode.Desc(), err)
		return log.NewError("fail render font map: %s", err)
	}
	err = mode.texture.LoadRGBA(rgba)
	if err != nil {
		log.Error("%s fail load font map: %s", mode.Desc(), err)
		return log.NewError("fail to load font map: %s", err)
	}
	mode.texture.TexImage()

	return nil
}

func (mode *LineMode) autoScale(camera *gfx.Camera, font *gfx.Font) float32 {

	fontRatio := font.Ratio()
	screenRatio := camera.Ratio()

	ratio := screenRatio / fontRatio

	scaleWidth := ratio * 2. / float32(mode.width)
	scaleHeight := 2. / float32(mode.height)

	if scaleWidth < scaleHeight {
		return scaleWidth
	} else {
		return scaleHeight
	}

	return float32(1.0)
}

func (mode *LineMode) Configure(lines *LineConfig, term *TermConfig, shader *ShaderConfig, camera *gfx.Camera, font *gfx.Font) {

	s := ""
	if shader != nil {
		s = " " + shader.Desc()
	}

	if term != nil && mode.termBuffer != nil {
		term.autoSize(camera.Ratio(), font.Ratio())
		if DEBUG_LINEMODE {
			log.Debug("%s configure %s%s", mode.Desc(), term.Desc(), s)
		}
		changed := false

		if term.GetSetWidth() && term.GetWidth() != 0 && uint(term.GetWidth()) != mode.width {
			mode.width = uint(term.GetWidth())
			changed = true
		}

		if term.GetSetHeight() && term.GetHeight() != 0 && uint(term.GetHeight()) != mode.height {
			mode.height = uint(term.GetHeight())
			changed = true
		}

		if changed {
			mode.termBuffer.Resize(mode.width, mode.height)
		}

		if term.GetSetFill() {
			fillStr := mode.fill(term.GetFill())
			mode.termBuffer.Fill(fillStr)
		}
		mode.ScheduleRefresh()
	}

	if lines != nil && mode.lineBuffer != nil {
		lines.autoSize(camera.Ratio(), font.Ratio())
		if DEBUG_LINEMODE {
			log.Debug("%s configure %s%s", mode.Desc(), lines.Desc(), s)
		}
		changed := false

		if lines.GetSetWidth() && lines.GetWidth() != 0 && uint(lines.GetWidth()) != mode.width {
			mode.width = uint(lines.GetWidth())
			changed = true
		}

		if lines.GetSetHeight() && lines.GetHeight() != 0 && uint(lines.GetHeight()) != mode.height {
			mode.height = uint(lines.GetHeight())
			changed = true
		}

		if lines != nil && lines.GetSetBuffer() && uint(lines.GetBuffer()) != mode.buffer {
			mode.buffer = uint(lines.GetBuffer())
			changed = true
		}

		if changed {
			mode.lineBuffer.Resize(mode.height, mode.buffer)
		}

		if lines.GetSetDownward() {
			mode.downward = lines.GetDownward()
		}
		if lines.GetSetSpeed() {
			mode.lineBuffer.SetSpeed(float32(lines.GetSpeed()))
		}

		if lines.GetSetFixed() {
			mode.lineBuffer.Fixed = lines.GetFixed()
		}
		if lines.GetSetDrop() {
			mode.lineBuffer.Drop = lines.GetDrop()
		}

		if lines.GetSetSmooth() {
			mode.lineBuffer.Smooth = lines.GetSmooth()
		}

		if lines.GetSetFill() {
			fillStr := mode.fill(lines.GetFill())
			mode.lineBuffer.Fill(fillStr)
		}

		mode.ScheduleRefresh()
	}

	if shader != nil {
		changed := false
		vert, frag := mode.vert, mode.frag

		if shader.GetSetVert() {
			changed = true
			mode.vert = shader.GetVert()
		}

		if shader.GetSetFrag() {
			changed = true
			mode.frag = shader.GetFrag()
		}

		if changed {
			err := mode.program.Link(mode.vert, mode.frag)
			if err != nil {
				mode.vert = vert
				mode.frag = frag
			}
			mode.ScheduleRefresh()
		}
	}

}

func (mode *LineMode) DumpBuffer() string {
	if mode.termBuffer != nil {
		return mode.termBuffer.Dump()
	}
	if mode.lineBuffer != nil {
		return mode.lineBuffer.Dump(mode.width)
	}
	return ""
}

func NewLineMode(lineBuffer *LineBuffer, termBuffer *TermBuffer) *LineMode {
	ret := &LineMode{}

	if lineBuffer != nil {
		ret.lineBuffer = lineBuffer
		ret.width = uint(LineDefaults.GetWidth())
		ret.height = uint(LineDefaults.GetHeight())
		ret.downward = LineDefaults.GetDownward()
		ret.buffer = uint(LineDefaults.GetBuffer())

	} else if termBuffer != nil {
		ret.termBuffer = termBuffer
		ret.width = uint(TermDefaults.GetWidth())
		ret.height = uint(TermDefaults.GetHeight())
	}

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()
	ret.refreshChan = make(chan bool, 1)

	return ret
}

func (mode *LineMode) Desc() string {
	s := "↑"
	if mode.downward {
		s = "↓"
	}
	if mode.termBuffer != nil {
		x, y := mode.termBuffer.GetCursor()
		return fmt.Sprintf("term[%dx%d %d,%d]", mode.width, mode.height, x, y)
	}
	if mode.lineBuffer != nil {
		return fmt.Sprintf("lines[%dx%d+%d %s]", mode.width, mode.height, mode.lineBuffer.GetBuffer(), s)
	}

	return "lines[]"
}

func (mode *LineMode) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: mode.vert,
		SetFrag: true, Frag: mode.frag,
	}
	return ret
}

func (mode *LineMode) LineConfig() *LineConfig {
	ret := &LineConfig{
		SetWidth: true, Width: uint64(mode.width),
		SetHeight: true, Height: uint64(mode.height),
		SetDownward: true, Downward: mode.downward,
		SetBuffer: true, Buffer: uint64(mode.buffer),
	}
	if mode.lineBuffer != nil {
		ret.SetSpeed = true
		ret.Speed = float64(mode.lineBuffer.Speed())
		ret.SetFixed = true
		ret.Fixed = mode.lineBuffer.Fixed
	}
	return ret
}

func (mode *LineMode) TermConfig() *TermConfig {
	ret := &TermConfig{
		SetWidth: true, Width: uint64(mode.width),
		SetHeight: true, Height: uint64(mode.height),
	}
	return ret
}
