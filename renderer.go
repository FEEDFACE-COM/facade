//go:build RENDERER
// +build RENDERER

package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"github.com/FEEDFACE-COM/piglet"
	"fmt"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"os"
	"strings"
	"sync"
	"time"
)

const RENDER_FRAME_RATE = 60.0

type Renderer struct {
	screen gfx.Size

	mode  facade.Mode
	debug bool

	terminal *facade.LineMode
	lines    *facade.LineMode
	words    *facade.WordMode
	chars    *facade.CharMode

	font   *gfx.Font
	camera *gfx.Camera
	mask   *gfx.Mask

	axis *gfx.Axis

	lineBuffer *facade.LineBuffer
	termBuffer *facade.TermBuffer
	wordBuffer *facade.WordBuffer
	charBuffer *facade.CharBuffer

	fontService    *gfx.FontService
	programService *gfx.ProgramService

	stateMutex *sync.Mutex
	directory  string

	refreshChan chan bool

	prevFrame gfx.ClockFrame

	tickChannel chan Tick
}

func NewRenderer(directory string, tickChannel chan Tick) *Renderer {
	ret := &Renderer{directory: directory, tickChannel: tickChannel}
	ret.stateMutex = &sync.Mutex{}
	ret.refreshChan = make(chan bool, 1)
	if ret.directory == "" {
		ret.fontService = gfx.NewFontService("", facade.FontAsset)
		ret.programService = gfx.NewProgramService("", facade.ShaderAsset)
	} else {
		if strings.HasPrefix(ret.directory, "~/") {
			ret.directory = os.Getenv("HOME") + ret.directory[1:]
		}
		ret.fontService = gfx.NewFontService(ret.directory+"/font", facade.FontAsset)
		ret.programService = gfx.NewProgramService(ret.directory+"/shader", facade.ShaderAsset)
	}
	return ret
}

func (renderer *Renderer) ScheduleRefresh() {

	select {
	case renderer.refreshChan <- true:
	default:
	}

}

func (renderer *Renderer) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-renderer.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}

func (renderer *Renderer) Init() error {
	var err error
	if DEBUG_RENDERER {
		log.Debug("%s init %s", renderer.Desc(), renderer.directory)
	}

	err = piglet.CreateContext()
	if err != nil {
		return log.NewError("fail to create context: %s", err)
	}

	w, h := piglet.GetDisplaySize()
	renderer.screen = gfx.Size{float32(w), float32(h)}
	log.Notice("%s screen %s", renderer.Desc(), renderer.screen.Desc())

	piglet.MakeCurrent()
	err = gl.InitWithProcAddrFunc(piglet.GetProcAddress)
	if err != nil {
		log.Error("%s fail to gl init: %s", renderer.Desc(), err)
		return log.NewError("fail to gl init: %s", err)
	}

	log.Info("%s renderer %s %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VENDOR))), gl.GoStr(gl.GetString((gl.RENDERER))))
	log.Info("%s version %s shader %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VERSION))), gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))))

	renderer.mode = facade.Defaults.Mode
	renderer.debug = facade.Defaults.Debug

	renderer.axis = gfx.NewAxis()

	renderer.font, err = renderer.fontService.GetFont(facade.DEFAULT_FONT)
	if err != nil {
		return log.NewError("fail to get default font %s: %s", facade.DEFAULT_FONT, err)
	}

	renderer.camera = gfx.NewCamera(float32(facade.CameraDefaults.Zoom), facade.CameraDefaults.Isometric, renderer.screen)
	renderer.mask = gfx.NewMask(facade.MaskDefaults.Name, renderer.screen)

	renderer.termBuffer = facade.NewTermBuffer(uint(facade.TermDefaults.Width), uint(facade.TermDefaults.Height))
	renderer.lineBuffer = facade.NewLineBuffer(uint(facade.TermDefaults.Height), uint(facade.LineDefaults.Buffer), renderer.refreshChan)
	renderer.wordBuffer = facade.NewWordBuffer(renderer.refreshChan)
	renderer.charBuffer = facade.NewCharBuffer(renderer.refreshChan)

	renderer.terminal = facade.NewLineMode(nil, renderer.termBuffer)
	renderer.lines = facade.NewLineMode(renderer.lineBuffer, nil)
	renderer.words = facade.NewWordMode(renderer.wordBuffer)
	renderer.chars = facade.NewCharMode(renderer.charBuffer)

	gfx.WorldClock().Reset()

	renderer.axis.Init(renderer.programService)
	renderer.camera.Init()
	renderer.mask.Init(renderer.programService)
	renderer.terminal.Init(renderer.programService, renderer.font)
	renderer.lines.Init(renderer.programService, renderer.font)
	renderer.words.Init(renderer.programService, renderer.font)
	renderer.chars.Init(renderer.programService, renderer.font)

	return err
}

func (renderer *Renderer) Finish() error {
	piglet.DestroyContext()
	log.Info("%s finished", renderer.Desc())
	return nil
}

func (renderer *Renderer) Configure(config *facade.Config) error {
	changed := false
	if config == nil {
		return log.NewError("renderer config nil")
	}

	log.Notice("%s configure %s", renderer.Desc(), config.Desc())
	var err error

	if config.GetSetDebug() {
		renderer.debug = config.GetDebug()
	} else {
		renderer.debug = false
	}

	if config.GetSetMode() {
		changed = true
		mode := config.GetMode()
		if renderer.mode != mode {
			log.Notice("%s switch to mode %s", renderer.Desc(), strings.ToLower(mode.String()))
			renderer.mode = mode
		}
	}

	var fill string = ""
	if config.GetSetFill() {
		fill = config.GetFill()
	}
	if cfg := config.GetFont(); cfg != nil {

		if cfg.GetSetName() {
			name := cfg.GetName()
			if name != renderer.font.GetName() {
				changed = true
				//				err = renderer.fontService.LoadFont(name) // REM, probably not needed here?
				//				if err != nil {
				//					log.Error("%s fail load font %s: %s", renderer.Desc(), name, err)
				//				}

				var fnt *gfx.Font
				fnt, err = renderer.fontService.GetFont(name)
				if err != nil {
					//					log.Error("%s fail get font %s: %s", renderer.Desc(), name, err)
				} else {
					log.Notice("%s switch to font %s", renderer.Desc(), name)
					renderer.font = fnt
					renderer.ScheduleRefresh()
					renderer.terminal.ScheduleRefresh()
					renderer.lines.ScheduleRefresh()
					renderer.words.ScheduleRefresh()
					renderer.chars.ScheduleRefresh()
				}
			}
		}
	}

	if cfg := config.GetCamera(); cfg != nil {
		changed = true
		if cfg.GetSetZoom() {
			renderer.camera.ConfigureZoom(float32(cfg.GetZoom()))
		}
		if cfg.GetSetIsometric() {
			renderer.camera.ConfigureIsometric(cfg.GetIsometric())
		}
	}

	if cfg := config.GetMask(); cfg != nil {
		changed = true
		if cfg.GetSetName() {
			renderer.mask.ConfigureName(cfg.GetName())
		}
	}

	if cfg := config.GetShader(); cfg != nil {
		changed = true
	}

	if config.GetLines() != nil || config.GetTerm() != nil || config.GetWords() != nil || config.GetChars() != nil {
		changed = true
	}

	switch renderer.mode {
	case facade.Mode_LINES:
		renderer.lines.Configure(config.GetLines(), nil, config.GetShader(), fill, renderer.camera, renderer.font)
	case facade.Mode_TERM:
		renderer.terminal.Configure(nil, config.GetTerm(), config.GetShader(), fill, renderer.camera, renderer.font)
	case facade.Mode_WORDS:
		renderer.words.Configure(config.GetWords(), config.GetShader(), fill, renderer.camera, renderer.font)
	case facade.Mode_CHARS:
		renderer.chars.Configure(config.GetChars(), config.GetShader(), fill, renderer.camera, renderer.font)
	}

	if changed && DEBUG_CHANGES {
		log.Debug("%s", renderer.InfoMode())
	}

	return nil
}

func (renderer *Renderer) tick() {

	for { //forever
		renderer.tickChannel <- TICK // wait until can send
		time.Sleep(time.Duration(int64(time.Second / RENDER_FRAME_RATE)))
	}

}

func (renderer *Renderer) tock() bool {
	var tick Tick

	// wait for message
	tick = <-renderer.tickChannel

	// clear all messages
	for {

		switch tick {
		case QUIT:
			return false
		case STOP:
			renderer.TogglePause()
		}

		select {

		case tick = <-renderer.tickChannel:
			continue // evaluate

		default:
			return true // indicate render one frame

		}
	}

	// indicate render one frame
	return true
}

func (renderer *Renderer) Render(confChan chan facade.Config, showDebugInfo bool) error {

	go renderer.tick()

	gl.Viewport(0, 0, int32(renderer.screen.W), int32(renderer.screen.H))
	gl.Disable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0., 0., 0., 1.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	//    gl.Enable(gl.CULL_FACE)
	//    gl.CullFace(gl.BACK)

	gfx.WorldClock().Tick()
	renderer.prevFrame = gfx.WorldClock().Frame()
	log.Notice("%s start render", renderer.Desc())

	renderFailed := false
	for piglet.Loop() {

		if DEBUG_DIAG {
			DiagStart()
		}

		gfx.WorldClock().Tick()

		verboseFrame := gfx.WorldClock().VerboseFrame()

		renderer.stateMutex.Lock()

		renderer.ProcessConf(confChan)
		if renderer.checkRefresh() {
			//            if DEBUG_RENDERER { log.Debug("%s refresh",renderer.Desc()) }
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.terminal.ScheduleRefresh()

			case facade.Mode_LINES:
				renderer.lines.ScheduleRefresh()

			case facade.Mode_WORDS:
				renderer.words.ScheduleRefresh()

			case facade.Mode_CHARS:
				renderer.chars.ScheduleRefresh()

			}

		}

		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		renderer.programService.CheckRefresh()

		gl.BlendEquationSeparate(gl.FUNC_ADD, gl.FUNC_ADD)
		gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ZERO, gl.ONE)
		switch renderer.mode {
		case facade.Mode_TERM:
			renderer.terminal.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		case facade.Mode_LINES:
			renderer.lines.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		case facade.Mode_WORDS:
			renderer.words.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		case facade.Mode_CHARS:
			renderer.chars.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		}

		if renderer.debug {
			renderer.axis.Render(renderer.camera, renderer.debug)
		}

		gl.BlendEquationSeparate(gl.FUNC_ADD, gl.FUNC_ADD)
		gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA, gl.ZERO, gl.ONE)
		renderer.mask.Render(renderer.debug)

		piglet.SwapBuffers()
		renderer.stateMutex.Unlock()

		e := gl.GetError()
		if e == gl.NO_ERROR {
			if renderFailed { //first success
				log.Notice("%s render success", renderer.Desc())
			}
			renderFailed = false
		} else {
			if renderFailed == false { // first failure
				log.Error("%s render error: %s", renderer.Desc(), piglet.ErrorString(e))
			} else if RENDERER_FIXUP_RASPI {
				//HACK: remove gles2 debug output 'glGetError 0x502'
				str := fmt.Sprintf("\033[1A\033[K")
				os.Stderr.Write([]byte(str))
			}
			renderFailed = true
		}

		if DEBUG_DIAG {
			DiagDone()
		}

		if showDebugInfo {
			if renderFailed {
				renderer.printDebugInfo("RENDER FAILURE: " + piglet.ErrorString(e))
			} else {
				renderer.printDebugInfo(gfx.WorldClock().Info(renderer.prevFrame))
			}
		}

		if renderFailed == false {
			renderer.prevFrame = gfx.WorldClock().Frame()
		}

		// wait for redraw timer
		if !renderer.tock() {
			break
		}

	}
	return nil
}

func (renderer *Renderer) TogglePause() {
	gfx.WorldClock().Toggle()
	if DEBUG_RENDERER {
		log.Debug("%s toggle pause", renderer.Desc())
	}
}

func (renderer *Renderer) ProcessConf(confChan chan facade.Config) {

	select {
	case conf := <-confChan:
		renderer.Configure(&conf)

	default:
		//nop
	}
}

func (renderer *Renderer) ProcessTextSeqs(textChan chan facade.TextSeq) error {

	if DEBUG_RENDERER {
		log.Debug("%s start process text sequences", renderer.Desc())
	}

	for {
		item := <-textChan
		text, seq := item.Text, item.Seq

		if gfx.WorldClock().Paused() {
			continue
		}

		if text != nil && len(text) > 0 {
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.termBuffer.ProcessRunes(text)
			case facade.Mode_LINES:
				renderer.lineBuffer.ProcessRunes(text)
			case facade.Mode_WORDS:
				renderer.wordBuffer.ProcessRunes(text)
			case facade.Mode_CHARS:
				renderer.charBuffer.ProcessRunes(text)
			}
			renderer.ScheduleRefresh()

		}
		if seq != nil {
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.termBuffer.ProcessSequence(seq)
			case facade.Mode_LINES:
				renderer.lineBuffer.ProcessSequence(seq)
			case facade.Mode_WORDS:
				renderer.wordBuffer.ProcessSequence(seq)
			case facade.Mode_CHARS:
				renderer.charBuffer.ProcessSequence(seq)
			}
			renderer.ScheduleRefresh()
		}
	}
	return nil
}

func (renderer *Renderer) InfoMode() string {
	mode := ""
	switch renderer.mode {
	case facade.Mode_TERM:
		mode = renderer.terminal.Desc() + " " + renderer.terminal.ShaderConfig().Desc()
	case facade.Mode_LINES:
		mode = renderer.lines.Desc() + " " + renderer.lines.ShaderConfig().Desc()
	case facade.Mode_WORDS:
		mode = renderer.words.Desc() + " " + renderer.words.ShaderConfig().Desc()
	case facade.Mode_CHARS:
		mode = renderer.chars.Desc() + " " + renderer.chars.ShaderConfig().Desc()
	}
	dbg := ""
	if renderer.debug {
		dbg = " DEBUG"
	}

	buffer := ""
	switch renderer.mode {
	case facade.Mode_TERM:
		buffer = renderer.termBuffer.Desc()
	case facade.Mode_LINES:
		buffer = renderer.lineBuffer.Desc()
	case facade.Mode_WORDS:
		buffer = renderer.wordBuffer.Desc()
	case facade.Mode_CHARS:
		buffer = renderer.charBuffer.Desc()
	}

	ret := ""
	ret += fmt.Sprintf("%s %s %s\n", renderer.font.Desc(), renderer.camera.Desc(), renderer.mask.Desc())
	ret += fmt.Sprintf("%s%s\n", mode, dbg)
	ret += fmt.Sprintf("%s", buffer)
	return ret
}

const DEBUGINFO_INFO = 4
const DEBUGINFO_BUFFER = 12
const DEBUGINFO_HEIGHT = DEBUGINFO_INFO + DEBUGINFO_BUFFER

func (renderer *Renderer) printDebugInfo(info string) {

	text := ""

	text += fmt.Sprintf("FACADE %s\n", info)
	text += renderer.InfoMode() + "\n"

	if DEBUG_BUFFER {
		tmp := ""
		switch renderer.mode {
		case facade.Mode_LINES:
			tmp = renderer.lineBuffer.Dump(80)
		case facade.Mode_TERM:
			tmp = renderer.termBuffer.Dump()
		case facade.Mode_WORDS:
			tmp = renderer.wordBuffer.Dump()
		case facade.Mode_CHARS:
			tmp = renderer.charBuffer.Dump()
		}
		buffer := strings.Split(tmp, "\n")
		count := len(buffer)
		if count <= DEBUGINFO_BUFFER {
			text += strings.Join(buffer, "\n")
			for i := count; i < DEBUGINFO_BUFFER; i++ {
				text += "\n"
			}
		} else {
			text += strings.Join(buffer[0:DEBUGINFO_BUFFER], "\n")
		}
	}

	text += "\n"
	//text += fmt.Sprintf("FACADE %s\n", info)

	text = strings.TrimRight(text, "\n")
	//text += fmt.Sprintf("\n## FACADE %s ##\n", info)

	lines := strings.Split(text, "\n")

	//log.Debug("have %d lines",len(lines))

	fmt.Fprintf(os.Stderr, "\0337")  // save cursor pos
	fmt.Fprintf(os.Stderr, "\033[H") // jump to origin

	for i := 0; i < DEBUGINFO_HEIGHT; i++ {
		fmt.Fprintf(os.Stderr, "\033[0K")
		if i < len(lines) {
			fmt.Fprintf(os.Stderr, "%s\n", lines[i]) // erase to eol and write single line
		} else {
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	fmt.Fprintf(os.Stderr, "\0338") // restore cursor pos

	//fmt.Fprintf(os.Stderr, "\033[%dA", DEBUGINFO_HEIGHT) // cursor up HEIGHT lines
	//
	//i := 0
	//for ; i < DEBUGINFO_HEIGHT && i < len(lines); i++ {
	//	fmt.Fprintf(os.Stderr, "\033[0K%s\n", lines[i]) // erase to eol and write single line
	//}
	//fmt.Fprintf(os.Stderr, "\033[%dB", DEBUGINFO_HEIGHT-i) // cursor down remaining lines

}

func (renderer *Renderer) Desc() string {
	ret := "renderer["
	ret += fmt.Sprintf("%dx%d ", int(renderer.screen.W), int(renderer.screen.H))
	ret += strings.ToLower(renderer.mode.String())
	if renderer.debug {
		ret += " DEBUG"
	}
	ret += "]"
	return ret
}

const RENDERER_AVAILABLE = true
