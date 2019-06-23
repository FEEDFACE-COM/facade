
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





const RENDERER_AVAILABLE = true


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
    
    mutex *sync.Mutex
    directory string
    
    refreshChan chan bool
    
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


func (renderer *Renderer) Init(config *facade.Config) error {
    var err error
    log.Info("init renderer[%s] %s",renderer.directory,config.Desc())
    if strings.HasPrefix(renderer.directory, "~/") {
        renderer.directory = os.Getenv("HOME") + renderer.directory[1:]
    }

    gfx.SetShaderDirectory(renderer.directory+"/shader")
    gfx.SetFontDirectory(renderer.directory+"/font")
    
    err = piglet.CreateContext()
    if err != nil {
        log.PANIC("fail to initialize renderer: %s",err)    
    }
    
    w,h := piglet.GetDisplaySize()
    renderer.screen = gfx.Size{float32(w),float32(h)}
    log.Info("got screen %s",renderer.screen.Desc())
    

    piglet.MakeCurrent()

    err = gl.InitWithProcAddrFunc( piglet.GetProcAddress )
    if err != nil {
        log.PANIC("fail to init GLES: %s",err)    
    }
    

    log.Info("got renderer %s %s", gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Info("got version %s %s", gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    //setup things 
    renderer.mode = facade.Defaults.Mode
    renderer.debug = facade.Defaults.Debug
    if config.GetSetMode()  { renderer.mode = config.GetMode() }
    if config.GetSetDebug() { renderer.debug = config.GetDebug() }
    



//    renderer.config.Clean()
    
    renderer.axis = &gfx.Axis{}


    {
    	var name = facade.DEFAULT_FONT
        if cfg := config.GetFont(); cfg!=nil {
            if cfg.GetSetName() {
                name = cfg.GetName()
            }
        }
    	renderer.font,err = gfx.GetFont( name )
        if err != nil {
            log.PANIC("no default font %s: %s",name,err)    
        }
    	renderer.font.Init()
    }
    

    {
        var zoom = facade.CameraDefaults.Zoom
        var iso = facade.CameraDefaults.Isometric
        if cfg:=config.GetCamera(); cfg!=nil {
            if cfg.GetSetZoom() {
                zoom = cfg.GetZoom()
            }
            if cfg.GetSetIsometric() {
                iso = cfg.GetIsometric()
            }
        }

        renderer.camera = gfx.NewCamera( float32(zoom), iso, renderer.screen)
        renderer.camera.Init()
	}


    {
        var name = facade.MaskDefaults.Name
        if cfg:=config.GetMask(); cfg!=nil {
            if cfg.GetSetName() {
                name = cfg.GetName()
            }
        }
        renderer.mask = gfx.NewMask(name,renderer.screen)
        renderer.mask.Init()
    }



    renderer.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width),uint(facade.GridDefaults.Height)) 
    renderer.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height),uint(facade.GridDefaults.Buffer),renderer.refreshChan) 

    renderer.grid = facade.NewGrid( renderer.lineBuffer, renderer.termBuffer )
    renderer.grid.Init(renderer.camera,renderer.font)
    renderer.grid.Configure(config.GetGrid(),renderer.camera,renderer.font)

    renderer.axis.Init()


    gfx.ClockReset()
    return err
}


func (renderer *Renderer) Configure(config *facade.Config) error {
    
    if config == nil { log.Error("renderer config nil") ;return nil }
    
    log.Info("renderer config %s",config.Desc())
    
    
    
    
//    if tmp,ok := config.Font(); ok {
//		newFont, err := gfx.GetFont(&tmp)
//		if err != nil {
//			log.Error("fail to get font %s",tmp.Desc())
//		} else {
//			newFont.Init()
//			renderer.font = newFont
//			if renderer.grid != nil { 
//				cfg := make(facade.GridConfig)
//				cfg.SetHeight( renderer.grid.Height() )
//				renderer.grid.Configure(&cfg,renderer.camera,renderer.font)
//			}
//		}
//	}
//    
//    if tmp,ok := config.Camera(); ok {
//		renderer.camera.Configure(&tmp)    
//	}
//    

    if cfg := config.GetFont(); cfg!=nil {
        if cfg.GetSetName() {
            name := cfg.GetName()
            newFont, err := gfx.GetFont(name)
    		if err != nil {
    	   		log.Error("fail to switch font %s",name)
    	    } else {
        	   newFont.Init()
        	   renderer.font = newFont
        	   log.Info("switch font %s",renderer.font.Desc())
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
	
	if config.GetSetDebug() {
		renderer.debug = config.GetDebug()
	} else {
		renderer.debug = false	
	}
    
    
    if config.GetSetMode() {
        mode := config.GetMode()
        if renderer.mode != mode {
        
            log.Info("switch mode[] to mode[%s]",renderer.mode.String(),mode.String())
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
    var prev gfx.Clock = *gfx.NewClock()
    log.Info("render %s",renderer.Desc())
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

            renderer.printDebug(prev); 
            prev = *gfx.NewClock() 
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
                log.Error("post render gl error: %s",gl.ErrorString(e)) 
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
            if DEBUG_MEMORY { log.Info("mem now %s",MemUsage())}
        
        
        default:
            //nop    
    }
}









func (renderer *Renderer) ProcessTextSeqs(textChan chan facade.TextSeq) error {

    for {
        item := <- textChan    
        text, seq := item.Text, item.Seq
        if text != nil && len(text) > 0 {
            renderer.lineBuffer.ProcessRunes( text )
            renderer.termBuffer.ProcessRunes( text )    
            renderer.ScheduleRefresh()
//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
//                log.Debug( renderer.grid.DumpBuffer() )
//            }
        }
        if seq != nil {
            renderer.lineBuffer.ProcessSequence( seq )
            renderer.termBuffer.ProcessSequence( seq )
            renderer.ScheduleRefresh()
//            if DEBUG_BUFFER && renderer.mode == facade.Mode_GRID {
//                log.Debug( renderer.grid.DumpBuffer() )
//            }
        }
    }
    return nil
}











func (renderer *Renderer) printDebug(prev gfx.Clock) {

    if DEBUG_CLOCK||DEBUG_MODE||DEBUG_BUFFER {
        log.Debug("")
    }


    if DEBUG_CLOCK { log.Info("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(prev)) }
    
    if DEBUG_DIAG { log.Info("%s", MemUsage() ) }
        
    if DEBUG_MODE {
        tmp := ""
        switch renderer.mode { 
            case facade.Mode_GRID:
                tmp = renderer.grid.Desc()
        }
		tmp2 := ""
		if renderer.debug {
			tmp2 = " DEBUG"	
		}
        log.Info("%s %s %s %s%s",tmp,renderer.camera.Desc(),renderer.font.Desc(),renderer.mask.Desc(),tmp2)
    }

    if DEBUG_BUFFER && !log.DebugLogging() { log.Info( "%s", renderer.grid.DescBuffer() ) }
    if DEBUG_BUFFER &&  log.DebugLogging() { renderer.dumpBuffer() }
    
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




func (renderer *Renderer) Desc() string { 
    return fmt.Sprintf("renderer[%dx%d]",int(renderer.screen.W),int(renderer.screen.H))
}



