
package render

import (
    "fmt"
)

type Buffer struct {
    head *Line
    tail *Line
    size uint
}

type Line struct {
    text string
    next *Line
    prev *Line
}

func NewBuffer(size uint) *Buffer {
    buffer := &Buffer{head:nil, tail: nil, size: size}
    if size == 0 { return buffer }
    buffer.tail = &Line{text: ""}
    buffer.head = buffer.tail
    for i:=uint(1); i<buffer.size; i++ {
        line := &Line{text: "", next:buffer.head}
        buffer.head = line    
    }
    return buffer
}


func (buffer *Buffer) AddTail(text string) {
    line := &Line{text: text}
    
    if buffer.tail == nil {
        buffer.head = line
        buffer.tail = line
        return    
    }
    
    line.prev = buffer.tail
    buffer.tail.next = line
    buffer.tail = line
    
    buffer.size += 1
}

func (buffer *Buffer) DelHead() {
    if buffer.head == nil { return }
    ptr := buffer.head

    buffer.head = ptr.next
    buffer.head.prev = nil
    
    ptr.next = nil
    ptr.prev = nil

    buffer.size -= 1        
}

func (buffer *Buffer) Desc() string { return fmt.Sprintf("buffer(%d)",buffer.size) }

func (buffer *Buffer) Debug() {
    
}



