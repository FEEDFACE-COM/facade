
// +build linux,arm

package render

import (
    "time"   
    "fmt"
    math "../math32"
)


const CLOCK_RATE = 1.0

var startTime time.Time

type Clock struct {
    frame uint
    time  float32
    count uint    // increases monotonously once per second
    fader float32 // ramps from 0..1 once per second
    cycle float32 // cycles through 0..2π once per second
}

func NewClock() *Clock {
    ret := new(Clock)
    ret.Reset()
    return ret
}

func (clock *Clock) Reset() {
    startTime = time.Now()    
    clock.frame = 0
    clock.time = 0.0
}

func InitClock() {
    startTime = time.Now()
}

func (clock *Clock) Time() float32 {
    return clock.time    
}

func (clock *Clock) Tick() {
    prev := clock.time
    clock.frame += 1
    clock.time = CLOCK_RATE *  float32( time.Now().Sub(startTime) ) / (1000.*1000.*1000.)
    
    /***   x%y  ==  mod(x,y)  ==  x - y*floor(x/y)  ***/
    clock.cycle = clock.time - math.TAU * math.Floor(clock.time/math.TAU) // time % TAU
    clock.fader = clock.time - math.Floor(clock.time)                     // time % 1
    
    if math.Floor(clock.time) > math.Floor(prev) {
        clock.count += 1
    }
    
}

func (clock *Clock) Desc() string {
    return fmt.Sprintf("%7.2fs %4.2f↺ %4.2f⤢ %d#",clock.time,clock.cycle,clock.fader,clock.count)
}
