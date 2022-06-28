// +build RENDERER

package facade

import (
    
	"FEEDFACE.COM/facade/log"
	"sync"

	//	"fmt"
	"github.com/pborman/ansi"
	"strings"    
)

const DEBUG_CHARBUFFER = true
const DEBUG_CHARBUFFER_DUMP = true


const maxCharCount = 32;

type CharBuffer struct {

	line []rune
	charCount uint




	refreshChan chan bool
	mutex       *sync.Mutex

    
}

func (buffer *CharBuffer) Chars() string {
	return "FIXME"
}

func (buffer *CharBuffer) CharCount() uint {
	return uint( len("FIXME") )
}

func NewCharBuffer(refreshChan chan bool) *CharBuffer {

    ret := &CharBuffer{
		charCount: maxCharCount,
		refreshChan: refreshChan,
		mutex:       &sync.Mutex{},
	}
    ret.line = make([]rune, maxCharCount)
    ret.refreshChan = refreshChan


	if DEBUG_CHARBUFFER {
		log.Debug("%s created", ret.Desc())
	}
    
    return ret

}


func (buffer *CharBuffer) Clear() {
	buffer.mutex.Lock()
	buffer.line = make([]rune, buffer.charCount)
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


func (buffer *CharBuffer) Fill(fill string) {

	buffer.Clear()

	idx := 0
	for idx = 0; uint(idx)<buffer.charCount && idx<len(fill); idx++ {
		buffer.line[idx] = rune( fill[idx] )
	}


	if DEBUG_CHARBUFFER_DUMP {
		log.Debug("%s filled %d chars:\n%s", buffer.Desc(), idx, buffer.Dump())
	} else if DEBUG_CHARBUFFER {
		log.Debug("%s filled %d chars", buffer.Desc(),idx)
	}

}


func (buffer *CharBuffer) Desc() string {

	ret := "charbuffer["
	ret = strings.TrimSuffix(ret, " ")
	ret += "]"
	return ret
}

func (buffer *CharBuffer) Dump() string {
	ret := ""
	ret += string(buffer.line)
	ret += "\n"
	return ret
}


