
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
    count int
    max int
    
	refreshChan chan bool
	mutex *sync.Mutex
    
}



func NewSetBuffer(refreshChan chan bool) *SetBuffer {
    ret := &SetBuffer{
        duration: float32(TagDefaults.Duration),
    }
    ret.count = 0
    ret.max = int(TagDefaults.Slot)
    ret.buf = make(map[string] *SetItem)
	ret.refreshChan = refreshChan
	ret.mutex = &sync.Mutex{}
    return ret
}


func (buffer *SetBuffer) ScheduleRefresh() {
	select {
	case buffer.refreshChan <- true:
	default:
	}
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
    
    if len(text) <= 0 {
        log.Debug("%s not adding empty string",buffer.Desc())
        return
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

    } else if len(buffer.buf) >= buffer.max {

        log.Debug("%s at %d/%d capacity, drop item '%s'",buffer.Desc(),len(buffer.buf),buffer.max,tag)
            
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
        buffer.count = len(buffer.buf)
        if DEBUG_SETBUFFER {
            log.Debug("%s added '%s': %.1f",buffer.Desc(),tag,item.timer.Fader())
        }
    }
	buffer.mutex.Unlock()
    buffer.ScheduleRefresh()
}

func (buffer *SetBuffer) deleteItem(idx string) {
    buffer.mutex.Lock()
    delete(buffer.buf,idx)
    buffer.count = len(buffer.buf)
	buffer.mutex.Unlock()
    if DEBUG_SETBUFFER {
        log.Debug("%s item expired '%s'",buffer.Desc(),idx)
    }
    buffer.ScheduleRefresh()
}

func (buffer *SetBuffer) Clear() {
    old := []*gfx.Timer{}
    
    buffer.mutex.Lock()
    for _,item := range buffer.buf {
        old = append(old,item.timer)
    }
    
    buffer.buf = make(map[string] *SetItem)
	buffer.mutex.Unlock()
	
	for _,timer := range old {
    	gfx.WorldClock().DeleteTimer(timer)
    }

    buffer.ScheduleRefresh()
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
    ret := fmt.Sprintf("setbuffer[%d/%d %.1f ]",buffer.count,buffer.max,buffer.duration)
    return ret
}

func (buffer *SetBuffer) Dump() string {
    ret := ""
    
    items := buffer.Items(buffer.max)
    for _,item := range items {    
    
        txt := string(item.tag)
        rem := "    "
        if item.timer != nil {
            rem = fmt.Sprintf("%4.1f",item.timer.Fader())
        }
        ret += fmt.Sprintf("    #%05d %s %5d# %s\n",item.index,rem,item.count,txt) 
    }
    return ret
}

func (buffer *SetBuffer) Max() int { return buffer.max }

func (buffer *SetBuffer) Duration() float32 { return buffer.duration }

func (buffer *SetBuffer) SetDuration(duration float32) {
	buffer.duration = duration
}

func (buffer *SetBuffer) Resize(max int) {
    if max <= 0 {
        return
    }
    if DEBUG_SETBUFFER {
        log.Debug("%s resize %d",buffer.Desc(),max)
    }
    

    items := buffer.Items(max)
    buffer.mutex.Lock()
    old := buffer.buf
    buffer.max = max
    buffer.buf = make(map[string] *SetItem, max)

    for _,item := range items {
        tag := item.tag
        buffer.buf[tag] = item
        delete(old, tag)
    }
    buffer.count = len(buffer.buf)
    buffer.mutex.Unlock()    

    for _,item := range old {
        gfx.WorldClock().DeleteTimer(item.timer)
    }
    
    if DEBUG_SETBUFFER {
        log.Debug("%s resized",buffer.Desc())
    }
    
    buffer.ScheduleRefresh()
}

