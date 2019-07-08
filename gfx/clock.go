
package gfx


import (
    "time"
    "fmt"
    log "../log"
)



const DEBUG_CLOCK = true


const ClockRate = 1.0
const VerboseFrames = 60

type ClockFrame struct {
    Frame uint
    Time float32
}






func Now() float32 { return worldClock.time }
func WorldClock() *Clock { return worldClock }


func (clock *Clock) NewTimer(duration float32, repeat bool, valueFun func(float32)float32, triggerFun func()) *Timer {

    timer := &Timer{
        start: clock.time,
        duration: duration, 
        repeat: repeat, 
        valueFun: valueFun,
        triggerFun: triggerFun,
    }

    if timer.valueFun == nil {
        timer.valueFun = func(x float32) float32 { return x }
    }
    
    
    clock.timers[timer] = timer
    
    return timer
    
}


func (clock *Clock) DeleteTimer(timer *Timer) {

    tmp := clock.timers[timer]
    
    if tmp != nil {
        if DEBUG_CLOCK { log.Debug("%s delete %s",clock.Desc(),tmp.Desc()) }
        clock.timers[timer] = nil    
    } else {
        log.Debug("%s fail delete timer %s",clock.Desc(), timer.Desc())
    }

}


func (clock *Clock) Reset() {

    if DEBUG_CLOCK { log.Debug("%s reset",clock.Desc()) }

    clock.start = time.Now()
    clock.frame = 0
    clock.time = 0.0
    
    for k,_ := range( clock.timers ) {
        delete(clock.timers,k)
    }
    
}

func (clock *Clock) VerboseFrame() bool {
    return clock.frame % VerboseFrames == 0  
}



func (clock *Clock) Tick() {
 
    clock.frame += 1
    clock.time = ClockRate * float32( time.Now().Sub(clock.start).Seconds() )
    
    
    for _,timer := range( clock.timers ) {
        
        if (*timer).Tick( clock.time ) == false {
            
            delete(clock.timers, timer)
            
        }
    }
    
}


func (clock *Clock) Desc() string {
    return fmt.Sprintf("clock[#%05d %.2fs]",clock.frame,clock.time)    
}

func (clock *Clock) Info(prev ClockFrame) string {
    fps := float32(clock.frame - prev.Frame) / (clock.time-prev.Time) 
    return fmt.Sprintf("#%05d %.2fs %.2ffps",clock.frame,clock.time,fps)
}

func (clock *Clock) Frame() ClockFrame { return ClockFrame{Frame: clock.frame, Time: clock.time } }
    










type Clock struct {
    frame uint
    time float32


    start time.Time
    timers map[*Timer]*Timer
}






var worldClock *Clock = &Clock{start: time.Now(), timers: map[*Timer]*Timer{} }



