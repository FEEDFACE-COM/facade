package gfx

import (
	log "../log"
	"fmt"
	"sync"
	"time"
)

const DEBUG_CLOCK = false

const ClockRate = 1.0
const VerboseFrames = 60

type ClockFrame struct {
	Frame uint
	Time  float32
}


func Now() float32       { return worldClock.time }
func WorldClock() *Clock { return worldClock }


func (clock *Clock) NewTimer(duration float32, repeat bool, valueFun func(float32) float32, triggerFun func()) *Timer {

	timer := &Timer{
		start:      clock.time,
		duration:   duration,
		repeat:     repeat,
		valueFun:   valueFun,
		triggerFun: triggerFun,
	}

	clock.mux.Lock()
	clock.timers[timer] = timer
	clock.mux.Unlock()

	if DEBUG_CLOCK {
		log.Debug("%s add %p:%s", clock.Desc(),timer,timer.Desc())
	}

	return timer

}

func (clock *Clock) DeleteTimer(timer *Timer) {
    if timer == nil {
        return
    }

	clock.mux.Lock()

	tmp, ok := clock.timers[timer]

	if ok {
		if DEBUG_CLOCK {
			log.Debug("%s delete %p:%s", clock.Desc(), timer, tmp.Desc())
		}
		delete(clock.timers, timer)
	} 
	clock.mux.Unlock()
}

func (clock *Clock) Toggle() {
    clock.paused = !clock.paused 
    clock.pausetime = clock.running() - clock.pausetime
    if DEBUG_CLOCK {
        s := "unpaused"
        if clock.paused { s = "paused" }
        log.Debug("%s %s",clock.Desc(),s)
    }
}



func (clock *Clock) Reset() {
	clock.mux.Lock()

	if DEBUG_CLOCK {
		log.Debug("%s reset", clock.Desc())
	}

	clock.start = time.Now()
	clock.frame = 0
	clock.time = 0.0

	for k, _ := range clock.timers {
		delete(clock.timers, k)
	}
	clock.mux.Unlock()
}

func (clock *Clock) VerboseFrame() bool {
	return clock.frame%VerboseFrames == 0
}

func (clock *Clock) running() float32 {
    return ClockRate * float32(time.Now().Sub(clock.start).Seconds())
}

func (clock *Clock) Tick() {

    clock.mux.Lock()

    running := clock.running()
    const HOUR = 60. * 60. 
    if running > 2. * HOUR {
        if DEBUG_CLOCK {
            log.Debug("%s rewind", clock.Desc())
        }
        add,err := time.ParseDuration( fmt.Sprintf("%.0fs",HOUR) )
        if err != nil {
            log.PANIC("fail to parse duration '%.0fs': %s",HOUR,err)
        }
        clock.hours += 1
        clock.start = clock.start.Add(add)
    	for _, timer := range clock.timers {
        	timer.start -= HOUR 
        }
        running = clock.running()
    }

    if clock.paused {
        clock.frame += 1
        clock.time = clock.pausetime
    } else {
        clock.frame += 1
        clock.time = running - clock.pausetime 
    }	   

    triggers := []func(){}
	for _, timer := range clock.timers {

        keep, fun := (*timer).Tick(clock.time)    
		if keep == false {
    		delete(clock.timers,timer)
		}
		if fun != nil { //keep note of triggers
    		triggers = append(triggers, fun)
        }
	}
	clock.mux.Unlock()

    // now that we're unlocked, run triggers
    for _,fun := range triggers {
        fun()
    }
}

func (clock *Clock) Now() float32 { return clock.time }
func (clock *Clock) Paused() bool { return clock.paused }

func (clock *Clock) Desc() string {
    s := ""
    if clock.paused { s = " PAUSED" }
	h := fmt.Sprintf("%.2fs",clock.time)
	if clock.hours > 0 {
    	h = fmt.Sprintf("%dh %.0fs",clock.hours,clock.time)
    }
	return fmt.Sprintf("clock[#%05d %s%s]", clock.frame,h,s)
}

func (clock *Clock) Info(prev ClockFrame) string {
    s := ""
    if clock.paused { s = " PAUSED" }
	fps := float32(clock.frame-prev.Frame) / (clock.time - prev.Time)
	h := fmt.Sprintf("%.2fs",clock.time)
	if clock.hours > 0 {
    	h = fmt.Sprintf("%dh %.0fs",clock.hours,clock.time)
    }
	return fmt.Sprintf("#%05d %s %.2ffps%s", clock.frame,h,fps,s)
}

func (clock *Clock) Frame() ClockFrame { return ClockFrame{Frame: clock.frame, Time: clock.time} }

type Clock struct {
	frame uint
	time  float32
	hours uint
	
	paused bool
	pausetime float32

	start  time.Time
	timers map[*Timer]*Timer
	mux    sync.Mutex
}

var worldClock *Clock = &Clock{start: time.Now(), timers: map[*Timer]*Timer{}}
