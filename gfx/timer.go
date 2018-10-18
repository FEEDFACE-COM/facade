
// +build linux,arm

package gfx


import (
    "fmt"
//    math "../math32"
)




type Timer struct {
    Name string

    Count uint
    Fader float32
    Fun func()
    
    
    
    start float32
    duration float32
    repeat bool
}


func Clamp(x float32) float32 {
    if x < 0.0 { return 0.0 }
    if x > 1.0 { return 1.0 }
    return x   
}


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
    timer.Fader = 0.0
    timer.Count = 0    
}

func (timer *Timer) Start() {
    timer.start = NOW()
    timer.Fader = 0.0
    timer.Count = 0    
}




func (timer *Timer) Update() bool {
    t := NOW() - timer.start
    d := timer.duration
    
    timer.Fader = Clamp( t/d )
    timer.Count += 1
    
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
    return fmt.Sprintf("t[%s "+/*%4.2f↺ */"%4.2f⤢ %d#]",timer.Name,timer.Fader,timer.Count)
}
