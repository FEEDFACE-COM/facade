
// +build linux,arm

package render

import (
    "fmt"
    "time"
    "sync"
//    "runtime"
    log "../log"
    conf "../conf"
    modes "../modes"
    gfx "../gfx"
    "src.feedface.com/gfx/piglet"
    gl "src.feedface.com/gfx/piglet/gles2"
)





const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

const BUFFER_SIZE = 80

type Renderer struct {
    screen gfx.Size

    config conf.Config

    grid *modes.Grid
    lines *modes.Lines
    test *modes.Test
    font *gfx.Font
    camera *gfx.Camera
    mask *gfx.Mask
    
    axis *gfx.Axis

    now Clock
    buffer *gfx.Buffer
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    ret.buffer = gfx.NewBuffer(BUFFER_SIZE)
    return ret
}

const DEBUG_CLOCK  = false
const DEBUG_MODE   = false
const DEBUG_BUFFER = true
const DEBUG_DIAG   = false
 

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init(config *conf.Config) error {
    var err error
    log.Debug("initializing renderer")
    
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
    renderer.axis = &gfx.Axis{}
    
    renderer.grid = modes.NewGrid(config.Grid)
    renderer.lines = modes.NewLines(config.Lines)
    renderer.test = modes.NewTest(config.Test)
    renderer.font = gfx.GetFont(config.Font,conf.DIRECTORY)
    renderer.camera = gfx.NewCamera(config.Camera,renderer.screen)
    renderer.mask = gfx.NewMask(config.Mask,renderer.screen)

    renderer.font.Configure(config.Font)
    renderer.camera.Configure(config.Camera)

    InitClock()
    renderer.now = Clock{}

    return err
}


func (renderer *Renderer) Configure(config *conf.Config) error {
    
    if config == nil {
        log.Error("renderer config nil")
        return nil
    }
    
    log.Debug("configure %s",config.Desc())
    
    if renderer.config.Mode != config.Mode {
        log.Debug("switch mode %s",string(config.Mode))
    }
    renderer.config = *config
    
    if config.Font != nil {
        newFont := gfx.GetFont(config.Font, conf.DIRECTORY)
        newFont.Init()
        renderer.font.Close()
        renderer.font = nil
//        oldFont := renderer.font
        renderer.font = newFont
//        oldFont.Close()
    }
    
    
    renderer.font.Configure(config.Font)
    renderer.lines.Configure(config.Lines)
    renderer.grid.Configure(config.Grid,renderer.font)
    renderer.test.Configure(config.Test)
    renderer.camera.Configure(config.Camera)
    renderer.mask.Configure(config.Mask)
    return nil
}

func (renderer *Renderer) Desc() string { 
    return fmt.Sprintf("renderer[%s]",renderer.screen.Desc())
}

func (renderer *Renderer) Render(confChan chan conf.Config, textChan chan conf.Text) error {

    
    var now *Clock = &renderer.now
    var prev Clock = *now

    log.Debug("renderer start")
    gl.ClearColor(0.5,0.5,0.5,1)
    gl.Viewport(0, 0, int32(renderer.screen.W),int32(renderer.screen.H))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

//    gl.Enable(gl.CULL_FACE)
//    gl.CullFace(gl.BACK)


    renderer.font.Init()
    renderer.camera.Init()
    renderer.grid.Init(renderer.camera,renderer.font)
    renderer.lines.Init(renderer.camera,renderer.font)
    renderer.test.Init(renderer.camera,renderer.font)
    renderer.mask.Init()
    renderer.axis.Init()


    renderer.grid.FillTest("coord",renderer.font)

    for {
//        if e := gl.GetError(); e != gl.NO_ERROR && debug { log.Error("pre render gl error: %s",gl.ErrorString(e)) }
        
        verbose := now.frame % DEBUG_FRAMES == 0
        
//        if verbose { log.Debug("render %s",renderer.Desc()) }
        
        now.Tick()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.ReadChannels(confChan,textChan)
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )


        switch renderer.config.Mode {
            case conf.GRID:
                renderer.grid.Render(renderer.camera, renderer.font, renderer.config.Debug, verbose )
            case conf.LINES:
                renderer.lines.Render(renderer.camera, renderer.config.Debug, verbose )
            case conf.TEST:
                renderer.test.Render(renderer.camera, renderer.config.Debug, verbose)
        }
      
        gl.Disable(gl.DEPTH_TEST)
        renderer.mask.Render()
        if renderer.config.Debug {renderer.axis.Render(renderer.camera) }
        
        if verbose { renderer.PrintDebug(now,&prev); prev = *now }

        if verbose {
            log.Debug("draw %s %s %s %s ",renderer.Desc(),renderer.grid.Desc(),renderer.camera.Desc(),renderer.font.Desc())    
        }    

        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??

        if e := gl.GetError(); e != gl.NO_ERROR && verbose { log.Error("post render gl error: %s",gl.ErrorString(e)) }
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}






func (renderer *Renderer) ReadChannels(confChan chan conf.Config, textChan chan conf.Text) {

    select {
        case config := <-confChan:
            renderer.Configure(&config)
        default:
    }
    
    select {
        case text := <-textChan:
//            renderer.buffer.Queue( gfx.NewText(string(text)) )
            log.Debug("queue %s",text)
            renderer.lines.Queue( string(text), renderer.font )
            renderer.grid.Queue(string(text), renderer.font)
            renderer.test.Queue(string(text))
        default:
    }
    
}


func (renderer *Renderer) PrintDebug(now *Clock, prev *Clock) {

    if DEBUG_CLOCK   {
        fps := float32(now.frame - prev.frame) / (now.time - prev.time)
        log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
    }
    
    if DEBUG_DIAG {
        log.Debug( MemUsage() )    
    }
    
    
    if DEBUG_BUFFER {
//        log.Debug(renderer.buffer.Dump())    
        switch renderer.config.Mode { 
            case conf.LINES:
                log.Debug( renderer.lines.Dump() )
            case conf.GRID:
                log.Debug( renderer.grid.Dump() )
        } 
    }
    
    if DEBUG_MODE {
        switch renderer.config.Mode { 
            case conf.LINES:
                log.Debug( renderer.lines.Desc() + " " +renderer.font.Desc() )
            case conf.GRID:
                log.Debug( renderer.grid.Desc() + " " +renderer.font.Desc() )
        }
    }
    
}

//
//func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
//    for {
//        text := <-textChan
//        log.Debug("read: %s",text)
//        newText := gfx.NewText(string(text))
//        renderer.mutex.Lock()
//        renderer.buffer.Queue(newText)
////        renderer.lines.Queue(string(text))
////        renderer.grid.Queue( string(text) )
////        renderer.lines.Queue( string(text) )
//        renderer.mutex.Unlock()
//    }
//    return nil
//    
//}
//
//
//func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
//    for {
//        config := <-confChan
//        log.Debug("conf: %s",config.Desc())    
//        renderer.mutex.Lock()
////        renderer.Configure(&config)
//        renderer.mutex.Unlock()
//    }
//    return nil
//}
//
//



