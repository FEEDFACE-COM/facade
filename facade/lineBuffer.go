
package facade

import(
    "fmt"
    gfx "../gfx"
    math "../math32"
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_LINEBUFFER = true
const DEBUG_LINEBUFFER_DUMP = false


type LineBuffer struct {
    rows uint // lines on screen, min 1
    off  uint // lines off screen, min 0
    buf []*Line
    rem []rune
    
    timer *gfx.Timer
    Speed float32


    refreshChan chan bool
    
}




func NewLineBuffer(rows,off uint, refreshChan chan bool) *LineBuffer {
    if rows == 0 { rows = 1 }
    total := rows + off
    ret := &LineBuffer{}
    ret.rows = rows
    ret.off = off
    ret.buf = make( []*Line, total )
    ret.rem = []rune{}
    ret.refreshChan = refreshChan
    if DEBUG_LINEBUFFER { log.Debug("%s created",ret.Desc()) }
    return ret
}


func (buffer *LineBuffer) GetLine(idx uint) Line {
    // REM probably should lock mutex?
    if idx == buffer.rows && buffer.buf[idx] == nil {
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

func (buffer *LineBuffer) GetScroller() float32 {
    
    if buffer.timer != nil {
        return buffer.timer.Custom()
    }
    
    return float32(0.0)  
}

func (buffer *LineBuffer) dequeueLine() {
    // probably should lock mutex?
    
    head := ""
    if buffer.buf[0] != nil {
        head = string( *buffer.buf[0] )
    }
    if DEBUG_LINEBUFFER { log.Debug("%s dequeue %s",buffer.Desc(),head) }
    total := buffer.rows + buffer.off
    idx := uint(0)
    for ; idx<total-1; idx++ {
        buffer.buf[idx] = buffer.buf[idx+1]
    }
    buffer.buf[idx] = nil
    
    if off > 0 && buffer.buf[buffer.rows] != nil {
        more := false
        if off > 1 && buffer.buf[buffer.rows+1] != nil {
            more = true
        }
        buffer.scrollOnce(true, more)
    }
    
        
    select { case buffer.refreshChan <- true: ; default: ; }
    
}


func (buffer *LineBuffer) scrollOnce(fromDequeue, moreToCome bool) {
    if buffer.timer != nil {
        log.Error("%s refuse scroll with existing timer",buffer.Desc())
        return    
    }
    
    //most lines are scrolled ease in / ease out
    custom := math.Identity
    
    
    
    if fromDequeue && moreToCome{  
    
        if DEBUG_LINEBUFFER { log.Debug("%s start timer with Identity",buffer.Desc()) }
        custom = math.Identity
    
    
    } else if fromDequeue { 

        if DEBUG_LINEBUFFER { log.Debug("%s start timer with Identity",buffer.Desc()) }
        custom = math.Identity
        
    } else  { 

        if DEBUG_LINEBUFFER { log.Debug("%s start timer with EaseInEaseOut",buffer.Desc()) }
        custom = math.EaseInEaseOut
        
    }

        
    buffer.timer = gfx.NewTimer( float32(buffer.Speed), false, custom )
    buffer.timer.Fun = func() {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
        buffer.dequeueLine()

    }
    buffer.timer.Start()
}


func (buffer *LineBuffer) pushLine(row Line) {
    // lock lock lock
    
    //dont want timer to mess with new buffer
    if buffer.timer != nil {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
    }

    total := buffer.rows + buffer.off // buffer.off should be zero
    
    r := uint(0)
    for ; r < total-1; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[r] = &   row
    if DEBUG_LINEBUFFER { log.Debug("%s push #%d %s",buffer.Desc(),r,string(row)) }

    select { case buffer.refreshChan <- true: ; default: ; }

}


func (buffer *LineBuffer) queueLine(row Line) {
    
    if buffer.off == 0 {
        buffer.pushLine(row)
        return
    }
    
    
    
    // REM probably should lock mutex?
    total := buffer.rows + buffer.off

    idx := buffer.rows

    if buffer.buf[idx] == nil { //first offscreen slot available
        
        if DEBUG_LINEBUFFER { log.Debug("%s next #%d %s",buffer.Desc(),idx,string(row)) }
        buffer.buf[idx] = &row
        buffer.scrollOnce(false) 
        
        
    } else { // first offscreen slot full, find next available
     
        for ;idx<total;idx++ {
            if buffer.buf[idx] == nil {
                break
            }    
        }
        
        if idx >= total {
            log.Debug("%s overflow!!",buffer.Desc())
            return
        }

        if DEBUG_LINEBUFFER { log.Debug("%s queue #%d %s",buffer.Desc(),idx,string(row)) }
        buffer.buf[idx] = &row
        
    }
    
    select { case buffer.refreshChan <- true: ; default: ; }

}



func (buffer *LineBuffer) Clear() {
    // probably should lock mutex?
    if DEBUG_LINEBUFFER { log.Debug("%s clear",buffer.Desc()) }
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
            break

    }    
}

func (buffer *LineBuffer) ProcessRunes(runes []rune) {

    rem := append(buffer.rem, runes ...)
    tmp := []rune{}
    
    
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
//                if DEBUG_LINEBUFFER { log.Debug("BEL") }
            
            case '\b':
                l := len(tmp)
                if l >= 1 {
                    tmp = tmp[0:l-1];
                }
//                if DEBUG_LINEBUFFER { log.Debug("BS") }
            
            default:
                tmp = append(tmp, r)
            
            
        }

        
    }
    buffer.rem = tmp
}

func (buffer *LineBuffer) Fill(fill []string) {
    
    // lock lock lock

    //dont want timer to mess with new buffer
    if buffer.timer != nil {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
    }

    
    rows := uint( len(fill) )
    if DEBUG_LINEBUFFER { log.Debug("%s fill %d lines",buffer.Desc(),rows) }
    
    r := uint(0)
    for  ; r<rows && r < buffer.rows; r++ {

        line := Line( fill[r] )
        buffer.buf[r] = &line
        
    }
    for ; r < buffer.rows+buffer.off; r++ {
        
//        line := Line( []rune{} )
        buffer.buf[r] = nil//&line
        
    }
    
}


func (buffer *LineBuffer) Resize(newRows,newOff uint) {
    //lock lock lock
    if DEBUG_LINEBUFFER { log.Debug("%s resize %d+%d",buffer.Desc(),newRows,newOff) }

    if newRows == 0 { newRows = 1 }

    oldTotal := buffer.rows + buffer.off
    newTotal := newRows + newOff
    
    oldBuf := buffer.buf
    newBuf := make( []*Line, newTotal )


    //dont want timer to mess with new buffer
    if buffer.timer != nil {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
    }




    //start with newest buffered lines in old buffer
    //copy to new buffer starting at newest visible line
    var oidx int = int(oldTotal-1)
    var nidx int = int(newRows-1)
    
    for ; oidx >=0 && nidx >= 0;  {
        
        if oldBuf[oidx] == nil { //skip nil lines
            oidx -= 1
            continue    
        }
        
        newBuf[nidx] = oldBuf[oidx]
        oidx -= 1
        nidx -= 1
        
    }

                
    
    //adjust buffer info
    buffer.rows = newRows
    buffer.off  = newOff
    buffer.buf = newBuf
    
    
}


func (buffer *LineBuffer) Desc() string { 
    tmp := ""
    if buffer.timer != nil {
        tmp = fmt.Sprintf(" %.1f",buffer.timer.Fader())
    }
    return fmt.Sprintf("linebuffer[%d+%d%s]",buffer.rows,buffer.off,tmp )
}


func (buffer *LineBuffer) Dump(width uint) string {
    ret := ""
    for i := uint(0); i<buffer.rows+buffer.off;i++ {
        
        ret += fmt.Sprintf("%2d |",i)
        
        line := buffer.buf[ i ]
        if line != nil {
            for c:=uint(0); c<width && c<uint(len(*line)); c++ {
                ret += fmt.Sprintf("%c",(*line)[c]) 
            }
        } 
        ret += "\n"
        
            
        if i == buffer.rows-1 {
            ret += "---+"
            for c:=uint(0); c<width; c++ { ret += "-" }
            ret += "\n"
        }

    }
    return ret
}


