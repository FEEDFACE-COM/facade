

package gfx

import(
    "fmt"
    "strings"
    log "../log"
)




type TextBuffer struct {
    count uint
    index uint
    items []*Text
}


func NewTextBuffer(count uint) *TextBuffer {
    ret := &TextBuffer{}
    if count == 0 { count = 1 }
    ret.count = count
    ret.index = 0
    ret.items = make( []*Text, ret.count )
    return ret
}


func (buffer *TextBuffer) WriteText(text *Text) {
    if buffer.index >= buffer.count {
        log.Error("buffer overflow!")
        return   
    }
    buffer.items[buffer.index] = text
    buffer.index += 1
}


func (buffer *TextBuffer) WriteString(str string) {
    
    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        text := NewText( line )
        buffer.WriteText( text )
        
    }
    
}

func (buffer *TextBuffer) Resize(newCount uint) {
    log.Debug("resize %s -> %d",buffer.Desc(),newCount)
    if newCount == 0 { newCount = 1 }
    
    newItems := make( []*Text, newCount )
    if newCount < buffer.count {

//        d := buffer.count - newCount

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newCount && idx<buffer.count; idx++ {
            newItems[idx] = buffer.items[idx]
        }

        if buffer.index >= newCount {
            buffer.index = newCount-1
        }
                
        //cleanup remaining items
        for j:=idx; j<buffer.count; j++ {
            if buffer.items[j] != nil {
                buffer.items[j].Close()    
            }
        }


    } else if newCount > buffer.count {

        // copy all items
        d := newCount - buffer.count
        for idx:= uint(0); idx<buffer.count; idx++ {
            newItems[ (idx+d) % newCount ] = buffer.items[idx]
        } 
        
        buffer.index = buffer.count
    }        
    
    //adjust buffer info
    buffer.count = newCount
    buffer.items = newItems
    
    
}

func (buffer *TextBuffer) Desc() string { 
    return fmt.Sprintf("textbuffer[%d]",buffer.count )
}


func (buffer *TextBuffer) Dump(width,height uint) string {
    ret := ""
//    ret += fmt.Sprintf("text[%d]\n",buffer.count)
    for i := uint(0); i<buffer.count;i++ {
        
        if i == height {
            for j:=uint(0); j<width; j++ {
                ret += "-"    
            }
            ret += "\n"
        }
        
        item := buffer.items[ i ]
        if item == nil {
            ret += fmt.Sprintf("#%02d\n",i)
            continue
        }
        
        txt := (*item).Desc()
        txt0 := txt
        txt1 := ""
        if l := uint(len(txt0)); l >= width {
            txt0 = txt[:width]
            txt1 = ""//txt[width:]
        }
        ret += fmt.Sprintf("#%02d [%s]%s\n",i,txt0,txt1)
    }
    return ret
}


