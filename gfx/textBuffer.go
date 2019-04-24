

package gfx

import(
    "fmt"
    "strings"
//    log "../log"
)



type bufferItem struct {
    item *Text
    next *bufferItem
    prev *bufferItem    
}

type TextBuffer struct {
    count uint
    items []*Text
}


func NewTextBuffer(count uint) *TextBuffer {
    ret := &TextBuffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.items = make( []*Text, ret.count )
    return ret
}


func (buffer *TextBuffer) WriteText(text *Text) {
        
}


func (buffer *TextBuffer) WriteString(str string) {
    
    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        text := NewText( line )
        buffer.WriteText( text )
        
    }
    
}


func (buffer *TextBuffer) Dump() string {
    ret := ""
    ret += fmt.Sprintf("text[%d]\n",buffer.count)
    for i := uint(0); i<buffer.count;i++ {
        
        item := buffer.items[ i ]
        if item == nil {
            break
        }
        
        txt := (*item).Desc()
        ret += fmt.Sprintf("#%02d %s\n",i,txt)
    }
    return ret
}


