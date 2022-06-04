//go:build (darwin && amd64) || (darwin && arm64)
// +build darwin,amd64 darwin,arm64

package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"fmt"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
	"os"
	"strings"
	"sync"
	"time"
)

type Renderer struct {
	screen gfx.Size

	mode  facade.Mode
	debug bool

	terminal *facade.Grid
	lines    *facade.Grid

	tags  *facade.Set
	words *facade.Set

	font   *gfx.Font
	camera *gfx.Camera
	mask   *gfx.Mask

	axis *gfx.Axis

	lineBuffer *facade.LineBuffer
	termBuffer *facade.TermBuffer
	wordBuffer *facade.WordBuffer
	tagBuffer  *facade.WordBuffer

	fontService    *gfx.FontService
	programService *gfx.ProgramService

	stateMutex *sync.Mutex
	directory  string

	window  *glfw.Window
	monitor *glfw.Monitor
	vidmode *glfw.VidMode
	winpos  *gfx.Frame

	refreshChan chan bool

	prevFrame gfx.ClockFrame

	tickChannel chan bool
}

func NewRenderer(directory string, tickChannel chan bool) *Renderer {
	ret := &Renderer{directory: directory, tickChannel: tickChannel}
	ret.stateMutex = &sync.Mutex{}
	ret.refreshChan = make(chan bool, 1)
	if strings.HasPrefix(ret.directory, "~/") {
		ret.directory = os.Getenv("HOME") + ret.directory[1:]
	}
	ret.fontService = gfx.NewFontService(ret.directory+"/font", facade.FontAsset)
	ret.programService = gfx.NewProgramService(ret.directory+"/shader", facade.ShaderAsset)
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

	/*
		go get -u -tags=gles3 github.com/go-gl/glfw/v3.3/glfw
	*/

	const WINDOW_WIDTH = 864
	const WINDOW_HEIGHT = 540

	var err error
	log.Notice("%s init %s", renderer.Desc(), renderer.directory)

	err = glfw.Init()
	if err != nil {
		return log.NewError("fail to initialize renderer: %s", err)
	}

	renderer.monitor = glfw.GetPrimaryMonitor()
	renderer.vidmode = renderer.monitor.GetVideoMode()
	log.Debug("%s mode %dx%d @%d fps", renderer.Desc(), renderer.vidmode.Width, renderer.vidmode.Height, renderer.vidmode.RefreshRate)
//	{
//		for _, mode := range renderer.monitor.GetVideoModes() {
//	   	w, h, fps := mode.Width, mode.Height, mode.RefreshRate
//	   	log.Debug("%s mode %dx%d @%d fps", renderer.Desc(), w, h, fps)
//	   }
//	}
	renderer.window, err = glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "FACADE by FEEDFACE.COM", nil, nil)
	if err != nil {
		glfw.Terminate()
		return log.NewError("fail to glfw create window: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	renderer.window.SetAspectRatio(WINDOW_WIDTH, WINDOW_HEIGHT)
	renderer.window.SetSizeLimits(WINDOW_WIDTH/2., WINDOW_WIDTH/2., gl.DONT_CARE, gl.DONT_CARE)
	renderer.window.SetSizeCallback(func(win *glfw.Window, w int, h int) { renderer.SizeFun(w, h) })
	renderer.window.SetFramebufferSizeCallback(func(win *glfw.Window, w int, h int) { renderer.FramebufferSizeFun(w, h) })
	renderer.window.SetKeyCallback(func(win *glfw.Window, k glfw.Key, c int, a glfw.Action, m glfw.ModifierKey) { renderer.KeyFun(k, a, m) })
	//renderer.window.SetRefreshCallback( func(win *glfw.Window) { renderer.RefreshFun() } )

	w, h := renderer.window.GetFramebufferSize()
	renderer.screen = gfx.Size{W: float32(w), H: float32(h)}
	log.Notice("%s screen %s", renderer.Desc(), renderer.screen.Desc())

	renderer.window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		return log.NewError("fail to init GLES: %s", err)
	}

	log.Info("%s renderer %s %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VENDOR))), gl.GoStr(gl.GetString((gl.RENDERER))))
	log.Info("%s version %s shader %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VERSION))), gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))))

	renderer.mode = facade.Defaults.Mode
	renderer.debug = facade.Defaults.Debug

	renderer.axis = gfx.NewAxis()
	renderer.axis.Init(renderer.programService)

	renderer.font, err = renderer.fontService.GetFont(facade.DEFAULT_FONT)
	if err != nil {
		return log.NewError("fail to get default font %s: %s", facade.DEFAULT_FONT, err)
	}

	renderer.camera = gfx.NewCamera(float32(facade.CameraDefaults.Zoom), facade.CameraDefaults.Isometric, renderer.screen)
	renderer.camera.Init()

	renderer.mask = gfx.NewMask(facade.MaskDefaults.Name, renderer.screen)
	renderer.mask.Init(renderer.programService)

	renderer.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width), uint(facade.GridDefaults.Height))
	renderer.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height), uint(facade.LineDefaults.Buffer), renderer.refreshChan)
	renderer.wordBuffer = facade.NewWordBuffer(renderer.refreshChan)
	renderer.tagBuffer = facade.NewTagBuffer(renderer.refreshChan)

	renderer.terminal = facade.NewGrid(nil, renderer.termBuffer)
	renderer.terminal.Init(renderer.programService, renderer.font)

	renderer.lines = facade.NewGrid(renderer.lineBuffer, nil)
	renderer.lines.Init(renderer.programService, renderer.font)

	renderer.words = facade.NewSet(renderer.wordBuffer)
	renderer.words.Init(renderer.programService, renderer.font)

	renderer.tags = facade.NewSet(renderer.tagBuffer)
	renderer.tags.Init(renderer.programService, renderer.font)

	gfx.WorldClock().Reset()

	return err
}

func (renderer *Renderer) RefreshFun() {
	log.Debug("%s refresh", renderer.Desc())
}

func (renderer *Renderer) FramebufferSizeFun(width int, height int) {
	renderer.screen = gfx.Size{W: float32(width), H: float32(height)}
	log.Notice("%s framebuffer size %s", renderer.Desc(), renderer.screen.Desc())
}

func (renderer *Renderer) SizeFun(width int, height int) {
	log.Notice("%s size %s", renderer.Desc(), renderer.screen.Desc())
}

func (renderer *Renderer) KeyFun(key glfw.Key, action glfw.Action, mod glfw.ModifierKey) {

	if mod == 0x2 && action == 0x1 && key == 0x43 {
		log.Notice("%s key ctrl-c", renderer.Desc())
		renderer.window.SetShouldClose(true)
		return
	}

	if key == 0x20 && action == 0x1 {
		log.Notice("%s key space", renderer.Desc())
		renderer.ToggleFullScreen()
		return
	}

	if key == 0x100 && action == 0x1 {
		log.Notice("%s key escape", renderer.Desc())
		renderer.window.SetShouldClose(true)
		return
	}

	log.Debug("%s key 0x%02x action %d mod %x", renderer.Desc(), key, action, mod)

}

func (renderer *Renderer) Finish() error {
	glfw.Terminate()
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
					renderer.tags.ScheduleRefresh()
					renderer.words.ScheduleRefresh()

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

	if cfg := config.GetLines(); cfg != nil {
		changed = true
		renderer.lines.Configure(cfg, nil, renderer.camera, renderer.font)
	}

	if cfg := config.GetTerminal(); cfg != nil {
		changed = true
		renderer.terminal.Configure(nil, cfg, renderer.camera, renderer.font)
	}

	if cfg := config.GetWords(); cfg != nil {
		changed = true
		renderer.words.Configure(cfg, nil, renderer.camera, renderer.font)
	}

	if cfg := config.GetTags(); cfg != nil {
		changed = true
		renderer.tags.Configure(nil, cfg, renderer.camera, renderer.font)
	}

	if config.GetSetMode() {
		changed = true
		mode := config.GetMode()
		if renderer.mode != mode {
			log.Notice("%s switch to mode %s", renderer.Desc(), strings.ToLower(mode.String()))
			renderer.mode = mode
		}
	}

	if changed && DEBUG_CHANGES {
		renderer.printDebug()
		//		renderer.prevFrame = gfx.WorldClock().Frame()
	}

	return nil
}

func (renderer *Renderer) tick() {

	for { //forever
		renderer.tickChannel <- true // wait until can send
		time.Sleep(time.Duration(int64(time.Second / RENDER_FRAME_RATE)))
	}

}

func (renderer *Renderer) tock() bool {
	var tick bool

	// wait for message
	tick = <-renderer.tickChannel
	if !tick {
		return false // indicate stop render
	}

	// clear all messages
	for {
		select {
		case tick = <-renderer.tickChannel:
			if !tick {
				return false // indicate stop render
			}

		default:
			return true // indicate render one frame
		}
	}

	// indicate render one frame
	return true
}

func (renderer *Renderer) Render(confChan chan facade.Config) error {

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
	for !renderer.window.ShouldClose() {
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

			case facade.Mode_TAGS:
				renderer.tags.ScheduleRefresh()

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
		case facade.Mode_TAGS:
			renderer.tags.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		}

		if renderer.debug {
			renderer.axis.Render(renderer.camera, renderer.debug)
		}

		gl.BlendEquationSeparate(gl.FUNC_ADD, gl.FUNC_ADD)
		gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA, gl.ZERO, gl.ONE)
		renderer.mask.Render(renderer.debug)

		if verboseFrame {
			if DEBUG_PERIODIC {
				renderer.printDebug()
			}
			renderer.prevFrame = gfx.WorldClock().Frame()
		}

		renderer.window.SwapBuffers()
		glfw.PollEvents()
		renderer.stateMutex.Unlock()

		e := gl.GetError()
		if e == gl.NO_ERROR {
			if renderFailed { //first success
				log.Notice("%s render success", renderer.Desc())
			}
			renderFailed = false
		} else {
			if renderFailed == false { // first failure
				log.Error("%s render error: %s", renderer.Desc(), glfw.ErrorCode(e).String())
			}
			renderFailed = true
		}

		if DEBUG_DIAG {
			DiagDone()
		}

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

func (renderer *Renderer) ToggleFullScreen() {
	if renderer.winpos == nil {
    	const SCREEN_WIDTH, SCREEN_HEIGHT = 3840,2160
		x, y := renderer.window.GetPos()
		w, h := renderer.window.GetSize()
		fps := renderer.vidmode.RefreshRate
		frame := gfx.Frame{P: gfx.Point{X: float32(x), Y: float32(y)}, S: gfx.Size{W: float32(w), H: float32(h)}}
		renderer.winpos = &frame
		log.Info("%s fullscreen %dx%d @%d fps", renderer.Desc(),SCREEN_WIDTH,SCREEN_HEIGHT,fps)
		renderer.window.SetMonitor(renderer.monitor, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, fps)
	} else {
		x, y := int(renderer.winpos.P.X), int(renderer.winpos.P.Y)
		w, h := int(renderer.winpos.S.W), int(renderer.winpos.S.H)
		fps := renderer.vidmode.RefreshRate
		renderer.winpos = nil
		log.Info("%s window %dx%d @%d fps", renderer.Desc(),w,h,fps)
		renderer.window.SetMonitor(nil, x, y, w, h, fps)
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

func (renderer *Renderer) ProcessQueries(queryChan chan (chan string)) error {

	if DEBUG_RENDERER {
		log.Debug("%s start process info queries", renderer.Desc())
	}

	for {

		chn := <-queryChan
		info := renderer.Info()

		select {
		case chn <- info:
			continue

		case <-time.After(1000. * time.Millisecond):
			continue
		}

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

			case facade.Mode_TAGS:
				renderer.tagBuffer.ProcessRunes(text)

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
			case facade.Mode_TAGS:
				renderer.tagBuffer.ProcessSequence(seq)
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
		mode = "term " + renderer.terminal.Desc() + " " + renderer.terminal.ShaderConfig().Desc()
	case facade.Mode_LINES:
		mode = "lines " + renderer.lines.Desc() + " " + renderer.lines.ShaderConfig().Desc()
	case facade.Mode_WORDS:
		mode = "words " + renderer.words.Desc() + " " + renderer.words.ShaderConfig().Desc()
	case facade.Mode_TAGS:
		mode = "tags " + renderer.tags.Desc() + " " + renderer.tags.ShaderConfig().Desc()
	}
	dbg := ""
	if renderer.debug {
		dbg = " DEBUG"
	}
	return fmt.Sprintf("%s%s\n  %s %s %s", mode, dbg, renderer.font.Desc(), renderer.camera.Desc(), renderer.mask.Desc())

}

//func (renderer *Renderer) InfoClock() string {
//    return fmt.Sprintf("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(renderer.prevFrame)  )
//}

func (renderer *Renderer) printDebug() {

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Info("")
	}

	if DEBUG_MEMORY {
		log.Info("memory usage %s", MemUsage())
	}

	if DEBUG_DIAG {
		log.Info("%s    %s", gfx.WorldClock().Info(renderer.prevFrame), InfoDiag())
	}

	if DEBUG_CLOCK {
		log.Info("%s", gfx.WorldClock().Info(renderer.prevFrame))
	}

	if DEBUG_MODE {
		log.Info("  %s", renderer.InfoMode())
		switch renderer.mode {
		case facade.Mode_LINES:
			log.Info("  %s", renderer.lineBuffer.Desc())
		case facade.Mode_TERM:
			log.Info("  %s", renderer.termBuffer.Desc())
		case facade.Mode_WORDS:
			log.Info("  %s", renderer.wordBuffer.Desc())
		case facade.Mode_TAGS:
			log.Info("  %s", renderer.tagBuffer.Desc())
		}
	}

	if DEBUG_FONT {
		log.Info("  %s", renderer.fontService.Desc())
		if renderer.font != nil {
			log.Info("  %s", renderer.font.Desc())
		}
	}

	if DEBUG_BUFFER && log.DebugLogging() {
		renderer.dumpBuffer()
	}

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Info("")
	}

}

func (renderer *Renderer) dumpBuffer() {
	if !DEBUG_BUFFER {
		return
	}
	if renderer.mode == facade.Mode_TERM {
		os.Stdout.Write([]byte(renderer.terminal.DumpBuffer()))
	} else if renderer.mode == facade.Mode_LINES {
		os.Stdout.Write([]byte(renderer.lines.DumpBuffer()))
	} else if renderer.mode == facade.Mode_WORDS {
		os.Stdout.Write([]byte(renderer.wordBuffer.Dump()))
	} else if renderer.mode == facade.Mode_TAGS {
		os.Stdout.Write([]byte(renderer.tagBuffer.Dump()))
	}
	os.Stdout.Write([]byte("\n"))
	os.Stdout.Sync()
}

func (renderer *Renderer) Info() string {
	ret := ""

	ret += InfoAuthor()
	ret += InfoVersion()
	ret += "\n\n"

	ret += gfx.WorldClock().Info(renderer.prevFrame)
	ret += "\n  " + renderer.InfoMode()
	ret += "\n  " + renderer.lineBuffer.Desc()
	ret += "\n  " + renderer.termBuffer.Desc()
	ret += "\n  " + renderer.wordBuffer.Desc()
	ret += "\n  " + renderer.tagBuffer.Desc()
	ret += "\n\n"

	return ret
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
