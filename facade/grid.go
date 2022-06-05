//go:build (linux && arm) || DARWIN_GUI
// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"fmt"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

const DEBUG_GRID = false

type Grid struct {
	width, height uint

	downward bool

	buffer   uint
	terminal bool

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
	DOWNWARD   gfx.UniformName = "downward"
)

const (
	TILECOORD gfx.AttribName = "tileCoord"
	GRIDCOORD gfx.AttribName = "gridCoord"
)

func (grid *Grid) ScheduleRefresh() {

	select {
	case grid.refreshChan <- true:
	default:
	}

}

func (grid *Grid) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-grid.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}

func (grid *Grid) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	if grid.checkRefresh() {
		if DEBUG_GRID {
			log.Debug("%s refresh", grid.Desc())
		}
		grid.generateData(font)
		grid.renderMap(font) // TODO: should not renderMap anytime anything changed!
	}

	gl.ActiveTexture(gl.TEXTURE0)
	grid.program.UseProgram(debug)
	grid.object.BindBuffer()

	tileCount := mgl32.Vec2{float32(grid.width), float32(grid.height)}
	grid.program.Uniform2fv(TILECOUNT, 1, &tileCount[0])

	tileSize := mgl32.Vec2{font.MaxSize().W / font.MaxSize().H, font.MaxSize().H / font.MaxSize().H}
	grid.program.Uniform2fv(TILESIZE, 1, &tileSize[0])

	tileOffset := mgl32.Vec2{-1., 0.0}
	if grid.width%2 == 0 { //even columns
		tileOffset[0] = 0.5
	}
	if grid.height%2 == 0 { //even rows
		tileOffset[1] = -0.5
	}
	if grid.downward && !grid.terminal {
		tileOffset[1] += 1.
	}
	grid.program.Uniform2fv(TILEOFFSET, 1, &tileOffset[0])

	cursorPos := mgl32.Vec2{-1., -1.}
	if grid.terminal {
		x, y := grid.termBuffer.GetCursor()
		cursorPos[0] = float32(x)
		cursorPos[1] = float32(y)
	}
	grid.program.Uniform2fv(CURSORPOS, 1, &cursorPos[0])

	grid.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	grid.program.Uniform1f(gfx.FONTRATIO, font.Ratio())

	clocknow := float32(gfx.Now())
	grid.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow)

	scroller := float32(0.0)
	if !grid.terminal {
		scroller = float32(grid.lineBuffer.GetScroller())
		if grid.downward {
			scroller *= -1.
		}
	}
	grid.program.Uniform1f(SCROLLER, scroller)

	camera.Uniform(grid.program)

	grid.texture.Uniform(grid.program)

	scale := float32(1.0)
	scale = grid.autoScale(camera, font)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])

	grid.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+2+2)*4, (0)*4)
	grid.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+2+2)*4, (3)*4)
	grid.program.VertexAttribPointer(TILECOORD, 2, (3+2+2+2)*4, (3+2)*4)
	grid.program.VertexAttribPointer(GRIDCOORD, 2, (3+2+2+2)*4, (3+2+2)*4)

	count := int32(grid.width * (grid.height + 1))
	offset := int32(0)

	if !debug || debug {
		grid.program.SetDebug(false)
		grid.texture.BindTexture()
		gl.DrawArrays(gl.TRIANGLES, int32(offset*2*3), (count)*(2*3))
		grid.program.SetDebug(debug)
	}

	if debug {
		gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		off := offset
		// REM, use single gl.DrawElements call instead (create indice array before)
		for r := 0; r < int(grid.height+1); r++ {
			for c := 0; c < int(grid.width); c++ {
				gl.DrawArrays(gl.LINE_STRIP, int32(off*2*3), int32(1*2*3))
				off += int32(1)
			}
			//    	   gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(grid.Width*2*3) )
			//    	   off += int32(grid.Width)
		}

	}
}

func (grid *Grid) fill(name string) []string {
	ret := []string{}

	switch name {
	case "title":
		title := []string{""}

		if grid.width >= 80 && grid.height >= 16 {
			title = strings.Split(`
 _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _      
|_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|
|   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |
`, "\n")[1:]
		} else if grid.width >= 40 && grid.height >= 8 {
			title = strings.Split(`
 _   _   _   _   _   _
|_  |_| /   |_| | \ |_
|   | | \_  | | |_/ |_
                      
       by FEEDFACE.COM
`, "\n")[1:]
		} else if grid.width >= 13 {
			title = []string{"F A C A D E", ""}
		} else if grid.width >= 8 {
			title = []string{"FACADE", ""}
		}

		gw, gh := int(grid.width), int(grid.height)
		tw, th := len(title[0]), len(title)-1
		dw, dh := ((gw - tw) / 2.), ((gh - th) / 2.)

		if DEBUG_GRID {
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

	case "grid":
		w, h := int(grid.width), int(grid.height)
		for r := 0; r < h; r++ {
			tmp := ""
			for c := 0; c < w; c++ {
				d := "."
				if c%5 == 0 {
					d = fmt.Sprintf("%d", r%10)
				}
				if r%5 == 0 {
					d = fmt.Sprintf("%d", c%10)
				}
				if c%5 == 0 && r%5 == 0 {
					d = "#"
				}
				tmp += fmt.Sprintf("%s", d)
			}
			ret = append(ret, tmp)
		}
		return ret

	case "alpha":
		w, h := int(grid.width), int(grid.height)
		alpha := strings.Repeat(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", w)
		s := 0
		for r := 0; r < h; r++ {
			ret = append(ret, alpha[s:s+w])
			s += w
		}
		return ret

	case "clear":
		w, h := int(grid.width), int(grid.height)
		for r := 0; r < h; r++ {
			s := ""
			for c := 0; c < w; c++ {
				s += " "
			}
			ret = append(ret, s)
		}
		return ret

	}

	log.Warning("%s no such fill pattern: %s", grid.Desc(), name)
	return ret

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (grid *Grid) vertices(
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

func (grid *Grid) generateData(font *gfx.Font) {
	grid.data = []float32{}

	//    if DEBUG_GRID { log.Debug("%s generate vertex data %s",grid.Desc(),font.Desc()) }

	width, height := int(grid.width), int(grid.height)

	for r := 0; r <= height; r++ {
		y := -1 * (r - height/2)

		row := r
		if grid.downward && !grid.terminal {
			row = height - r
		}

		line := Line("")
		if grid.terminal {
			line = grid.termBuffer.GetLine(uint(row))
		} else {
			line = grid.lineBuffer.GetLine(uint(row))
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

			grid.data = append(grid.data, grid.vertices(tileSize, glyphSize, maxGlyphSize, gridCoord, tilePos, texOffset)...)

		}

	}
	grid.object.BufferData(len(grid.data)*4, grid.data)

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

func (grid *Grid) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Info("%s init", grid.Desc())

	grid.object = gfx.NewObject("grid")
	grid.object.Init()

	grid.texture = gfx.NewTexture("grid")
	grid.texture.Init()

	grid.program = programService.GetProgram("grid", "grid/")
	err := grid.program.Link(grid.vert, grid.frag)
	if err != nil {
		log.Error("%s fail link program: %s", grid.Desc(), err)
	}

	grid.renderMap(font)

	grid.ScheduleRefresh()

}

func (grid *Grid) renderMap(font *gfx.Font) error {

	if DEBUG_GRID {
		log.Debug("%s render texture map %s", grid.Desc(), font.Desc())
	}

	rgba, err := font.RenderMap(false)
	if err != nil {
		log.Error("%s fail render font map: %s", grid.Desc(), err)
		return log.NewError("fail render font map: %s", err)
	}
	err = grid.texture.LoadRGBA(rgba)
	if err != nil {
		log.Error("%s fail load font map: %s", grid.Desc(), err)
		return log.NewError("fail to load font map: %s", err)
	}
	grid.texture.TexImage()

	return nil
}

func (grid *Grid) autoScale(camera *gfx.Camera, font *gfx.Font) float32 {

	fontRatio := font.Ratio()
	screenRatio := camera.Ratio()

	ratio := screenRatio / fontRatio

	scaleWidth := ratio * 2. / float32(grid.width)
	scaleHeight := 2. / float32(grid.height)

	if scaleWidth < scaleHeight {
		return scaleWidth
	} else {
		return scaleHeight
	}

	return float32(1.0)
}

func (grid *Grid) Configure(lines *LineConfig, terminal *TermConfig, camera *gfx.Camera, font *gfx.Font) {
	var shader *ShaderConfig = nil
	var config *GridConfig = nil

	if grid.terminal && terminal != nil {

		if DEBUG_GRID {
			log.Debug("%s configure %s", grid.Desc(), terminal.Desc())
		}
		shader = terminal.GetShader()
		config = terminal.GetGrid()

	} else if !grid.terminal && lines != nil {

		if DEBUG_GRID {
			log.Debug("%s configure %s", grid.Desc(), lines.Desc())
		}
		shader = lines.GetShader()
		config = lines.GetGrid()

		if lines.GetSetDownward() {
			grid.downward = lines.GetDownward()
		}
		if lines.GetSetSpeed() {
			grid.lineBuffer.SetSpeed(float32(lines.GetSpeed()))
		}

		if lines.GetSetFixed() {
			grid.lineBuffer.Fixed = lines.GetFixed()
		}
		if lines.GetSetDrop() {
			grid.lineBuffer.Drop = lines.GetDrop()
		}

		if lines.GetSetSmooth() {
			grid.lineBuffer.Smooth = lines.GetSmooth()
		}

	} else {
		log.Warning("%s cannot configure", grid.Desc())
		return
	}

	// if true { //optimize!!
	// 	log.Debug("%s rendermap %s", grid.Desc(), font.Desc())
	// 	grid.renderMap(font)
	// }

	config.autoSize(camera.Ratio(), font.Ratio())

	{
		changed := false
		if config.GetSetWidth() && config.GetWidth() != 0 && uint(config.GetWidth()) != grid.width {
			grid.width = uint(config.GetWidth())
			changed = true
		}

		if config.GetSetHeight() && config.GetHeight() != 0 && uint(config.GetHeight()) != grid.height {
			grid.height = uint(config.GetHeight())
			changed = true
		}

		if lines != nil && lines.GetSetBuffer() && uint(lines.GetBuffer()) != grid.buffer {
			grid.buffer = uint(lines.GetBuffer())
			changed = true
		}

		if changed {
			if grid.termBuffer != nil {
				grid.termBuffer.Resize(grid.width, grid.height)
			}
			if grid.lineBuffer != nil {
				grid.lineBuffer.Resize(grid.height, grid.buffer)
			}
		}
	}

	{
		changed := false
		vert, frag := grid.vert, grid.frag

		if shader != nil {

			if shader.GetSetVert() {
				changed = true
				grid.vert = shader.GetVert()
			}

			if shader.GetSetFrag() {
				changed = true
				grid.frag = shader.GetFrag()
			}
		}

		if changed {
			err := grid.program.Link(grid.vert, grid.frag)
			//			err := grid.LoadShaders()
			if err != nil {
				grid.vert = vert
				grid.frag = frag
			}
		}
	}

	if config.GetSetFill() {
		fillStr := grid.fill(config.GetFill())
		if grid.termBuffer != nil {
			grid.termBuffer.Fill(fillStr)
		}
		if grid.lineBuffer != nil {
			grid.lineBuffer.Fill(fillStr)
		}
	}

	grid.ScheduleRefresh()

}

func (grid *Grid) DescBuffer() string {
	ret := ""
	if grid.termBuffer != nil {
		ret += grid.termBuffer.Desc()
	}
	if grid.lineBuffer != nil {
		ret += grid.lineBuffer.Desc()
	}
	return ret
}

func (grid *Grid) DumpBuffer() string {
	ret := ""
	if grid.termBuffer != nil {
		ret += grid.termBuffer.Dump()
	}
	if grid.lineBuffer != nil {
		ret += grid.lineBuffer.Dump(grid.width)
	}
	return ret
}

func NewGrid(lineBuffer *LineBuffer, termBuffer *TermBuffer) *Grid {
	ret := &Grid{}
	ret.width = uint(GridDefaults.GetWidth())
	ret.height = uint(GridDefaults.GetHeight())

	ret.downward = LineDefaults.GetDownward()

	ret.buffer = uint(LineDefaults.GetBuffer())

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()

	ret.refreshChan = make(chan bool, 1)
	ret.lineBuffer = lineBuffer
	ret.termBuffer = termBuffer

	if termBuffer != nil {
		ret.terminal = true
	} else {
		ret.terminal = false
	}

	return ret
}

func (grid *Grid) Desc() string {
	return grid.Config().Desc()
}

func (grid *Grid) Config() *GridConfig {
	return &GridConfig{
		SetWidth: true, Width: uint64(grid.width),
		SetHeight: true, Height: uint64(grid.height),
	}
}

func (grid *Grid) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: grid.vert,
		SetFrag: true, Frag: grid.frag,
	}
	return ret
}

func (grid *Grid) LineConfig() *LineConfig {
	ret := &LineConfig{
		SetDownward: true, Downward: grid.downward,
		SetBuffer: true, Buffer: uint64(grid.buffer),
	}
	if grid.lineBuffer != nil {
		ret.SetSpeed = true
		ret.Speed = float64(grid.lineBuffer.Speed())
		ret.SetFixed = true
		ret.Fixed = grid.lineBuffer.Fixed
	}
	return ret
}
