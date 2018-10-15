
// +build linux,arm

package gfx

import (
    "time"
    "fmt"
    log "../log"
)

const DEBUG_FRAMES = 90
const CLOCK_RATE = 1.0

var clockStart time.Time

var clockFrame uint
var clockTime float32

var clockTimers []*Timer

func NOW() float32 { return clockTime }


type Clock struct {
    frame uint
    time  float32
}

func NewClock() *Clock {
    return &Clock{frame: clockFrame, time: clockTime}
}

func ClockReset() {
    clockStart = time.Now()    
    clockFrame = 0
    clockTime = 0.0
}



func RegisterTimer(timer *Timer) {
    clockTimers = append(clockTimers, timer )    
    log.Debug("register %s",(*timer).Desc())
}





func ClockDebug() bool { return clockFrame % DEBUG_FRAMES == 0 }


func ClockDelta(prev Clock) float32 { 
    return float32(clockFrame - prev.frame) / (clockTime-prev.time) 
}



func ClockTick() {
    
    clockFrame += 1
    clockTime = CLOCK_RATE *  float32( time.Now().Sub(clockStart) ) / (1000.*1000.*1000.)
    
    
    for _,timer := range( clockTimers ) {
        trigger := (*timer).Update(clockTime)
        if trigger {
//            log.Debug("trigger %s",timer.Desc())
        }
    }
    
    
}

func ClockDesc() string {
    return fmt.Sprintf("frame #%05d %7.2fs",clockFrame,clockTime)
}
