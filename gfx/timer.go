

package gfx


import (
    "fmt"
    math "../math32"
    log "../log"
)


const DEBUG_TIMER = true


type Timer struct {

    count uint
    fader float32
    Fun func()
    
    start float32
    duration float32
    repeat bool
}

func (timer *Timer) Count() uint    { return timer.count }
func (timer *Timer) Fader() float32 { return timer.fader }





func NewTimer(duration float32, repeat bool) *Timer {
    ret := &Timer{start: NOW(), duration: duration, repeat: repeat}
    RegisterTimer(ret)
    return ret
}



func (timer *Timer) Close() {
    //todo: dereg
}


func (timer *Timer) Reset() {
    timer.start = NOW()
    timer.fader = 0.0
    timer.count = 0    
}

func (timer *Timer) Start() {
    timer.start = NOW()
    timer.fader = 0.0
    timer.count = 0    
    if DEBUG_TIMER { log.Debug("start %s",timer.Desc()) }
}




func (timer *Timer) Update() bool {
    t := NOW() - timer.start
    d := timer.duration
    
    timer.fader = math.Clamp( t/d )
    
    //triggered?
    if NOW() > timer.start+timer.duration {

        timer.count += 1

        if DEBUG_TIMER { log.Debug("trigger %s",timer.Desc()) }


        if timer.Fun != nil {
            timer.Fun()
        }
     
        if timer.repeat {
            timer.start = NOW()
        }
     
        return true   
    }        
    return false
}



func (timer *Timer) Desc() string { 
    return fmt.Sprintf("timer[%4.2f⤢ #%d]",timer.fader,timer.count)
}
