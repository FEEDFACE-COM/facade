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

const DEBUG_SET = true

const HARD_MAX_LENGTH = 80.0

type Set struct {
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
	CHAROFFSET gfx.AttribName = "charOffset"
	CHARINDEX  gfx.AttribName = "charIndex"
)

func (set *Set) ScheduleRefresh() {

	select {
	case set.refreshChan <- true:
	default:
	}

}

func (set *Set) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-set.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}

func NewSet(buffer *WordBuffer) *Set {
	ret := &Set{
		wordBuffer: buffer,
		maxWord:    Word{},
	}

	ret.vert = ShaderDefaults.GetVert()
	ret.frag = ShaderDefaults.GetFrag()

	ret.refreshChan = make(chan bool, 1)
	return ret
}

func (set *Set) generateData(font *gfx.Font) {

	//setup vertex + bind order arrays
	set.data = []float32{}

	words := set.wordBuffer.Words()
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
		set.data = append(set.data, set.vertices(word.text, float32(word.index), float32(word.length), float32(word.width), font)...)
	}

	if set.maxLength == 0 { // actually counted max values
		set.maxWord = Word{length: maxLength, width: maxWidth}
	} else { // config max length and font ratio
		set.maxWord = Word{length: set.maxLength, width: float32(set.maxLength) * font.Ratio()}
	}

	set.object.BufferData(len(set.data)*4, set.data)
	if DEBUG_SET {
		log.Debug("%s generate words:%d chars:%d floats:%d, max length:%d width:%.1f", set.Desc(), len(words), charCount, len(set.data), maxLength, maxWidth)
	}

}

func (set *Set) vertices(
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

	if DEBUG_SET {
		//log.Debug("%s data generate '%s'", set.Desc(), text)
	}

	return ret
}

func (set *Set) autoScale(camera *gfx.Camera) float32 {
	//scaleHeight := float32(1.) / math32.Sqrt( float32(set.wordBuffer.SlotCount() ) )
	//scaleHeight := float32(1.) / float32(set.wordBuffer.SlotCount())
	//scaleHeight := float32(1.) / float32(8.)
	scaleHeight := float32(1.)
	return scaleHeight //* 2.
}

func (set *Set) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {
	set.wordBuffer.mutex.Lock()

	if set.checkRefresh() {
		if DEBUG_SET {
			log.Debug("%s refresh", set.Desc())
		}
		set.generateData(font)
		set.renderMap(font)
	}

	words := set.wordBuffer.Words()
	gl.ActiveTexture(gl.TEXTURE0)

	set.program.UseProgram(debug)
	set.object.BindBuffer()

	set.program.Uniform1f(WORDCOUNT, float32(set.wordBuffer.SlotCount()))

	set.program.Uniform1f(WORDMAXLENGTH, float32(set.maxWord.length))
	set.program.Uniform1f(WORDMAXWIDTH, set.maxWord.width)
	set.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
	set.program.Uniform1f(gfx.FONTRATIO, font.Ratio())
	set.program.Uniform1f(gfx.CLOCKNOW, float32(gfx.Now()))

	camera.Uniform(set.program)
	scale := float32(1.0)
	scale = set.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	set.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])
	set.program.VertexAttribPointer(gfx.VERTEX, 3, (3+2+1+1+1+1+1)*4, (0)*4)
	set.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2+1+1+1+1+1)*4, (0+3)*4)
	set.program.VertexAttribPointer(CHARINDEX, 1, (3+2+1+1+1+1+1)*4, (0+3+2)*4)
	set.program.VertexAttribPointer(CHAROFFSET, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1)*4)
	set.program.VertexAttribPointer(WORDINDEX, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1)*4)
	set.program.VertexAttribPointer(WORDWIDTH, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1+1)*4)
	set.program.VertexAttribPointer(WORDLENGTH, 1, (3+2+1+1+1+1+1)*4, (0+3+2+1+1+1+1)*4)

	count := int32(1)
	offset := int32(0)

	for _, word := range words {

		text := word.text

		if len(text) <= 0 {
			continue
		}

		set.texture.Uniform(set.program)
		count = int32(utf8.RuneCountInString(word.text))

		fader := float32(1.0)
		if word.fader != nil {
			fader = word.fader.Value()
		}
		set.program.Uniform1f(WORDFADER, fader)

		age := float32(0.0)
		if word.timer != nil {
			age = word.timer.Value()
		}
		set.program.Uniform1f(WORDAGE, age)

		if !debug || debug {
			set.program.SetDebug(false)
			set.texture.BindTexture()
			gl.DrawArrays(gl.TRIANGLES, int32(offset*(2*3)), int32(count*2*3))
			set.program.SetDebug(debug)
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
	if DEBUG_SET && verbose {
		//log.Debug("%s render %d words", set.Desc(), len(words) )
	}
	set.wordBuffer.mutex.Unlock()
}

func (set *Set) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", set.Desc())

	set.object = gfx.NewObject("set")
	set.object.Init()

	set.texture = gfx.NewTexture("set")
	set.texture.Init()

	set.program = programService.GetProgram("set", "set/")
	set.program.Link(set.vert, set.frag)

	set.renderMap(font)

	set.ScheduleRefresh()

}

func (set *Set) renderMap(font *gfx.Font) error {

	if DEBUG_SET {
		log.Debug("%s render texture map %s", set.Desc(), font.Desc())
	}

	rgba, err := font.RenderMap(false)
	if err != nil {
		log.Error("%s fail render font map: %s", set.Desc(), err)
		return log.NewError("fail render font map: %s", err)
	}
	err = set.texture.LoadRGBA(rgba)
	if err != nil {
		log.Error("%s fail load font map: %s", set.Desc(), err)
		return log.NewError("fail to load font map: %s", err)
	}
	set.texture.TexImage()

	return nil
}

func (set *Set) Configure(words *WordConfig, camera *gfx.Camera, font *gfx.Font) {
	var shader *ShaderConfig = nil
	var config *SetConfig = nil

	log.Debug("%s configure %s", set.Desc(), words.Desc())
	shader = words.GetShader()
	config = words.GetSet()

	{
		changed := false
		vert, frag := set.vert, set.frag

		if shader != nil {

			if shader.GetSetVert() {
				changed = true
				set.vert = shader.GetVert()
			}

			if shader.GetSetFrag() {
				changed = true
				set.frag = shader.GetFrag()
			}
		}

		if changed {
			err := set.program.Link(set.vert, set.frag)
			//			err := set.LoadShaders()
			if err != nil {
				set.vert = vert
				set.frag = frag
			}
		}
	}

	if config.GetSetSlots() {
		set.wordBuffer.Resize(int(config.GetSlots()))
	}

	if config.GetSetLifetime() {
		set.wordBuffer.SetLifetime(float32(config.GetLifetime()))
	}

	if config.GetSetWatermark() {
		set.wordBuffer.SetWatermark(float32(config.GetWatermark()))
	}

	if config.GetSetShuffle() {
		set.wordBuffer.SetShuffle(config.GetShuffle())
	}

	if config.GetSetMaxLength() {
		set.maxLength = uint(config.GetMaxLength())
	}

	if config.GetSetAging() {
		set.wordBuffer.SetAging(config.GetAging())
	}

	if config.GetSetFill() {
		fillStr := set.fill(config.GetFill())
		if fillStr != nil {
			set.wordBuffer.Fill(fillStr)
		}
	}

	set.ScheduleRefresh()

}

func (set *Set) fill(name string) []string {
	switch name {
	case "title":
		ret := []string{"FACADE"}
		if set.maxLength >= 11. {
			ret = []string{"F A C A D E"}
		}
		if set.maxLength >= 15. {
			ret = append(ret, "by FEEDFACE.COM")
		}
		return ret

	case "index":
		maxLength := 8
		if set.maxLength >= 1.0 {
			maxLength = int(set.maxLength)
		}
		ret := []string{}
		for i := 0; i < set.wordBuffer.slotCount; i++ {
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

func (set *Set) Desc() string {
	ret := "set["
	ret += fmt.Sprintf("#%d/%d", set.wordBuffer.WordCount(), set.wordBuffer.SlotCount())
	ret += fmt.Sprintf("≤%d", set.maxLength)
	if set.wordBuffer.Lifetime() > 0.0 {
		ret += fmt.Sprintf(" l%.1f", set.wordBuffer.Lifetime())
	}
	if set.wordBuffer.Watermark() > 0.0 {
		ret += fmt.Sprintf(" m%0.1f", set.wordBuffer.Watermark())
	}
	if set.wordBuffer.shuffle || set.wordBuffer.aging {
		ret += " "
		if set.wordBuffer.shuffle {
			ret += "⧢"
		}
		if set.wordBuffer.aging {
			ret += "å"
		}
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (set *Set) Config() *SetConfig {
	ret := &SetConfig{
		SetSlots: true, Slots: uint64(set.wordBuffer.SlotCount()),
		SetMaxLength: true, MaxLength: uint64(set.maxLength),
		SetLifetime: true, Lifetime: float64(set.wordBuffer.Lifetime()),
		SetWatermark: true, Watermark: float64(set.wordBuffer.Watermark()),
		SetShuffle: true, Shuffle: bool(set.wordBuffer.Shuffle()),
		SetAging: true, Aging: bool(set.wordBuffer.Aging()),
	}
	return ret

}

func (set *Set) ShaderConfig() *ShaderConfig {
	ret := &ShaderConfig{
		SetVert: true, Vert: set.vert,
		SetFrag: true, Frag: set.frag,
	}
	return ret
}
