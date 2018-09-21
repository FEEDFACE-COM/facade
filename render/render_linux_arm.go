
// +build linux,arm

package render

import (
//    "fmt"
    "time"
    "sync"
    "runtime"
    log "../log"
    conf "../conf"
    modes "../modes"
    gfx "../gfx"
    "src.feedface.com/gfx/piglet"
    gl "src.feedface.com/gfx/piglet/gles2"
)





const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

const BUFFER_SIZE = 40

type Renderer struct {
    size struct{width int32; height int32}

    mode conf.Mode
    grid *modes.Grid
    lines *modes.Lines
    font *gfx.Font

    now Clock
    buffer *Buffer
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    ret.buffer = &Buffer{}
    ret.buffer.Resize(BUFFER_SIZE)
    return ret
}

const DEBUG_CLOCK  = false
const DEBUG_MODE   = false
const DEBUG_BUFFER = false
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
    renderer.size = struct{width int32; height int32} {int32(w),int32(h)}
    log.Info("got display %dx%d",renderer.size.width,renderer.size.height)
    

    piglet.MakeCurrent()

    err = gl.InitWithProcAddrFunc( piglet.GetProcAddress )
    if err != nil {
        log.PANIC("fail to init GLES: %s",err)    
    }
    

    log.Debug("got renderer %s %s", gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Debug("got version %s %s", gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    //setup things    
    renderer.mode = config.Mode
    renderer.grid = modes.NewGrid(config.Grid)
    renderer.lines = modes.NewLines(config.Line)
    renderer.font = gfx.NewFont(config.Font)

    renderer.font.Configure(config.Font,conf.DIRECTORY)

    InitClock()
    renderer.now = Clock{frame: 0}

    return err
}


func (renderer *Renderer) Configure(config *conf.Config) error {
    
    if config == nil {
        return nil
    }
    
    log.Debug("configure %s",config.Desc())
    
    if renderer.mode != config.Mode {
        log.Debug("switch mode to %s",string(config.Mode))
    }
    
    renderer.mode = config.Mode
    renderer.font.Configure(config.Font,conf.DIRECTORY)
    renderer.lines.Configure(config.Line)
    renderer.grid.Configure(config.Grid)
    return nil
}



func (renderer *Renderer) Render(confChan chan conf.Config, textChan chan conf.Text) error {

    
    var now *Clock = &renderer.now
    var prev Clock = *now

    log.Debug("renderer start")
    gl.ClearColor(0xFE,0xED,0xFA,0xCE)
    gl.Viewport(0, 0, renderer.size.width,renderer.size.height)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

//    gl.Enable(gl.CULL_FACE)
//    gl.CullFace(gl.BACK)


	camera := gfx.NewCamera( float32(renderer.size.width), float32(renderer.size.height) )
    renderer.lines.Init(camera)
    renderer.grid.Init(camera)

    for {
        now.Tick()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.ReadChannels(confChan,textChan)
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )

        switch renderer.mode {
            case conf.GRID:
                renderer.grid.Render(camera)
            case conf.LINE:
                renderer.lines.Render(camera)
        }
        
        if now.frame % DEBUG_FRAMES == 0 { renderer.PrintDebug(now,&prev) }

        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}






func (renderer *Renderer) ReadChannels(confChan chan conf.Config, textChan chan conf.Text) {

    select {
        case config := <-confChan:
            log.Debug("conf: %s",config.Desc())
            renderer.Configure(&config)
            log.Debug( renderer.lines.Dump() )
            runtime.GC()
        default:
    }
    
    select {
        case text := <-textChan:
//                log.Debug("read: %s",text)
//                renderer.buffer.Queue(renderer.now.Time(),text)
            renderer.lines.Queue( string(text), renderer.font )
            renderer.grid.Queue( string(text) )
            log.Debug( renderer.lines.Dump() )
            runtime.GC()
        default:
    }
    
}


func (renderer *Renderer) PrintDebug(now *Clock, prev *Clock) {

    if DEBUG_CLOCK   {
        fps := float32(now.frame - prev.frame) / (now.time - prev.time)
        log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
        *prev = *now
    }
    
    if DEBUG_DIAG {
        log.Debug( MemUsage() )    
    }
    
    
    if DEBUG_BUFFER {
//        log.Debug("%5.1f %5.1f %s",renderer.now.Time(),now.Time(),renderer.buffer.Dump())    
        switch renderer.mode { 
            case conf.LINE:
                log.Debug( renderer.lines.Dump() )
            case conf.GRID:
                log.Debug( renderer.grid.Dump() )
        } 
    }
    
    if DEBUG_MODE {
        switch renderer.mode { 
            case conf.LINE:
                log.Debug( renderer.lines.Desc() + " " +renderer.font.Desc() )
            case conf.GRID:
                log.Debug( renderer.grid.Desc() + " " +renderer.font.Desc() )
        }
    }
    
}


func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    for {
        text := <-textChan
        log.Debug("read: %s",text)
        renderer.mutex.Lock()
        renderer.buffer.Queue(renderer.now.Time(),text)
//        renderer.lines.Queue(string(text))
//        renderer.grid.Queue( string(text) )
//        renderer.lines.Queue( string(text) )
        renderer.mutex.Unlock()
    }
    return nil
    
}


func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    for {
        config := <-confChan
        log.Debug("conf: %s",config.Desc())    
        renderer.mutex.Lock()
//        renderer.Configure(&config)
        renderer.mutex.Unlock()
    }
    return nil
}




var quadVertices = []float32{
	//  X, Y, Z, U, V

//	// Front
    -1.0,  1.0, 0.0, 0.0, 0.0,
    -1.0, -1.0, 0.0, 0.0, 1.0,
     1.0, -1.0, 0.0, 1.0, 1.0,
    
    -1.0,  1.0, 0.0, 0.0, 0.0,
     1.0, -1.0, 0.0, 1.0, 1.0,
     1.0,  1.0, 0.0, 1.0, 0.0,

}

