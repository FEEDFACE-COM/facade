//go:build RENDERER
// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"sync"
	"unicode/utf8"

	"fmt"
	"github.com/pborman/ansi"
	"strings"
)

const DEBUG_CHARBUFFER = true
const DEBUG_CHARBUFFER_DUMP = false

type CharBuffer struct {
	line Line
	last Line
	rem  []rune

	timer *gfx.Timer
	speed float32

	charCount uint

	refreshChan chan bool
	mutex       *sync.Mutex
}

func NewCharBuffer(refreshChan chan bool) *CharBuffer {

	ret := &CharBuffer{
		charCount:   uint(CharDefaults.CharCount),
		speed:       float32(CharDefaults.Speed),
		refreshChan: refreshChan,
		mutex:       &sync.Mutex{},
	}
	ret.line = Line{}
	ret.last = Line{}

	ret.timer = gfx.WorldClock().NewTimer(0., false, nil, nil)

	ret.refreshChan = refreshChan
	if DEBUG_CHARBUFFER {
		log.Debug("%s created", ret.Desc())
	}

	return ret

}

func (buffer *CharBuffer) scroll() {

	if len(buffer.line) <= 0 {
		return
	}

	buffer.timer.Restart()
	buffer.line = buffer.line[1:]

	//if DEBUG_CHARBUFFER_DUMP {
	//	log.Debug("%s scroll: %s\n%s",buffer.Desc(),buffer.timer.Desc(),buffer.Dump())
	//} else if DEBUG_CHARBUFFER {
	//	log.Debug("%s scroll: %s",buffer.Desc(),buffer.timer.Desc())
	//}
	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) Clear() {
	buffer.mutex.Lock()
	buffer.line = Line{}
	buffer.last = Line{}
	buffer.mutex.Unlock()

	if DEBUG_CHARBUFFER {
		log.Debug("%s clear", buffer.Desc())
	}
	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) ScheduleRefresh() {
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *CharBuffer) ProcessSequence(seq *ansi.S) {
	sequence, ok := ansi.Table[seq.Code]
	if !ok {
		return
		//unlcok tho?
	}

	switch sequence {

	case ansi.Table[ansi.ED]:
		buffer.Clear()

	default:
		break

	}
}

func (buffer *CharBuffer) ProcessRunes(runes []rune) {
	rem := append(buffer.rem, runes...)
	tmp := []rune{}

	for _, r := range rem {

		switch r {
		case '\n':
			buffer.addLine(tmp)
			tmp = []rune{}

		default:
			tmp = append(tmp, r)
		}

	}
	buffer.rem = tmp

}

func (buffer *CharBuffer) addLine(line Line) {

}

func (buffer *CharBuffer) GetScroller() float32 {
	if buffer.timer != nil {
		return buffer.timer.Value()
	}
	return float32(0.0)
}

func (buffer *CharBuffer) GetLine() Line {
	ret := buffer.line
	if len(ret) > int(buffer.charCount) {
		ret = ret[:buffer.charCount]
	}
	return ret
}

func (buffer *CharBuffer) CharCount() uint {
	return buffer.charCount
}

func (buffer *CharBuffer) Repeat() bool {
	return false
}

func (buffer *CharBuffer) Fill(fill string) {
	log.Debug("fill %s", fill)
	n := utf8.RuneCountInString(fill)
	if n > int(buffer.charCount) {
		n = int(buffer.charCount)
	}
	buffer.line = make([]rune, n)

	off := 0
	run, s := utf8.DecodeRuneInString(fill[off:])
	for i := 0; i < n && s > 0; i++ {
		buffer.line[i] = run
		off += s
		run, s = utf8.DecodeRuneInString(fill[off:])
	}

	if DEBUG_CHARBUFFER_DUMP {
		log.Debug("%s filled %d chars:\n%s", buffer.Desc(), n, buffer.Dump())
	} else if DEBUG_CHARBUFFER {
		log.Debug("%s filled %d chars", buffer.Desc(), n)
	}
	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) Desc() string {

	ret := "charbuffer["
	ret += fmt.Sprintf("#%d:%d %.1f s%.1f", buffer.charCount, len(buffer.line), buffer.timer.Value(), buffer.timer.Duration())
	ret = strings.TrimSuffix(ret, " ")
	ret += "]"
	return ret
}

func (buffer *CharBuffer) Dump() string {
	ret := ""
	//n := len(buffer.line)
	c := int(buffer.charCount)
	if c <= 1 {
		ret += "+\n0\n" + string(buffer.line) + "\n+\n"
		return ret
	}

	pad := "%" + fmt.Sprintf("%d", c-1) + "d"
	ret += "    0" + fmt.Sprintf(pad, c-1) + "\n"
	ret += "    " + string(buffer.line) + "\n"
	ret += "    +" + strings.Repeat(" ", c-2) + "+\n"
	ret += "lst " + fmt.Sprintf("%s", string(buffer.last)) + "\n"
	ret += "rem " + fmt.Sprintf("%s", string(buffer.rem)) + "\n"
	return ret
}

func (buffer *CharBuffer) Resize(charCount uint) {
	if charCount <= 0 {
		return
	}

	buffer.mutex.Lock()
	buffer.charCount = charCount
	buffer.mutex.Unlock()

	if DEBUG_CHARBUFFER_DUMP {
		log.Debug("%s resize %d:\n%s", buffer.Desc(), charCount, buffer.Dump())
	} else if DEBUG_CHARBUFFER {
		log.Debug("%s resize %d", buffer.Desc(), charCount)
	}

	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) SetSpeed(speed float32) {

	if speed > 0. {
		buffer.timer.SetValueFun(func(x float32) float32 { return x })
		buffer.timer.SetTriggerFun(func() { buffer.scroll() })
		buffer.timer.SetDuration(speed)
	} else {
		buffer.timer.SetValueFun(func(x float32) float32 { return 0. })
		buffer.timer.SetTriggerFun(nil)
	}

	if DEBUG_CHARBUFFER {
		log.Debug("%s set speed %.1f: %s", buffer.Desc(), speed, buffer.timer.Desc())
	}
}
