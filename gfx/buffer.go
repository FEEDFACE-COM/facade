
// +build linux,arm

package gfx

import(
    "fmt"
    log "../log"
)

        

type Buffer struct {
    count uint
    head  uint
    tail  uint
    items []*Text
}


func NewBuffer(count uint) *Buffer {
    ret := &Buffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.head = 0
    ret.tail = count-1
    ret.items = make( []*Text, ret.count )
    return ret    
}





func (buffer *Buffer) Resize(newCount uint) {
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






func (buffer *Buffer) Queue(newItem *Text) {

    if buffer.items[buffer.tail] == nil && newItem != nil {
        buffer.items[buffer.tail] = newItem
        if newItem == nil {
            log.Debug("jump nil %s\n%s",buffer.Desc(),buffer.Dump())
        } else {
            log.Debug("jump '%s' %s\n%s",newItem.Desc(),buffer.Desc(),buffer.Dump())
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
    
    if newItem == nil {
        log.Debug("queue nil %s\n%s",buffer.Desc(),buffer.Dump())
    } else {
        log.Debug("queue '%s' %s\n%s",newItem.Desc(),buffer.Desc(),buffer.Dump())
    }
}







func (buffer *Buffer) Tail(off uint) *Text {
    // off /= buffer.count   // probably?
    idx := buffer.count + buffer.tail - off
    return buffer.items[idx % buffer.count]    
}

func (buffer *Buffer) Head(off uint) *Text {
    idx := buffer.count + buffer.head + off
    return buffer.items[idx % buffer.count]    
}






func (buffer *Buffer) Desc() string { 
    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    return ret
}


func (buffer *Buffer) Dump() string {
    ret := ""
    
    for i:= uint(0);i<buffer.count;i++ {
        idx := ( i ) % buffer.count 
        item := buffer.items[ idx ]

        s0 := fmt.Sprintf("%d", ((idx-buffer.head)%buffer.count)%10)
        s1 := fmt.Sprintf("%d", ((buffer.count+buffer.tail-idx)%buffer.count)%10)
        s2 := ""
        
        if buffer.head == i { s0 = "H" }
        if buffer.tail == i { s1 = "T" }
        if item != nil { s2 = (*item).Desc() }
        ret += fmt.Sprintf("  %s%s %s\n",s0,s1,s2) 
    }
    return ret
}






