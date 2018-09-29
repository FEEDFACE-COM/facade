
// +build linux,arm

package main

import (
    "fmt"
    "strings"
    "time"
    "sync"
    "os"
//    "runtime"
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
    lines *facade.Lines
    test *facade.Test
    font *gfx.Font
    camera *gfx.Camera
    mask *gfx.Mask
    
    axis *gfx.Axis

    now gfx.Clock
    buffer *gfx.Buffer
    mutex *sync.Mutex
    directory string
}

const DEBUG_CLOCK  =  true
const DEBUG_MODE   = false
const DEBUG_BUFFER = false
const DEBUG_DIAG   = false
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
    renderer.axis = &gfx.Axis{}

    renderer.font,err = gfx.GetFont(config.Font)
    if err != nil {
        log.PANIC("no default font: %s",err)    
    }
    renderer.camera = gfx.NewCamera(config.Camera,renderer.screen)
    renderer.mask = gfx.NewMask(config.Mask,renderer.screen)

    renderer.font.Configure(config.Font)
    renderer.camera.Configure(config.Camera)

    
    renderer.grid = facade.NewGrid(config.Grid)
    renderer.lines = facade.NewLines(config.Lines)
    renderer.test = facade.NewTest(config.Test)


    gfx.StartClock()
    renderer.now = gfx.Clock{}

    return err
}


func (renderer *Renderer) Configure(config *facade.Config) error {
    
    if config == nil { log.Error("renderer config nil") ;return nil }
    
    log.Debug("conf %s",config.Desc())
    
    
    old := renderer.config

    renderer.config = *config
    if renderer.config.Mode != old.Mode {
        log.Debug("switch mode -> %s",string(config.Mode))
    }
    
    if config.Font != nil && config.Font != old.Font {
        newFont,err := gfx.GetFont(config.Font)
        if err == nil {
            log.Debug("switch font -> %s",string(config.Font.Name))
            newFont.Init()
            renderer.font = newFont
            renderer.font.Configure(config.Font)
//        oldFont.Close() //still, leaks?
        }

    }
    
    
    renderer.lines.Configure(config.Lines)
    renderer.grid.Configure(config.Grid,renderer.camera,renderer.font)
    renderer.test.Configure(config.Test)
    renderer.camera.Configure(config.Camera)
    renderer.mask.Configure(config.Mask)
    return nil
}


func (renderer *Renderer) Render(confChan chan facade.Config, textChan chan string) error {

    
    var now *gfx.Clock = &renderer.now
    var prev gfx.Clock = *now

    now.Tick()

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
    renderer.grid.Init(now, renderer.camera,renderer.font)
    renderer.lines.Init(renderer.camera,renderer.font)
    renderer.test.Init(renderer.camera,renderer.font)
    renderer.mask.Init()
    renderer.axis.Init()


    renderer.grid.Configure(renderer.config.Grid,renderer.camera,renderer.font)

    renderer.grid.FillTest("author",renderer.font)



    log.Debug("renderer start")
    for {
        
        verbose := now.DebugFrame()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        renderer.QueueTexts(textChan)
        renderer.QueueConfs(confChan)
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )


        switch facade.Mode(renderer.config.Mode) {
            case facade.GRID:
                renderer.grid.Render(renderer.camera, renderer.font, renderer.config.Debug, verbose )
            case facade.LINES:
                renderer.lines.Render(renderer.camera, renderer.config.Debug, verbose )
            case facade.TEST:
                renderer.test.Render(renderer.camera, renderer.config.Debug, verbose)
        }
      
        gl.Disable(gl.DEPTH_TEST)
        renderer.mask.Render()
        if renderer.config.Debug {renderer.axis.Render(renderer.camera) }
        
        if verbose { renderer.PrintDebug(*now,prev); prev = *now }

        if verbose {
            log.Debug("draw %s %s %s %s ",renderer.Desc(),renderer.grid.Desc(),renderer.camera.Desc(),renderer.font.Desc())    
        }    

        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??

        if e := gl.GetError(); e != gl.NO_ERROR && verbose { log.Error("post render gl error: %s",gl.ErrorString(e)) }
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        now.Tick()
    }
    return nil
}


func (renderer *Renderer) QueueTexts(textChan chan string) {

    select {
        case text := <-textChan:
            
            renderer.lines.Queue(text, renderer.font )
            renderer.grid.Queue(text, renderer.font)
            renderer.test.Queue(text)
        
        default:
            //nop    
    }
    
}



func (renderer *Renderer) QueueConfs(confChan chan facade.Config) {
    
    select {
        case conf := <-confChan:
        
            renderer.Configure(&conf)
        
        
        default:
            //nop    
    }    
}







func (renderer *Renderer) PrintDebug(now gfx.Clock, prev gfx.Clock) {

    if DEBUG_CLOCK   {
        log.Debug("%s    %4.1ffps",now.Desc(),now.Delta(prev))
    }
    
    if DEBUG_DIAG {
        log.Debug( MemUsage() )    
    }
    
    
    if DEBUG_BUFFER {
//        log.Debug(renderer.buffer.Dump())    
        switch facade.Mode(renderer.config.Mode) { 
            case facade.LINES:
                log.Debug( renderer.lines.Dump() )
            case facade.GRID:
                log.Debug( renderer.grid.Dump() )
        } 
    }
    
    if DEBUG_MODE {
            tmp := ""
        switch facade.Mode(renderer.config.Mode) { 
            case facade.LINES:
                tmp = renderer.lines.Desc()
            case facade.GRID:
                tmp = renderer.grid.Desc()
        }
        log.Debug("%s %s %s %s ",tmp,renderer.camera.Desc(),renderer.font.Desc(),renderer.Desc())    
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



