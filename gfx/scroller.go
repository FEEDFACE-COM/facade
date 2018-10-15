
// +build linux,arm

package gfx


import (
    "fmt"
    math "../math32"
)


type Scroller struct {
    
    Scroll bool
    Speed float32
    
    Timer *Timer
}

func (scroller *Scroller) Uniform(program *Program, downward bool) {
    var val float32
//    if (downward) { val = -0.0 }
    
    if scroller.Scroll {
        val = scroller.Timer.Fader
    } else {
        val = 1.0;
    }

    if (downward) {
        val *= -1.
    }
    program.Uniform1f(SCROLLER, val)
}

func NewScroller(scroll bool,speed float32) *Scroller {
    ret := &Scroller{Scroll: scroll, Speed: speed}
    return ret
}



func (scroller *Scroller) Init() {
    scroller.Timer = NewTimer(math.Abs(scroller.Speed) )
}



func (scroller *Scroller) Desc() string {
    tmp := ""
    if scroller.Scroll { tmp += fmt.Sprintf("%.2f",scroller.Speed) }
    if scroller.Timer != nil { tmp += " " + scroller.Timer.Desc() }
    return fmt.Sprintf("scroller[%s]",tmp)
}

func (scroller *Scroller) SetScrollSpeed(scroll bool, speed float32) {
    scroller.Scroll = scroll
    scroller.Speed = speed
    scroller.Timer.duration=speed
}
