//go:build RENDERER
// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/log"
	"sync"
	"unicode/utf8"

	"fmt"
	"github.com/pborman/ansi"
	"strings"
)

const DEBUG_CHARBUFFER = true
const DEBUG_CHARBUFFER_DUMP = true

type CharBuffer struct {
	line      Line
	charCount uint

	refreshChan chan bool
	mutex       *sync.Mutex
}

func NewCharBuffer(refreshChan chan bool) *CharBuffer {

	ret := &CharBuffer{
		charCount:   uint(CharDefaults.CharCount),
		refreshChan: refreshChan,
		mutex:       &sync.Mutex{},
	}
	ret.line = make(Line, 1)
	ret.line[0] = ' '
	ret.refreshChan = refreshChan

	if DEBUG_CHARBUFFER {
		log.Debug("%s created", ret.Desc())
	}

	return ret

}

func (buffer *CharBuffer) Clear() {
	buffer.mutex.Lock()
	buffer.line = make([]rune, 1)
	buffer.line[0] = ' '
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
}

func (buffer *CharBuffer) ProcessRunes(runes []rune) {
}

func (buffer *CharBuffer) GetLine() Line {
	return buffer.line
}

func (buffer *CharBuffer) CharCount() uint {
	return buffer.charCount
}

func (buffer *CharBuffer) Repeat() bool {
	return false
}

func (buffer *CharBuffer) Fill(fill string) {
	log.Debug("fill " + fill)
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
	ret += fmt.Sprintf("#%d", len(buffer.line))
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
	ret += "+" + strings.Repeat(" ", c-2) + "+\n"
	ret += string(buffer.line) + "\n"
	ret += "+" + strings.Repeat(" ", c-2) + "+\n"
	ret += "\n"
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
