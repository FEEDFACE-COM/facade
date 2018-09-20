//
package modes

import(
    "fmt"
    log "../log"
)

type BufferItem interface {
    Desc() string
}
        

type Buffer struct {
    count uint
    index uint
    items []*BufferItem
}



func NewBuffer(count uint) Buffer {
    ret := Buffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.index = 0
    ret.items = make( []*BufferItem, ret.count )
    return ret    
}


func (buffer *Buffer) Resize(count uint) {
    log.Debug("%s resize(%d)",buffer.Desc(),count)
    if count == 0 { count = 1 }
    newItems := make( []*BufferItem, count )
    var idx uint = 0
    for ; idx<count && idx<buffer.count; idx++ {
        oidx := buffer.count + buffer.index - idx -1 
        newItems[idx] = buffer.items[(buffer.count+oidx)%buffer.count]    
    }
    buffer.count = count
    buffer.index = idx % count
    buffer.items = newItems
}


func (buffer *Buffer) Queue(newItem BufferItem) {
    buffer.items[ buffer.index ] = &newItem
    buffer.index = ( buffer.count + buffer.index + 1 ) % buffer.count
}

func (buffer *Buffer) Item(off uint) *BufferItem {
    return buffer.items[ (buffer.count+buffer.index+off) % buffer.count ]
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
        if buffer.index == i { s0 = ">" }
        s1 := ""
        if item != nil { s1 = (*item).Desc() }
        ret += fmt.Sprintf("  %s%02d %s\n",s0,idx,s1) 
    }
    return ret
}






