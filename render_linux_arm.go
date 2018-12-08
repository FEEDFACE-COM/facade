
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

const FRAME_RATE = 60.0

const BUFFER_SIZE = 80

type Renderer struct {
    screen gfx.Size

    config facade.Config

    grid *facade.Grid
//    lines *facade.Lines
//    test *facade.Test
    font *gfx.Font
    camera *gfx.Camera
    mask *gfx.Mask
    
    axis *gfx.Axis

    buffer *gfx.Buffer
    mutex *sync.Mutex
    directory string
    
    mode facade.Mode
    debug bool
    
}

const DEBUG_CLOCK  = false
const DEBUG_MODE   = true
const DEBUG_BUFFER = false
const DEBUG_DIAG   = false
const DEBUG_MEMORY = true
const DEBUG_MESSAGES = false



func NewRenderer(directory string) *Renderer {
    ret := &Renderer{directory: directory}
    ret.mutex = &sync.Mutex{}
    ret.buffer = gfx.NewBuffer(BUFFER_SIZE)
    return ret
}


func (renderer *Renderer) Init(config *facade.Config) error {
    var err error
    log.Debug("init renderer[%s]",renderer.directory)
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
    renderer.config = *config
    renderer.config.Clean()
    
    renderer.mode,_ = config.Mode()
    renderer.debug,_ = config.Debug()
       
    renderer.axis = &gfx.Axis{}

    fontConfig,_ := config.Font()
    renderer.font,err = gfx.GetFont( &fontConfig)
    if err != nil {
        log.PANIC("no default font: %s",err)    
    }
    renderer.font.Init()
 


    cameraConfig,_ := config.Camera()
    renderer.camera = gfx.NewCamera(&cameraConfig,renderer.screen)
    renderer.camera.Init()

    maskConfig,_ := config.Mask()
    renderer.mask = gfx.NewMask(&maskConfig,renderer.screen)
    renderer.mask.Init()
    


    gridDelta,_ := config.Grid()    
    renderer.grid = facade.NewGrid(nil)
    renderer.grid.Init(renderer.camera,renderer.font)
    renderer.grid.Configure(&gridDelta,renderer.camera,renderer.font)
    
//    renderer.lines = facade.NewLines(config.Lines)
//    renderer.lines.Init(renderer.camera,renderer.font)

//    renderer.test = facade.NewTest(config.Test)
  //  renderer.test.Init(renderer.camera,renderer.font)

    renderer.axis.Init()


    gfx.ClockReset()
    return err
}


func (renderer *Renderer) Configure(config *facade.Config) error {
    
    if config == nil { log.Error("renderer config nil") ;return nil }
    
    config.Clean()
    log.Debug("configure %s",config.Desc())
    
    renderer.config = *config
    
    if mode,ok := config.Mode(); ok && mode != renderer.mode {
        renderer.mode = mode
        log.Debug("switch mode %s",string(mode))
    }
    
    
    if fontConfig,ok := config.Font(); ok && string(fontConfig) != renderer.font.Name() {
        log.Debug("switch %s",fontConfig)
        newFont,err := gfx.GetFont(&fontConfig)
        if err == nil {
            newFont.Init()
            renderer.font = newFont
        } else {
            log.Error("fail to get font %s",fontConfig.Desc())
        } 
    }
    
    if gridConfig,ok := config.Grid(); ok {
        renderer.grid.Configure(&gridConfig,renderer.camera,renderer.font)
    }
//    renderer.lines.Configure(config.Lines)
//    renderer.test.Configure(config.Test)

    if cameraConfig,ok := config.Camera(); ok { renderer.camera.Configure(&cameraConfig) }
    if maskConfig, ok := config.Mask(); ok { renderer.mask.Configure(&maskConfig) }
    
    return nil
}


func (renderer *Renderer) Render(confChan chan facade.Config, textChan chan string) error {

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
    log.Debug("renderer start")
    for {
        
        verbose := gfx.ClockDebug()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.ProcessTexts(textChan)
        renderer.ProcessConfs(confChan)

        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )

        gfx.RefreshPrograms()


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA,gl.ZERO,gl.ONE)
        switch renderer.mode {
            case facade.GRID:
                renderer.grid.Render(renderer.camera, renderer.font, renderer.debug, verbose )
//            case facade.LINES:
//                renderer.lines.Render(renderer.camera, renderer.debug, verbose )
//            case facade.TEST:
//                renderer.test.Render(renderer.camera, renderer.debug, verbose)
        }
      
        if renderer.debug {renderer.axis.Render(renderer.camera, renderer.debug) }


        gl.BlendEquationSeparate(gl.FUNC_ADD,gl.FUNC_ADD)
        gl.BlendFuncSeparate(gl.ONE, gl.SRC_ALPHA,gl.ZERO,gl.ONE)
        renderer.mask.Render(renderer.debug)
        
        if verbose { 
            renderer.PrintDebug(prev); 
            prev = *gfx.NewClock() 
        }

        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??


        e := gl.GetError()
//        e := uint32(gl.NO_ERROR)
        if e != gl.NO_ERROR && verbose { 
            log.Error("post render gl error: %s",gl.ErrorString(e)) 
        }
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        gfx.ClockTick()
    }
    return nil
}


func (renderer *Renderer) ProcessTexts(textChan chan string) {

    select {
        case text := <-textChan:
            
            renderer.grid.Queue(text, renderer.font)
//            renderer.lines.Queue(text, renderer.font )
//            renderer.test.Queue(text)
            if DEBUG_MEMORY { log.Debug("mem now %s",MemUsage())}
        
        default:
            //nop    
    }
}



func (renderer *Renderer) ProcessConfs(confChan chan facade.Config) {
    
    select {
        case conf := <-confChan:
        
            renderer.Configure(&conf)
            if DEBUG_MEMORY { log.Debug("mem now %s",MemUsage())}
        
        
        default:
            //nop    
    }
}







func (renderer *Renderer) PrintDebug(prev gfx.Clock) {

    if DEBUG_CLOCK   {
        log.Debug("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(prev))
    }
    
    if DEBUG_DIAG {
        log.Debug( MemUsage() )    
    }
    
    
    if DEBUG_BUFFER {
//        log.Debug(renderer.buffer.Dump())    
        switch renderer.mode { 
//            case facade.LINES:
//                log.Debug( renderer.lines.Dump() )
            case facade.GRID:
                log.Debug( renderer.grid.Dump() )
        } 
    }
    
    if DEBUG_MODE {
        tmp := ""
        switch renderer.mode { 
//            case facade.LINES:
//                tmp = renderer.lines.Desc()
            case facade.GRID:
                tmp = renderer.grid.Desc()
        }
        tmp2 := ""
        if renderer.mask.Mask() {
            tmp2 = " " + renderer.mask.Desc()
        }
        log.Debug("%s %s %s%s",tmp,renderer.camera.Desc(),renderer.font.Desc(),tmp2)
    }
    
}

//


func (renderer *Renderer) SanitizeText(raw facade.RawText) string {
    const TABWIDTH = 8
    ret := string(raw)
    ret = strings.Replace(ret, "	", strings.Repeat(" ", TABWIDTH), -1)
    return ret
}


func (renderer *Renderer) SanitizeConfig(raw facade.Config) facade.Config {
    ret := raw
    return ret
}




func (renderer *Renderer) ProcessText(rawChan chan facade.RawText, textChan chan string) error {

    for {
        rawText := <-rawChan
        if DEBUG_MESSAGES {
            log.Debug("process raw text: %s",string(rawText))
        }
        text := renderer.SanitizeText(rawText)
        
        renderer.mutex.Lock()
        renderer.buffer.Queue( gfx.NewText(text) )
        renderer.mutex.Unlock()
        
        textChan <- text
        
    }
    return nil    
    
}





func (renderer *Renderer) ProcessConf(rawChan chan facade.Config, confChan chan facade.Config) error {
    for {
        rawConf := <-rawChan
        if DEBUG_MESSAGES {
            log.Debug("process raw conf: %s",rawConf.Desc())
        }
        conf := renderer.SanitizeConfig(rawConf)

        renderer.mutex.Lock()
        // prep some stuff i guess?
        renderer.mutex.Unlock()
        
        confChan <- conf

    }
    return nil
}



func (renderer *Renderer) Desc() string { 
    return fmt.Sprintf("renderer[%dx%d]",int(renderer.screen.W),int(renderer.screen.H))
}



