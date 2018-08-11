
package render

import (
    "time"
    gfx "../gfx"
    log "../log"
)

const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

type Renderer struct {
    buffer *Buffer
    textChan chan Text
}

func NewRenderer() *Renderer {
    return &Renderer{}    
}

const DEBUG_CLOCK  = false
const DEBUG_BUFFER = false

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init() error {
    log.Debug("initialize renderer")
    renderer.buffer = NewBufferDebug(4)
    var err error
    if err != nil {
        log.Fatal("could not initialize renderer: %s",err)    
    }
    return err
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
        renderer.buffer.mutex.Lock()
        


        if now.frame % DEBUG_FRAMES == 0 {
            if DEBUG_CLOCK   {
                fps := float64(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
                prev = now
            }
        
            if DEBUG_BUFFER {
                log.Debug("\n%.1fs %s\n%s",now.time,renderer.buffer.Desc(),renderer.buffer.Debug(gfx.PageDown)) 
            }
        }

        renderer.buffer.mutex.Unlock()
        
        // wait for next frame
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}


func (renderer *Renderer) ReadText(textChan chan Text) error {
    for {
        text := <-textChan
//        log.Debug(">> %s",text)
        renderer.buffer.Queue(text,1.0)
    }
    return nil
    
}
