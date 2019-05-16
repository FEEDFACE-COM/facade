
package facade

import(
    "fmt"
    log "../log"
    "github.com/pborman/ansi"
)


type Line []rune

type LineBuffer struct {
    rows uint // lines on screen, min 1
    off  uint // lines off screen, min 1
    buf []*Line
    timer *Timer
    rem []rune
}


func NewLineBuffer(rows,off uint) *LineBuffer {
    if rows == 0 { rows = 1 }
    if off == 0 { off = 1 }
    total := rows + off
    ret := &LineBuffer{}
    ret.rows = rows
    ret.off = off
    ret.buf = make( []*Line, total )
    ret.rem = []rune{}
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
    buffer.timer = NewTimer( float32(duration), false )
    buffer.timer.Fun = func() {
        UnRegisterTimer(buffer.timer)
        buffer.timer = nil
        buffer.dequeueLine()

    }
    buffer.timer.Start()
}


func (buffer *LineBuffer) queueLine(row Line) {
    // probably should lock mutex?
    total := buffer.rows + buffer.off

    idx := buffer.rows



    if buffer.buf[idx] == nil { //first offscreen slot available
        
        log.Debug("queue #%d %s %s",idx,buffer.Desc(),string(row))
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

        log.Debug("queue #%d %s %s",idx,buffer.Desc(),string(row))
        buffer.buf[idx] = &row
        
    }
    
    // REM, schedule refresh

}


func (buffer LineBuffer) processRunes(runes []rune) {
    runes = append(buffer.rem, runes ...)
    
    tmp := []rune{}
    
    for _,r := range(runes) {
        
        switch r {
            
            case '\n':
                buffer.queueLine ( Line(tmp) )
                tmp = []rune{}
            
            default:
                tmp = append(tmp, r)
            
            
        }

        
    }
    buffer.rem = tmp
}


func (buffer *LineBuffer) ProcessBytes(raw []byte) {
    var err error
    var seq *ansi.S

    var ptr []byte = raw
    var rem []byte = raw
    
    var tmp []rune = []rune{}
    
    for rem != nil {


        rem,seq,err = ansi.Decode(ptr)
        if err != nil {
            log.Error("fail ansi decode: %s\n%s",err,log.Dump(ptr,0,0)) 
            break    
        }

        if seq == nil {
            log.Error("ansi sequence nil")
            break
        }
        
        switch seq.Type {
    
            case "":  // no sequence
                s := []rune( seq.String() )
                tmp = append(tmp, []rune(s) ... )
                if DEBUG_ANSI { log.Debug("plain %s",string(s)) }

            case "C0":
                if DEBUG_ANSI { log.Debug("ansi C0 byte: '%c' %02x",ptr[0],ptr[0]) }
            case "C1":
                if DEBUG_ANSI { log.Debug("ansi C1 byte: '%c' %02x",ptr[0],ptr[0]) }
                // The C1 control set has both a two byte and a single byte representation.  The
                // two byte representation is an Escape followed by a byte in the range of 0x40
                // to 0x5f.  They may also be specified by a single byte in the range of 0x80 -
                // 0x9f. 
                if ptr[0] >= 0x80 && ptr[0] <= 0x9f {
                    
                }
                log.Debug("codes: %x",seq.Code)
            case "CSI", "IF":
                params := ""
                for _,v := range(seq.Params) { 
                    params += string(v) + ", "
                }
                sequence, ok := ansi.Table[seq.Code]
                if !ok {
                    log.Error("ansi %s 0x%x not in table",seq.Type,seq.Code)    
                } else {
                    if DEBUG_ANSI { log.Debug("ansi %s %s: %s(%s)",seq.Type,sequence.Desc,sequence.Name,params) }
                }
            
            default:
                log.Error("ansi unknown sequence type %s",seq.Type)
        }

        ptr = rem
        
    }
    
    buffer.processRunes(tmp)
    
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


