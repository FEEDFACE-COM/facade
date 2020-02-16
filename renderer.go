// +build linux,arm

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	facade "./facade"
	gfx "./gfx"
	log "./log"
	"github.com/FEEDFACE-COM/piglet"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
)

type Renderer struct {
	screen gfx.Size

	mode  facade.Mode
	debug bool

	terminal *facade.Grid
	lines    *facade.Grid
	
	tags     * facade.Set

	font   *gfx.Font
	camera *gfx.Camera
	mask   *gfx.Mask

	axis *gfx.Axis

	lineBuffer *facade.LineBuffer
	termBuffer *facade.TermBuffer
	setBuffer *facade.SetBuffer

	fontService    *gfx.FontService
	programService *gfx.ProgramService

	stateMutex *sync.Mutex
	directory  string

	refreshChan chan bool

    paused bool

	prevFrame gfx.ClockFrame

	tickChannel chan bool
}

func NewRenderer(directory string) *Renderer {
	ret := &Renderer{directory: directory}
	ret.stateMutex = &sync.Mutex{}
	ret.refreshChan = make(chan bool, 1)
	ret.tickChannel = make(chan bool, 1)
	ret.fontService = gfx.NewFontService(directory+"/font", facade.FontAsset)
	ret.programService = gfx.NewProgramService(directory+"/shader", facade.ShaderAsset)
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
	log.Info("init renderer[%s]", renderer.directory)
	if strings.HasPrefix(renderer.directory, "~/") {
		renderer.directory = os.Getenv("HOME") + renderer.directory[1:]
	}

	err = piglet.CreateContext()
	if err != nil {
		return log.NewError("fail to initialize renderer: %s", err)
	}

	w, h := piglet.GetDisplaySize()
	renderer.screen = gfx.Size{float32(w), float32(h)}
	log.Info("%s got screen %s", renderer.Desc(), renderer.screen.Desc())

	piglet.MakeCurrent()

	err = gl.InitWithProcAddrFunc(piglet.GetProcAddress)
	if err != nil {
		return log.NewError("fail to init GLES: %s", err)
	}

	log.Debug("%s got renderer %s %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VENDOR))), gl.GoStr(gl.GetString((gl.RENDERER))))
	log.Debug("%s got version %s %s", renderer.Desc(), gl.GoStr(gl.GetString((gl.VERSION))), gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))))

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
	renderer.setBuffer = facade.NewSetBuffer( renderer.refreshChan )

	renderer.terminal = facade.NewGrid(nil, renderer.termBuffer)
	renderer.terminal.Init(renderer.programService, renderer.font)

	renderer.lines = facade.NewGrid(renderer.lineBuffer, nil)
	renderer.lines.Init(renderer.programService, renderer.font)
	
	
	renderer.tags = facade.NewSet(renderer.setBuffer)
	renderer.tags.Init(renderer.programService, renderer.font)

	gfx.WorldClock().Reset()
	return err
}

func (renderer *Renderer) Configure(config *facade.Config) error {
	changed := false
	if config == nil {
		return log.NewError("renderer config nil")
	}

	log.Info("%s configure %s", renderer.Desc(), config.Desc())
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
				err = renderer.fontService.LoadFont(name)
				if err != nil {
					log.Error("%s fail load font %s: %s", renderer.Desc(), name, err)
				}

				var fnt *gfx.Font
				fnt, err = renderer.fontService.GetFont(name)
				if err != nil {
					log.Error("%s fail get font %s: %s", renderer.Desc(), name, err)
				} else {
					if DEBUG_RENDERER {
						log.Debug("%s switch to font %s", renderer.Desc(), name)
					}
					renderer.font = fnt
					renderer.ScheduleRefresh()
					renderer.terminal.ScheduleRefresh()
					renderer.lines.ScheduleRefresh()
					renderer.tags.ScheduleRefresh()

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
	
	if cfg := config.GetTags(); cfg != nil {
    	changed = true
    	renderer.tags.Configure(cfg, renderer.camera, renderer.font)
    }

	if config.GetSetMode() {
		changed = true
		mode := config.GetMode()
		if renderer.mode != mode {
			log.Info("%s switch to mode %s", renderer.Desc(), strings.ToLower(mode.String()))
			renderer.mode = mode
		}
	}

	if changed && DEBUG_EVENTS { 
		renderer.printDebug()
		renderer.prevFrame = gfx.WorldClock().Frame()
	}

	return nil
}

func (renderer *Renderer) tick() {

	for { //forever
		renderer.tickChannel <- true // wait until can send
		time.Sleep(time.Duration(int64(time.Second / FRAME_RATE)))
	}

}

func (renderer *Renderer) tock() {

	// wait for message
	<-renderer.tickChannel

	// clear all messages
	for {
		select {
		case <-renderer.tickChannel:
		default:
			return
		}
	}

	// return to render one frame
}

func (renderer *Renderer) Render(confChan chan facade.Config, pauseChan chan bool) error {

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
	log.Info("%s start render", renderer.Desc())

	renderFailed := false
	for {
		if DEBUG_DIAG {
			DiagStart()
		}

		gfx.WorldClock().Tick()

		verboseFrame := gfx.WorldClock().VerboseFrame()

		renderer.stateMutex.Lock()
		piglet.MakeCurrent()

		renderer.ProcessConf(confChan,pauseChan)
		if renderer.checkRefresh() {
			//            if DEBUG_RENDERER { log.Debug("%s refresh",renderer.Desc()) }
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.terminal.ScheduleRefresh()

			case facade.Mode_LINE:
				renderer.lines.ScheduleRefresh()

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
		case facade.Mode_LINE:
			renderer.lines.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
        case facade.Mode_TAGS:
            renderer.tags.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame)
		}

		if renderer.debug && renderer.paused {
			renderer.axis.Render(renderer.camera, renderer.debug)
		}

		gl.BlendEquationSeparate(gl.FUNC_ADD, gl.FUNC_ADD)
		gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA, gl.ZERO, gl.ONE)
		renderer.mask.Render(renderer.debug)

		if DEBUG_PERIODIC && verboseFrame {
			renderer.printDebug()
			renderer.prevFrame = gfx.WorldClock().Frame()
		}

		piglet.SwapBuffers()
		renderer.stateMutex.Unlock()

		e := gl.GetError()
		if e == gl.NO_ERROR {
			if renderFailed { //first success
				log.Notice("%s render success", renderer.Desc())
			}
			renderFailed = false
		} else {

			//HACK: remove gles2 debug output 'glGetError 0x502'
			os.Stderr.Write([]byte("\b\r"))

			if renderFailed == false { // first failure
				log.Error("%s render error: %s", renderer.Desc(), gl.ErrorString(e))
			}
			renderFailed = true
		}

		if DEBUG_DIAG {
			DiagDone()
		}

		renderer.tock()

	}
	return nil
}

func (renderer *Renderer) TogglePause() {
    gfx.WorldClock().Toggle()
    if DEBUG_RENDERER {
        log.Debug("%s toggle pause",renderer.Desc())
    }
}

func (renderer *Renderer) ProcessConf(confChan chan facade.Config, pauseChan chan bool) {

	select {
	case conf := <-confChan:
		renderer.Configure(&conf)
    
    case <-pauseChan:
        renderer.TogglePause()

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
		
		if renderer.paused {
    		continue
        }
		
		if text != nil && len(text) > 0 {
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.termBuffer.ProcessRunes(text)

			case facade.Mode_LINE:
				renderer.lineBuffer.ProcessRunes(text)
			
			case facade.Mode_TAGS:
                renderer.setBuffer.ProcessRunes(text)

			}

			renderer.ScheduleRefresh()

            if DEBUG_EVENTS {
    			renderer.printDebug()
                renderer.prevFrame = gfx.WorldClock().Frame()
			//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
			//                log.Debug( "%s", renderer.grid.DumpBuffer() )
			//            }
            }
	   		
		}
		if seq != nil {
			switch renderer.mode {
			case facade.Mode_TERM:
				renderer.termBuffer.ProcessSequence(seq)
			case facade.Mode_LINE:
				renderer.lineBuffer.ProcessSequence(seq)
            case facade.Mode_TAGS:
                renderer.setBuffer.ProcessSequence(seq)
			}
			renderer.ScheduleRefresh()
			if DEBUG_EVENTS {
    			renderer.printDebug()
	       		renderer.prevFrame = gfx.WorldClock().Frame()
			//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
			//                log.Debug( "%s", renderer.grid.DumpBuffer() )
			//            }
            }
		}
	}
	return nil
}

func (renderer *Renderer) InfoMode() string {
	mode := ""
	switch renderer.mode {
	case facade.Mode_TERM:
		mode = "term " + renderer.terminal.Desc() + " " + renderer.terminal.ShaderConfig().Desc()
	case facade.Mode_LINE:
		mode = "line " + renderer.lines.Desc() + " " + renderer.lines.ShaderConfig().Desc()
    case facade.Mode_TAGS:
        mode = "tags " + renderer.tags.Desc() + " " + renderer.tags.ShaderConfig().Desc()
	}
	dbg := ""
	if renderer.debug {
		dbg = " DEBUG"
	}
	return fmt.Sprintf("%s%s\n  %s %s %s", mode ,dbg, renderer.font.Desc(), renderer.camera.Desc(), renderer.mask.Desc())

}

//func (renderer *Renderer) InfoClock() string {
//    return fmt.Sprintf("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(renderer.prevFrame)  )
//}

func (renderer *Renderer) printDebug() {

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Debug("")
	}

	if DEBUG_MEMORY {
		log.Debug("memory usage %s", MemUsage())
	}

	if DEBUG_DIAG {
		log.Debug("%s    %s", gfx.WorldClock().Info(renderer.prevFrame), InfoDiag())
	}

	if DEBUG_CLOCK {
		log.Info("%s", gfx.WorldClock().Info(renderer.prevFrame))
	}

	if DEBUG_MODE {
		log.Debug("  %s", renderer.InfoMode())
		log.Debug("  %s", renderer.lineBuffer.Desc())
		log.Debug("  %s", renderer.termBuffer.Desc())
		log.Debug("  %s", renderer.setBuffer.Desc())
	}

	if DEBUG_FONT {
		log.Info("  %s", renderer.fontService.Desc())
		if renderer.font != nil {
			log.Debug("  %s", renderer.font.Desc())
		}
	}

	if DEBUG_BUFFER && log.DebugLogging() {
		renderer.dumpBuffer()
	}

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Debug("")
	}

}

func (renderer *Renderer) dumpBuffer() {
	if !DEBUG_BUFFER {
		return
	}
	if renderer.mode == facade.Mode_TERM {
		os.Stdout.Write([]byte(renderer.terminal.DumpBuffer()))
	} else if renderer.mode == facade.Mode_LINE {
		os.Stdout.Write([]byte(renderer.lines.DumpBuffer()))
	} else if renderer.mode == facade.Mode_TAGS {
    	os.Stdout.Write([]byte(renderer.setBuffer.Dump() ))
    }
	os.Stdout.Write([]byte("\n"))
	os.Stdout.Sync()
}

func (renderer *Renderer) Info() string {
	ret := ""

	ret += InfoVersion()
	ret += InfoAssets(renderer.programService.GetAvailableNames(), renderer.fontService.GetAvailableNames())
	ret += "\n\n"

	ret += gfx.WorldClock().Info(renderer.prevFrame)
	ret += "\n  " + renderer.InfoMode()
	ret += "\n  " + renderer.lineBuffer.Desc()
	ret += "\n  " + renderer.termBuffer.Desc()
	ret += "\n\n"

	return ret
}

func (renderer *Renderer) Desc() string {
	ret := "renderer["
	ret += strings.ToLower(renderer.mode.String())
	if renderer.debug {
		ret += " DEBUG"
	}
	ret += "]"
	return ret
}

const RENDERER_AVAILABLE = true
