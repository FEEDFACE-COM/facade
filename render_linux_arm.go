
// +build linux,arm

package main

import (
    "fmt"
    "strings"
    "time"
    "sync"
    "os"
    log "./log"
    facade "./facade"
    gfx "./gfx"
    "src.feedface.com/gfx/piglet"
    gl "src.feedface.com/gfx/piglet/gles2"
)


const DEBUG_RENDERER = true




type Renderer struct {
    screen gfx.Size

    mode facade.Mode
    debug bool

    grid *facade.Grid

    font *gfx.Font
    camera *gfx.Camera
    mask *gfx.Mask
    
    axis *gfx.Axis

    
    lineBuffer *facade.LineBuffer
    termBuffer *facade.TermBuffer

    fontService *gfx.FontService
    shaderService *gfx.ShaderService
    
    mutex *sync.Mutex
    directory string
    
    refreshChan chan bool
    
    prevClock gfx.Clock
    
}



func NewRenderer(directory string) *Renderer {
    ret := &Renderer{directory: directory}
    ret.mutex = &sync.Mutex{}
    ret.refreshChan = make( chan bool, 1 )
    return ret
}


func (renderer *Renderer) ScheduleRefresh() {

    select { case renderer.refreshChan <- true: ; default: ; }
	
}


func (renderer *Renderer) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
			case refresh := <- renderer.refreshChan:
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
    log.Info("init renderer[%s]",renderer.directory)
    if strings.HasPrefix(renderer.directory, "~/") {
        renderer.directory = os.Getenv("HOME") + renderer.directory[1:]
    }

    renderer.fontService = gfx.NewFontService(renderer.directory+"/font")
    renderer.shaderService = gfx.NewShaderService(renderer.directory+"/shader")

    err = piglet.CreateContext()
    if err != nil {
        log.PANIC("%s fail to initialize renderer: %s",renderer.Desc(),err)    
    }
    
    w,h := piglet.GetDisplaySize()
    renderer.screen = gfx.Size{float32(w),float32(h)}
    log.Info("%s got screen %s",renderer.Desc(),renderer.screen.Desc())
    

    piglet.MakeCurrent()

    err = gl.InitWithProcAddrFunc( piglet.GetProcAddress )
    if err != nil {
        log.PANIC("%s fail to init GLES: %s",renderer.Desc(),err)    
    }
    

    log.Debug("%s got renderer %s %s", renderer.Desc(),gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Debug("%s got version %s %s", renderer.Desc(),gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    renderer.mode = facade.Defaults.Mode
    renderer.debug = facade.Defaults.Debug
    


    renderer.axis = gfx.NewAxis()
    renderer.axis.Init(renderer.shaderService)



    renderer.font,err = renderer.fontService.GetFont( facade.DEFAULT_FONT )
    if err != nil {
        log.PANIC("%s fail to get default font %s: %s",renderer.Desc(),facade.DEFAULT_FONT,err)
    }


    renderer.camera = gfx.NewCamera( float32(facade.CameraDefaults.Zoom), facade.CameraDefaults.Isometric, renderer.screen)
    renderer.camera.Init()


    renderer.mask = gfx.NewMask(facade.MaskDefaults.Name,renderer.screen)
    renderer.mask.Init(renderer.shaderService)


    renderer.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width),uint(facade.GridDefaults.Height)) 
    renderer.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height),uint(facade.GridDefaults.Buffer),renderer.refreshChan) 

    renderer.grid = facade.NewGrid( renderer.lineBuffer, renderer.termBuffer )
    renderer.grid.Init(renderer.shaderService,renderer.font)


    gfx.ClockReset()
    return err
}


func (renderer *Renderer) Configure(config *facade.Config) error {
    if config == nil { return log.NewError("renderer config nil") }
    
    log.Info("%s config %s",renderer.Desc(),config.Desc())
    var err error

	if config.GetSetDebug() {
		renderer.debug = config.GetDebug()
	} else {
		renderer.debug = false	
	}


    if cfg := config.GetFont(); cfg!=nil {
    
        if cfg.GetSetName() {
            name := cfg.GetName()
            if name != renderer.font.GetName() {
        
                err = renderer.fontService.LoadFont( name )
                if err != nil {
                    log.Error("%s fail load font %s: %s",renderer.Desc(),name,err)
                }
                
                var fnt *gfx.Font
                fnt,err = renderer.fontService.GetFont( name )
                if err != nil {
                    log.Error("%s fail get font %s: %s",renderer.Desc(),name,err)
                } else {
                    if DEBUG_RENDERER { log.Debug("%s switch to font %s",renderer.Desc(),name) }
                    renderer.font = fnt
                    renderer.ScheduleRefresh()

                }                
            }        
        }
    }


    if cfg := config.GetCamera(); cfg!=nil {
        if cfg.GetSetZoom() {
            renderer.camera.ConfigureZoom( float32(cfg.GetZoom()) )
        }
        if cfg.GetSetIsometric() {
            renderer.camera.ConfigureIsometric( cfg.GetIsometric() )
        }
    }


    if cfg := config.GetMask(); cfg!=nil {
        if cfg.GetSetName() { 
            renderer.mask.ConfigureName( cfg.GetName() )
        }
    }


    if cfg := config.GetGrid(); cfg!=nil {
        renderer.grid.Configure(cfg,renderer.camera,renderer.font)
    }

	
    if config.GetSetMode() {
        mode := config.GetMode()
        if renderer.mode != mode {
        
            log.Info("%s switch to mode %s",renderer.Desc(),strings.ToLower(mode.String()))
            renderer.mode = mode
        
        }
        
    }
    
    
    return nil
}


func (renderer *Renderer) Render(confChan chan facade.Config) error {

    gl.Viewport(0, 0, int32(renderer.screen.W),int32(renderer.screen.H))
    gl.Disable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
    gl.ClearColor(0., 0., 0., 1.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
//    gl.Enable(gl.CULL_FACE)
//    gl.CullFace(gl.BACK)


    gfx.ClockTick()
    renderer.prevClock = *gfx.NewClock()
    log.Info("%s start render",renderer.Desc())
    for {
        
        verboseFrame := gfx.ClockVerboseFrame()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.ProcessConf(confChan)
        if renderer.checkRefresh() {
            switch renderer.mode {
                case facade.Mode_GRID:
                    renderer.grid.GenerateData(renderer.font)
            }
        }
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )

        gfx.RefreshPrograms()


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA,gl.ZERO,gl.ONE)
        switch renderer.mode {
            case facade.Mode_GRID:
                renderer.grid.Render(renderer.camera, renderer.font, renderer.debug, verboseFrame )
        }
      
        if renderer.debug {renderer.axis.Render(renderer.camera, renderer.debug) }


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA,gl.ZERO,gl.ONE)
        renderer.mask.Render(renderer.debug)
        
        
        
        if verboseFrame { 

            renderer.printDebug(); 
            renderer.prevClock = *gfx.NewClock()
        }

        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??


//        if e != gl.NO_ERROR && verboseFrame { 

        e := uint32(gl.NO_ERROR)
        if verboseFrame { 
            e = gl.GetError()
            if e != gl.NO_ERROR {
                log.Error("%s post render gl error: %s",renderer.Desc(),gl.ErrorString(e)) 
            }
        }
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        gfx.ClockTick()
    }
    return nil
}

func (renderer *Renderer) ProcessConf(confChan chan facade.Config) {
    
    select {
        case conf := <-confChan:
            renderer.Configure(&conf)
        
        default:
            //nop    
    }
}






func (renderer *Renderer) ProcessQueries(queryChan chan (chan string) ) error {

    if DEBUG_RENDERER { log.Debug("%s start process info queries",renderer.Desc(),) }

    for {
    
        chn := <- queryChan
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

    if DEBUG_RENDERER { log.Debug("%s start process text sequences",renderer.Desc(),) }

    for {
        item := <- textChan    
        text, seq := item.Text, item.Seq
        if text != nil && len(text) > 0 {
            renderer.lineBuffer.ProcessRunes( text )
            renderer.termBuffer.ProcessRunes( text )    
            renderer.ScheduleRefresh()
//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
//                log.Debug( "%s", renderer.grid.DumpBuffer() )
//            }
        }
        if seq != nil {
            renderer.lineBuffer.ProcessSequence( seq )
            renderer.termBuffer.ProcessSequence( seq )
            renderer.ScheduleRefresh()
//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
//                log.Debug( "%s", renderer.grid.DumpBuffer() )
//            }
        }
    }
    return nil
}







func (renderer *Renderer) InfoMode() string {
        mode := ""
        switch renderer.mode { 
            case facade.Mode_GRID:
                mode = renderer.grid.Desc()
        }
		dbg := ""
		if renderer.debug {
			dbg = " DEBUG"	
		}
        return fmt.Sprintf("%s\n  %s %s %s%s",mode,renderer.camera.Desc(),renderer.font.Desc(),renderer.mask.Desc(),dbg)
    
}


func (renderer *Renderer) InfoClock() string {
    return fmt.Sprintf("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(renderer.prevClock)  )
}


func (renderer *Renderer) printDebug() {

    if DEBUG_CLOCK||DEBUG_MODE||DEBUG_BUFFER {
        log.Debug("")
    }


    if DEBUG_CLOCK { log.Info( "%s", renderer.InfoClock() ) }
    
    if DEBUG_DIAG { log.Info("  %s", MemUsage() ) }
        
    if DEBUG_MODE { 
            log.Info("  %s", renderer.InfoMode() ) 
            log.Info("  %s", renderer.lineBuffer.Desc() )
            log.Info("  %s", renderer.termBuffer.Desc() )
    }

    if DEBUG_FONT {
        log.Info("  %s",renderer.fontService.Desc())
        if renderer.font != nil {
            log.Info("  %s",renderer.font.Desc())
        }
    }    
    

    if DEBUG_BUFFER &&  log.DebugLogging() { renderer.dumpBuffer() }
 
    if DEBUG_CLOCK||DEBUG_MODE||DEBUG_BUFFER {
        log.Debug("")
    }

    
}

func (renderer *Renderer) dumpBuffer() {
//    if ! DEBUG_BUFFER {
//        return
//    }
//    if renderer.mode  == facade.Mode_GRID {
//        os.Stdout.Write( []byte( renderer.grid.DumpBuffer() ) )        
//    }
//    os.Stdout.Write( []byte( "\n" ) )
//    os.Stdout.Sync()
}

func (renderer *Renderer) Info() string { 
    ret := ""
    
    ret += InfoVersion()
    ret += InfoAssets(nil,nil)
    ret += "\n\n"


    ret += renderer.InfoClock()
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



