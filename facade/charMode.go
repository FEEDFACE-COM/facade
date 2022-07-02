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
	//    "fmt"
)

const DEBUG_CHARMODE = true

type CharMode struct {
	charBuffer *CharBuffer
	charCount  uint    // chars in data
	charWidth  float32 // total width of chars in data

	texture *gfx.Texture
	program *gfx.Program
	object  *gfx.Object
	data    []float32

	vert, frag  string
	refreshChan chan bool
}

const (
	CHARCOUNT gfx.UniformName = "charCount"
	CHARWIDTH gfx.UniformName = "charWidth"
)

const (
	CHAROFFSET gfx.AttribName = "charOffset"
	CHARINDEX  gfx.AttribName = "charIndex"
)

func (mode *CharMode) ScheduleRefresh() {

	select {
	case mode.refreshChan <- true:
	default:
	}

}

func (mode *CharMode) checkRefresh() bool {
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

func NewCharMode(buffer *CharBuffer) *CharMode {
	ret := &CharMode{
		charBuffer: buffer,
		charCount:  0,
	}

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()

	ret.refreshChan = make(chan bool, 1)
	return ret
}

func (mode *CharMode) generateData(font *gfx.Font) {
	mode.data = []float32{}

	line := mode.charBuffer.GetLine()

	index := float32(0.)
	offset := float32(0.)

	data := []float32{}
	w := float32(0.)
	for _, run := range line {

		data, w = mode.vertices(run, index, offset, font)
		mode.data = append(mode.data, data...)
		index += 1.
		offset += w
	}

	mode.object.BufferData(len(mode.data)*4, mode.data)
	mode.charCount = uint(index)
	mode.charWidth = float32(offset)
	if DEBUG_CHARMODE {
		log.Debug("%s generate chars:%d width:%.1f verts:%d floats:%d", mode.Desc(), mode.charCount, mode.charWidth, 6*mode.charCount, len(mode.data))
	}
}

func (mode *CharMode) vertices(
	run rune,
	index float32,
	offset float32,
	font *gfx.Font,
) ([]float32, float32) {

	glyphCoord := getGlyphCoord(run)
	glyphSize := font.Size(glyphCoord.X, glyphCoord.Y)
	maxGlyphSize := font.MaxSize()

	texOffset := gfx.Point{
		X: float32(glyphCoord.X) / (gfx.GlyphMapCols),
		Y: float32(glyphCoord.Y) / (gfx.GlyphMapRows),
	}

	ox, oy := texOffset.X, texOffset.Y
	th := 1. / float32(gfx.GlyphMapRows)
	tw := glyphSize.W / (maxGlyphSize.W * float32(gfx.GlyphMapCols))

	w := glyphSize.W / glyphSize.H
	h := float32(1.)
	//log.Debug("got coord %d/%d size %fx%f wh %fx%f",glyphCoord.X,glyphCoord.Y,glyphSize.W,glyphSize.H,w,h)

	/*

	     A          D
	 -w/2,h/2____w/2,h/2
	     |          |
	     |          |
	 -w/2,-h/2___w/2,-h/2
	     B          C


	     A          D
	    0,0________1,0
	     |          |
	     |          |
	    0,1________1,1
	     B          C


	  A     A_D
	  |\    \ |
	  |_\    \|
	  B C     C

	*/

	data := []float32{
		//  x,     y,   z,      tx,      ty, idx, offset,
		-w / 2., +h / 2., 0.0, 0. + ox, 0. + oy, index, offset, // A
		-w / 2., -h / 2., 0.0, 0. + ox, th + oy, index, offset, // B
		+w / 2., -h / 2., 0.0, tw + ox, th + oy, index, offset, // C
		+w / 2., -h / 2., 0.0, tw + ox, th + oy, index, offset, // C
		+w / 2., +h / 2., 0.0, tw + ox, 0. + oy, index, offset, // D
		-w / 2., +h / 2., 0.0, 0. + ox, 0. + oy, index, offset, // A
	}

	return data, w

}

func (mode *CharMode) autoScale(camera *gfx.Camera) float32 {

	return 1.

}

func (mode *CharMode) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	mode.charBuffer.mutex.Lock()

	if mode.checkRefresh() {
		if DEBUG_CHARMODE {
			log.Debug("%s refresh", mode.Desc())
		}
		mode.generateData(font)
		mode.renderMap(font)
	}

	//line := mode.charBuffer.GetLine()
	charCount := mode.charCount
	charWidth := mode.charWidth
	mode.charBuffer.mutex.Unlock()

	gl.ActiveTexture(gl.TEXTURE0)
	mode.program.UseProgram(debug)
	mode.object.BindBuffer()

	mode.program.Uniform1f(CHARCOUNT, float32(charCount))
	mode.program.Uniform1f(CHARWIDTH, float32(charWidth))

	mode.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	mode.program.Uniform1f(gfx.FONTRATIO, font.Ratio())
	mode.program.Uniform1f(gfx.CLOCKNOW, float32(gfx.Now()))

	camera.Uniform(mode.program)
	scale := float32(1.0)
	scale = mode.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	mode.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])

	mode.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+1+1)*4, (0)*4)
	mode.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+1+1)*4, (0+3)*4)
	mode.program.VertexAttribPointer(CHARINDEX, 1, (3+2+1+1)*4, (0+3+2)*4)
	mode.program.VertexAttribPointer(CHAROFFSET, 1, (3+2+1+1)*4, (0+3+2+1)*4)

	count := int32(charCount)
	offset := 0

	//if DEBUG_CHARMODE && verbose {
	//	log.Debug("%s render chars:%d verts:%d  ", mode.Desc(), count, count*(2*3))
	//}

	if charCount <= 0. {
		return
	}
	mode.texture.Uniform(mode.program)

	if !debug || debug {
		mode.program.SetDebug(false)
		mode.texture.BindTexture()
		gl.DrawArrays(gl.TRIANGLES, int32(offset*2*3), (count)*(2*3))
		mode.program.SetDebug(debug)
	}

	if debug {
		mode.program.SetDebug(true)
		gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		off := offset
		// REM, use single gl.DrawElements call instead (create indice array before)
		for c := int32(0); c < count; c++ {
			gl.DrawArrays(gl.LINE_STRIP, int32(off*2*3), int32(1*2*3))
			off += 1
		}

	}

}

func (mode *CharMode) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", mode.Desc())

	mode.object = gfx.NewObject("chars")
	mode.object.Init()

	mode.texture = gfx.NewTexture("chars")
	mode.texture.Init()

	mode.program = programService.GetProgram("chars", "chars/")
	mode.program.Link(mode.vert, mode.frag)

	mode.renderMap(font)
	mode.ScheduleRefresh()

}

func (mode *CharMode) renderMap(font *gfx.Font) error {

	if DEBUG_CHARMODE {
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

func (mode *CharMode) Configure(config *CharConfig, shader *ShaderConfig, camera *gfx.Camera, font *gfx.Font) {
	s := ""
	if shader != nil {
		s = " " + shader.Desc()
	}
	log.Debug("%s configure %s%s", mode.Desc(), config.Desc(), s)

	if shader != nil {
		changed := false
		vert, frag := mode.vert, mode.frag

		if shader != nil {

			if shader.GetSetVert() {
				changed = true
				mode.vert = shader.GetVert()
			}

			if shader.GetSetFrag() {
				changed = true
				mode.frag = shader.GetFrag()
			}
		}

		if changed {
			err := mode.program.Link(mode.vert, mode.frag)
			if err != nil {
				mode.vert = vert
				mode.frag = frag
			}
		}
		mode.ScheduleRefresh()
	}

	if config != nil {

		if config.GetSetCharCount() {
			mode.charBuffer.Resize(uint(config.GetCharCount()))
		}

		if config.GetSetFill() {
			fillStr := mode.fill(config.GetFill())
			if fillStr != "" {
				mode.charBuffer.Fill(fillStr)
			}
		}
		mode.ScheduleRefresh()
	}

}

func (mode *CharMode) fill(name string) string {
	switch name {
	case "title":
		return "FACADE"
	case "index":
		ret := ""
		for i := 0; uint(i) < mode.charBuffer.charCount; i++ {
			if i%10 == 0 {
				ret += "#"
			} else {
				ret += fmt.Sprintf("%1d", i%10)
			}
		}
		return ret
	case "alpha":
		ret := ""
		alpha := " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
		d := uint(len(alpha))
		for i := uint(0); i < mode.charBuffer.charCount; i++ {
			ret += fmt.Sprintf("%c", alpha[i%d])
		}
		return ret

	default:
		log.Error("no such charbuffer fill pattern: '%s'", name)
	}
	return ""

}

func (mode *CharMode) Desc() string {
	ret := "chars["
	ret += mode.charBuffer.Desc()
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (mode *CharMode) Config() *CharConfig {
	ret := &CharConfig{
		SetRepeat: true, Repeat: bool(mode.charBuffer.Repeat()),
		SetCharCount: true, CharCount: uint64(mode.charBuffer.CharCount()),
	}
	return ret

}

func (mode *CharMode) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: mode.vert,
		SetFrag: true, Frag: mode.frag,
	}
	return ret
}
