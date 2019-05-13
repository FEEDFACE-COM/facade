
package gfx

import(
    "fmt"
    "strings"
    log "../log"
)


type Line []rune

type LineBuffer struct {
    count uint
    delim uint // index to first line not on screen
    next uint  // index to add next line 
    buf []*Line
    raw []byte
}


func NewLineBuffer(count,delim uint) *LineBuffer {
    ret := &LineBuffer{}
    if count == 0 { count = 1 }
    if delim >= count { delim = count - 1 }
    ret.count = count
    ret.delim = delim
    ret.buf = make( []*Line, ret.count )
    ret.next = ret.delim
    return ret
}



func (buffer *LineBuffer) dequeueLine() {
    // REM, implement me (triggered by scroll timer)    
}

func (buffer *LineBuffer) queueLine(row Line) {
    if buffer.next >= buffer.count {
        log.Error("buffer overflow!")
        return   
    }
    buffer.buf[buffer.next] = &row
    buffer.next += 1
}

// REM THIS IS BAD AND DIRTY AND NEEDS REWRITING
func (buffer *LineBuffer) ProcessBytes(raw []byte) {
    
    
    //rem, we will need to split bytes by newline and append resulting rows
    //but then also keep remaining bytes around until next time we're called??
    
    str := string(raw)    

    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        row := []rune( line )
        buffer.queueLine( Line(row) )
        
    }

}


func (buffer *LineBuffer) Resize(newCount,newDelim uint) {
    log.Debug("resize %d,%d %s",newCount,newDelim,buffer.Desc())
    if newCount == 0 { return }
    
    newBuf := make( []*Line, newCount )
    if newCount < buffer.count {

//        d := buffer.row - newCount

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newCount && idx<buffer.count; idx++ {
            newBuf[idx] = buffer.buf[idx]
        }

        if buffer.next >= newCount {
            buffer.next = newCount-1
        }
                

    } else if newCount > buffer.count {

        // copy all items
        d := newCount - buffer.count
        for idx:= uint(0); idx<buffer.count; idx++ {
            newBuf[ (idx+d) % newCount ] = buffer.buf[idx]
        } 
        
        buffer.next = buffer.count
    }        
    
    //adjust buffer info
    buffer.count = newCount
    buffer.buf = newBuf
    
    
}

func (buffer *LineBuffer) Desc() string { 
    return fmt.Sprintf("linebuffer[%d]",buffer.count )
}


func (buffer *LineBuffer) Dump(width,height uint) string {
    ret := ""
    for i := uint(0); i<buffer.count;i++ {
        
        ret += fmt.Sprintf(" %02d | ",i)
        
        row := buffer.buf[ i ]
        if row != nil {
            for c:=uint(0); c<width && c<uint(len(*row)); c++ {
                ret += fmt.Sprintf("%c",(*row)[c]) 
            }
        } 
        ret += "\n"
        
            
        if i == buffer.delim-1 {
            ret += " ---+-"
            for c:=uint(0); c<width; c++ { ret += "-" }
            ret += "\n"
        }

    }
    return ret
}


