// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"

	//    "fmt"
)

const DEBUG_SCROLL = true



type Scroll struct {

	charBuffer *CharBuffer


	texture *gfx.Texture
	program *gfx.Program
	object  *gfx.Object
	data    []float32
    


	vert, frag string
	refreshChan chan bool
    
}

const (


)

const (


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

	//ret.line = Line{}
	//    ret.textures = make( map[string] *gfx.Texture, ret.wordBuffer.SlotCount() )

	ret.refreshChan = make(chan bool, 1)
	return ret
}


func (scroll *Scroll) generateData(font *gfx.Font) {

	FIXME := 0

	scroll.data = make([]float32,FIXME)


	//    old := set.textures
	//
	//    set.textures = make( map[string] *gfx.Texture, set.wordBuffer.SlotCount())



//	chars.line = chars.lineBuffer.GetLine(0)

	if DEBUG_SCROLL {
		log.Debug("%s generate chars", scroll.Desc())
	}

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


	chars := scroll.charBuffer.Chars()
	scroll.charBuffer.mutex.Unlock()

	gl.ActiveTexture(gl.TEXTURE0)

	scroll.program.UseProgram(debug)
	scroll.object.BindBuffer()

	// FIXME: verify nothing missing
	scroll.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	scroll.program.Uniform1f(gfx.FONTRATIO, font.Ratio())
	scroll.program.Uniform1f(gfx.CLOCKNOW, float32(gfx.Now()))

	camera.Uniform(scroll.program)
	scale := float32(1.0)
	scale = scroll.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	scroll.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])
	scroll.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+1+1+1+1+1)*4, (0)*4)
	scroll.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+1+1+1+1+1)*4, (0+3)*4)
	scroll.program.VertexAttribPointer(CHARINDEX, 1, (3+2+1+1+1+1+1)*4, (0+3+2)*4)



	if DEBUG_SCROLL && verbose {
			log.Debug("%s render %d chars.", scroll.Desc(),len(chars))
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
			//			err := set.LoadShaders()
			if err != nil {
				scroll.vert = vert
				scroll.frag = frag
			}
		}
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
		for i:=0; i<maxCharCount;i++ {
			if i%10 == 0 {
				ret += "#"
			} else {
				ret += fmt.Sprintf("%1d",i)
			}
		}
		return ret
	default:
		log.Error("no such charbuffer fill pattern: '%s'", name)
	}
	return ""

	}

func (scroll *Scroll) Desc() string {
	ret := "scroll["

	//if scroll.charBuf.wordBuffer.shuffle || set.wordBuffer.aging {
	//	ret += " "
	//	if set.wordBuffer.shuffle {
	//		ret += "⧢"
	//	}
	//	if set.wordBuffer.aging {
	//		ret += "å"
	//	}
	//}


	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (scroll *Scroll) Config() *ScrollConfig {
	ret := &ScrollConfig{
		//SetRepeat: true, Repeat: bool(scroll.Repeat()),
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


