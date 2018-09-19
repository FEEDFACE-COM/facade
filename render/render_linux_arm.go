
package render

import (
//    "fmt"
    "time"
    "sync"
    log "../log"
    conf "../conf"
    grid "../grid"
    gfx "../gfx"
    "src.feedface.com/gfx/piglet"
    gl "src.feedface.com/gfx/piglet/gles2"
)





const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

type Renderer struct {
    size struct{width int32; height int32}
    mode conf.Mode
    grid *grid.Grid
    font *gfx.Font
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    return ret
}

const DEBUG_CLOCK  = true
const DEBUG_CONF   = false
const DEBUG_TEXT   = false

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init() error {
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
    

    log.Notice("got renderer %s %s", gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Notice("got version %s %s", gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    //setup things    
    renderer.mode = conf.DEFAULT_MODE
    renderer.grid = grid.NewGrid()
    renderer.font = gfx.NewFont()

    return err
}


func (renderer *Renderer) Configure(config *conf.Config) error {
    log.Debug("configure %s",config.Describe())
    
    if renderer.mode != config.Mode {
        log.Debug("switch mode to %s",string(config.Mode))
    }
    
    renderer.mode = config.Mode
    switch (config.Mode) {
        case conf.GRID:
            renderer.grid.Configure(config.Grid)
    }
    if config.Font != nil {
        renderer.font.Configure(config.Font,conf.DIRECTORY)
    }
    return nil
}



func (renderer *Renderer) Render() error {

    InitClock()

    var now  Clock = Clock{frame: 0}
    var prev Clock = Clock{frame: 0}
    


    log.Debug("renderer start")
//    gl.ClearColor(0x0, 0x0, 0x0, 1.0)
    gl.ClearColor(0xFE,0xED,0xFA,0xCE)
    gl.Viewport(0, 0, renderer.size.width,renderer.size.height)


    for {
        now.Tick()

        //draw
        // TBD
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )


        switch renderer.mode {
            case conf.GRID:
                renderer.grid.Render()
                
        }


        if now.frame % DEBUG_FRAMES == 0 {
            if DEBUG_CLOCK   {
                fps := float64(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Describe(),fps)
                prev = now
            }
            
        }
        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}


func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    for {
        text := <-textChan
//        log.Debug("read text: %s",text)
        renderer.mutex.Lock()
        renderer.grid.Queue( string(text) )
        renderer.mutex.Unlock()
    }
    return nil
    
}


func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    for {
        config := <-confChan
//        log.Debug("read config: %s",config.Describe())    
        renderer.mutex.Lock()
        renderer.Configure(&config)
        renderer.mutex.Unlock()
    }
    return nil
}

