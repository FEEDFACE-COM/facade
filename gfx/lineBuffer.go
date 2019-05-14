
package gfx

import(
    "fmt"
    "strings"
    log "../log"
)


type Line []rune

type LineBuffer struct {
    rows uint // lines on screen, min 1
    off  uint // lines off screen, min 1
    buf []*Line
    timer *Timer
}


func NewLineBuffer(rows,off uint) *LineBuffer {
    if rows == 0 { rows = 1 }
    if off == 0 { off = 1 }
    total := rows + off
    ret := &LineBuffer{}
    ret.rows = rows
    ret.off = off
    ret.buf = make( []*Line, total )
    return ret
}



func (buffer *LineBuffer) dequeueLine() {
    // probably should lock mutex?
    
    head := ""
    if buffer.buf[0] != nil {
        head = string( *buffer.buf[0] )
    }
    log.Debug("dequeue %s %s",buffer.Desc(),head)
    total := buffer.rows + buffer.off
    idx := uint(0)
    for ; idx<total-1; idx++ {
        buffer.buf[idx] = buffer.buf[idx+1]
    }
    buffer.buf[idx] = nil
    
    if buffer.buf[buffer.rows] != nil {
        buffer.scrollOnce( 1.0 ) // REM take from config
    }
    
    
    // REM, schedule refresh
    
}


func (buffer *LineBuffer) GetLine(idx uint) Line {
    // probably should lock mutex?
    if idx > buffer.rows {
        log.Error("no line %d in %s",idx,buffer.Desc())
    }
    return *(buffer.buf[idx])
}

func (buffer *LineBuffer) scrollOnce(duration float64) {
    if buffer.timer != nil {
        log.Error("SCROLLING WITH EXISTING TIMER")
        return    
    }
    buffer.timer = NewTimer( duration, false )
    buffer.timer.Fun = func() {
        UnRegisterTimer(buffer.timer)
        buffer.timer = nil
        buffer.dequeueLine()

    }
    buffer.timer.Start()
}


func (buffer *LineBuffer) queueLine(row Line) {
    // probably should lock mutex?
    log.Debug("queue %s %s",buffer.Desc(),string(row))
    total := buffer.rows + buffer.off

    idx := buffer.rows



    if buffer.buf[idx] == nil { //first offscreen slot available
        
        buffer.buf[idx] = &row
        buffer.scrollOnce( 1.0 ) // REM take from config
        
        
    } else { // first offscreen slot full, find next available
        
     
        for ;idx<total;idx++ {
            if buffer.buf[idx] == nil {
                break
            }    
        }
        
        if idx >= total {
            log.Error("overflow %s",buffer.Desc())
            return
        }
        
        buffer.buf[idx] = &row
        
    }
    
    // REM, schedule refresh

}

// REM THIS IS BAD AND DIRTY AND NEEDS REWRITING
func (buffer *LineBuffer) ProcessBytes(raw []byte) {
    
    
    //rem, we will need to split bytes by newline and append resulting rows
    //but then also keep remaining bytes around until next time we're called??
    
    //rem, also remove ansi colors and most ansi controls
    //and obvs remove any vt2x0 chars
    
    str := string(raw)    

    lines := strings.Split(str, "\n")
    for _,line := range(lines) {
    
        tmp := []rune( line )
        if len(tmp) > 0 {
            buffer.queueLine( Line(tmp) )
        }
        
    }

}


func (buffer *LineBuffer) Resize(newRows,newOff uint) {
    log.Debug("resize %d+%d %s",newRows,newOff,buffer.Desc())

    if newRows == 0 { newRows = 1 }
    if newOff == 0 { newOff = 1 }

    oldTotal := buffer.rows + buffer.off
    newTotal := newRows + newOff
    newBuf := make( []*Line, newTotal )



    if newTotal < oldTotal {

        // copy as many items as fit
        var idx uint = 0
        for ; idx<newTotal && idx<oldTotal; idx++ {
            newBuf[idx] = buffer.buf[idx]
        }

                

    } else if newTotal >= oldTotal {

        // copy all items
        d := newTotal - oldTotal
        for idx:= uint(0); idx<oldTotal; idx++ {
            newBuf[ (idx+d) % newTotal ] = buffer.buf[idx]
        } 
        
    }        
    
    //adjust buffer info
    buffer.rows = newRows
    buffer.off  = newOff
    buffer.buf = newBuf
    
    
}

func (buffer *LineBuffer) Desc() string { 
    return fmt.Sprintf("linebuffer[%d+%d]",buffer.rows,buffer.off )
}


func (buffer *LineBuffer) Dump(width uint) string {
    ret := ""
    for i := uint(0); i<buffer.rows+buffer.off;i++ {
        
        ret += fmt.Sprintf(" %02d | ",i)
        
        line := buffer.buf[ i ]
        if line != nil {
            for c:=uint(0); c<width && c<uint(len(*line)); c++ {
                ret += fmt.Sprintf("%c",(*line)[c]) 
            }
        } 
        ret += "\n"
        
            
        if i == buffer.rows-1 {
            ret += " ---+-"
            for c:=uint(0); c<width; c++ { ret += "-" }
            ret += "\n"
        }

    }
    return ret
}


