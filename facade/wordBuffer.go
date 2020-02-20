
package facade

import(
    "fmt"
    "sync"
    "math/bits"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_WORDBUFFER = true

const maxWordLength = 32 // found experimentally


type Word struct {
    tag string
    count uint
    index uint
    timer *gfx.Timer
}


type WordBuffer struct {
    
    tags  map[string] *Word
    words []*Word

    indices []uint
    
    slotCount int
    duration float32
    shuffle bool

    tagmode bool
    
    rem []rune
	refreshChan chan bool
	mutex *sync.Mutex
    
}



func NewWordBuffer(tagmode bool, refreshChan chan bool) *WordBuffer {
    ret := &WordBuffer{
        tagmode: tagmode,
        duration: float32(SetDefaults.Duration),
        slotCount: int(SetDefaults.Slot),
    }
    if tagmode {
        ret.tags = make( map [string]*Word ) 
    } else {
        ret.words = make( []*Word, ret.slotCount )
    }
    ret.indices = make( []uint, ret.slotCount)
    for i:=0;i<ret.slotCount;i++ {
        ret.indices[i] = uint(i)
    }
	ret.refreshChan = refreshChan
	ret.mutex = &sync.Mutex{}
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



    
    
func (buffer *WordBuffer) Words(max int) []*Word {
    
    buffer.mutex.Lock()
    
    ret := []*Word{}
    tmp := make( map[string] *Word)
    
    for k,v := range buffer.tags {
        
        tmp[k] = v
        
    }
    
    mx := max
    if len(tmp) < mx {
        mx = len(tmp)
    }
    
    // pick max words from tmp
    for i := 0; i<mx; i++ {
        
        var min uint = 1 << bits.UintSize -1 
        var key string
        
        // look at each remaining word
        for k,v := range tmp {
            
            // if smallest seen, keep note of key
            if v.index < min {
                min = v.index
                key = k
            } 
            
        }
        ret = append(ret, tmp[key])
        delete(tmp, key)
        
            
    }
    
    
    buffer.mutex.Unlock()
    return ret
}    

func (buffer *WordBuffer) addWord(text []rune) {
    
    if len(text) <= 0 {
        log.Debug("%s not adding empty string",buffer.Desc())
        return
    }

    tag := string(text)
    if len(tag) > maxWordLength {
        tag = tag[0: maxWordLength-1]
    }

    buffer.mutex.Lock()
    word, ok := buffer.tags[tag]
    if ok {
        word.count += 1
        r := word.timer.Extend( gfx.Now() )
        if DEBUG_WORDBUFFER {
            if r {
                log.Debug("%s extended: '%s'",word.timer.Desc(),tag)
            }
        }

    } else if len(buffer.tags) >= buffer.slotCount {

        log.Debug("%s at %d/%d capacity, drop word '%s'",buffer.Desc(),len(buffer.tags),buffer.slotCount,tag)
            
    } else {

        triggerFun := func() {
            buffer.deleteWord(tag)
        }
        word = &Word{}
        word.tag = tag
        word.count = 1
        if buffer.shuffle {

        } else {
            word.index = buffer.indices[0]
            buffer.indices = buffer.indices[1:]
        }
        if DEBUG_WORDBUFFER {
            s := ""
            for _,r := range(buffer.indices) {
                s += fmt.Sprintf(" %d",r)
            }
            log.Debug("%s used #%d indices:%s",buffer.Desc(),word.index,s)
        }
        word.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.tags[tag] = word
        if DEBUG_WORDBUFFER {
            log.Debug("%s added '%s': %.1f",buffer.Desc(),tag,word.timer.Fader())
        }
    }
	buffer.mutex.Unlock()
    buffer.ScheduleRefresh()
}

func (buffer *WordBuffer) deleteWord(tag string) {
    word := buffer.tags[tag]
    buffer.mutex.Lock()
    delete(buffer.tags,tag)
	buffer.mutex.Unlock()
    buffer.indices = append(buffer.indices, word.index)

    if DEBUG_WORDBUFFER {
        log.Debug("%s word expired '%s'",buffer.Desc(),tag)
    }

    if DEBUG_WORDBUFFER {
        s := ""
        for _,r := range(buffer.indices) {
            s += fmt.Sprintf(" %d",r)
        }
        log.Debug("%s free #%d indices:%s",buffer.Desc(),word.index,s)
    }
    buffer.ScheduleRefresh()
}

func (buffer *WordBuffer) Clear() {
    old := []*gfx.Timer{}
    
    buffer.mutex.Lock()
    for _,word := range buffer.tags {
        old = append(old,word.timer)
    }
    
    buffer.tags = make(map[string] *Word)
	buffer.mutex.Unlock()
	
	for _,timer := range old {
    	gfx.WorldClock().DeleteTimer(timer)
    }

    buffer.ScheduleRefresh()
}


func (buffer *WordBuffer) Fill(fill []string) {

    buffer.Clear()

    rows := uint(len(fill))
    if DEBUG_WORDBUFFER {
        log.Debug("%s fill %s words",buffer.Desc(), rows)
    }
    
    for r:=uint(0); r<rows; r++  {
        
        buffer.addWord( []rune( fill[r] ) )
            
    }
    
}

func (buffer *WordBuffer) Desc() string {
    b := "words"
    if buffer.tagmode {
        b = "tags"
    }
    s := ""
    if buffer.shuffle {
        s = "⧢"
    }
    ret := fmt.Sprintf("wordbuffer[%d/%d %s %.1f%s]",len(buffer.tags),buffer.slotCount,b,buffer.duration,s)
    return ret
}

func (buffer *WordBuffer) Dump() string {
    ret := ""
    
    words := buffer.Words(buffer.slotCount)
    for _,word := range words {
    
        txt := string(word.tag)
        rem := "    "
        if word.timer != nil {
            rem = word.timer.Desc()
//            rem = fmt.Sprintf("%4.1f %4.1f",word.timer.Fader(),word.timer.Remaining(gfx.Now()))
        }
        ret += fmt.Sprintf("    #%05d %s %5d# %s\n",word.index,rem,word.count,txt)
    }
    return ret
}

func (buffer *WordBuffer) SlotCount() int    { return buffer.slotCount }
func (buffer *WordBuffer) WordCount() int     { return len(buffer.tags) }
func (buffer *WordBuffer) Duration() float32 { return buffer.duration }
func (buffer *WordBuffer) Shuffle() bool { return buffer.shuffle }

func (buffer *WordBuffer) SetDuration(duration float32) {
	buffer.duration = duration
    for _,word := range(buffer.tags) {
        word.timer.SetDuration(duration)
    }
}

func (buffer *WordBuffer) Resize(slotCount int) {
    if slotCount <= 0 {
        return
    }
    if DEBUG_WORDBUFFER {
        log.Debug("%s resize %d",buffer.Desc(),slotCount)
    }
    

    words := buffer.Words(slotCount)
    buffer.mutex.Lock()
    old := buffer.tags
    buffer.slotCount = slotCount
    buffer.tags = make(map[string] *Word, slotCount)
    buffer.indices = make( []uint, buffer.slotCount)
    for i:=0;i<buffer.slotCount;i++ {
        buffer.indices[i] = uint(i)
    }


    for _,word := range words {
        tag := word.tag
        word.index = buffer.indices[0]
        buffer.indices = buffer.indices[1:]
        buffer.tags[tag] = word
        delete(old, tag)
    }
    buffer.mutex.Unlock()    

    for _,word := range old {
        gfx.WorldClock().DeleteTimer(word.timer)
    }
    
    if DEBUG_WORDBUFFER {
        log.Debug("%s resized",buffer.Desc())
    }
    
    buffer.ScheduleRefresh()
}

