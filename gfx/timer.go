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

	bias   float32
}


func (timer *Timer) SetDuration(duration float32) {
    timer.duration = duration
}

func (timer *Timer) Count() uint    { return timer.count }
func (timer *Timer) Fader() float32 { return timer.fader }

func (timer *Timer) Remaining(now float32) float32 { 
    return timer.duration - (now - timer.start)
}

func (timer *Timer) Value() float32 {
	if timer.valueFun != nil {
		return timer.valueFun(timer.fader)
	}
	return timer.fader

}

//func (timer *Timer) Restart(now float32) {
//    timer.start = now
//    timer.amp = 1.
//    timer.bias = 0.
//}

func (timer *Timer) Extend(now float32) bool {
    
    
    const MAX = float32( 0.95 )
    
    if timer.fade(now) >= MAX {
        timer.bias = 0.
        timer.start = now
        return false

    } 

    timer.bias = timer.fade(now)
    timer.start = now
    
    return true
        
}
    


func (timer *Timer) fade(now float32) float32 {
    t := now - timer.start
    d := timer.duration
    b := timer.bias

    return b + (1.-b) * (t/d) 
    
}

func (timer *Timer) Tick(now float32) (bool, func()) {
    
    timer.fader = math.Clamp( timer.fade(now) )

	//triggered?
	if now > timer.start+timer.duration {

		timer.count += 1

		if DEBUG_CLOCK {
			log.Debug("%s trigger", timer.Desc())
		}

		if timer.repeat { //keep triggered repeating timer
			timer.start = now
			return true, timer.triggerFun
		}

		return false, timer.triggerFun // remove triggered single time
	}
	return true, nil // keep untriggered timer
}

func (timer *Timer) Desc() string {
    run := Now() - timer.start
	ret := fmt.Sprintf("t[%.1f/%.1f", run,timer.duration)
	if timer.repeat {
		ret += fmt.Sprintf(" #%d", timer.count)
    }
	if timer.bias == 0.0 {
    	ret += fmt.Sprintf(" →%3.1f    ", timer.fader)
    } else {
        ret += fmt.Sprintf(" →%3.1f+%3.1f",(1.-timer.bias)*timer.fader,timer.bias)    
    }
	if timer.valueFun != nil {
		ret += fmt.Sprintf(" ↑%4.2f", timer.valueFun(timer.fader))
	}
	ret += "]"
	return ret
}
