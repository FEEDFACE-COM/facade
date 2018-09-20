//
package render

import(
    "fmt"
    conf "../conf"
)

type BufferItem struct {
    time float32
    text conf.Text
}
        

type Buffer struct {
    count uint
    index uint
    items []*BufferItem
}



func NewBuffer(count uint) Buffer {
    ret := Buffer{}
    ret.count = count
    ret.index = ret.count
    ret.items = make( []*BufferItem, ret.count )
    return ret    
}


func (buffer *Buffer) Resize(count uint) {

    newItems := make( []*BufferItem, count )
    var idx uint = 0
    for ; idx<count && idx<buffer.count; idx++ {
        newItems[idx] = buffer.items[idx]    
    }
    buffer.count = count
    buffer.index = idx
    buffer.items = newItems
}


func (buffer *Buffer) Queue(time float32, text conf.Text) {
    newItem := &BufferItem{time: time, text: text}
    buffer.items[buffer.index] = newItem
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
    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    for i:= uint(0);i<buffer.count;i++ {
        item := buffer.items[i]
//        if item == nil { continue }
        s0 := "#"
        if buffer.index == i { s0 = ">" }
        s1 := ""
        if item != nil { s1 = fmt.Sprintf("%5.1f %s",item.time,item.text) }
        ret += fmt.Sprintf("\n  %s%02d %s",s0,i,s1) 
    }
    return ret
}






