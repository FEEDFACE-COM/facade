
package gfx

import(
    "fmt"
    "strings"
    log "../log"
)


type Row []rune

type TextBuffer struct {
    rows uint
    index uint
    buf []*Row
}


func NewTextBuffer(rows uint) *TextBuffer {
    ret := &TextBuffer{}
    if rows == 0 { rows = 1 }
    ret.rows = rows
    ret.index = 0
    ret.buf = make( []*Row, ret.rows )
    return ret
}



func (buffer *TextBuffer) dequeueRow() {
    // REM, implement me (triggered by scroll timer)    
}

func (buffer *TextBuffer) queueRow(row Row) {
    if buffer.index >= buffer.rows {
        log.Error("buffer overflow!")
        return   
    }
    buffer.buf[buffer.index] = &row
    buffer.index += 1
}

// REM THIS IS BAD AND DIRTY AND NEEDS REWRITING
func (buffer *TextBuffer) WriteBytes(raw []byte) {
    
    
    //rem, we will need to split bytes by newline and append resulting rows
    //but then also keep remaining bytes around until next time we're called??
    
    str := string(raw)    

    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        row := []rune( line )
        buffer.queueRow( Row(row) )
        
    }

}


func (buffer *TextBuffer) Resize(newRows uint) {
    log.Debug("resize %d %s",newRows,buffer.Desc())
    if newRows == 0 { newRows = 1 }
    
    newBuf := make( []*Row, newRows )
    if newRows < buffer.rows {

//        d := buffer.row - newRows

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newRows && idx<buffer.rows; idx++ {
            newBuf[idx] = buffer.buf[idx]
        }

        if buffer.index >= newRows {
            buffer.index = newRows-1
        }
                

    } else if newRows > buffer.rows {

        // copy all items
        d := newRows - buffer.rows
        for idx:= uint(0); idx<buffer.rows; idx++ {
            newBuf[ (idx+d) % newRows ] = buffer.buf[idx]
        } 
        
        buffer.index = buffer.rows
    }        
    
    //adjust buffer info
    buffer.rows = newRows
    buffer.buf = newBuf
    
    
}

func (buffer *TextBuffer) Desc() string { 
    return fmt.Sprintf("textbuffer[%d]",buffer.rows )
}


func (buffer *TextBuffer) Dump(width,height uint) string {
    ret := ""
//    ret += fmt.Sprintf("text[%d]\n",buffer.rows)
    for i := uint(0); i<buffer.rows;i++ {
        
        
        ret += fmt.Sprintf(" %02d | ",i)
        
        row := buffer.buf[ i ]
        if row != nil {
            for c:=uint(0); c<width && c<uint(len(*row)); c++ {
                ret += fmt.Sprintf("%c",(*row)[c]) 
            }
//            txt := (*item).Desc()
//            ret += fmt.Sprintf("%s",txt[:width])
        } 
        ret += "\n"
        
            
        

        if i == height-1 {
            ret += " ---+-"
            for c:=uint(0); c<width; c++ { ret += "-" }
            ret += "\n"
//            continue
        }
        

//
//        txt0 := txt
//        txt1 := ""
//        if l := uint(len(txt0)); l >= width {
//            txt0 = txt[:width]
//            txt1 = ""//txt[width:]
//        }
//        ret += fmt.Sprintf("#%02d [%s]%s\n",i,txt0,txt1)
    }
    return ret
}


