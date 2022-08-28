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

const DEBUG_CHARBUFFER = false
const DEBUG_CHARBUFFER_DUMP = false

const MAX_CHARCOUNT = 128

type CharBuffer struct {
	line Line
	next Line
	rem  []rune

	timer  *gfx.Timer
	speed  float32
	repeat bool

	charCount uint

	refreshChan chan bool
	mutex       *sync.Mutex
}

func NewCharBuffer(refreshChan chan bool) *CharBuffer {

	ret := &CharBuffer{
		charCount:   uint(CharDefaults.CharCount),
		speed:       float32(CharDefaults.Speed),
		repeat:      CharDefaults.Repeat,
		line:        Line{},
		next:        Line{},
		refreshChan: refreshChan,
		mutex:       &sync.Mutex{},
	}
	ret.timer = gfx.WorldClock().NewTimer(0., false, func(x float32) float32 { return 0. }, nil)
	if ret.speed > 0. {
		ret.timer.SetValueFun(func(x float32) float32 { return x })
		ret.timer.SetTriggerFun(func() { ret.scroll() })
		ret.timer.SetDuration(ret.speed)
	}
	if DEBUG_CHARBUFFER {
		log.Debug("%s created", ret.Desc())
	}
	return ret
}

func (buffer *CharBuffer) scroll() {

	buffer.timer.Restart()

	if len(buffer.line) <= 0 {
		return
	}

	if len(buffer.line) == 1 { // scroll out last char, check for refill

		if buffer.Repeat() && len(buffer.next) > 0 { // refill with last line received

			buffer.line = EmptyLine(int(buffer.charCount))
			buffer.line = append(buffer.line, buffer.next...)

		} else { // empty buffer
			buffer.line = Line{}
		}

		//if DEBUG_CHARBUFFER {
		//	log.Debug("%s scroll",buffer.Desc())
		//}

	} else {

		buffer.line = buffer.line[1:]

	}

	buffer.ScheduleRefresh()

}

func (buffer *CharBuffer) Clear() {
	buffer.mutex.Lock()
	buffer.line = Line{}
	buffer.next = Line{}
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

	if len(line) > MAX_CHARCOUNT {
		if DEBUG_CHARMODE {
			log.Debug("%s drop line %d > %d max", buffer.Desc(), len(line), MAX_CHARCOUNT)
		}
		return
	}

	buffer.next = line

	if len(line)+len(buffer.line) > MAX_CHARCOUNT {

		return // too long total line

		//if DEBUG_CHARMODE {
		//	log.Debug("%s hold line %d+%d > %d max", buffer.Desc(), len(line), len(buffer.line), MAX_CHARCOUNT)
		//}
	}

	if uint(len(buffer.line)) > buffer.charCount {

		return // already out of bounds

	}

	//fillup with spaces?
	n := int(buffer.charCount) - len(buffer.line)
	if n > 0 {
		buffer.line = append(buffer.line, EmptyLine(n-1)...)
	}
	buffer.line = append(buffer.line, EmptyLine(1)...)
	buffer.line = append(buffer.line, line...)
	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) GetScroller() float32 {
	if buffer.timer != nil {
		return buffer.timer.Value()
	}
	return float32(0.0)
}

func (buffer *CharBuffer) GetLine() Line {
	n := uint(len(buffer.line))
	if n > buffer.charCount {
		n = buffer.charCount
	}
	return buffer.line[:n]
}

func (buffer *CharBuffer) CharCount() uint {
	return buffer.charCount
}

func (buffer *CharBuffer) Repeat() bool {
	return buffer.repeat
}

func (buffer *CharBuffer) SetRepeat(repeat bool) {
	buffer.repeat = repeat
}

func (buffer *CharBuffer) Fill(fill string) {
	n := utf8.RuneCountInString(fill)
	if n > int(MAX_CHARCOUNT) {
		n = int(MAX_CHARCOUNT)
	}
	buffer.next = Line{}
	buffer.line = make([]rune, n)

	off := 0
	run, s := utf8.DecodeRuneInString(fill[off:])
	for i := 0; i < n && s > 0; i++ {
		buffer.line[i] = run
		off += s
		run, s = utf8.DecodeRuneInString(fill[off:])
	}

	buffer.next = buffer.line

	if DEBUG_CHARBUFFER_DUMP {
		log.Debug("%s filled %d chars:\n%s", buffer.Desc(), n, buffer.Dump())
	} else if DEBUG_CHARBUFFER {
		log.Debug("%s filled %d chars", buffer.Desc(), n)
	}
	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) Desc() string {

	ret := "charbuffer["
	ret += fmt.Sprintf("#%d/%d s%.1f %.1f/%.1f ", len(buffer.line), buffer.charCount, buffer.speed, buffer.timer.Value(), buffer.timer.Duration())
	if buffer.repeat {
		ret += "⟳ "
	}
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
	ret += "0" + fmt.Sprintf(pad, c-1) + "\n"
	for _, run := range buffer.line {
		//if run == ' ' {
		//	run = '.'
		//}
		ret += fmt.Sprintf("%c", run)
	}
	ret += "\n"
	ret += "+" + strings.Repeat(" ", c-2) + "+\n"
	ret += fmt.Sprintf("%s", string(buffer.next)) + "\n"
	return ret
}

func (buffer *CharBuffer) Resize(charCount uint) {
	if charCount <= 0 {
		return
	}

	if charCount > MAX_CHARCOUNT {
		charCount = MAX_CHARCOUNT
	}

	buffer.mutex.Lock()
	buffer.charCount = charCount
	buffer.timer.SetDuration(buffer.speed)
	buffer.mutex.Unlock()

	if DEBUG_CHARBUFFER_DUMP {
		log.Debug("%s resize %d:\n%s", buffer.Desc(), charCount, buffer.Dump())
	} else if DEBUG_CHARBUFFER {
		log.Debug("%s resize %d", buffer.Desc(), charCount)
	}

	buffer.ScheduleRefresh()
}

func (buffer *CharBuffer) SetSpeed(speed float32) {

	buffer.speed = speed // time (seconds) to travel 1 char horizontally

	if buffer.speed > 0. {

		buffer.timer.SetValueFun(func(x float32) float32 { return x })
		buffer.timer.SetTriggerFun(func() { buffer.scroll() })
		buffer.timer.SetDuration(buffer.speed)

	} else {
		buffer.timer.SetValueFun(func(x float32) float32 { return 0. })
		buffer.timer.SetTriggerFun(nil)
	}

	if DEBUG_CHARBUFFER {
		log.Debug("%s set speed %.1f: %s", buffer.Desc(), speed, buffer.timer.Desc())
	}
}
