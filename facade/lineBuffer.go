
package facade

import(
    "fmt"
    "strings"
    gfx "../gfx"
    math "../math32"
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_LINEBUFFER = false
const DEBUG_LINEBUFFER_DUMP = false


type Mark int
const (
    LOWER Mark = -1.0
    LEVEL Mark =  0.0
    UPPER Mark = +1.0
)
    


type LineBuffer struct {
    rows uint // lines on screen, min 1
    offs  uint // lines off screen, min 0
    buf []*Line
    rem []rune
    
    timer *gfx.Timer
    speed float32
    Adaptive bool 
    Drop bool
    Smooth bool

    refreshChan chan bool

    //metering
    mark Mark
    meterBuffer *gfx.RB
    meterTimer *gfx.Timer
    meterTimestamp float32
    
}

const METER_SAMPLES  = 10   
const METER_INTERVAL = 1.


func NewLineBuffer(rows,offs uint, refreshChan chan bool) *LineBuffer {
    if rows == 0 { rows = 1 }
    total := rows + offs
    ret := &LineBuffer{
            speed: float32(GridDefaults.Speed), 
            Adaptive: GridDefaults.Adaptive, 
            Drop: GridDefaults.Drop, 
            Smooth: GridDefaults.Smooth,
    }
    ret.rows = rows
    ret.offs = offs
    ret.buf = make( []*Line, total )
    ret.rem = []rune{}
    ret.refreshChan = refreshChan

    ret.mark = LOWER
    ret.meterBuffer = gfx.NewRB( METER_SAMPLES )
    for i:=0;i<METER_SAMPLES;i++ { ret.meterBuffer.Add(1.0) }
    ret.meterTimestamp = gfx.Now()
    ret.meterTimer = gfx.WorldClock().NewTimer(METER_INTERVAL, true, nil, func() {
        if ret.meterTimestamp + METER_INTERVAL < gfx.Now() { // no new lines since last check
//            log.Debug("%s no line since one sec",ret.Desc())
            ret.meterBuffer.Add( 1.0 )
        }        
    })



    if DEBUG_LINEBUFFER { log.Debug("%s created",ret.Desc()) }
    return ret
}


func (buffer *LineBuffer) GetLine(idx uint) Line {
    // REM probably should lock mutex?
    
    if buffer.offs == 0 && idx >= buffer.rows {
        return Line{}
    }
    
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
        return buffer.timer.Value()
    }
    
    return float32(0.0)  
}

func (buffer *LineBuffer) dequeueLine(scrollNext bool) {
    // probably should lock mutex?
    
    total := buffer.rows + buffer.offs
    idx := uint(0)
    for ; idx<total-1; idx++ {
        buffer.buf[idx] = buffer.buf[idx+1]
    }
    buffer.buf[idx] = nil
    
    if scrollNext {
        if buffer.offs > 0 && buffer.buf[buffer.rows] != nil {
            buffer.scrollOnce(false)
        }
    }    
        
    select { case buffer.refreshChan <- true: ; default: ; }
    
}


func (buffer *LineBuffer) scrollOnce(freshLine bool) {
    if buffer.timer != nil {
        log.Error("%s refuse scroll with existing %s",buffer.Desc(),buffer.timer.Desc())
        return    
    }
    
    //most lines are scrolled ease in / ease out
    valueFun := math.EaseInEaseOut
    speed := float32(buffer.speed)    
    
    tmp := ""


    if freshLine  {
        valueFun = math.EaseInEaseOut
        tmp += " ease"
    }

    if !freshLine  { 

        if buffer.Smooth {
            valueFun = math.Identity
            tmp += " smooth"
        }
        if buffer.Adaptive {
            speed = buffer.adaptedSpeed()
            tmp += " adapted"
        }
    }
    
    
    if DEBUG_LINEBUFFER { log.Debug("%s scroll %.2f%s",buffer.Desc(),speed,tmp) }


    triggerFun := func() {
        buffer.timer = nil
        buffer.dequeueLine(true)
    }

    buffer.timer = gfx.WorldClock().NewTimer( speed, false, valueFun, triggerFun)

}



func (buffer *LineBuffer) pushLine(row Line) {
    // lock lock lock
    

    total := buffer.rows + buffer.offs // buffer.offs should be zero
    
    r := uint(0)
    for ; r < total-1; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[r] = &row
    if DEBUG_LINEBUFFER { log.Debug("%s push #%d",buffer.Desc(),r) }

    select { case buffer.refreshChan <- true: ; default: ; }

}




func (buffer *LineBuffer) queueLine(row Line) {

    // measure speed
    metered := gfx.Now()-buffer.meterTimestamp
    buffer.meterBuffer.Add( metered )
    buffer.meterTimestamp = gfx.Now()
    
    
    // no buffering, just push line    
    if buffer.offs == 0 {
        buffer.pushLine(row)
        return
    }

    

    if buffer.buf[buffer.rows] == nil { //first offscreen slot available

        buffer.buf[buffer.rows] = &row
        buffer.scrollOnce(true) 
        
        
    } else { // first offscreen slot full, find next available

        // REM probably should lock mutex?
        total := buffer.rows + buffer.offs
        idx := buffer.rows

     
        for ;idx<total;idx++ {
            if buffer.buf[idx] == nil {
                break
            }    
        }

        
        if idx < total {
            
            buffer.buf[idx] = &row
            
        } else { //all slots full

            
            // restart timer in case its gone somehow
            if buffer.timer == nil {
                log.Debug("%s restart nil timer",buffer.Desc())
                buffer.scrollOnce(false)    
            }            

            
            if buffer.Drop {

                log.Debug("%s overflow !! line dropped !!",buffer.Desc())
                
            } else {

                if DEBUG_LINEBUFFER { log.Debug("%s overflow, line jumped",buffer.Desc()) }
                buffer.pushLine(row)

            }
            
            return
            
        }
        
    }
    
    select { case buffer.refreshChan <- true: ; default: ; }

}



func (buffer *LineBuffer) Clear() {
    // probably should lock mutex?
    if DEBUG_LINEBUFFER { log.Debug("%s clear",buffer.Desc()) }
    buffer.buf = make( []*Line, buffer.rows + buffer.offs)    
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
                for c:=0; c<TabWidth ; c++ {
                    tmp = append(tmp, rune(' '))
                }
                            
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

    // kill timer
    if buffer.timer != nil {
        gfx.WorldClock().DeleteTimer( buffer.timer )
        buffer.timer = nil
        buffer.dequeueLine(false)
    }


    
    rows := uint( len(fill) )
    if DEBUG_LINEBUFFER { log.Debug("%s fill %d lines",buffer.Desc(),rows) }
    
    r := uint(0)
    for  ; r<rows && r < buffer.rows; r++ {

        line := Line( fill[r] )
        buffer.buf[r] = &line
        
    }
    for ; r < buffer.rows+buffer.offs; r++ {
        
//        line := Line( []rune{} )
        buffer.buf[r] = nil//&line
        
    }
    
}


func (buffer *LineBuffer) Resize(newRows,newOffs uint) {
    //lock lock lock
    if DEBUG_LINEBUFFER { log.Debug("%s resize %d+%d",buffer.Desc(),newRows,newOffs) }

    if newRows == 0 { newRows = 1 }

    oldTotal := buffer.rows + buffer.offs
    newTotal := newRows + newOffs
    
    oldBuf := buffer.buf
    newBuf := make( []*Line, newTotal )


    // kill timer
    if buffer.timer != nil {
        gfx.WorldClock().DeleteTimer( buffer.timer )
        buffer.timer = nil
        buffer.dequeueLine(false)
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
    buffer.offs  = newOffs
    buffer.buf = newBuf
    
    
}

func (buffer *LineBuffer) Speed() float32 { return buffer.speed }

func (buffer *LineBuffer) SetSpeed(speed float32) {
    buffer.speed = speed
    
    // kill timer
    if buffer.timer != nil {
        gfx.WorldClock().DeleteTimer( buffer.timer )
        buffer.timer = nil
        buffer.dequeueLine(false)
    }

}

func (buffer *LineBuffer) fillage() float32 { return float32(buffer.buffered()) / float32(buffer.offs) }
func (buffer *LineBuffer) buffered() uint { 
    cnt := uint(0)
    for i:=buffer.rows; i<buffer.rows+buffer.offs; i++ {
        if buffer.buf[i] == nil {
            break
        }
        cnt += 1
        
    }
    return cnt
}


func (buffer *LineBuffer) adaptedSpeed() float32 {

    fillage := buffer.fillage()
    speed := buffer.speed
    average := buffer.meterBuffer.Average()


    const FinalBound = 0.88
    const UpperBound = 0.75
    const CenterMark = 0.50
    const LowerBound = 0.25

    
    const SpeedDelta = 0.2
    const SpeedDouble = 0.5            


    //adjust mark
    if buffer.mark == UPPER && fillage <= CenterMark { 

        buffer.mark = LEVEL 

    } else if buffer.mark == LOWER && fillage >= CenterMark { 
    
        buffer.mark = LEVEL 

    } else if buffer.mark == LEVEL && fillage > UpperBound {
        
        buffer.mark = UPPER
        
    } else if buffer.mark == LEVEL && fillage < LowerBound {
    
        buffer.mark = LOWER

    }

    if average <= buffer.speed /*|| buffer.speed == 0.0*/ {


        if fillage >= FinalBound {

            speed = average * SpeedDouble
            
        } else if fillage > UpperBound || ( fillage > CenterMark && buffer.mark == UPPER ) {
    
            speed = average * (1. - SpeedDelta)
    
        } else if fillage < LowerBound || ( fillage < CenterMark && buffer.mark == LOWER ) {
    
            speed = average * (1. + SpeedDelta)
    
        } else {
            
            speed = average
                
        }

    }    
    

    
    

    
    return speed


}

func (buffer *LineBuffer) Desc() string { 

    ret := "linebuffer["
    
    
    if buffer.offs == 0 {
        ret += fmt.Sprintf("%d ",buffer.rows)
    } else {
        buffered := buffer.buffered()
        fillage := buffer.fillage()
        ret += fmt.Sprintf("%d+%d ",buffer.rows,buffer.offs)
        ret += fmt.Sprintf("%3.0f%%=%d/%d ",100.*fillage,buffered,buffer.offs)
    }

    {
        ret += fmt.Sprintf("spd%.2f ",buffer.speed)
        if buffer.Adaptive { 
            ret += fmt.Sprintf("adp%.2f ",buffer.adaptedSpeed())
        }
        
        ret += fmt.Sprintf("avg%.2f ",buffer.meterBuffer.Average())
    }    
    
    if ! buffer.Drop { ret += "!drop " }
    if ! buffer.Smooth { ret += "!smooth " }

    if DEBUG_LINEBUFFER {
        mp := map[int] string{
                -1 : "lower",
                 0 : "level",
                +1 : "upper",
        }
        ret += " " + mp[int(buffer.mark)]
    }

    ret = strings.TrimSuffix(ret," ")
    ret += "]"
    return ret
}


func (buffer *LineBuffer) Dump(width uint) string {
    ret := ""
    for i := uint(0); i<buffer.rows+buffer.offs;i++ {
        
        
        if i > buffer.rows + 2 && i < buffer.rows+buffer.offs - 2 {
            continue
        }
        
        
        ret += fmt.Sprintf("%3d |",i)
        
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




func (buffer *LineBuffer) GetBuffer() uint64 { return uint64(buffer.offs) }
func (buffer *LineBuffer) GetHeight() uint64 { return uint64(buffer.rows) }


