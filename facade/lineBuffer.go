
package facade

import(
    "fmt"
    "strings"
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
    speed float32
    Adaptive bool 
    Drop bool
    Smooth bool

    refreshChan chan bool
    
    checker *gfx.Timer

    lastTimestamp float32
    average float32
    
    packetCount uint
}


const CHECK_INTERVAL = 1.


func NewLineBuffer(rows,off uint, refreshChan chan bool) *LineBuffer {
    if rows == 0 { rows = 1 }
    total := rows + off
    ret := &LineBuffer{
            speed: float32(GridDefaults.Speed), 
            Adaptive: GridDefaults.Adaptive, 
            Drop: GridDefaults.Drop, 
            Smooth: GridDefaults.Smooth,
    }
    ret.rows = rows
    ret.off = off
    ret.buf = make( []*Line, total )
    ret.rem = []rune{}
    ret.refreshChan = refreshChan


    ret.average = float32(GridDefaults.Speed);
    ret.lastTimestamp = gfx.NOW()
    ret.checker = gfx.NewTimer(CHECK_INTERVAL, true, nil)
    ret.checker.Fun = func() {
        if ret.packetCount > 0 {
            ret.average = CHECK_INTERVAL  / float32(ret.packetCount)
//            if DEBUG_LINEBUFFER { log.Debug("%s checkd %.2f average %d/%.1fs ",ret.Desc(),ret.average,ret.packetCount,CHECK_INTERVAL) }
        }
        ret.packetCount = 0
    }
    
    
    
    ret.packetCount = 0
    if DEBUG_LINEBUFFER { log.Debug("%s created",ret.Desc()) }
    return ret
}


func (buffer *LineBuffer) GetLine(idx uint) Line {
    // REM probably should lock mutex?
    
    if buffer.off == 0 && idx >= buffer.rows {
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
        return buffer.timer.Custom()
    }
    
    return float32(0.0)  
}

func (buffer *LineBuffer) dequeueLine() {
    // probably should lock mutex?
    
//    head := ""
//    if buffer.buf[0] != nil {
//        head = string( *buffer.buf[0] )
//    }
//    if DEBUG_LINEBUFFER { log.Debug("%s dequeue",buffer.Desc()) }
    total := buffer.rows + buffer.off
    idx := uint(0)
    for ; idx<total-1; idx++ {
        buffer.buf[idx] = buffer.buf[idx+1]
    }
    buffer.buf[idx] = nil
    
    if buffer.off > 0 && buffer.buf[buffer.rows] != nil {
//        more := false
//        if buffer.off > 1 && buffer.buf[buffer.rows+1] != nil { ///fillage>0?
//            more = true
//        }
        buffer.scrollOnce(true, nil)
    }
    
        
    select { case buffer.refreshChan <- true: ; default: ; }
    
}


func (buffer *LineBuffer) scrollOnce(fromDequeue bool, withSpeed *float32) {
    if buffer.timer != nil {
        log.Error("%s refuse scroll with existing timer",buffer.Desc())
        return    
    }
    
    //most lines are scrolled ease in / ease out
    custom := math.Identity
    speed := float32(buffer.speed)    
    
    tmp := "speed"

    if fromDequeue && buffer.Smooth  { 

        custom = math.Identity
        tmp = "smooth"
        if buffer.Adaptive {
            speed = buffer.adaptedSpeed()
            tmp += " adapted"
        }
    }
    
    if ! fromDequeue  {

        custom = math.EaseInEaseOut
        tmp = "ease"

//        if buffer.Adaptive {
//            speed = buffer.adaptedSpeed()
//            tmp = " adapted"
//        }
//        if withSpeed != nil && *withSpeed < speed {
//            speed = *withSpeed
//            tmp = "given"
//        }

        
    }
    
    if DEBUG_LINEBUFFER { log.Debug("%s scroll %.2f %s",buffer.Desc(),speed,tmp) }

    buffer.timer = gfx.NewTimer( speed, false, custom )
    buffer.timer.Fun = func() {
        gfx.UnRegisterTimer(buffer.timer)
        buffer.timer = nil
        buffer.dequeueLine()

    }
    buffer.timer.Start()
}



func (buffer *LineBuffer) pushLine(row Line) {
    // lock lock lock
    
//    //dont want timer to mess with new buffer
//    if buffer.timer != nil {
//        gfx.UnRegisterTimer(buffer.timer)
//        buffer.timer = nil
//    }

    total := buffer.rows + buffer.off // buffer.off should be zero
    
    r := uint(0)
    for ; r < total-1; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[r] = &   row
    if DEBUG_LINEBUFFER { log.Debug("%s push #%d",buffer.Desc(),r) }

    select { case buffer.refreshChan <- true: ; default: ; }

}


func (buffer *LineBuffer) queueLine(row Line) {
    
    now := gfx.NOW()
//    buffer.average = now - buffer.lastTimestamp 
//    buffer.lastTimestamp = now
    buffer.packetCount += 1

    
    
    if buffer.off == 0 {
        buffer.pushLine(row)
        return
    }

    
    
    
    // REM probably should lock mutex?
    total := buffer.rows + buffer.off

    idx := buffer.rows

    if buffer.buf[buffer.rows] == nil { //first offscreen slot available

        //change speed preemptively
        buffer.lastTimestamp = now
        
//        if DEBUG_LINEBUFFER { log.Debug("%s first #%d",buffer.Desc(),idx) }

        buffer.buf[idx] = &row
    
//        speed := now - buffer.lastTimestamp
//        buffer.scrollOnce(false,&speed) 

        buffer.scrollOnce(false,nil) 
        
        
    } else { // first offscreen slot full, find next available

        buffer.lastTimestamp = now
     
        for ;idx<total;idx++ {
            if buffer.buf[idx] == nil {
                break
            }    
        }

        
        if idx >= total {
            
            if buffer.Drop {

                log.Debug("%s overflow !! line dropped !!",buffer.Desc())
                return
                
            } else {

                if DEBUG_LINEBUFFER { log.Debug("%s overflow, line jumped",buffer.Desc()) }
                buffer.pushLine(row)
                return
            
            }
        }

//        if DEBUG_LINEBUFFER { log.Debug("%s queue #%d",buffer.Desc(),idx) }
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

func (buffer *LineBuffer) Speed() float32 { return buffer.speed }

func (buffer *LineBuffer) SetSpeed(speed float32) {
    buffer.speed = speed
    
    if buffer.timer != nil {
        buffer.timer.Fun() //execute timer, will unregister itself
    }    
}

func (buffer *LineBuffer) fillage() float32 { return float32(buffer.buffered()) / float32(buffer.off) }
func (buffer *LineBuffer) buffered() uint { 
    cnt := uint(0)
    for i:=buffer.rows; i<buffer.rows+buffer.off; i++ {
        if buffer.buf[i] == nil {
            break
        }
        cnt += 1
        
    }
    return cnt
}


func (buffer *LineBuffer) adaptedSpeed() float32 {

    fillage := buffer.fillage()

    ret := buffer.speed
    if buffer.average <= buffer.speed {
        ret = buffer.average
        
//        if buffer.buf[buffer.rows+buffer.off-1.] != nil {
//            ret /= 2.;
//        }
        
//        log.Debug("fillage %f",fillage)
        if fillage >= 0.75 { //speed up
            ret *= 0.5
        } else if fillage <= 0.25 { //slow down
            ret *= 1.2
        } else {
            ret *= 0.9    
        }
        
//        if fillage > 0.8 {
//            ret *= 1. - fillage
//        }

//        if fillage > 0.125 { ret /= 2.  } 
//        if fillage > 0.25 { ret /= 2.  }
        
        
         
//        if fillage > 0.5 { ret /= 2.  } 
//        if fillage > 0.75 { ret /= 2.  } 
//        if fillage > 0.95 { ret /= 2.  } 

//            log.Debug("using last delta + fillage %f!",ret);

//        ret *= 1. - math.Log(1. - fillage)

        
//        if fillage > .5  {
//            ret /= 2.    
//            log.Debug("using last delta %f HALVED!",ret);
//        } else {
//                    log.Debug("using last delta %f",ret);
//        }
    }
    
    return ret


    
////    pageDuration := float32(buffer.rows) * float32(buffer.speed)
////    bufferDuration := float32(buffered) * float32(buffer.speed)
//    
//    //ratio := pageDuration / bufferDuration
//    if buffered == 0 {
//        return buffer.speed
//    }
//    
//    return buffer.speed * float32( 1. - float32(buffered)/float32(buffer.off) )
//
//
////    if fillage > .25 { speed /= 2. }
////    if fillage > .50 { speed /= 2. }
////    if fillage > .75 { speed /= 2. }
////    if buffered > 8 { speed /= 2. }
////    if buffered > 16 { speed /= 2. }
////    if buffered > 32 { speed /= 2. }
    

}

func (buffer *LineBuffer) Desc() string { 

    ret := "linebuffer["
    
    
    if buffer.off == 0 {
        ret += fmt.Sprintf("%d ",buffer.rows)
    } else {
        buffered := buffer.buffered()
        fillage := buffer.fillage()
        ret += fmt.Sprintf("%d+%d ",buffer.rows,buffer.off)
        ret += fmt.Sprintf("%3.0f%%=%d/%d ",100.*fillage,buffered,buffer.off)
    }

    {
        ret += fmt.Sprintf("spd%.2f ",buffer.speed)
        if buffer.Adaptive { 
            ret += fmt.Sprintf("adp%.2f ",buffer.adaptedSpeed())
        }
        
        
        ret += fmt.Sprintf("avg%.2f ",buffer.average)
        
    }    
    
    if ! buffer.Drop { ret += "!drop " }
    if ! buffer.Smooth { ret += "!smooth " }

    ret = strings.TrimSuffix(ret," ")
    ret += "]"
    return ret
}


func (buffer *LineBuffer) Dump(width uint) string {
    ret := ""
    for i := uint(0); i<buffer.rows+buffer.off;i++ {
        
        
        if i > buffer.rows + 2 && i < buffer.rows+buffer.off - 2 {
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


