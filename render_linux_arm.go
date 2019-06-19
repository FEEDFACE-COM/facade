
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

    state facade.State

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
    log.Debug("init renderer[%s] %s",renderer.directory,config.Desc())
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
    

    log.Debug("got renderer %s %s", gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Debug("got version %s %s", gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    //setup things 
	renderer.state = facade.Defaults
    renderer.state.ApplyConfig(config)



//    renderer.config.Clean()
    
    renderer.axis = &gfx.Axis{}

	fontConfig := gfx.FontDefaults.Config()
	if cfg,ok := config.Font(); ok {
		fontConfig.ApplyConfig( &cfg )	
	}
	renderer.font,err = gfx.GetFont( fontConfig )
    if err != nil {
        log.PANIC("no default font: %s",err)    
    }
	renderer.font.Init()


    cameraConfig := gfx.CameraDefaults.Config()
    if cfg,ok := config.Camera(); ok {
		cameraConfig.ApplyConfig( &cfg )    
	}
    renderer.camera = gfx.NewCamera(cameraConfig,renderer.screen)
    renderer.camera.Init(cameraConfig)

    maskConfig := gfx.MaskDefaults.Config()
    if cfg,ok := config.Mask(); ok {
		maskConfig.ApplyConfig(&cfg)    
	}
    renderer.mask = gfx.NewMask(maskConfig,renderer.screen)
    renderer.mask.Init()



    gridConfig := facade.GridDefaults.Config()
    if cfg,ok := config.Grid(); ok {
        gridConfig.ApplyConfig(&cfg)
    }	

    width,_ := gridConfig.Width()
    height,_ := gridConfig.Height()
    buflen,_ := gridConfig.BufLen()


    renderer.termBuffer = facade.NewTermBuffer(width,height) 
    renderer.lineBuffer = facade.NewLineBuffer(height,buflen,renderer.refreshChan) 

    //initialize mode, REM this should probably init all modes
	switch renderer.state.Mode {
		case facade.GRID:
			renderer.grid = facade.NewGrid( gridConfig, renderer.lineBuffer, renderer.termBuffer )
			renderer.grid.Init(renderer.camera,renderer.font)
			renderer.grid.Configure(gridConfig,renderer.camera,renderer.font)
	}

    renderer.axis.Init()


    gfx.ClockReset()
    return err
}


func (renderer *Renderer) Configure(config *facade.Config) error {
    
    if config == nil { log.Error("renderer config nil") ;return nil }
    if len(*config) <= 0 { return nil }
    
    log.Debug("renderer config %s",config.Desc())
    
    if tmp,ok := config.Font(); ok {
		newFont, err := gfx.GetFont(&tmp)
		if err != nil {
			log.Error("fail to get font %s",tmp.Desc())
		} else {
			newFont.Init()
			renderer.font = newFont
			if renderer.grid != nil { 
				cfg := make(facade.GridConfig)
				cfg.SetHeight( renderer.grid.Height() )
				renderer.grid.Configure(&cfg,renderer.camera,renderer.font)
			}
		}
	}
    
    if tmp,ok := config.Camera(); ok {
		renderer.camera.Configure(&tmp)    
	}
    
    if tmp,ok := config.Mask(); ok {
		renderer.mask.Configure(&tmp)    
	}

    if tmp,ok := config.Grid(); ok {
		renderer.grid.Configure(&tmp,renderer.camera,renderer.font)    
	}
	
	if debug,ok := config.Debug(); ok {
		renderer.state.Debug = debug	
	} else {
		renderer.state.Debug = false	
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
    log.Debug("render %s",renderer.state.Desc())
    for {
        
        verboseFrame := gfx.ClockVerboseFrame()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.ProcessConf(confChan)
        if renderer.checkRefresh() {
            switch renderer.state.Mode {
                case facade.GRID:
                    renderer.grid.GenerateData(renderer.font)
            }
        }
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )

        gfx.RefreshPrograms()


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA,gl.ZERO,gl.ONE)
        switch renderer.state.Mode {
            case facade.GRID:
                renderer.grid.Render(renderer.camera, renderer.font, renderer.state.Debug, verboseFrame )
        }
      
        if renderer.state.Debug {renderer.axis.Render(renderer.camera, renderer.state.Debug) }


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA,gl.ZERO,gl.ONE)
        renderer.mask.Render(renderer.state.Debug)
        
        
        
        if verboseFrame { 

            if DEBUG_BUFFER {
                renderer.dumpBuffers()    
            }

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
            if DEBUG_MEMORY { log.Debug("mem now %s",MemUsage())}
        
        
        default:
            //nop    
    }
}









func (renderer *Renderer) ProcessBufferItems(bufChan chan facade.BufferItem) error {

    for {
        item := <- bufChan    
        if DEBUG_MESSAGES { log.Debug("buffer %s",item.Desc()) }
        text, seq := item.Text, item.Seq
        if text != nil && len(text) > 0 {
            renderer.lineBuffer.ProcessRunes( text )
            renderer.termBuffer.ProcessRunes( text )    
            renderer.ScheduleRefresh()
            if DEBUG_BUFFER {
                renderer.dumpBuffers()    
            }
        }
        if seq != nil {
            renderer.lineBuffer.ProcessSequence( seq )
            renderer.termBuffer.ProcessSequence( seq )
            renderer.ScheduleRefresh()
            if DEBUG_BUFFER {
                renderer.dumpBuffers()    
            }
        }
    }
    return nil
}









func (renderer *Renderer) ProcessRawConfs(rawChan chan facade.Config, confChan chan facade.Config) error {
    for {
        rawConf := <-rawChan
        if DEBUG_MESSAGES { log.Debug("process %s",rawConf.Desc()) }


//        renderer.mutex.Lock()
//        // prep some stuff i guess?
//        renderer.mutex.Unlock()
        
        confChan <- rawConf

    }
    return nil
}




func (renderer *Renderer) printDebug(prev gfx.Clock) {

    if DEBUG_CLOCK { log.Debug("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(prev)) }
    
    if DEBUG_DIAG { log.Debug( MemUsage() ) }
        
    if DEBUG_MODE {
        tmp := ""
        switch renderer.state.Mode { 
            case facade.GRID:
                tmp = renderer.grid.Desc()
        }
		tmp2 := ""
		if renderer.state.Debug {
			tmp2 = " DEBUG"	
		}
        log.Debug("%s %s %s %s%s",tmp,renderer.camera.Desc(),renderer.font.Desc(),renderer.mask.Desc(),tmp2)
    }
    
}

func (renderer *Renderer) dumpBuffers() {
    if renderer.state.Mode  == facade.GRID {
        os.Stdout.Write( []byte( renderer.grid.Dump() ) )        
    }
    os.Stdout.Write( []byte( "\n" ) )
    os.Stdout.Sync()
}




func (renderer *Renderer) Desc() string { 
    return fmt.Sprintf("renderer[%dx%d]",int(renderer.screen.W),int(renderer.screen.H))
}



