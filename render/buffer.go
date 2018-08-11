
package render

import (
    "fmt"
    "sync"
//    log "../log"
    gfx "../gfx"
)

type Buffer struct {
    head *Line
    tail *Line
    mutex *sync.Mutex
    size uint
    twice bool
}

type Line struct {
    text string
    next *Line
    prev *Line
}

type Text string

    

func NewBuffer(size uint) *Buffer {
    buffer := &Buffer{}
    buffer.mutex = &sync.Mutex{}
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


func (buffer *Buffer) Queue(text Text, fade float64) {
    buffer.mutex.Lock()
    if fade < 0.5 && buffer.size >= 2 && buffer.head.text == "" && buffer.head.next.text == "" {
        buffer.head.text = string(text)
        buffer.twice = true
    } else if buffer.head.text == "" {
        buffer.head.text = string(text)
    } else {
        buffer.addTail(string(text))
        buffer.delHead()
    }
    buffer.mutex.Unlock()
//    log.Debug("queue %s",text)
}

func (buffer *Buffer) Desc() string { return fmt.Sprintf("gridBuffer[%d]",buffer.size) }

func (buffer *Buffer) Debug(dir gfx.PageDirection) string {
    ret := ""
    var ptr *Line
    var idx uint
    if dir == gfx.PageDown { ptr = buffer.head ; idx = 0 }
    if dir == gfx.PageUp   { ptr = buffer.tail ; idx = buffer.size }

    for ptr != nil {
        h,t := " ", " "
        if ptr == buffer.head { h = "h" }
        if ptr == buffer.tail { t = "t" }
        ret = ret + fmt.Sprintf("#%02d %s%s %s\n",idx,h,t,ptr.text)
        if dir == gfx.PageDown { ptr = ptr.next; idx = (idx + 1)   }
        if dir == gfx.PageUp   { ptr = ptr.prev; idx = (idx - 1 + buffer.size) % buffer.size }
    }
    return ret
    
}



