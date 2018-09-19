
package grid

import (
    "fmt"
    conf "../conf"
    log "../log"
)




type Buffer struct {
    head *Line
    tail *Line
    size uint
    twice bool
}

type Line struct {
    text string
    next *Line
    prev *Line
}


func NewBuffer(size uint) *Buffer {
    buffer := &Buffer{}
    for i:=uint(0); i<size; i++ {
        buffer.addTail("")
    }
    return buffer
}


func NewBufferDebug(size uint) *Buffer {
    buffer := NewBuffer(size)
    for i:=uint(0); i<size; i++ {
        buffer.addTail( fmt.Sprintf("buffered %d/%d",i+1,size) )
        buffer.delHead()
    }
    return buffer

}



func (buffer *Buffer) Queue(text string, fade float64) {
    if buffer.tail.text == "" {
        buffer.tail.text = text
    } else {
        buffer.addTail(text)
        buffer.delHead()
    }
}






func (buffer *Buffer) addTail(text string) {
    line := &Line{text: text}
    
    if buffer.tail == nil && buffer.head == nil {
        buffer.head = line
        buffer.tail = line
        buffer.size = 1
        return    
    }
    
    line.prev = buffer.tail
    buffer.tail.next = line
    buffer.tail = line
    buffer.size += 1
}


func (buffer *Buffer) addHead(text string) {
    line := &Line{text: text}
    
    if buffer.tail == nil && buffer.head == nil {
        buffer.head = line
        buffer.tail = line
        buffer.size = 1   
        return 
    }
    
    line.next = buffer.head
    buffer.head.prev = line
    buffer.head = line
    buffer.size += 1
    
}



func (buffer *Buffer) delTail() {
    if buffer.tail == nil { return }
    
    if buffer.tail == buffer.head {
        ptr := buffer.tail
        buffer.tail = nil
        buffer.head = nil
        buffer.size = 0
        ptr.next = nil
        ptr.prev = nil
    }
    
    ptr := buffer.tail
    
    buffer.tail = buffer.tail.prev
    buffer.tail.next = nil
    buffer.size -= 1
    
    ptr.next = nil
    ptr.prev = nil
    
}

func (buffer *Buffer) delHead() {
    if buffer.head == nil { return }

    if buffer.head == buffer.tail {
        ptr := buffer.head
        buffer.tail = nil
        buffer.head = nil
        buffer.size = 0    
        ptr.next = nil
        ptr.prev = nil
    }


    ptr := buffer.head

    buffer.head = buffer.head.next
    buffer.head.prev = nil
    buffer.size -= 1        
    
    ptr.next = nil
    ptr.prev = nil

}














func (buffer *Buffer) Configure(config *conf.GridConfig) {
    log.Debug("configure buffer: %s",config.Describe())
    // add/del as needed
}

func (buffer *Buffer) Describe() string { return fmt.Sprintf("buffer[%d]",buffer.size) }

func (buffer *Buffer) Debug(dir conf.PageDirection) string {
    ret := ""
    var ptr *Line
    var idx uint
    if dir == conf.PageDown { ptr = buffer.head ; idx = 0 }
    if dir == conf.PageUp   { ptr = buffer.tail ; idx = buffer.size }

    for ptr != nil {
        h,t := " ", " "
        if ptr == buffer.head { h = "h" }
        if ptr == buffer.tail { t = "t" }
        ret = ret + fmt.Sprintf("#%02d %s%s %s\n",idx,h,t,ptr.text)
        if dir == conf.PageDown { ptr = ptr.next; idx = (idx + 1)   }
        if dir == conf.PageUp   { ptr = ptr.prev; idx = (idx - 1 + buffer.size) % buffer.size }
    }
    return ret
    
}



