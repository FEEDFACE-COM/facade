
package render

import (
    "fmt"
)

type Buffer struct {
    head *Line
    tail *Line
    size uint
    page PageDirection
    twice bool
}

type Line struct {
    text string
    next *Line
    prev *Line
}



type PageDirection string
const (
    PageUp   PageDirection = "up"
    PageDown PageDirection = "down"
    DefaultPageDirection PageDirection = "up"
)
    

func EmptyBuffer(size uint) *Buffer {
    buffer := &Buffer{page:DefaultPageDirection}
    for i:=uint(0); i<=size; i++ {
        buffer.addTail("")
    }
    return buffer
}

func DebugBuffer(size uint) *Buffer {
    buffer := &Buffer{page: DefaultPageDirection}
    for i:=uint(0); i<size; i++ {
        buffer.addTail(fmt.Sprintf("buffered %d/%d",i+1,size) )
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


func (buffer *Buffer) Queue(timer *Timer, text string) {
    if timer.fade() < 0.5 && buffer.size >= 2 && buffer.head.text == "" && buffer.head.next.text == "" {
        buffer.head.text = text
        buffer.twice = true
    } else if buffer.head.text == "" {
        buffer.head.text = text
    } else {
        buffer.addTail(text)
        buffer.delHead()
    }
}

func (buffer *Buffer) Desc() string { return fmt.Sprintf("gridBuffer[%d]",buffer.size) }

func (buffer *Buffer) Debug() string {
    ret := ""
    var ptr *Line
    var idx uint
    if buffer.page == PageDown { ptr = buffer.head ; idx = 0 }
    if buffer.page == PageUp   { ptr = buffer.tail ; idx = buffer.size }

    for ptr != nil {
        h,t := " ", " "
        if ptr == buffer.head { h = "h" }
        if ptr == buffer.tail { t = "t" }
        ret = ret + fmt.Sprintf("#%02d %s%s %s\n",idx,h,t,ptr.text)
        if buffer.page == PageDown { ptr = ptr.next; idx = (idx + 1)   }
        if buffer.page == PageUp   { ptr = ptr.prev; idx = (idx - 1 + buffer.size) % buffer.size }
    }
    return ret
    
}



