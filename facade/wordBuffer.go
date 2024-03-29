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

const DEBUG_WORDBUFFER = false
const DEBUG_WORDBUFFER_DUMP = false

const maxWordLength = 64 // found experimentally

const FadeDuration = 0.5

type WordState string

const (
	WORD_FADEIN  WordState = "fadein"
	WORD_ALIVE   WordState = "alive"
	WORD_FADEOUT WordState = "fadeout"
)

type Word struct {
	text   string
	index  uint
	length uint
	width  float32
	state  WordState
	timer  *gfx.Timer
	fader  *gfx.Timer
}

type WordBuffer struct {
	words     []*Word
	nextIndex int

	slotCount int
	lifetime  float32
	watermark float32
	shuffle   bool
	aging     bool
	unique    bool

	rem         []rune
	refreshChan chan bool
	mutex       *sync.Mutex
}

func NewWordBuffer(refreshChan chan bool) *WordBuffer {
	ret := &WordBuffer{
		slotCount:   int(WordDefaults.Slots),
		lifetime:    float32(WordDefaults.Lifetime),
		watermark:   float32(WordDefaults.Watermark),
		shuffle:     WordDefaults.Shuffle,
		aging:       WordDefaults.Aging,
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
	sequence, ok := ansi.Table[seq.Code]
	if !ok {
		return
	}

	switch sequence {

	case ansi.Table[ansi.ED]:
		buffer.Clear()

	default:
		break

	}

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

func (buffer *WordBuffer) GetWords() []*Word {
	ret := []*Word{}
	for _, word := range buffer.words {
		if word != nil {
			ret = append(ret, word)
		}
	}
	return ret
}

func (buffer *WordBuffer) checkWatermark() {
	allowed := int(math.Floor(buffer.watermark * float32(buffer.slotCount)))
	alive := []int{}

	for i := 0; i < buffer.slotCount; i++ {
		idx := (buffer.nextIndex + i) % buffer.slotCount
		if buffer.words[idx] != nil && buffer.words[idx].state == WORD_ALIVE {
			alive = append(alive, idx)
		}
	}
	used := len(alive)

	if used > allowed {

		var word *Word = nil

		for i := 0; i < used-allowed; i++ {

			u := len(alive)

			if u <= 0 {
				if DEBUG_WORDBUFFER {
					log.Debug("%s high water mark, no free slots", buffer.Desc())
				}
				break
			}

			if buffer.shuffle {
				r := int(rand.Int31n(int32(u)))
				tmp := alive[0]
				alive[0] = alive[r]
				alive[r] = tmp
			}
			index := alive[0]
			alive = alive[1:] // remove first entry

			word = buffer.words[index]
			if word != nil {

				buffer.fadeoutWord(word)

				if DEBUG_WORDBUFFER {
					log.Debug("%s high water mark, fade out: %s", buffer.Desc(), word.Desc())
				}

			}

		}

	}
}

func (buffer *WordBuffer) addWord(raw []rune) {
	if len(raw) <= 0 {
		//log.Debug("%s not adding empty string", buffer.Desc())
		return
	}

	text := string(raw)
	if len(text) > maxWordLength {
		text = text[0 : maxWordLength-1]
	}

	buffer.mutex.Lock()
	buffer.checkWatermark()

	// check if already in buffer
	if buffer.unique {
		for i := 0; i < buffer.slotCount; i++ {
			idx := (buffer.nextIndex + i) % buffer.slotCount
			if buffer.words[idx] != nil && buffer.words[idx].text == text {
				uniq := buffer.words[idx]
				if DEBUG_WORDBUFFER {
					log.Debug("%s word already in #%d: %s", buffer.Desc(), uniq.index, text)
				}
				buffer.mutex.Unlock()
				return
			}
		}
	}


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
	word.fader = gfx.WorldClock().NewTimer(
		FadeDuration,
		false,
		func(x float32) float32 { return math.EaseIn(x) },
		func() { buffer.fadedinWord(word) },
	)
	if buffer.lifetime > 0.0 {
		word.timer = gfx.WorldClock().NewTimer(
			buffer.lifetime,
			false,
			func(x float32) float32 { return x },
			nil,
		)
	}
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
	gfx.WorldClock().DeleteTimer(word.fader)
	word.fader = nil
	word.state = WORD_ALIVE
	lifetime := buffer.lifetime
	if buffer.watermark <= 0. {

		buffer.fadeoutWord(word)

	} else if lifetime <= 0.0 {

		// nop

	} else if lifetime <= 2.*FadeDuration {

		buffer.fadeoutWord(word)

	} else {
		word.fader = gfx.WorldClock().NewTimer(
			lifetime-2.*FadeDuration,
			false,
			func(x float32) float32 { return 1.0 },
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
	if word.fader != nil {
		val = word.fader.Value()
	}

	gfx.WorldClock().DeleteTimer(word.fader)
	word.fader = nil
	word.state = WORD_FADEOUT
	word.fader = gfx.WorldClock().NewTimer(
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
	//buffer.checkWatermark()
	buffer.mutex.Unlock()
	gfx.WorldClock().DeleteTimer(word.fader)
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
		word.width = float32(utf8.RuneCountInString(word.text))
		word.state = WORD_ALIVE
		if buffer.lifetime > 0.0 {

			word.timer = gfx.WorldClock().NewTimer(
				buffer.lifetime,
				false,
				func(x float32) float32 { return x },
				nil,
			)
			if buffer.lifetime <= FadeDuration {

				buffer.fadeoutWord(word)

			} else {
				word.fader = gfx.WorldClock().NewTimer(
					buffer.lifetime-FadeDuration,
					false,
					func(x float32) float32 { return 1.0 },
					func() { buffer.fadeoutWord(word) },
				)
			}
		}
		buffer.words[idx] = word

	}

	if DEBUG_WORDBUFFER {
		log.Debug("%s filled %d words", buffer.Desc(), buffer.slotCount)
	}

	if buffer.watermark != 0.0 {
		buffer.checkWatermark()
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
	//	{
	//		slots := float32(buffer.slotCount)
	//		count := float32(0.)
	//		for i := 0; i < buffer.slotCount; i++ {
	//			idx := (buffer.nextIndex + i) % buffer.slotCount
	//			if buffer.words[idx] != nil && buffer.words[idx].state == WORD_ALIVE {
	//				count += 1.
	//			}
	//		}
	//		used := count / slots
	//		ret += fmt.Sprintf("%0.1f:%0.1f", used, buffer.watermark)
	//	}
	if buffer.watermark <= 1. {
		ret += fmt.Sprintf("%.1fm ", buffer.watermark)
	}
	if buffer.lifetime > 0. {
		ret += fmt.Sprintf("%.1fl ", buffer.lifetime)
	}

	if buffer.shuffle {
		ret += "⧢"
	}
	if buffer.aging {
		ret += "å"
	}
	if buffer.unique {
		ret += "û"
	}
	ret = strings.TrimRight(ret, "\n")
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
	ret := ""
	max := 32
	min := func(a, b int) int {
		if a <= b {
			return a
		} else {
			return b
		}
	}

	for idx, word := range buffer.words {
		if word != nil {
			ret += fmt.Sprintf("%-2d |", idx)
			ret += fmt.Sprintf(" #%-2d", word.index)
			ret += fmt.Sprintf(" [%.*s]", max, word.text[:min(max, len(word.text))])
			if word.timer != nil {
				ret += " " + word.timer.Desc()
			}
			if word.fader != nil {
				ret += " " + word.fader.Desc("f")
			}
			ret += "\n"
		} else {
			ret += fmt.Sprintf("%-2d |\n", idx)
		}
	}
	return strings.TrimRight(ret, "\n")
}

func (buffer *WordBuffer) WordCount() int {
	ret := 0
	for _, w := range buffer.words {
		if w != nil {
			ret += 1
		}
	}
	return ret
}

func (buffer *WordBuffer) SlotCount() int     { return buffer.slotCount }
func (buffer *WordBuffer) Lifetime() float32  { return buffer.lifetime }
func (buffer *WordBuffer) Watermark() float32 { return buffer.watermark }
func (buffer *WordBuffer) Shuffle() bool      { return buffer.shuffle }
func (buffer *WordBuffer) Aging() bool        { return buffer.aging }
func (buffer *WordBuffer) Unique() bool       { return buffer.unique }
func (buffer *WordBuffer) SetShuffle(shuffle bool) {
	buffer.shuffle = shuffle
}

func (buffer *WordBuffer) SetAging(aging bool) {
	buffer.aging = aging
}

func (buffer *WordBuffer) SetUnique(unique bool) {
	buffer.unique = unique
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
			buffer.words[word.index] = word
		} else {
			gfx.WorldClock().DeleteTimer(word.timer)
			gfx.WorldClock().DeleteTimer(word.fader)
		}
	}

	buffer.mutex.Unlock()

	if DEBUG_WORDBUFFER_DUMP {
		log.Debug("%s resize %d:\n%s", buffer.Desc(), slotCount, buffer.Dump())
	} else if DEBUG_WORDBUFFER {
		log.Debug("%s resize %d", buffer.Desc(), slotCount)
	}

	buffer.ScheduleRefresh()
}
