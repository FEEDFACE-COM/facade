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

const DEBUG_SCROLL = true

type Scroll struct {
	charBuffer *CharBuffer
	charCount  uint

	texture *gfx.Texture
	program *gfx.Program
	object  *gfx.Object
	data    []float32

	vert, frag  string
	refreshChan chan bool
}

const (
	CHARCOUNT gfx.UniformName = "charCount"
)

const (
	CHAROFFSET gfx.AttribName = "charOffset"
	CHARINDEX  gfx.AttribName = "charIndex"
)

func (scroll *Scroll) ScheduleRefresh() {

	select {
	case scroll.refreshChan <- true:
	default:
	}

}

func (scroll *Scroll) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-scroll.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}

func NewScroll(buffer *CharBuffer) *Scroll {
	ret := &Scroll{
		charBuffer: buffer,
	}

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()

	ret.refreshChan = make(chan bool, 1)
	return ret
}

func (scroll *Scroll) generateData(font *gfx.Font) {
	scroll.data = []float32{}

	line := scroll.charBuffer.GetLine()

	index := float32(0.)
	offset := float32(0.)

	data := []float32{}
	w := float32(0.)
	for _, run := range line {

		data, w = scroll.vertices(run, index, offset, font)
		scroll.data = append(scroll.data, data...)
		index += 1.
		offset += w
	}

	if DEBUG_SCROLL {
		log.Debug("%s generate %1.0f chars, wide:%.2f verts:%d floats:%d", scroll.Desc(), index, offset, len(scroll.data)/len(data), len(scroll.data))
	}
}

func (scroll *Scroll) vertices(
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

func (chars *Scroll) autoScale(camera *gfx.Camera) float32 {

	return 1.

}

func (scroll *Scroll) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	scroll.charBuffer.mutex.Lock()

	if scroll.checkRefresh() {
		if DEBUG_SCROLL {
			log.Debug("%s refresh", scroll.Desc())
		}
		scroll.generateData(font)
		scroll.renderMap(font)
	}

	//line := scroll.charBuffer.GetLine()
	charCount := scroll.charBuffer.charCount
	scroll.charBuffer.mutex.Unlock()

	gl.ActiveTexture(gl.TEXTURE0)

	scroll.program.UseProgram(debug)
	scroll.object.BindBuffer()

	// FIXME: verify nothing missing
	scroll.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	scroll.program.Uniform1f(gfx.FONTRATIO, font.Ratio())
	scroll.program.Uniform1f(gfx.CLOCKNOW, float32(gfx.Now()))
	scroll.program.Uniform1f(CHARCOUNT, float32(charCount))

	camera.Uniform(scroll.program)
	scale := float32(1.0)
	scale = scroll.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	scroll.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])

	scroll.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+1+1)*4, (0)*4)
	scroll.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+1+1)*4, (0+3)*4)
	scroll.program.VertexAttribPointer(CHARINDEX, 1, (3+2+1+1)*4, (0+3+2)*4)
	scroll.program.VertexAttribPointer(CHAROFFSET, 1, (3+2+1+1)*4, (0+3+2+1)*4)

	count := int32(charCount)
	offset := 0

	if DEBUG_SCROLL && verbose {
		log.Debug("%s render chars:%d verts:%d  ", scroll.Desc(), count, count*(2*3))
	}

	if !debug {
		scroll.program.SetDebug(false)
		scroll.texture.BindTexture()
		//gl.DrawArrays(gl.TRIANGLES, int32(offset*2*3), (count)*(2*3))
		scroll.program.SetDebug(debug)
	}

	if debug {
		scroll.program.SetDebug(true)
		gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		off := offset
		// REM, use single gl.DrawElements call instead (create indice array before)
		for c := int32(0); c < count; c++ {
			//gl.DrawArrays(gl.LINE_STRIP, int32(off*2*3), int32(1*2*3))
			off += 1
		}

	}

}

func (scroll *Scroll) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", scroll.Desc())

	scroll.object = gfx.NewObject("scroll")
	scroll.object.Init()

	scroll.texture = gfx.NewTexture("scroll")
	scroll.texture.Init()

	scroll.program = programService.GetProgram("scroll", "scroll/")
	scroll.program.Link(scroll.vert, scroll.frag)

	scroll.renderMap(font)

	scroll.ScheduleRefresh()

}

func (scroll *Scroll) renderMap(font *gfx.Font) error {

	if DEBUG_SCROLL {
		log.Debug("%s render texture map %s", scroll.Desc(), font.Desc())
	}

	rgba, err := font.RenderMap(false)
	if err != nil {
		log.Error("%s fail render font map: %s", scroll.Desc(), err)
		return log.NewError("fail render font map: %s", err)
	}
	err = scroll.texture.LoadRGBA(rgba)
	if err != nil {
		log.Error("%s fail load font map: %s", scroll.Desc(), err)
		return log.NewError("fail to load font map: %s", err)
	}
	scroll.texture.TexImage()

	return nil
}

func (scroll *Scroll) Configure(chars *CharConfig, camera *gfx.Camera, font *gfx.Font) {
	var shader *ShaderConfig = nil
	var config *ScrollConfig = nil

	log.Debug("%s configure", scroll.Desc())
	shader = chars.GetShader()
	config = chars.GetScroll()

	{
		changed := false
		vert, frag := scroll.vert, scroll.frag

		if shader != nil {

			if shader.GetSetVert() {
				changed = true
				scroll.vert = shader.GetVert()
			}

			if shader.GetSetFrag() {
				changed = true
				scroll.frag = shader.GetFrag()
			}
		}

		if changed {
			err := scroll.program.Link(scroll.vert, scroll.frag)
			if err != nil {
				scroll.vert = vert
				scroll.frag = frag
			}
		}
	}

	if config.GetSetCharCount() {
		scroll.charBuffer.Resize(uint(config.GetCharCount()))
	}

	if config.GetSetFill() {
		fillStr := scroll.fill(config.GetFill())
		if fillStr != "" {
			scroll.charBuffer.Fill(fillStr)
		}
	}

	scroll.ScheduleRefresh()

}

func (scroll *Scroll) fill(name string) string {
	switch name {
	case "title":
		return "FACADE"
	case "index":
		ret := ""
		for i := 0; uint(i) < scroll.charBuffer.charCount; i++ {
			if i%10 == 0 {
				ret += "#"
			} else {
				ret += fmt.Sprintf("%1d", i%10)
			}
		}
		return ret
	case "alpha":
		ret := ""
		alpha := []string{
			"alpha",
			"beta",
			"gamma",
			"delta",
			"epsilon",
			"zeta",
			"eta",
			"theta",
			"iota",
			"kappa",
			"lambda",
			"mu",
			"nu",
			"xi",
			"omicron",
			"pi",
			"rho",
			"sigma",
			"tau",
			"upsilon",
			"phi",
			"chi",
			"psi",
			"omega",
		}
		l := uint(0)
		for i := 0; l < scroll.charBuffer.charCount && i < len(alpha); i++ {
			ret += alpha[i] + " "
		}
		return ret

	default:
		log.Error("no such charbuffer fill pattern: '%s'", name)
	}
	return ""

}

func (scroll *Scroll) Desc() string {
	ret := "scroll["
	ret += scroll.charBuffer.Desc()
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (scroll *Scroll) Config() *ScrollConfig {
	ret := &ScrollConfig{
		SetRepeat: true, Repeat: bool(scroll.charBuffer.Repeat()),
		SetCharCount: true, CharCount: uint64(scroll.charBuffer.CharCount()),
	}
	return ret

}

func (scroll *Scroll) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: scroll.vert,
		SetFrag: true, Frag: scroll.frag,
	}
	return ret
}
