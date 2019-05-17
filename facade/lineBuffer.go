
package facade

import(
    "fmt"
    gfx "../gfx"
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_LINEBUFFER = false


type LineBuffer struct {
    rows uint // lines on screen, min 1
    off  uint // lines off screen, min 1
    buf []*Line
    timer *gfx.Timer
    rem []rune

    refreshChan chan bool
    
}


func NewLineBuffer(rows,off uint, refreshChan chan bool) *LineBuffer {
    if rows == 0 { rows = 1 }
    if off == 0 { off = 1 }
    total := rows + off
    ret := &LineBuffer{}
    ret.rows = rows
    ret.off = off
    ret.buf = make( []*Line, total )
    ret.rem = []rune{}
    ret.refreshChan = refreshChan
    return ret
}


func (buffer *LineBuffer) GetLine(idx uint) Line {
    // REM probably should lock mutex?
    if idx == buffer.rows {
        return Line{}    
    } else if idx > buffer.rows {
        log.Error("no line %d in %s",idx,buffer.Desc())
        return Line{}
    }
    ret := buffer.buf[idx]
    if ret == nil {
        return Line{}
    }
    return *ret
}



func (buffer *LineBuffer) dequeueLine() {
    // probably should lock mutex?
    
    head := ""
    if buffer.buf[0] != nil {
        head = string( *buffer.buf[0] )
    }
    if DEBUG_LINEBUFFER { log.Debug("dequeue %s %s",buffer.Desc(),head) }
    total := buffer.rows + buffer.off
    idx := uint(0)
    for ; idx<total-1; idx++ {
        buffer.buf[idx] = buffer.buf[idx+1]
    }
    buffer.buf[idx] = nil
    
    if buffer.buf[buffer.rows] != nil {
        buffer.scrollOnce( 1.0 ) // REM take from config
    }
    
        
    select { case buffer.refreshChan <- true: ; default: ; }
    
}




func (buffer *LineBuffer) scrollOnce(duration float64) {
    if buffer.timer != nil {
        log.Error("refuse scroll with existing timer")
        return    
    }
    buffer.timer = gfx.NewTimer( float32(duration), false )
    buffer.timer.Fun = func() {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
        buffer.dequeueLine()

    }
    buffer.timer.Start()
}


func (buffer *LineBuffer) queueLine(row Line) {
    // REM probably should lock mutex?
    total := buffer.rows + buffer.off

    idx := buffer.rows

    if buffer.buf[idx] == nil { //first offscreen slot available
        
        if DEBUG_LINEBUFFER { log.Debug("queue #%d %s %s",idx-buffer.rows,buffer.Desc(),string(row)) }
        buffer.buf[idx] = &row
        buffer.scrollOnce( 1.0 ) // REM take from config
        
        
    } else { // first offscreen slot full, find next available
     
        for ;idx<total;idx++ {
            if buffer.buf[idx] == nil {
                break
            }    
        }
        
        if idx >= total {
            log.Warning("overflow %s",buffer.Desc())
            return
        }

        if DEBUG_LINEBUFFER { log.Debug("queue #%d %s %s",idx-buffer.rows,buffer.Desc(),string(row)) }
        buffer.buf[idx] = &row
        
    }
    
    select { case buffer.refreshChan <- true: ; default: ; }

}



func (buffer *LineBuffer) Clear() {
    // probably should lock mutex?
    if DEBUG_LINEBUFFER { log.Debug("clear %s",buffer.Desc()) }
    buffer.buf = make( []*Line, buffer.rows + buffer.off)    
}



func (buffer *LineBuffer) ProcessSequence(seq *ansi.S) {
    //lock mutex?

    sequence, ok := ansi.Table[seq.Code]
    if !ok {
        return
        //unlcok tho?
    }

    switch sequence {

        case ansi.Table[ansi.ED]: 
            buffer.Clear()
            
        default:
            if DEBUG_LINEBUFFER {             
                tmp := ""
                for _,v := range(seq.Params) { 
                    tmp += string(v) + ", "
                }
                log.Debug("sequence unhandled: %s %s(%s)",sequence.Desc,sequence.Name,tmp)
            }

    }    
}

func (buffer *LineBuffer) ProcessRunes(runes []rune) {

    rem := append(buffer.rem, runes ...)
    tmp := []rune{}
//    log.Debug("REM %s",string(rem))
    
    
    for _,r := range(rem) {
        
        switch r {
            
            case '\n':
                buffer.queueLine ( Line(tmp) )
                tmp = []rune{}
            
            case '\t':
//                if DEBUG_LINEBUFFER { log.Debug("TAB") }
            
            case '\r':
//                if DEBUG_LINEBUFFER { log.Debug("CR") }
            
            case '\a':
                if DEBUG_LINEBUFFER { log.Debug("BEL") }
            
            case '\b':
//                if DEBUG_LINEBUFFER { log.Debug("BS") }
            
            default:
                tmp = append(tmp, r)
            
            
        }

        
    }
    buffer.rem = tmp
}




func (buffer *LineBuffer) Resize(newRows,newOff uint) {
    if DEBUG_LINEBUFFER { log.Debug("resize %d+%d %s",newRows,newOff,buffer.Desc()) }

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


