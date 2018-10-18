
// +build linux,arm

package gfx


import (
    "fmt"
    math "../math32"
)




type Timer struct {
    Name string

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
}




func (timer *Timer) Update() bool {
    t := NOW() - timer.start
    d := timer.duration
    
    timer.fader = math.Clamp( t/d )
    timer.count += 1
    
    //triggered?
    if NOW() > timer.start+timer.duration {

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
    return fmt.Sprintf("t[%s "+/*%4.2f↺ */"%4.2f⤢ %d#]",timer.Name,timer.fader,timer.count)
}
