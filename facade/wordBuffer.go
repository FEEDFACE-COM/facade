
package facade

import(
    "fmt"
    "sync"
    "sort"
    "math/rand"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_WORDBUFFER = false

const maxWordLength = 80 // found experimentally


type Word struct {
    text string
    count uint
    index uint
    timer *gfx.Timer
}


type WordBuffer struct {
    
    tags  map[string] *Word
    words []*Word
    nextIndex int
    
    slotCount int
    duration float32
    shuffle bool

    rem []rune
	refreshChan chan bool
	mutex *sync.Mutex
    
}

func newWordBuffer(tag bool, refreshChan chan bool) *WordBuffer {
    return &WordBuffer{
        duration: float32(SetDefaults.Duration),
        slotCount: int(SetDefaults.Slot),
        refreshChan: refreshChan,
        mutex: &sync.Mutex{},
    }
}

func NewTagBuffer(refreshChan chan bool) *WordBuffer {
    ret := newWordBuffer(true, refreshChan)
    ret.tags = make( map [string]*Word ) 
    return ret
    
}

func NewWordBuffer(refreshChan chan bool) *WordBuffer {
    ret := newWordBuffer(false, refreshChan)
    ret.words = make( []*Word, ret.slotCount )
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

    rem := append(buffer.rem, runes ... )
    tmp := []rune{}

    for _,r := range(rem) {

        switch r {
            case '\n':
                buffer.addWord( tmp )
                tmp = []rune{}

            default:
                tmp = append(tmp, r)
        }

    }
    buffer.rem = tmp
}


func (buffer *WordBuffer) Words() []Word {
    ret := []Word{}
    buffer.mutex.Lock()
    if buffer.tags != nil {
        for _,w := range(buffer.tags) {
            if w != nil {
                ret = append(ret, *w)
            }
        }
    }
    if buffer.words != nil {
        for _,w := range(buffer.words) {
            if w != nil {
                ret = append(ret, *w)
            }
        }
    }
    buffer.mutex.Unlock()    
    if len(ret) != buffer.WordCount() {
        log.Debug("mismatch buffer tags: expected %d got %d",buffer.WordCount(),len(ret))    
    }
    return ret    
} 
    


func (buffer *WordBuffer) addWordWord(text string) {
    buffer.mutex.Lock()
    var index int = -1
    //find index
//    if buffer.words[buffer.nextIndex] == nil {
//        index = buffer.nextIndex
//    } else {


    //find next empty slot
    empty := []int{}
    for i:=0; i < buffer.slotCount; i++ {
        idx := (buffer.nextIndex+i)%buffer.slotCount
        if buffer.words[idx] == nil {
            empty = append(empty, idx)
        }
    }
    if len(empty) > 0 {
        if buffer.shuffle {
            r := rand.Int31n( int32(len(empty)))
            index = empty[r]
        } else {
            index = empty[0]
        }
    }

//    }
    
    if index >= 0 {
        word := &Word{}
        word.index = uint(index)
        word.text = text
        triggerFun := func() {
            buffer.deleteWord(*word)
        }
        word.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.words[index] = word
        buffer.nextIndex = (index+1) % buffer.slotCount
        if DEBUG_WORDBUFFER {
            log.Debug("%s word add #%d: %s",buffer.Desc(),index,text)
        }
    } else if DEBUG_WORDBUFFER {
        log.Debug("%s word drop: %s",buffer.Desc(),text)
    }
    buffer.mutex.Unlock()
}

func (buffer *WordBuffer) addWordTag(text string) {

    if len(text) <= 0 {
        return    
    }

    buffer.mutex.Lock()
    word, ok := buffer.tags[text]
    if ok {
        word.count += 1
        r := word.timer.Extend( gfx.Now() )
        if DEBUG_WORDBUFFER {
            if r {
                log.Debug("%s tag extend: %s",word.timer.Desc(),word.text)
            }
        }

    } else if len(buffer.tags) >= buffer.slotCount {

        log.Debug("%s tag drop: %s",buffer.Desc(),text)
            
    } else {
        word = &Word{}
        word.text = text
        word.count = 1

        triggerFun := func() {
            buffer.deleteWord(*word)
        }
        
        
        { // find random free index
            idxmap := make( map [int]bool )
            for i:=0; i<buffer.slotCount; i++ {
                idxmap[i] = true
            }
            for _,word := range(buffer.tags) {
                idx := word.index
                delete(idxmap,int(idx))
            }
            indices := []uint{}
            for i,_ := range(idxmap) {
                indices = append(indices,uint(i))
            }
            r := rand.Int31n( int32(len(indices)) )
            word.index = indices[r]
        }        
        word.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.tags[text] = word
        if DEBUG_WORDBUFFER {
            log.Debug("%s tag add: %s",buffer.Desc(),text)
        }
    }
	buffer.mutex.Unlock()
}


func (buffer *WordBuffer) addWord(raw []rune) {
    
//    if len(raw) <= 0 {
//        log.Debug("%s not adding empty string",buffer.Desc())
//        return
//    }

    text := string(raw)
    if len(text) > maxWordLength {
        text = text[0: maxWordLength-1]
    }

    if buffer.tags != nil {
        buffer.addWordTag(text)
    }
    if buffer.words != nil {
        buffer.addWordWord(text)
    }

    buffer.ScheduleRefresh()
}



func (buffer *WordBuffer) deleteWord(word Word) {
    if buffer.words != nil {
        buffer.mutex.Lock()
        old := buffer.words[word.index]
        buffer.words[word.index] = nil
        buffer.mutex.Unlock()
        if old != nil {
            gfx.WorldClock().DeleteTimer(old.timer)
        }
        if DEBUG_WORDBUFFER {
            log.Debug("%s word delete: #%d %s",buffer.Desc(),word.index,word.text)
        }
    }
    if buffer.tags != nil {
        buffer.mutex.Lock()
        old := buffer.tags[word.text]
        delete(buffer.tags,word.text)
        buffer.mutex.Unlock()
        if old != nil {
            gfx.WorldClock().DeleteTimer(old.timer)
        }
        if DEBUG_WORDBUFFER {
            log.Debug("%s tag delete: %s",buffer.Desc(),word.text)
        }
    }
    buffer.ScheduleRefresh()
}


func (buffer *WordBuffer) Clear() {
    old := []*gfx.Timer{}
    
    buffer.mutex.Lock()
    if buffer.tags != nil {
        for _,word := range buffer.tags {
            if word != nil {
                old = append(old,word.timer)
            }
        }
        buffer.tags = make( map[string] *Word, buffer.slotCount )
    } 
    if buffer.words != nil {
        for _,word := range buffer.words {
            if word != nil {
                old = append(old,word.timer)
            }
        }    
        buffer.words = make([]*Word, buffer.slotCount)
    }
    buffer.nextIndex = 0
	buffer.mutex.Unlock()
	
	for _,timer := range old {
    	gfx.WorldClock().DeleteTimer(timer)
    }

    if DEBUG_WORDBUFFER {
        log.Debug("%s clear",buffer.Desc())
    }
    buffer.ScheduleRefresh()
}


func (buffer *WordBuffer) Fill(fill []string) {

    buffer.Clear()

    rows := uint(len(fill))
    
    for r:=uint(0); r<rows; r++  {
        
        buffer.addWord( []rune( fill[r] ) )
            
    }
    
}

func (buffer *WordBuffer) Desc() string {
    ret := ""
    if buffer.tags != nil {
        ret += "tag"
    }
    if buffer.words != nil {
        ret += "word"
    }
    ret += "buffer["
    if buffer.tags != nil {
        ret += fmt.Sprintf("%d/%d ",len(buffer.tags),buffer.slotCount)
    }
    if buffer.words != nil {
        c := 0
        for _,w := range( buffer.words ) {
            if w != nil {
                c += 1
            }
        }
        ret += fmt.Sprintf("%d/%d ",c,buffer.slotCount)
    }
    ret += fmt.Sprintf("%.1f",buffer.duration)
    if buffer.shuffle {
        ret += "⧢"
    }
    ret += "]"
    return ret
}

func (buffer *WordBuffer) Dump() string {
    ret := ""
    if buffer.words != nil {
        for _,word := range buffer.words {
            if word != nil {
                ret += fmt.Sprintf("    %2d |  ",word.index)
                if word.timer != nil {
                    ret += word.timer.Desc()
                }
                ret += fmt.Sprintf(" %s\n",word.text)
            }
        }
    }
    if buffer.tags != nil {
        indices := []int{}
        words := make( map [int]*Word )
        for _,word := range buffer.tags {
            indices = append(indices,int(word.index))
            words[int(word.index)] = word
        }
        sort.Ints(indices)
        for _,idx := range indices {
            word := words[idx]
            if word != nil {
                ret += fmt.Sprintf("    %2d[#%2d] ",word.index,word.count)        
                if word.timer != nil {
                    ret += word.timer.Desc()
                }
                ret += fmt.Sprintf(" %s\n",word.text)
            }
        }
    }
    return ret
}


func (buffer *WordBuffer) WordCount() int     { 
    if buffer.tags != nil {
        return len(buffer.tags) 
    }
    if buffer.words != nil {
        r := 0
        for _,w := range(buffer.words) {
            if w != nil {
                r += 1
            }
        }
        return r    
    }
    return 0
}

func (buffer *WordBuffer) SlotCount() int    { return buffer.slotCount }
func (buffer *WordBuffer) Duration() float32 { return buffer.duration }
func (buffer *WordBuffer) Shuffle() bool { return buffer.shuffle }

func (buffer *WordBuffer) SetShuffle(shuffle bool) {
    buffer.shuffle = shuffle
}


func (buffer *WordBuffer) SetDuration(duration float32) {
	buffer.duration = duration
	if buffer.tags != nil {
        for _,word := range(buffer.tags) {
            if word != nil {
                word.timer.SetDuration(duration)
            }
        }
    }
    if buffer.words != nil {
        for _,word := range(buffer.words) {
            if word != nil {
                word.timer.SetDuration(duration)
            }
        }
    }
}




func (buffer *WordBuffer) Resize(slotCount int) {
    if slotCount <= 0 {
        return
    }
    
    old := []*Word{}
    
    if buffer.tags != nil {
        buffer.mutex.Lock()
        for _,w := range(buffer.tags) {
            if w != nil {
                old = append(old,w)
            }
        }
        buffer.tags = make(map[string] *Word, slotCount)
        buffer.mutex.Unlock()
    }
    
    if buffer.words != nil {
        buffer.mutex.Lock()
        for _,w := range(buffer.words) {
            if w != nil {
                old = append(old,w)    
            }    
        }        
        buffer.words = make([]*Word, slotCount)
        buffer.mutex.Unlock()
    }

    buffer.mutex.Lock()
    buffer.slotCount = slotCount
    buffer.nextIndex = buffer.nextIndex % buffer.slotCount
    buffer.mutex.Unlock()

    for _,word := range old {
        gfx.WorldClock().DeleteTimer(word.timer)
    }
    
    if DEBUG_WORDBUFFER {
        log.Debug("%s resize %d",buffer.Desc(),slotCount)
    }
    
    buffer.ScheduleRefresh()
}

