
package grid

import (
    "fmt"
//    log "../log"
)


type PageDirection string
const (
    PageUp   PageDirection = "up"
    PageDown PageDirection = "down"
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


func (buffer *Buffer) addTail(text string) {
    line := &Line{text: text}
    
    if buffer.tail == nil || buffer.head == nil {
        buffer.head = line
        buffer.tail = line
        buffer.size += 1
        return    
    }
    
    line.prev = buffer.tail
    buffer.tail.next = line
    buffer.tail = line
    buffer.size += 1
}

func (buffer *Buffer) delHead() {
    if buffer.head == nil { return }
    ptr := buffer.head

    buffer.head = ptr.next
    buffer.head.prev = nil
    
    ptr.next = nil
    ptr.prev = nil

    buffer.size -= 1        
}


func (buffer *Buffer) Queue(text string, fade float64) {
    if buffer.tail.text == "" {
        buffer.tail.text = text
    } else {
        buffer.addTail(text)
        buffer.delHead()
    }
}

func (buffer *Buffer) Describe() string { return fmt.Sprintf("buffer[%d]",buffer.size) }

func (buffer *Buffer) Debug(dir PageDirection) string {
    ret := ""
    var ptr *Line
    var idx uint
    if dir == PageDown { ptr = buffer.head ; idx = 0 }
    if dir == PageUp   { ptr = buffer.tail ; idx = buffer.size }

    for ptr != nil {
        h,t := " ", " "
        if ptr == buffer.head { h = "h" }
        if ptr == buffer.tail { t = "t" }
        ret = ret + fmt.Sprintf("#%02d %s%s %s\n",idx,h,t,ptr.text)
        if dir == PageDown { ptr = ptr.next; idx = (idx + 1)   }
        if dir == PageUp   { ptr = ptr.prev; idx = (idx - 1 + buffer.size) % buffer.size }
    }
    return ret
    
}



