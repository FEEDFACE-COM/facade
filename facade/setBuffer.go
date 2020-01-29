
package facade

import(
    "fmt"
    "sync"
    "math/bits"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_SETBUFFER = true

const maxTagLength = 32 // found experimentally


type SetItem struct {
    tag string
    count uint
    index uint
    timer *gfx.Timer
}


type SetBuffer struct {
    buf map[string] *SetItem
    rem []rune
    duration float32
    nextIndex uint
    count uint
    
	refreshChan chan bool
	mutex *sync.Mutex
    
}



func NewSetBuffer(refreshChan chan bool) *SetBuffer {
    ret := &SetBuffer{
        duration: float32(TagDefaults.Duration),
    }
    ret.count = 0
    ret.buf = make(map[string] *SetItem)
	ret.refreshChan = refreshChan
	ret.mutex = &sync.Mutex{}
    return ret
}

func (buffer *SetBuffer) ProcessSequence(seq *ansi.S) {
    return
}

func (buffer *SetBuffer) ProcessRunes(runes []rune) {

    rem := append(buffer.rem, runes ... )
    tmp := []rune{}

    for _,r := range(rem) {

        switch r {
            case '\n':
                buffer.addItem( tmp )
                tmp = []rune{}

            default:
                tmp = append(tmp, r)
        }

    }
    buffer.rem = tmp
}



    
    
func (buffer *SetBuffer) Items(max int) []*SetItem {
    
    buffer.mutex.Lock()
    
    ret := []*SetItem{}
    tmp := make( map[string] *SetItem )
    
    for k,v := range buffer.buf {
        
        tmp[k] = v
        
    }
    
    mx := max
    if len(tmp) < mx {
        mx = len(tmp)
    }
    
    // pick max items from tmp
    for i := 0; i<mx; i++ {
        
        var min uint = 1 << bits.UintSize -1 
        var key string
        
        // look at each remaining item
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

func (buffer *SetBuffer) addItem(text []rune) {
    desc := buffer.Desc()
    
    if len(text) <= 0 {
        log.Debug("%s not adding empty string",buffer.Desc())
    }

    tag := string(text)
    if len(tag) > maxTagLength {
        tag = tag[0:maxTagLength-1]
    }

    buffer.mutex.Lock()
    item, ok := buffer.buf[tag]
    if ok {
        item.count += 1
        item.timer.Restart( gfx.Now() )
//        if DEBUG_SETBUFFER {
//            log.Debug("%s refreshed: '%s'",desc,tag)
//        }

    } else {
        triggerFun := func() {
            buffer.deleteItem(tag)
        }
        item = &SetItem{}
        item.tag = tag
        item.count = 1
        item.index = buffer.nextIndex
        buffer.nextIndex += 1
        item.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.buf[tag] = item
        buffer.count += 1
        if DEBUG_SETBUFFER {
            log.Debug("%s item added: '%s'",desc,tag)
        }
    }
	buffer.mutex.Unlock()
    
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *SetBuffer) deleteItem(idx string) {
    buffer.mutex.Lock() //thread #1 waits here
    delete(buffer.buf,idx)
    buffer.count -= 1
	buffer.mutex.Unlock()
    if DEBUG_SETBUFFER {
        log.Debug("%s item expired: '%s'",buffer.Desc(),idx)
    }
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *SetBuffer) Clear() {
    buffer.mutex.Lock()
    buffer.buf = make(map[string] *SetItem)
	select {
	case buffer.refreshChan <- true:
	default:
	}
	buffer.mutex.Unlock()
}


func (buffer *SetBuffer) Fill(fill []string) {

    buffer.Clear()

    rows := uint(len(fill))
    if DEBUG_SETBUFFER {
        log.Debug("%s fill %s items",buffer.Desc(), rows)
    }
    
    for r:=uint(0); r<rows; r++  {
        
        buffer.addItem( []rune( fill[r] ) )
            
    }
    
}

func (buffer *SetBuffer) Desc() string {
    ret := fmt.Sprintf("setbuffer[%.1f #%d]",buffer.duration,buffer.count)
    return ret
}

func (buffer *SetBuffer) Dump() string {
    ret := ""
    
    items := buffer.Items(32)
    for _,item := range items {    
    
        txt := string(item.tag)
        rem := "    "
        if item.timer != nil {
            rem = fmt.Sprintf("%4.1f",item.timer.Fader())
        }
        ret += fmt.Sprintf("#%05d %s %5d# %s\n",item.index,rem,item.count,txt) 
    }
    return ret
}

func (buffer *SetBuffer) Duration() float32 { return buffer.duration }

func (buffer *SetBuffer) SetDuration(duration float32) {
	buffer.duration = duration
}

