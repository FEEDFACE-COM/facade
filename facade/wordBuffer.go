//go:build RENDERER
// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	math "FEEDFACE.COM/facade/math32"
	"fmt"
	"github.com/pborman/ansi"
	"math/rand"
	"strings"
	"sync"
)

const DEBUG_WORDBUFFER = true
const DEBUG_WORDBUFFER_DUMP = true

const maxWordLength = 80 // found experimentally

const FadeDuration = 0.5;

type WordState string
const (
	WORD_FADEIN  WordState = "fadein"
	WORD_ALIVE   WordState = "alive"
	WORD_FADEOUT WordState = "fadeout"
)

type Word struct {
	text  string
	count uint
	index uint
	width float32
	state WordState
	timer *gfx.Timer
}

type WordBuffer struct {
	words     []*Word
	nextIndex int

	slotCount int
	lifetime  float32
	watermark float32
	shuffle   bool

	rem         []rune
	refreshChan chan bool
	mutex       *sync.Mutex
}



func NewWordBuffer(refreshChan chan bool) *WordBuffer {
	ret := &WordBuffer{
		slotCount:   int(SetDefaults.Slots),
		lifetime:    float32(SetDefaults.Lifetime),
		watermark:   float32(SetDefaults.Watermark),
		shuffle:     SetDefaults.Shuffle,
		refreshChan: refreshChan,
		mutex:       &sync.Mutex{},
	}
	ret.words = make([]*Word, ret.slotCount)
	return ret
}

func (buffer *WordBuffer) ScheduleRefresh() {
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *WordBuffer) ProcessSequence(seq *ansi.S) {
	return
}

func (buffer *WordBuffer) ProcessRunes(runes []rune) {

	rem := append(buffer.rem, runes...)
	tmp := []rune{}

	for _, r := range rem {

		switch r {
		case '\n':
			buffer.addWord(tmp)
			tmp = []rune{}

		default:
			tmp = append(tmp, r)
		}

	}
	buffer.rem = tmp
}

func (buffer *WordBuffer) Words() []*Word {
	ret := []*Word{}
	buffer.mutex.Lock()
	for _, w := range buffer.words {
		if w != nil {
			ret = append(ret, w)
		}
	}
	buffer.mutex.Unlock()
	if len(ret) != buffer.WordCount() {
		log.Warning("mismatch buffer tags: expected %d got %d", buffer.WordCount(), len(ret))
	}
	return ret
}

func (buffer *WordBuffer) addWord(raw []rune) {
	if len(raw) <= 0 {
		log.Debug("%s not adding empty string", buffer.Desc())
		return
	}

	text := string(raw)
	if len(text) > maxWordLength {
		text = text[0 : maxWordLength-1]
	}

	buffer.mutex.Lock()
	var index int = -1

	//find next empty slot
	empty := []int{}
	for i := 0; i < buffer.slotCount; i++ {
		idx := (buffer.nextIndex + i) % buffer.slotCount
		if buffer.words[idx] == nil {
			empty = append(empty, idx)
		}
	}
	if len(empty) > 0 {
		if buffer.shuffle {
			r := rand.Int31n(int32(len(empty)))
			index = empty[r]
		} else {
			index = empty[0]
		}
	}

	if index < 0 {
		if DEBUG_WORDBUFFER {
			log.Debug("%s word drop: %s", buffer.Desc(), text)
		}
		buffer.mutex.Unlock()
		return
	}

	word := &Word{}
	word.index = uint(index)
	word.text = text
	word.state = WORD_FADEIN
	word.timer = gfx.WorldClock().NewTimer(
		FadeDuration,
		false,
		func(x float32) float32 { return math.EaseIn(x) },
		func() { buffer.fadedinWord(word) },
	)

	buffer.words[index] = word
	buffer.nextIndex = (index + 1) % buffer.slotCount
	if DEBUG_WORDBUFFER_DUMP {
		log.Debug("%s word add: %s\n%s", buffer.Desc(), word.Desc(),buffer.Dump())

	} else if DEBUG_WORDBUFFER {
		log.Debug("%s word add: %s", buffer.Desc(), word.Desc())
	}
	buffer.mutex.Unlock()
	buffer.ScheduleRefresh()
}



func (buffer *WordBuffer) fadedinWord(word *Word) {
	word.state = WORD_ALIVE
	word.timer = nil
	lifetime := buffer.lifetime
	if lifetime <= 0.0 {


	} else if lifetime <= 2. * FadeDuration {

		buffer.fadeoutWord(word)

	} else {
		word.timer = gfx.WorldClock().NewTimer(
			lifetime - 2. * FadeDuration,
			false,
			func(x float32) float32 { return 1. },
			func() { buffer.fadeoutWord(word) },
		)
		log.Debug("%s word alife: %s",buffer.Desc(),word.Desc())
	}

}

func (buffer *WordBuffer) fadeoutWord(word *Word) {
	word.timer = nil
	word.state = WORD_FADEOUT
	word.timer = gfx.WorldClock().NewTimer(
		FadeDuration,
		false,
		func(x float32) float32 { return 1.- math.EaseOut(x) },
		func() { buffer.deleteWord(word) },
	)
	log.Debug("%s word fade out: %s",buffer.Desc(),word.Desc())
}


func (buffer *WordBuffer) deleteWord(word *Word) {
	buffer.mutex.Lock()
	old := buffer.words[word.index]
	buffer.words[word.index] = nil
	buffer.mutex.Unlock()
	if old != nil {
		old.timer = nil
		old.state = WORD_ALIVE
	}
	if DEBUG_WORDBUFFER_DUMP {
		log.Debug("%s word delete: %s\n%s", buffer.Desc(), word.Desc(),buffer.Dump())
	} else if DEBUG_WORDBUFFER {
		log.Debug("%s word delete: %s", buffer.Desc(), word.Desc())
	}

	buffer.ScheduleRefresh()
}

func (buffer *WordBuffer) Clear() {
	old := []*gfx.Timer{}

	buffer.mutex.Lock()
	for _, word := range buffer.words {
		if word != nil {
			old = append(old, word.timer)
		}
	}
	buffer.words = make([]*Word, buffer.slotCount)
	buffer.nextIndex = 0
	buffer.mutex.Unlock()

	for _, timer := range old {
		gfx.WorldClock().DeleteTimer(timer)
	}

	if DEBUG_WORDBUFFER {
		log.Debug("%s clear", buffer.Desc())
	}
	buffer.ScheduleRefresh()
}

func (buffer *WordBuffer) Fill(fill []string) {

	buffer.Clear()

	rows := uint(len(fill))

	for r := uint(0); r < rows; r++ {

		buffer.addWord([]rune(fill[r]))

	}

}

func (buffer *WordBuffer) Desc() string {
	ret := ""
	ret += "wordbuffer["
	if buffer.words != nil {
		c := 0
		for _, w := range buffer.words {
			if w != nil {
				c += 1
			}
		}
		ret += fmt.Sprintf("%d/%d ", c, buffer.slotCount)
	}
	ret += fmt.Sprintf("%.1fl ", buffer.lifetime)
	ret += fmt.Sprintf("%.1fm ", buffer.watermark)
	if buffer.shuffle {
		ret += "⧢"
	}
	ret += "]"
	return ret
}

func (word *Word) Desc() string {
	ret := ""
	ret += fmt.Sprintf("#%2d ",word.index)
	if word.timer != nil {
		ret += word.timer.Desc() + " "
	}
	ret += word.text
	return ret
}

func (buffer *WordBuffer) Dump() string {
	ret := ""
	if buffer.words != nil {
		for _, word := range buffer.words {
			if word != nil {
				ret += fmt.Sprintf("    %2d |  ", word.index)
				if word.timer != nil {
					ret += word.timer.Desc()
				}
				ret += fmt.Sprintf(" %s\n", word.text)
			}
		}
	}
	return strings.TrimRight(ret,"\n")
}

func (buffer *WordBuffer) WordCount() int {
	r := 0
	for _, w := range buffer.words {
		if w != nil {
			r += 1
		}
	}
	return r
}

func (buffer *WordBuffer) SlotCount() int    { return buffer.slotCount }
func (buffer *WordBuffer) Lifetime() float32 { return buffer.lifetime }
func (buffer *WordBuffer) Watermark() float32 { return buffer.watermark }
func (buffer *WordBuffer) Shuffle() bool     { return buffer.shuffle }

func (buffer *WordBuffer) SetShuffle(shuffle bool) {
	buffer.shuffle = shuffle
}

func (buffer *WordBuffer) SetLifetime(lifetime float32) {
	buffer.lifetime = lifetime
	if buffer.words != nil {
		for _, word := range buffer.words {
			if word != nil {
				if word.state == WORD_ALIVE {
					word.timer.SetDuration(lifetime)
				}
			}
		}
	}
}

func (buffer *WordBuffer) SetWatermark(watermark float32) {
	buffer.watermark = watermark
}

func (buffer *WordBuffer) Resize(slotCount int) {
	if slotCount <= 0 {
		return
	}

	old := []*Word{}


	buffer.mutex.Lock()
	for _, w := range buffer.words {
		if w != nil {
			old = append(old, w)
		}
	}
	buffer.words = make([]*Word, slotCount)
	buffer.mutex.Unlock()

	buffer.mutex.Lock()
	buffer.slotCount = slotCount
	buffer.nextIndex = buffer.nextIndex % buffer.slotCount
	buffer.mutex.Unlock()

	for _, word := range old {
		gfx.WorldClock().DeleteTimer(word.timer)
	}

	if DEBUG_WORDBUFFER {
		log.Debug("%s resize %d", buffer.Desc(), slotCount)
	}

	buffer.ScheduleRefresh()
}
