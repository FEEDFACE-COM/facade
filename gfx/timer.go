package gfx

import (
	log "../log"
	math "../math32"
	"fmt"
)

type Timer struct {
	count      uint
	fader      float32
	valueFun   func(float32) float32
	triggerFun func()

	start    float32
	duration float32
	repeat   bool
}

func (timer *Timer) Count() uint    { return timer.count }
func (timer *Timer) Fader() float32 { return timer.fader }

func (timer *Timer) Value() float32 {
	if timer.valueFun != nil {
		return timer.valueFun(timer.fader)
	}
	return timer.fader

}

func (timer *Timer) Restart(now float32) {
    timer.start = now
}

//
//func NewTimer(duration float32, repeat bool, fun func(), custom func(float32) float32 ) *Timer {
//    ret := &Timer{start: NOW(), duration: duration, repeat: repeat}
//    ret.custom = func(x float32) float32 { return x }
//    if custom != nil {
//        ret.custom = custom
//    }
//    ret.Fun = fun
//    RegisterTimer(ret)
//    return ret
//}
//

//func (timer *Timer) Close() {
//    //todo: dereg
//}

//func (timer *Timer) Reset() {
//    timer.start = NOW()
//    timer.fader = 0.0
//    timer.count = 0
//}
//
//func (timer *Timer) Start() {
//    timer.start = NOW()
//    timer.fader = 0.0
//    timer.count = 0
//    if DEBUG_CLOCK { log.Debug("%s start",timer.Desc()) }
//}

func (timer *Timer) Tick(now float32) bool {
	t := now - timer.start
	d := timer.duration

	timer.fader = math.Clamp(t / d)

	//triggered?
	if now > timer.start+timer.duration {

		timer.count += 1

		if DEBUG_CLOCK {
			log.Debug("%s trigger", timer.Desc())
		}

		if timer.triggerFun != nil {
			timer.triggerFun()
		}

		if timer.repeat { //keep triggered repeating timer
			timer.start = now
			return true
		}

		return false // remove triggered single time
	}
	return true // keep untriggered timer
}

func (timer *Timer) Desc() string {
	ret := fmt.Sprintf("timer[%.2f", timer.duration)
	if timer.repeat {
		ret += fmt.Sprintf(" #%d", timer.count)
	}
	ret += fmt.Sprintf(" →%4.2f", timer.fader)
	if timer.valueFun != nil {
		ret += fmt.Sprintf(" ↑%4.2f", timer.valueFun(timer.fader))
	}
	ret += "]"
	return ret
}
