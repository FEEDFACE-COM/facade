
package render

import (
    "time"
    "sync"
    log "../log"
    proto "../proto"
    gfx "../gfx"
)


const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

type Renderer struct {
    mode proto.Mode
    pager *gfx.Pager
    font *gfx.Font
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    return ret
}

const DEBUG_CLOCK  = false
const DEBUG_CONF = true
const DEBUG_TEXT = true

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init() error {
    var err error
    log.Debug("initializing renderer")
    if err != nil {
        log.Fatal("could not initialize renderer: %s",err)    
    }
    
    
    config := proto.NewConfig(proto.PAGER)
    renderer.Config(config)
    return err
}


func (renderer *Renderer) Config(config *proto.Config) {
    renderer.mode = config.Mode
    switch (config.Mode) {
        case proto.PAGER:
            renderer.pager = gfx.NewPager(config.Pager)    
    }
    if config.Font != nil {
        renderer.font = gfx.NewFont(config.Font)
    }
    if renderer.font == nil {
        renderer.font = gfx.NewFont(proto.NewFont())    
    }
}



func (renderer *Renderer) Render() error {

    InitClock()

    var now  Clock = Clock{frame: 0}
    var prev Clock = Clock{frame: 0}
    


    log.Debug("renderer start")
    for {
        now.Tick()

        //draw
        // TBD
        renderer.mutex.Lock()
        


        switch renderer.mode {
            case proto.PAGER:
                renderer.pager.Render()
        }


        if now.frame % DEBUG_FRAMES == 0 {
            if DEBUG_CLOCK   {
                fps := float64(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
                prev = now
            }

            if DEBUG_CONF {
                str := string(renderer.mode)
                str += " " + renderer.font.Desc()
                switch renderer.mode {
                    case proto.PAGER:
                        str += " " + renderer.pager.Desc()
                    }
                log.Debug(str)
            }
            
            if DEBUG_TEXT {
                log.Debug("%s\n%s",renderer.pager.Buffer().Desc(),renderer.pager.Buffer().Debug(gfx.PageDown))
            }
            
            
//            if DEBUG_RENDER {
//                log.Debug("%s\n%s",renderer.Buffer().Desc(),renderer.Buffer().Debug(gfx.PageDown)) 
//            }
//            
//            if DEBUG_GFX {
//                if renderer.mode == Pager  { log.Debug(renderer.pager.Desc())  }
//            }
            
                
            
        }

        renderer.mutex.Unlock()
        
        // wait for next frame
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}


func (renderer *Renderer) ReadText(textChan chan proto.Text) error {
    for {
        text := <-textChan
        log.Debug("read %s",text)
        renderer.mutex.Lock()        
        renderer.pager.Buffer().Queue(string(text),1.0)
        renderer.mutex.Unlock()
    }
    return nil
    
}


func (renderer *Renderer) ReadConf(confChan chan proto.Config) error {
    for {
        config := <-confChan
        log.Debug("read %s",config.Desc())    
        renderer.mutex.Lock()
        renderer.Config(&config)
        renderer.mutex.Unlock()
    }
    return nil
}

