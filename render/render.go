
package render

import (
    "time"
    gfx "../gfx"
    log "../log"
)

const FRAME_RATE = 60.0

type Renderer struct {
}

func NewRenderer() *Renderer {
    return &Renderer{}    
}

const DEBUG_CLOCK  = false
const DEBUG_BUFFER = true


func (renderer *Renderer) Init() error {
    return nil
}


func (renderer *Renderer) Start() error {

    InitClock()

    var now  Clock = Clock{frame: 0}
    var prev Clock = Clock{frame: 0}
    

//    buffer := EmptyBuffer(4)
    buffer := DebugBuffer(4)

    log.Debug("renderer start")
    for {
        now.Tick()

        //draw
        // TBD


        if now.frame % 50 == 0 {
            if DEBUG_CLOCK   {
                fps := float64(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
                prev = now
            }
        
            if DEBUG_BUFFER {
                log.Debug("\n%.1fs %s\n%s",now.time,buffer.Desc(),buffer.Debug(gfx.PageDown)) 
            }
        }
        
        // wait for next frame
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}

