
package facade

import(
//    "fmt"
    "sync"
    "math/rand"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_TAGBUFFER = true

const maxTagLength = 32 // found experimentally



type TagBuffer struct {
    
    tags  map[string] *Word

    slotCount int
    duration float32
    shuffle bool

    
    rem []rune
	refreshChan chan bool
	mutex *sync.Mutex
    
}


func NewTagBuffer(refreshChan chan bool) *WordBuffer {
    ret := &WordBuffer{
        duration: float32(SetDefaults.Duration),
        slotCount: int(SetDefaults.Slot),
    }
    ret.tags = make( map [string]*Word ) 
	ret.refreshChan = refreshChan
	ret.mutex = &sync.Mutex{}
    return ret
}

func (buffer *TagBuffer) ScheduleRefresh() {
	select {
	case buffer.refreshChan <- true:
	default:
	}
}


func (buffer *TagBuffer) ProcessSequence(seq *ansi.S) {
    return
}

func (buffer *TagBuffer) ProcessRunes(runes []rune) {

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

func (buffer *TagBuffer) addWord(txt []rune) {

    text := string(txt)
    
    if len(text) <= 0 {
        log.Debug("%s not adding empty string",buffer.Desc())
        return
    }

    if len(text) > maxTagLength {
        text = text[0: maxTagLength-1]
    }

    buffer.mutex.Lock()
    word, ok := buffer.tags[text]
    if ok {
        word.count += 1
        r := word.timer.Extend( gfx.Now() )
        if DEBUG_TAGBUFFER {
            if r {
                log.Debug("%s extended: '%s'",word.timer.Desc(),word.text)
            }
        }

    } else if len(buffer.tags) >= buffer.slotCount {

        log.Debug("%s at %d/%d capacity, drop word '%s'",buffer.Desc(),len(buffer.tags),buffer.slotCount,text)
            
    } else {

        triggerFun := func() {
            buffer.deleteWord(tag)
        }
        word = &Word{}
        word.tag = tag
        word.count = 1
        
        {
            idxmap = make( map [int]bool )
            for i:=0; i<buffer.slotCount; i++ {
                idxmap[i] = true
            }
            for _,word := range(buffer.tags) {
                idx := word.index
                delete(idxmap,idx)
            }
            var indices []int
            for i,_ := range(idxmap) {
                indices = append(indices,i)
            }
            r := rand.Int31n( len(indices) )
            word.index = indices[r]
        }        
        word.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        
        buffer.tags[text] = word
        if DEBUG_WORDBUFFER {
            log.Debug("%s added '%s': %.1f",buffer.Desc(),text,word.timer.Fader())
        }
    }
	buffer.mutex.Unlock()
    buffer.ScheduleRefresh()
}



