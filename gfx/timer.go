
// +build linux,arm

package gfx


import (
    "fmt"
    math "../math32"
)




type Timer struct {
    Name string

    Count uint
    Fader float32
    Cycle float32
    Fun func()
    
    start float32
    duration float32
}


func Clamp(x float32) float32 {
    if x < 0.0 { return 0.0 }
    if x > 1.0 { return 1.0 }
    return x   
}


func NewTimer(now *Clock, duration float32) *Timer {
    ret := &Timer{start: now.time, duration: duration}
    now.Register(ret)
    return ret
}

func (timer *Timer) Close() {
    //todo: dereg
}




func (timer *Timer) Update(now float32) bool {
    t := now - timer.start
    d := timer.duration
    
    timer.Fader = Clamp( t/d )
    timer.Cycle = math.TAU * timer.Fader
    
    if now > timer.start+timer.duration {
        if timer.Fun != nil {
            timer.Fun()
        }
        timer.Count += 1
        timer.start += timer.duration
        return true
    }
    return false
}



func (timer *Timer) Desc() string { 
    return fmt.Sprintf("t[%s %4.2f↺ %4.2f⤢ %d#]",timer.Name,timer.Cycle,timer.Fader,timer.Count)
}
