
package gfx

import(
    "fmt"
    "strings"
    log "../log"
)


type Row []rune

type LineBuffer struct {
    rows uint
    index uint
    line []*Row
}


func NewLineBuffer(rows uint) *LineBuffer {
    ret := &LineBuffer{}
    if rows == 0 { rows = 1 }
    ret.rows = rows
    ret.index = 0
    ret.line = make( []*Row, ret.rows )
    return ret
}



func (buffer *LineBuffer) dequeueRow() {
    // REM, implement me (triggered by scroll timer)    
}

func (buffer *LineBuffer) queueRow(row Row) {
    if buffer.index >= buffer.rows {
        log.Error("buffer overflow!")
        return   
    }
    buffer.line[buffer.index] = &row
    buffer.index += 1
}

// REM THIS IS BAD AND DIRTY AND NEEDS REWRITING
func (buffer *LineBuffer) WriteBytes(raw []byte) {
    
    
    //rem, we will need to split bytes by newline and append resulting rows
    //but then also keep remaining bytes around until next time we're called??
    
    str := string(raw)    

    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        row := []rune( line )
        buffer.queueRow( Row(row) )
        
    }

}


func (buffer *LineBuffer) Resize(newRows uint) {
    log.Debug("resize %d %s",newRows,buffer.Desc())
    if newRows == 0 { newRows = 1 }
    
    newBuf := make( []*Row, newRows )
    if newRows < buffer.rows {

//        d := buffer.row - newRows

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newRows && idx<buffer.rows; idx++ {
            newBuf[idx] = buffer.line[idx]
        }

        if buffer.index >= newRows {
            buffer.index = newRows-1
        }
                

    } else if newRows > buffer.rows {

        // copy all items
        d := newRows - buffer.rows
        for idx:= uint(0); idx<buffer.rows; idx++ {
            newBuf[ (idx+d) % newRows ] = buffer.line[idx]
        } 
        
        buffer.index = buffer.rows
    }        
    
    //adjust buffer info
    buffer.rows = newRows
    buffer.line = newBuf
    
    
}

func (buffer *LineBuffer) Desc() string { 
    return fmt.Sprintf("linebuffer[%d]",buffer.rows )
}


func (buffer *LineBuffer) Dump(width,height uint) string {
    ret := ""
    for i := uint(0); i<buffer.rows;i++ {
        
        ret += fmt.Sprintf(" %02d | ",i)
        
        row := buffer.line[ i ]
        if row != nil {
            for c:=uint(0); c<width && c<uint(len(*row)); c++ {
                ret += fmt.Sprintf("%c",(*row)[c]) 
            }
        } 
        ret += "\n"
        
            
        if i == height-1 {
            ret += " ---+-"
            for c:=uint(0); c<width; c++ { ret += "-" }
            ret += "\n"
        }

    }
    return ret
}


