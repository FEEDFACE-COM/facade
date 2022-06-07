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
	"unicode/utf8"
)

const DEBUG_WORDBUFFER = true
const DEBUG_WORDBUFFER_DUMP = true

const maxWordLength = 64 // found experimentally

const FadeDuration = 0.5

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
	aging     bool

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
		aging:       SetDefaults.Aging,
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

func (buffer *WordBuffer) checkWatermark() {
	allowed := buffer.watermark * float32(buffer.slotCount)
	used := float32(0.)
	for i := 0; i < buffer.slotCount; i++ {
		idx := (buffer.nextIndex + i) % buffer.slotCount
		if buffer.words[idx] != nil && buffer.words[idx].state == WORD_ALIVE {
			used += 1.
		}
	}

	if used >= allowed {

		var word *Word = nil

		//find next alive slot
		alive := []int{}
		for i := 0; i < buffer.slotCount; i++ {
			idx := (buffer.nextIndex + i) % buffer.slotCount
			if buffer.words[idx] != nil && buffer.words[idx].state == WORD_ALIVE {
				alive = append(alive, idx)
			}
		}
		if len(alive) > 0 {
			if buffer.shuffle {
				r := rand.Int31n(int32(len(alive)))
				index := alive[r]
				word = buffer.words[index]
			} else {
				index := alive[0]
				word = buffer.words[index]
			}
		}

		if word != nil {

			buffer.fadeoutWord(word)
			if DEBUG_WORDBUFFER {
				log.Debug("%s watermark exceeded, fade out: %s", buffer.Desc(), word.Desc())
			}

		} else {
			if DEBUG_WORDBUFFER {
				log.Debug("%s watermark exceeded, but no fadeout", buffer.Desc())
			}
		}

	}
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

	buffer.checkWatermark()

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
		log.Debug("%s word add: %s\n%s", buffer.Desc(), word.Desc(), buffer.Dump())

	} else if DEBUG_WORDBUFFER {
		log.Debug("%s word add: %s", buffer.Desc(), word.Desc())
	}
	buffer.mutex.Unlock()

	buffer.ScheduleRefresh()
}

func (buffer *WordBuffer) fadedinWord(word *Word) {
	gfx.WorldClock().DeleteTimer(word.timer)
	word.timer = nil
	word.state = WORD_ALIVE
	lifetime := buffer.lifetime
	if lifetime <= 0.0 {

	} else if lifetime <= 2.*FadeDuration {

		buffer.fadeoutWord(word)

	} else {
		fun := func(x float32) float32 { return 1. }
		if buffer.aging {
			fun = func(x float32) float32 { return 1. - math.EaseInEaseOut(x) }
		}
		word.timer = gfx.WorldClock().NewTimer(
			lifetime-2.*FadeDuration,
			false,
			fun,
			func() { buffer.fadeoutWord(word) },
		)
		if DEBUG_WORDBUFFER_DUMP {
			log.Debug("%s word alife: %s\n%s", buffer.Desc(), word.Desc(), buffer.Dump())
		} else if DEBUG_WORDBUFFER {
			log.Debug("%s word alife: %s", buffer.Desc(), word.Desc())
		}
	}

}

func (buffer *WordBuffer) fadeoutWord(word *Word) {
	var val float32 = 1.

	// get most recent value
	if word.timer != nil {
		val = word.timer.Value()
	}

	gfx.WorldClock().DeleteTimer(word.timer)
	word.timer = nil
	word.state = WORD_FADEOUT
	word.timer = gfx.WorldClock().NewTimer(
		FadeDuration,
		false,
		func(x float32) float32 { return val * (1. - math.EaseIn(x)) },
		func() { buffer.deleteWord(word) },
	)
	if DEBUG_WORDBUFFER_DUMP {
		log.Debug("%s word fade out: %s\n%s", buffer.Desc(), word.Desc(), buffer.Dump())
	} else if DEBUG_WORDBUFFER {
		log.Debug("%s word fade out: %s", buffer.Desc(), word.Desc())
	}
}

func (buffer *WordBuffer) deleteWord(word *Word) {
	buffer.mutex.Lock()
	//old := buffer.words[word.index]
	buffer.words[word.index] = nil
	buffer.checkWatermark()
	buffer.mutex.Unlock()
	gfx.WorldClock().DeleteTimer(word.timer)
	if DEBUG_WORDBUFFER_DUMP {
		log.Debug("%s word delete: %s\n%s", buffer.Desc(), word.Desc(), buffer.Dump())
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

	//get array of indices
	slots := []int{}
	for i := 0; i < buffer.slotCount; i++ {
		idx := (buffer.nextIndex + i) % buffer.slotCount
		slots = append(slots, idx)
	}

	if buffer.shuffle {
		rand.Shuffle(len(slots), func(i, j int) {
			slots[i], slots[j] = slots[j], slots[i]
		})
	}

	for row := 0; uint(row) < rows && row < len(slots); row++ {
		idx := slots[row]
		word := &Word{}
		word.index = uint(idx)
		word.text = fill[row]
		word.width = float32( utf8.RuneCountInString( word.text ) )
		word.state = WORD_ALIVE
		if buffer.lifetime != 0. {
			fun := func(x float32) float32 { return 1. }
			if buffer.aging {
				fun = func(x float32) float32 { return 1. - math.EaseInEaseOut(x) }
			}
			word.timer = gfx.WorldClock().NewTimer(
				buffer.lifetime-1.*FadeDuration,
				false,
				fun,
				func() { buffer.fadeoutWord(word) },
			)
		}
		buffer.words[idx] = word

	}

	//if buffer.watermark != 0.0 {
	//	buffer.checkWatermark()
	//}

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
	{
		slots := float32(buffer.slotCount)
		count := float32(0.)
		for i := 0; i < buffer.slotCount; i++ {
			idx := (buffer.nextIndex + i) % buffer.slotCount
			if buffer.words[idx] != nil && buffer.words[idx].state == WORD_ALIVE {
				count += 1.
			}
		}
		used := count / slots
		ret += fmt.Sprintf("%.1f/%.1f", used, buffer.watermark)
	}
	if buffer.lifetime != 0. {
		ret += fmt.Sprintf(" %.1fs", buffer.lifetime)
	}
	if buffer.shuffle || buffer.aging {
		ret += " "
	}
	if buffer.shuffle {
		ret += "⧢"
	}
	if buffer.aging {
		ret += "a"
	}
	ret += "]"
	return ret
}

func (word *Word) Desc() string {
	ret := ""
	ret += fmt.Sprintf("#%02d ", word.index)
	if word.timer != nil {
		ret += word.timer.Desc() + " "
	}
	ret += word.text
	return ret
}

func (buffer *WordBuffer) Dump() string {
	max := 0
	for _, word := range buffer.words {
		if word != nil {
			if len(word.text) > max {
				max = len(word.text)
			}
		}
	}
	ret := ""
	for _, word := range buffer.words {
		if word != nil {
			ret += fmt.Sprintf("    %2d |  ", word.index)
			f := fmt.Sprintf("%d", max+1)
			ret += fmt.Sprintf("%"+f+"s", word.text)
			if word.timer != nil {
				ret += " " + word.timer.Desc()
			}
			ret += "\n"
		}
	}
	return strings.TrimRight(ret, "\n")
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

func (buffer *WordBuffer) SlotCount() int     { return buffer.slotCount }
func (buffer *WordBuffer) Lifetime() float32  { return buffer.lifetime }
func (buffer *WordBuffer) Watermark() float32 { return buffer.watermark }
func (buffer *WordBuffer) Shuffle() bool      { return buffer.shuffle }
func (buffer *WordBuffer) Aging() bool        { return buffer.aging }

func (buffer *WordBuffer) SetShuffle(shuffle bool) {
	buffer.shuffle = shuffle
}

func (buffer *WordBuffer) SetAging(aging bool) {
	buffer.aging = aging
}

func (buffer *WordBuffer) SetLifetime(lifetime float32) {
	buffer.lifetime = lifetime
	if buffer.words != nil {
		for _, word := range buffer.words {
			if word != nil {
				if word.state == WORD_ALIVE {
					//					word.timer.SetDuration(lifetime)
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
	buffer.slotCount = slotCount
	buffer.nextIndex = buffer.nextIndex % buffer.slotCount

	for _, word := range old {
		if word.index < uint(buffer.slotCount) {
		buffer.words[ word.index ] = word
	} else{
		gfx.WorldClock().DeleteTimer(word.timer)
	}
	}

	buffer.mutex.Unlock()


	if DEBUG_WORDBUFFER {
		log.Debug("%s resize %d", buffer.Desc(), slotCount)
	}

	buffer.ScheduleRefresh()
}
