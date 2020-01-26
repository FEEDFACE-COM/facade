
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
    timer *gfx.Timer
}


type SetBuffer struct {
    buf map[string] *SetItem
    rem []rune
}



func NewSetBuffer(refreshChan chan bool) *SetBuffer {
    ret := &SetBuffer{}
    ret.buf = make(map[string] *SetItem)
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

func (buffer *SetBuffer) addItem(text []rune) {
    //lock?
    idx := string(text)
    item, ok := buffer.buf[idx]
    if ok {
        if DEBUG_SETBUFFER {
            log.Debug("%s again seen item %s",buffer.Desc(),idx)
        }
        item.count += 1


    } else {
        if DEBUG_SETBUFFER {
            log.Debug("%s first seen item %s",buffer.Desc(),idx)
        }
        item = &SetItem{}
        item.text = text 
        item.count = 1
        buffer.buf[idx] = item
    }

}


func (buffer *SetBuffer) Desc() string {
    return fmt.Sprintf("setbuffer[%d]",len(buffer.buf))
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


