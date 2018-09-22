
// +build linux,arm

package gfx

import(
    "fmt"
    log "../log"
)

        

type Buffer struct {
    count uint
    index uint
    head  uint
    tail  uint
    items []*Text
}


func NewBuffer(count uint) *Buffer {
    ret := &Buffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.index = 0
    ret.head = 0
    ret.tail = count-1
    ret.items = make( []*Text, ret.count )
    return ret    
}


func (buffer *Buffer) Resize(count uint) {
    log.Debug("%s resize(%d)",buffer.Desc(),count)
    if count == 0 { count = 1 }
    newItems := make( []*Text, count )
    var idx uint = 0
    for ; idx<count && idx<buffer.count; idx++ {
        oidx := buffer.count + buffer.index - idx -1 
        newItems[idx] = buffer.items[(buffer.count+oidx)%buffer.count]    
    }
    //cleanup remains of old buffer
    for j:=idx; j<buffer.count; j++ {
        if buffer.items[j] != nil {
            buffer.items[j].Close()    
        }
    }
    buffer.count = count
    buffer.index = idx % count
    buffer.head = 0
    buffer.tail = count - 1
    buffer.items = newItems
}



func (buffer *Buffer) Queue(newItem *Text) {
    newIndex := (buffer.head)%buffer.count
    if buffer.items[newIndex] != nil {
        buffer.items[newIndex].Close()
    }
    buffer.items[ newIndex ] = newItem
    buffer.index = newIndex
    buffer.head = (buffer.head+1)%buffer.count
    buffer.tail = (buffer.tail+1)%buffer.count
}

func (buffer *Buffer) Item(off uint) *Text {
    return buffer.items[ (buffer.count+buffer.index+off) % buffer.count ]
}

func (buffer *Buffer) Tail(off uint) *Text {
    idx := buffer.count + buffer.tail - off
    return buffer.items[idx % buffer.count]    
}



func (buffer *Buffer) Items() []*Text {
    return buffer.items
}

func (buffer *Buffer) Index(off uint) uint {
    return (buffer.count+buffer.index+off) % buffer.count
}

func (buffer *Buffer) Desc() string { 
    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    return ret
}


func (buffer *Buffer) Dump() string {
//    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    ret := ""
    for i:= uint(0);i<buffer.count;i++ {
        idx := ( i ) % buffer.count 
        item := buffer.items[ idx ]
//        if item == nil { continue }
        s0 := "#"
        s2 := " "
        if buffer.head == i { s0 = "h" }
        if buffer.tail == i { s0 = "t" }
        if buffer.head == i && buffer.tail == i { s0 = "X" }
        if buffer.index == i { s2 = ">" }
        s1 := ""
        if item != nil { s1 = (*item).Desc() }
        ret += fmt.Sprintf("  %s%s%02d %s\n",s2,s0,idx,s1) 
    }
    return ret
}






