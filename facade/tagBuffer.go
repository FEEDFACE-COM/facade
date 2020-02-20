
package facade

import(
    "fmt"
    "sync"
    "math/bits"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_TAGBUFFER = true

const maxTagLength = 32 // found experimentally


type SetItem struct {
    tag string
    count uint
    index uint
    timer *gfx.Timer
}


type TagBuffer struct {

    buf map[string] *SetItem

    indices []uint
    
    slotCount int
    duration float32
    shuffle bool
    
    rem []rune
	refreshChan chan bool
	mutex *sync.Mutex
    
}



func NewTagBuffer(refreshChan chan bool) *TagBuffer {
    ret := &TagBuffer{
        duration: float32(TagDefaults.Duration),
    }
    ret.slotCount = int(TagDefaults.Slot)
    ret.buf = make(map[string] *SetItem)
    ret.indices = make( []uint, ret.slotCount)
    for i:=0;i<ret.slotCount;i++ {
        ret.indices[i] = uint(i)
    }
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
                buffer.addItem( tmp )
                tmp = []rune{}

            default:
                tmp = append(tmp, r)
        }

    }
    buffer.rem = tmp
}



    
    
func (buffer *TagBuffer) Items(max int) []*SetItem {
    
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

func (buffer *TagBuffer) addItem(text []rune) {
    
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
        r := item.timer.Extend( gfx.Now() )
        if DEBUG_TAGBUFFER {
            if r {
                log.Debug("%s extended: '%s'",item.timer.Desc(),tag)
            }
        }

    } else if len(buffer.buf) >= buffer.slotCount {

        log.Debug("%s at %d/%d capacity, drop item '%s'",buffer.Desc(),len(buffer.buf),buffer.slotCount,tag)
            
    } else {

        triggerFun := func() {
            buffer.deleteItem(tag)
        }
        item = &SetItem{}
        item.tag = tag
        item.count = 1
        if buffer.shuffle {

        } else {
            item.index = buffer.indices[0]
            buffer.indices = buffer.indices[1:]
        }
        if DEBUG_TAGBUFFER {
            s := ""
            for _,r := range(buffer.indices) {
                s += fmt.Sprintf(" %d",r)
            }
            log.Debug("%s used #%d indices:%s",buffer.Desc(),item.index,s)
        }
        item.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.buf[tag] = item
        if DEBUG_TAGBUFFER {
            log.Debug("%s added '%s': %.1f",buffer.Desc(),tag,item.timer.Fader())
        }
    }
	buffer.mutex.Unlock()
    buffer.ScheduleRefresh()
}

func (buffer *TagBuffer) deleteItem(tag string) {
    item := buffer.buf[tag]
    buffer.mutex.Lock()
    delete(buffer.buf,tag)
	buffer.mutex.Unlock()
    buffer.indices = append(buffer.indices, item.index)

    if DEBUG_TAGBUFFER {
        log.Debug("%s item expired '%s'",buffer.Desc(),tag)
    }

    if DEBUG_TAGBUFFER {
        s := ""
        for _,r := range(buffer.indices) {
            s += fmt.Sprintf(" %d",r)
        }
        log.Debug("%s free #%d indices:%s",buffer.Desc(),item.index,s)
    }
    buffer.ScheduleRefresh()
}

func (buffer *TagBuffer) Clear() {
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


func (buffer *TagBuffer) Fill(fill []string) {

    buffer.Clear()

    rows := uint(len(fill))
    if DEBUG_TAGBUFFER {
        log.Debug("%s fill %s items",buffer.Desc(), rows)
    }
    
    for r:=uint(0); r<rows; r++  {
        
        buffer.addItem( []rune( fill[r] ) )
            
    }
    
}

func (buffer *TagBuffer) Desc() string {
    s := ""
    if buffer.shuffle {
        s = " ⧢"
    }
    ret := fmt.Sprintf("tagbuffer[%d/%d %.1f%s]",len(buffer.buf),buffer.slotCount,buffer.duration,s)
    return ret
}

func (buffer *TagBuffer) Dump() string {
    ret := ""
    
    items := buffer.Items(buffer.slotCount)
    for _,item := range items {    
    
        txt := string(item.tag)
        rem := "    "
        if item.timer != nil {
            rem = item.timer.Desc()
//            rem = fmt.Sprintf("%4.1f %4.1f",item.timer.Fader(),item.timer.Remaining(gfx.Now()))
        }
        ret += fmt.Sprintf("    #%05d %s %5d# %s\n",item.index,rem,item.count,txt) 
    }
    return ret
}

func (buffer *TagBuffer) SlotCount() int { return buffer.slotCount }
func (buffer *TagBuffer) TagCount() int { return len(buffer.buf) }
func (buffer *TagBuffer) Duration() float32 { return buffer.duration }

func (buffer *TagBuffer) SetDuration(duration float32) {
	buffer.duration = duration
    for _,item := range(buffer.buf) {
        item.timer.SetDuration(duration)
    }
}

func (buffer *TagBuffer) Resize(slotCount int) {
    if slotCount <= 0 {
        return
    }
    if DEBUG_TAGBUFFER {
        log.Debug("%s resize %d",buffer.Desc(),slotCount)
    }
    

    items := buffer.Items(slotCount)
    buffer.mutex.Lock()
    old := buffer.buf
    buffer.slotCount = slotCount
    buffer.buf = make(map[string] *SetItem, slotCount)
    buffer.indices = make( []uint, buffer.slotCount)
    for i:=0;i<buffer.slotCount;i++ {
        buffer.indices[i] = uint(i)
    }


    for _,item := range items {
        tag := item.tag
        item.index = buffer.indices[0]
        buffer.indices = buffer.indices[1:]
        buffer.buf[tag] = item
        delete(old, tag)
    }
    buffer.mutex.Unlock()    

    for _,item := range old {
        gfx.WorldClock().DeleteTimer(item.timer)
    }
    
    if DEBUG_TAGBUFFER {
        log.Debug("%s resized",buffer.Desc())
    }
    
    buffer.ScheduleRefresh()
}

