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

	//    if timer.valueFun == nil {
	//        timer.valueFun = func(x float32) float32 { return x }
	//    }

	clock.mux.Lock()
	clock.timers[timer] = timer
	clock.mux.Unlock()
	if DEBUG_CLOCK {
		log.Debug("%s add %s", clock.Desc(),timer.Desc())
	}

	return timer

}

func (clock *Clock) DeleteTimer(timer *Timer) {
	clock.mux.Lock()

	tmp, ok := clock.timers[timer]

	if ok {
		if DEBUG_CLOCK {
			log.Debug("%s delete %s", clock.Desc(), tmp.Desc())
		}
		delete(clock.timers, timer)
	} else {
		log.Error("%s fail delete timer", clock.Desc())
	}

	clock.mux.Unlock()
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

func (clock *Clock) Tick() {

	clock.frame += 1
	clock.time = ClockRate * float32(time.Now().Sub(clock.start).Seconds())

	for _, timer := range clock.timers {

		if (*timer).Tick(clock.time) == false {

			clock.mux.Lock()
			delete(clock.timers, timer)
			clock.mux.Unlock()

		}
	}
}

func (clock *Clock) Now() float32 { return clock.time }

func (clock *Clock) Desc() string {
	return fmt.Sprintf("clock[#%05d %.2fs]", clock.frame, clock.time)
}

func (clock *Clock) Info(prev ClockFrame) string {
	fps := float32(clock.frame-prev.Frame) / (clock.time - prev.Time)
	return fmt.Sprintf("#%05d %.2fs %.2ffps", clock.frame, clock.time, fps)
}

func (clock *Clock) Frame() ClockFrame { return ClockFrame{Frame: clock.frame, Time: clock.time} }

type Clock struct {
	frame uint
	time  float32

	start  time.Time
	timers map[*Timer]*Timer
	mux    sync.Mutex
}

var worldClock *Clock = &Clock{start: time.Now(), timers: map[*Timer]*Timer{}}
