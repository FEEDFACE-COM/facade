
package facade

import(
    "fmt"
    log "../log"
    gfx "../gfx"
    "github.com/pborman/ansi"
)

const DEBUG_SETBUFFER = true


type SetItem struct {
    text []rune
    count uint
    index uint
    timer *gfx.Timer
}


type SetBuffer struct {
    buf map[string] *SetItem
    rem []rune
    duration float32
    nextIndex uint
    
	refreshChan chan bool
    
}



func NewSetBuffer(refreshChan chan bool) *SetBuffer {
    ret := &SetBuffer{
        duration: float32(TagDefaults.Duration),
    }
    ret.buf = make(map[string] *SetItem)
	ret.refreshChan = refreshChan
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



func (buffer *SetBuffer) Tags(max int) []string {
    ret := []string{}
    for idx := range buffer.buf {
        ret = append(ret,idx)
    }
    if len(ret) < max {
        return ret
    }
    return ret[0:max-1]
    }
    
    
func (buffer *SetBuffer) Items(max int) map[string] *SetItem {
    
    ret := make( map[string] *SetItem )
    
    for k,v := range buffer.buf {
        ret[k] = v
    }
    
    return ret
}    

func (buffer *SetBuffer) addItem(text []rune) {
    //lock?
    idx := string(text)
    item, ok := buffer.buf[idx]
    if ok {
        item.count += 1
        item.timer.Restart( gfx.Now() )
        if DEBUG_SETBUFFER {
            log.Debug("%s item refreshed: '%s'",buffer.Desc(),idx)
        }

    } else {
        triggerFun := func() {
            buffer.deleteItem(idx)
        }
        item = &SetItem{}
        item.text = text        
        item.count = 1
        item.index = buffer.nextIndex
        buffer.nextIndex += 1
        item.timer = gfx.WorldClock().NewTimer(buffer.duration, false, nil, triggerFun)
        buffer.buf[idx] = item
        if DEBUG_SETBUFFER {
            log.Debug("%s item added: '%s'",buffer.Desc(),idx)
        }
    }
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *SetBuffer) deleteItem(idx string) {
    delete(buffer.buf,idx)
    if DEBUG_SETBUFFER {
        log.Debug("%s item expired: '%s'",buffer.Desc(),idx)
    }
	select {
	case buffer.refreshChan <- true:
	default:
	}
}

func (buffer *SetBuffer) Clear() {
    buffer.buf = make(map[string] *SetItem)
	select {
	case buffer.refreshChan <- true:
	default:
	}
}


func (buffer *SetBuffer) Fill(fill []string) {

    buffer.buf = make(map[string] *SetItem)

    rows := uint(len(fill))
    if DEBUG_SETBUFFER {
        log.Debug("%s fill %s items",buffer.Desc(), rows)
    }
    
    for r:=uint(0); r<rows; r++  {
        
        buffer.addItem( []rune( fill[r] ) )
            
    }
    
}

func (buffer *SetBuffer) Desc() string {
    return fmt.Sprintf("setbuffer[%.1f #%d]",buffer.duration,len(buffer.buf))
}

func (buffer *SetBuffer) Dump() string {
    ret := ""
    for _,itm := range(buffer.buf) {
        txt := string(itm.text)
        rem := " "
        if itm.timer != nil {
            rem = fmt.Sprintf("%4.1f",itm.timer.Fader())
        }
        ret += fmt.Sprintf("%5d# %s %s\n",itm.count,rem,txt) 
    }
    return ret
}

func (buffer *SetBuffer) Duration() float32 { return buffer.duration }

func (buffer *SetBuffer) SetDuration(duration float32) {
	buffer.duration = duration
}
