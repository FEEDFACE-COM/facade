
// +build linux,arm

package gfx

import(
    "fmt"
    "strings"
    log "../log"
)

        

type RingBuffer struct {
    count uint
    head  uint
    tail  uint
    items []*Text
}


const DEBUG_BUFFER = false

func NewRingBuffer(count uint) *RingBuffer {
    ret := &RingBuffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.head = 0
    ret.tail = count-1
    ret.items = make( []*Text, ret.count )
    return ret    
}





func (buffer *RingBuffer) Resize(newCount uint) {
    log.Debug("resize %s -> %d",buffer.Desc(),newCount)
    if newCount == 0 { newCount = 1 }

    newItems := make( []*Text, newCount )


    if newCount < buffer.count {

        d := buffer.count - newCount

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newCount && idx<buffer.count; idx++ {
            newItems[idx] = buffer.Head( idx + d )
        }
        
        //cleanup remaining items
        for j:=idx; j<buffer.count; j++ {
            if buffer.items[j] != nil {
                buffer.items[j].Close()    
            }
        }

    } else if newCount > buffer.count {
        
        // copy all items
        d := newCount - buffer.count
        for idx:= uint(0); idx<buffer.count; idx++ {
            newItems[ (idx+d) % newCount ] = buffer.items[idx]
        } 

    
    } else {
        //nop
    }

    
    //adjust buffer info
    buffer.count = newCount
    buffer.head = 0
    buffer.tail = newCount - 1
    buffer.items = newItems
}


func (buffer *RingBuffer) WriteText(text *Text) {

    text.Text = strings.Trim(text.Text, "\n")

    buffer.queue( text )
    
}


func (buffer *RingBuffer) WriteString(str string) {

    s := strings.Trim(str, "\n")
    buffer.queue( NewText(s) )
    
}



func (buffer *RingBuffer) queue(newItem *Text) {
    if buffer.items[buffer.tail] == nil && newItem != nil {
        buffer.items[buffer.tail] = newItem
        if DEBUG_BUFFER {
            tmp := ""
            if newItem != nil { tmp = newItem.Desc() }
            log.Debug("jump %s\n%s",tmp,buffer.Dump())
        }        
        return
    }


    idx := (buffer.head)%buffer.count
    
    // clean up old item
    if buffer.items[idx] != nil {
        buffer.items[idx].Close()
    }

    //insert item
    buffer.items[ idx ] = newItem

    //adjust buffer info        
    buffer.head = (buffer.head+1)%buffer.count
    buffer.tail = (buffer.tail+1)%buffer.count

    if DEBUG_BUFFER {
        tmp := ""
        if newItem != nil { tmp = newItem.Desc() }
        log.Debug("queue %s\n%s",tmp,buffer.Dump())
    }        
}







func (buffer *RingBuffer) Tail(off uint) *Text {
    // off /= buffer.count   // probably?
    idx := buffer.count + buffer.tail - off
    return buffer.items[idx % buffer.count]    
}

func (buffer *RingBuffer) Head(off uint) *Text {
    idx := buffer.count + buffer.head + off
    return buffer.items[idx % buffer.count]    
}






func (buffer *RingBuffer) Desc() string { 
    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    return ret
}


func (buffer *RingBuffer) Dump() string {
    ret := ""
    ret += fmt.Sprintf("ring[%d]\n",buffer.count)
    for i:= buffer.count-1; i-1>uint(0);i-- {
//    for i:= uint(0); i<buffer.count;i++ {
        idx := ( i ) % buffer.count 
        item := buffer.items[ idx ]

        s0 := fmt.Sprintf("%d", ((idx-buffer.head)%buffer.count)%10)
        s1 := fmt.Sprintf("%d", ((buffer.count+buffer.tail-idx)%buffer.count)%10)
        s2 := ""
        
        t0,t1 := " "," "
        if buffer.head == i { t0 = "H" }
        if buffer.tail == i { t1 = "T" }
        if item != nil { s2 = (*item).Desc() }
        ret += fmt.Sprintf("  %s%s %s%s %s\n",t0,t1,s0,s1,s2) 
    }
    return ret
}






