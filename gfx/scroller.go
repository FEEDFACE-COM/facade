
// +build linux,arm

package gfx


import (
    "fmt"
//    log "../log"
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
    
    if scroller.Timer != nil {
        val =   1. - math.EaseInEaseOut( scroller.Timer.Fader() )
    } else {
        val = 0.0;
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



func (scroller *Scroller) Once(fun func()) bool {
    if ! scroller.Scroll {
        return false
    }
    if scroller.Timer != nil {
        return false
    }
    scroller.Timer = NewTimer(math.Abs(scroller.Speed),false )
    scroller.Timer.Start()
    scroller.Timer.Fun = func() {
        UnRegisterTimer(scroller.Timer)
        scroller.Timer = nil
        fun()
//        log.Debug("stop %s",scroller.Desc())

    }    
    return true	
//    log.Debug("start %s",scroller.Desc())
}



func (scroller *Scroller) Desc() string {
    tmp := ""
    if scroller.Scroll { tmp += fmt.Sprintf("%.2f",scroller.Speed) }
    if scroller.Timer != nil { tmp += " " + scroller.Timer.Desc() }
    return fmt.Sprintf("scroll[%s]",tmp)
}

func (scroller *Scroller) SetScrollSpeed(scroll bool, speed float32) {
    scroller.Scroll = scroll
    scroller.Speed = speed
}
