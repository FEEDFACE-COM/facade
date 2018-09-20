//
package render

import(
    "fmt"
)

type BufferItem interface {
    Desc() string
    Close()
}
        

type Buffer struct {
    count uint
    index uint
    item []*BufferItem
}



func NewBuffer(count uint) Buffer {
    ret := Buffer{}
    ret.count = count
    ret.index = ret.count
    ret.item = make( []*BufferItem, ret.count )
    return ret    
}


func (buffer *Buffer) Resize() {
    
    
}


func (buffer *Buffer) Queue(newItem *BufferItem) {
    tmp := buffer.item[buffer.index]
    buffer.item[buffer.index] = newItem
    buffer.index = ( buffer.count + buffer.index - 1 ) % buffer.count
    (*tmp).Close()
}

func (buffer *Buffer) Item(off uint) *BufferItem {
    return buffer.item[ (buffer.count+buffer.index+off) % buffer.count ]
}

func (buffer *Buffer) Desc() string { 
    ret := fmt.Sprintf("buffer[%d]",buffer.count)
    return ret
}


func (buffer *Buffer) Dump() string {
    ret := ""
    for i:= uint(0);i<buffer.count;i++ {
        s := ""
        if buffer.index == i { s = "*" }
        ret += fmt.Sprintf("  #%02d %s \n",i,s,(*buffer.item[i]).Desc()) 
    }
    return ret
}






