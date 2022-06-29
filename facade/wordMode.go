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
	"unicode/utf8"
)

const DEBUG_WORDMODE = false

const HARD_MAX_LENGTH = 80.0

type WordMode struct {
	maxLength  uint
	wordBuffer *WordBuffer

	maxWord Word

	texture *gfx.Texture
	program *gfx.Program
	object  *gfx.Object
	data    []float32

	vert, frag  string
	refreshChan chan bool
}

const (
	WORDCOUNT     gfx.UniformName = "wordCount"
	WORDMAXWIDTH  gfx.UniformName = "wordMaxWidth"
	WORDMAXLENGTH gfx.UniformName = "wordMaxLength"
	WORDFADER     gfx.UniformName = "wordFader"
	WORDAGE       gfx.UniformName = "wordAge"
)

const (
	WORDINDEX  gfx.AttribName = "wordIndex"
	WORDWIDTH  gfx.AttribName = "wordWidth"
	WORDLENGTH gfx.AttribName = "wordLength"
)

func (mode *WordMode) ScheduleRefresh() {

	select {
	case mode.refreshChan <- true:
	default:
	}

}

func (mode *WordMode) checkRefresh() bool {
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

func NewWordMode(buffer *WordBuffer) *WordMode {
	ret := &WordMode{
		wordBuffer: buffer,
		maxWord:    Word{},
	}

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()

	ret.refreshChan = make(chan bool, 1)
	return ret
}

func (mode *WordMode) generateData(font *gfx.Font) {

	//setup vertex + bind order arrays
	mode.data = []float32{}

	words := mode.wordBuffer.GetWords()
	maxWidth := float32(0.)
	maxLength := uint(0)
	charCount := 0
	for _, word := range words {

		word.width = float32(0.)
		word.length = 0
		for _, run := range word.text {
			glyphCoord := getGlyphCoord(run)
			glyphSize := font.Size(glyphCoord.X, glyphCoord.Y)
			word.width += glyphSize.W / glyphSize.H
			word.length += 1
			charCount += 1
		}
		if word.width > maxWidth {
			maxWidth = word.width
			maxLength = word.length
		}
		if word.length > maxLength {
			maxLength = word.length
			maxWidth = word.width
		}
		mode.data = append(mode.data, mode.vertices(word.text, float32(word.index), float32(word.length), float32(word.width), font)...)
	}

	if mode.maxLength == 0 { // actually counted max values
		mode.maxWord = Word{length: maxLength, width: maxWidth}
	} else { // config max length and font ratio
		mode.maxWord = Word{length: mode.maxLength, width: float32(mode.maxLength) * font.Ratio()}
	}

	mode.object.BufferData(len(mode.data)*4, mode.data)
	if DEBUG_WORDMODE {
		log.Debug("%s generate words:%d chars:%d floats:%d, max length:%d width:%.1f", mode.Desc(), len(words), charCount, len(mode.data), maxLength, maxWidth)
	}

}

func (mode *WordMode) vertices(
	text string,
	wordIndex float32,
	wordLength float32,
	wordWidth float32,
	font *gfx.Font,
) []float32 {

	var ret = []float32{}

	charIndex := 0
	offset := -wordWidth / 2.
	for _, run := range text {

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

		offset += w / 2.

		//        dx := offset + w/2.

		idx := float32(charIndex)
		off := float32(offset)

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
			//        x,       y,   z,      tx,      ty,
			-w/2. + off, +h / 2., 0.0, 0. + ox, 0. + oy, idx, off, wordIndex, wordWidth, wordLength, // A
			-w/2. + off, -h / 2., 0.0, 0. + ox, th + oy, idx, off, wordIndex, wordWidth, wordLength, // B
			+w/2. + off, -h / 2., 0.0, tw + ox, th + oy, idx, off, wordIndex, wordWidth, wordLength, // C
			+w/2. + off, -h / 2., 0.0, tw + ox, th + oy, idx, off, wordIndex, wordWidth, wordLength, // C
			+w/2. + off, +h / 2., 0.0, tw + ox, 0. + oy, idx, off, wordIndex, wordWidth, wordLength, // D
			-w/2. + off, +h / 2., 0.0, 0. + ox, 0. + oy, idx, off, wordIndex, wordWidth, wordLength, // A
		}
		ret = append(ret, data...)

		offset += w / 2.
		charIndex += 1
	}

	if DEBUG_WORDMODE {
		//log.Debug("%s data generate '%s'", mode.Desc(), text)
	}

	return ret
}

func (mode *WordMode) autoScale(camera *gfx.Camera) float32 {
	//scaleHeight := float32(1.) / math32.Sqrt( float32(mode.wordBuffer.SlotCount() ) )
	//scaleHeight := float32(1.) / float32(mode.wordBuffer.SlotCount())
	//scaleHeight := float32(1.) / float32(8.)
	scaleHeight := float32(1.)
	return scaleHeight //* 2.
}

func (mode *WordMode) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {
	mode.wordBuffer.mutex.Lock()

	if mode.checkRefresh() {
		if DEBUG_WORDMODE {
			log.Debug("%s refresh", mode.Desc())
		}
		mode.generateData(font)
		mode.renderMap(font)
	}

	words := mode.wordBuffer.GetWords()
	wordCount := float32(mode.wordBuffer.SlotCount())
	mode.wordBuffer.mutex.Unlock()

	gl.ActiveTexture(gl.TEXTURE0)

	mode.program.UseProgram(debug)
	mode.object.BindBuffer()

	mode.program.Uniform1f(WORDCOUNT, wordCount)

	mode.program.Uniform1f(WORDMAXLENGTH, float32(mode.maxWord.length))
	mode.program.Uniform1f(WORDMAXWIDTH, mode.maxWord.width)
	mode.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	mode.program.Uniform1f(gfx.FONTRATIO, font.Ratio())
	mode.program.Uniform1f(gfx.CLOCKNOW, float32(gfx.Now()))

	camera.Uniform(mode.program)
	scale := float32(1.0)
	scale = mode.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	mode.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])
	mode.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+1+1+1+1+1)*4, (0)*4)
	mode.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+1+1+1+1+1)*4, (0+3)*4)
	mode.program.VertexAttribPointer(CHARINDEX, 1, (3+2+1+1+1+1+1)*4, (0+3+2)*4)
	mode.program.VertexAttribPointer(CHAROFFSET, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1)*4)
	mode.program.VertexAttribPointer(WORDINDEX, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1)*4)
	mode.program.VertexAttribPointer(WORDWIDTH, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1+1)*4)
	mode.program.VertexAttribPointer(WORDLENGTH, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1+1+1)*4)

	count := int32(1)
	offset := int32(0)

	for _, word := range words {

		text := word.text

		if len(text) <= 0 {
			continue
		}

		mode.texture.Uniform(mode.program)
		count = int32(utf8.RuneCountInString(word.text))

		fader := float32(1.0)
		if word.fader != nil {
			fader = word.fader.Value()
		}
		mode.program.Uniform1f(WORDFADER, fader)

		age := float32(0.0)
		if word.timer != nil {
			age = word.timer.Value()
		}
		mode.program.Uniform1f(WORDAGE, age)

		if !debug || debug {
			mode.program.SetDebug(false)
			mode.texture.BindTexture()
			gl.DrawArrays(gl.TRIANGLES, int32(offset*(2*3)), int32(count*2*3))
			mode.program.SetDebug(debug)
		}

		if debug {
			gl.LineWidth(3.0)
			gl.BindTexture(gl.TEXTURE_2D, 0)
			off := offset
			for i := 0; i < int(count); i++ {
				gl.DrawArrays(gl.LINE_STRIP, int32(off*2*3), int32(1*2*3))
				off += int32(1)
			}
		}
		offset += count
	}
	if DEBUG_WORDMODE && verbose {
		//log.Debug("%s render %d words", mode.Desc(), len(words) )
	}
}

func (mode *WordMode) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", mode.Desc())

	mode.object = gfx.NewObject("words")
	mode.object.Init()

	mode.texture = gfx.NewTexture("words")
	mode.texture.Init()

	mode.program = programService.GetProgram("words", "words/")
	mode.program.Link(mode.vert, mode.frag)

	mode.renderMap(font)
	mode.ScheduleRefresh()

}

func (mode *WordMode) renderMap(font *gfx.Font) error {

	if DEBUG_WORDMODE {
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

func (mode *WordMode) Configure(config *WordConfig, shader *ShaderConfig, camera *gfx.Camera, font *gfx.Font) {

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
			mode.ScheduleRefresh()
		}
	}

	if config != nil {

		if config.GetSetSlots() {
			mode.wordBuffer.Resize(int(config.GetSlots()))
		}

		if config.GetSetLifetime() {
			mode.wordBuffer.SetLifetime(float32(config.GetLifetime()))
		}

		if config.GetSetWatermark() {
			mode.wordBuffer.SetWatermark(float32(config.GetWatermark()))
		}

		if config.GetSetShuffle() {
			mode.wordBuffer.SetShuffle(config.GetShuffle())
		}

		if config.GetSetMaxLength() {
			mode.maxLength = uint(config.GetMaxLength())
		}

		if config.GetSetAging() {
			mode.wordBuffer.SetAging(config.GetAging())
		}

		if config.GetSetFill() {
			fillStr := mode.fill(config.GetFill())
			if fillStr != nil {
				mode.wordBuffer.Fill(fillStr)
			}
		}
		mode.ScheduleRefresh()

	}

}

func (mode *WordMode) fill(name string) []string {
	switch name {
	case "title":
		ret := []string{"FACADE"}
		if mode.maxLength >= 11. {
			ret = []string{"F A C A D E"}
		}
		if mode.maxLength >= 15. {
			ret = append(ret, "by FEEDFACE.COM")
		}
		return ret

	case "index":
		maxLength := 8
		if mode.maxLength >= 1.0 {
			maxLength = int(mode.maxLength)
		}
		ret := []string{}
		for i := 0; i < mode.wordBuffer.slotCount; i++ {
			s := ""
			for j := 0; j < maxLength; j++ {
				s += fmt.Sprintf("%1x", i%0x10)
			}
			ret = append(ret, s)
		}
		return ret
	case "alpha":
		return strings.Split(`
alpha
bravo
charlie
delta
echo
foxtrot
golf
hotel
india
juliet
kilo
lima
mike
november
oscar
papa
quebec
romeo
sierra
tango
uniform
victor
whiskey
xray
yankee
zulu
`, "\n")[1:]
	default:
		log.Error("no such wordbuffer fill pattern: '%s'", name)
	}
	return []string{}
}

func (mode *WordMode) Desc() string {
	ret := "words["
	ret += fmt.Sprintf("#%d/%d ", mode.wordBuffer.WordCount(), mode.wordBuffer.SlotCount())
	ret += fmt.Sprintf("≤%d ", mode.maxLength)
	if mode.wordBuffer.Lifetime() > 0.0 {
		ret += fmt.Sprintf("l%.1f ", mode.wordBuffer.Lifetime())
	}
	if mode.wordBuffer.Watermark() > 0.0 {
		ret += fmt.Sprintf("m%0.1f ", mode.wordBuffer.Watermark())
	}
	if mode.wordBuffer.shuffle {
		ret += "⧢"
	}
	if mode.wordBuffer.aging {
		ret += "å"
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (mode *WordMode) Config() *WordConfig {
	ret := &WordConfig{
		SetSlots: true, Slots: uint64(mode.wordBuffer.SlotCount()),
		SetMaxLength: true, MaxLength: uint64(mode.maxLength),
		SetLifetime: true, Lifetime: float64(mode.wordBuffer.Lifetime()),
		SetWatermark: true, Watermark: float64(mode.wordBuffer.Watermark()),
		SetShuffle: true, Shuffle: bool(mode.wordBuffer.Shuffle()),
		SetAging: true, Aging: bool(mode.wordBuffer.Aging()),
	}
	return ret

}

func (mode *WordMode) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: mode.vert,
		SetFrag: true, Frag: mode.frag,
	}
	return ret
}
