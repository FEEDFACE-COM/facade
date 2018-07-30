
package render

import (
    "time"
    log "../log"
)

const FRAME_RATE = 60.0

var (
    
)


func Init() error {
    return nil
}


func Start() error {

    InitClock()

    var now  Clock = Clock{frame: 0}
    var prev Clock = Clock{frame: 0}
    

    log.Debug("render start")
    for {
        now.Tick()

        //draw
        // TBD


        // show benchmarks
        if (now.frame % 50 == 0) {
            fps := float64(now.frame - prev.frame) / (now.time - prev.time)
            log.Debug("frame #%05d %5.2ffps    %7.2fs  %4.2f↺  %4.2f⤢  %d#",
                now.frame,fps,now.time,now.cycle,now.fader,now.count)
            prev = now
        }
        // wait for next frame
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}

