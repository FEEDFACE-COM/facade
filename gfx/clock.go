

package gfx

import (
    "time"
    "fmt"
    log "../log"
)

const VERBOSE_FRAMES = 60
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
    //TODO: replace with hashmap!
    clockTimers = append(clockTimers, timer )    
    if DEBUG_TIMER { log.Debug("register %s",(*timer).Desc()) }
}


func UnRegisterTimer(timer *Timer) {
    //TODO: replace with hashmap!
    var ct []*Timer
    if DEBUG_TIMER { log.Debug("deregister %s",(*timer).Desc()) }
    for _,v := range( clockTimers ) {
        if v != timer {
            ct = append(ct,v)
        }
    }
    clockTimers = ct    
}




func ClockVerboseFrame() bool { return clockFrame % VERBOSE_FRAMES == 0 }


func ClockDelta(prev Clock) float32 { 
    return float32(clockFrame - prev.frame) / (clockTime-prev.time) 
}



func ClockTick() {
    
    clockFrame += 1
    clockTime = CLOCK_RATE *  float32( time.Now().Sub(clockStart) ) / (1000.*1000.*1000.)
    
    
    for _,timer := range( clockTimers ) {
        trigger := (*timer).Update()
        if trigger {
//            log.Debug("trigger %s",timer.Desc())
        }
    }
    
    
}

func ClockDesc() string {
    return fmt.Sprintf("frame #%05d %7.2fs",clockFrame,clockTime)
}
