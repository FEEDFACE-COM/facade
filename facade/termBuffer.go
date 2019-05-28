

package facade

import(
    "fmt"
//    "os"
//    "strings"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_TERMBUFFER = true
const DEBUG_TERMBUFFER_DUMP = true


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type pos struct {
    x,y uint
}

type region struct {
    from,to uint
}

type TermBuffer struct {
    cols uint    // runes per line
    rows  uint   // lines on screen
    buffer [][]rune  // cols+1 x rows+1
    max pos      // max row / column
    cursor pos    
    
    scroll region
    
    altBuffer *[][]rune
    altCursor *pos
}



func makeBuffer(cols,rows uint) [][]rune {
    ret := make([][]rune, rows+1)
    for r:=uint(0); r<=rows; r++ { 
        ret[r] = makeRow(cols)
    }
    return ret
}

func makeRow(cols uint) []rune {
    ret := make([]rune,cols+1)
    for c:=uint(0); c<=cols; c++ { 
        ret[c] = rune(' ')
    }
    return ret
}




func NewTermBuffer(cols, rows uint) *TermBuffer {
    ret := &TermBuffer{cols:cols, rows:rows}
    ret.max = pos{cols,rows}
    ret.buffer = makeBuffer(cols,rows)
    ret.cursor = pos{1,1}
    ret.scroll = region{1,rows}

    return ret
}


func (buffer *TermBuffer) Fill(fill []string) {

    // lock lock lock

    rows := uint( len(fill) )
    for r := uint(0); r<rows && r < buffer.rows; r++ {
        line := Line( fill[r] )
        cols := uint( len(line) )
        for c := uint(0); c<cols && c < buffer.cols; c++ {
            buffer.buffer[r+1][c+1] = line[c]
        }
    }

    buffer.cursor = pos{1,buffer.max.y}    
    
}


func (buffer *TermBuffer) Resize(cols, rows uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s resize %dx%d",buffer.Desc(),cols,rows) }
    max := pos{cols,rows}
    buf := makeBuffer(cols,rows)
    for r:=uint(1); r<max.y && r<buffer.max.y; r++ {
        for c:=uint(1); c<max.x && c<buffer.max.x; c++ {
            buf[r][c] = buffer.buffer[r][c]
        }
    }
    buffer.cols, buffer.rows = cols,rows
    buffer.buffer = buf
    buffer.max = max
    buffer.cursor = pos{1,1}
    buffer.scroll = region{1,rows}
    
    //throw away alternate buffer/cursor
    buffer.altBuffer = nil
    buffer.altCursor = nil
    
}


func (buffer *TermBuffer) GetCursor() (uint,uint) {
    return buffer.cursor.x-1,buffer.cursor.y-1
}

func (buffer *TermBuffer) GetLine(idx uint) Line {
    // REM probably should lock mutex?
    if idx == buffer.rows {
        return Line{}
    } else if idx >= buffer.rows {
        log.Error("no line %d in %s",idx,buffer.Desc())
        return Line{}
    }
    return buffer.buffer[idx+1][1:]
}







func (buffer *TermBuffer) ProcessRunes(runes []rune) {

    if DEBUG_TERMBUFFER_DUMP {
        log.Debug("%s process %d runes:\n%s",buffer.Desc(),len(runes),log.Dump( []byte(string(runes)), 0,0) )
    } else if DEBUG_TERMBUFFER { 
        log.Debug("%s process %d runes",buffer.Desc(),len(runes)) 
    }

//    tmp := []rune{}
    
    for _,run := range(runes) {
        
        switch (run) {
            
            case '\n':
                if DEBUG_TERMBUFFER { log.Debug("%s linefeed",buffer.Desc()) }
                buffer.cursor.x = 1
                buffer.cursor.y += 1
                if buffer.shouldScroll() {
                    buffer.scrollLine()
                    buffer.cursor.y = buffer.max.y
                }

            
            
            case '\t':
                if DEBUG_TERMBUFFER { log.Debug("%s tabulator",buffer.Desc()) }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {

                    if buffer.cursor.x > buffer.max.x {
                        buffer.cursor.x = 1
                        buffer.cursor.y += 1
                    }
                    
                    if buffer.shouldScroll() {
                        buffer.scrollLine()
                        buffer.cursor.x = 1
                        buffer.cursor.y = buffer.max.y
                    }


                    buffer.buffer[ buffer.cursor.y ][ buffer.cursor.x ] = rune(' ')
                    buffer.cursor.x += 1

                    if int(buffer.cursor.x) % TABWIDTH == 1 { //hit tab stop
                        break
                    }

//                    if buffer.cursor.x > max.x {
//                        break
//                    }
                }
//            }
                
            
            case '\r':
//                if DEBUG_TERMBUFFER { log.Debug("%s carriage return",buffer.Desc()) }
                buffer.cursor.x = 1
            
            case '\a':
                if DEBUG_TERMBUFFER { log.Debug("%s bell.",buffer.Desc()) }
            
            case '\b':
                if DEBUG_TERMBUFFER { log.Debug("%s backspace",buffer.Desc()) }
                buffer.cursor.x -= 1
                if buffer.cursor.x <= 1 { buffer.cursor.x = 1 }

            
            default:
                if buffer.cursor.x > buffer.max.x {
                    buffer.cursor.x = 1
                    buffer.cursor.y += 1
                }
                
                if buffer.shouldScroll() {
                    buffer.scrollLine()
                    buffer.cursor.x = 1
                    buffer.cursor.y = buffer.max.y
                }
            
//                if DEBUG_TERMBUFFER { log.Debug("rune %c %d,%d",run,cur.x,cur.y) }
                buffer.buffer[ buffer.cursor.y ][ buffer.cursor.x ] = run
                buffer.cursor.x += 1
//                tmp = append(tmp, run)
                

        }

//        buffer.trace()
        
    }
//    if DEBUG_TERMBUFFER { log.Debug("%s wrote %d runes:\n%s",buffer.Desc(),len(tmp),log.Dump([]byte(string(tmp)),0,0)) }
    
}



func (buffer *TermBuffer) saveBuffer() {
    if DEBUG_TERMBUFFER { log.Debug("%s save buffer",buffer.Desc()) }
    alt := makeBuffer(buffer.cols,buffer.rows)
    for r:=uint(0); r<=buffer.rows; r++ {
        for c:=uint(0); c<=buffer.cols; c++ {
            alt[r][c] = buffer.buffer[r][c]
        }
    }
    buffer.altBuffer = &alt
}

func (buffer *TermBuffer) restoreBuffer() {
    if buffer.altBuffer == nil {
        log.Warning("%s cannot restore nil buffer",buffer.Desc())
        return
    }
    // rem check for same size!!
    if DEBUG_TERMBUFFER { log.Debug("%s restore buffer",buffer.Desc()) }

    var alt [][]rune = *(buffer.altBuffer)
    for r:=uint(0); r<=buffer.rows; r++ {
        for c:=uint(0); c<=buffer.cols; c++ {
            buffer.buffer[r][c] = alt[r][c]
        }
    }
    buffer.altBuffer = nil
}


func (buffer *TermBuffer) saveCursor() {
    if DEBUG_TERMBUFFER { log.Debug("%s save cursor",buffer.Desc()) }
    alt := pos{ buffer.cursor.x, buffer.cursor.y }
    buffer.altCursor = &alt
}


func (buffer *TermBuffer) restoreCursor() {
    if buffer.altCursor == nil {
        log.Warning("%s cannot restore nil cursor",buffer.Desc())
        return
    }
    if DEBUG_TERMBUFFER { log.Debug("%s restore cursor",buffer.Desc()) }
    buffer.cursor = pos{buffer.altCursor.x, buffer.altCursor.y }
    buffer.altCursor = nil
}


func (buffer *TermBuffer) shouldScroll() bool {
    return buffer.cursor.y > buffer.max.y
}

func (buffer *TermBuffer) scrollLine() {
    if DEBUG_TERMBUFFER { log.Debug("%s scroll",buffer.Desc()) }
    for r:=uint(1); r<buffer.max.y; r++ {
        buffer.buffer[r] = buffer.buffer[r+1]
    }
    buffer.buffer[ buffer.max.y ] = makeRow(buffer.max.x)
}


//func (buffer *TermBuffer) lineFeed() {
//    if DEBUG_TERMBUFFER { log.Debug("%s scroll",buffer.Desc()) }
//    for r:=uint(1); r<buffer.max.y; r++ {
//        buffer.buffer[r] = buffer.buffer[r+1]
//    }
//    buffer.buffer[ buffer.max.y ] = makeRow(buffer.max.x)
//}

func (buffer *TermBuffer) reverseLineFeed() {
    if DEBUG_TERMBUFFER { log.Debug("%s scroll reverse",buffer.Desc()) }
    for r:=uint(buffer.max.y); r>1; r-- {
        buffer.buffer[r] = buffer.buffer[r-1]
    }
    buffer.buffer[ 1 ] = makeRow(buffer.max.x)
}





func (buffer *TermBuffer) clear() {
    if DEBUG_TERMBUFFER { log.Debug("%s clear",buffer.Desc()) }
    buffer.buffer = makeBuffer(buffer.cols,buffer.rows)
}

func (buffer *TermBuffer) erasePage(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s erase page %d",buffer.Desc(),val) }
    switch val {
        case 2:
            buffer.buffer = makeBuffer(buffer.cols,buffer.rows)
            buffer.cursor = pos{1,1}
        default:
            log.Warning("NOT IMPLEMENTED: erase page %d",val)
    }
}

func (buffer *TermBuffer) setScrollRegion(from,to uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s set scroll region %d-%d",buffer.Desc(),from,to) }
    if from <= 0 || to > buffer.max.y { from = 1 }
    if to   <= 0 || to > buffer.max.y {   to = buffer.max.y }
    if from < to {
        buffer.scroll = region{from,to}
    }

    
}


func (buffer *TermBuffer) setCursor(x,y uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s set cursor %d,%d",buffer.Desc(),x,y) }
    buffer.cursor = pos{x,y}
}

func (buffer *TermBuffer) cursorUp(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s cursor up %d",buffer.Desc(),val) }
    if int(buffer.cursor.y) - int(val) >= 1 {
        buffer.cursor.y = buffer.cursor.y - val
    }

}

func (buffer *TermBuffer) cursorRight(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s cursor up %d",buffer.Desc(),val) }
    if buffer.cursor.x + val <= buffer.max.x {
        buffer.cursor.x = buffer.cursor.x + val
    }

}

func (buffer *TermBuffer) deleteLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s delete line %d",buffer.Desc(),val) }
    log.Warning("NOT IMPLEMENTED: delete line %d",val)
}


func (buffer *TermBuffer) deleteCharacter(val uint) {
    log.Warning("NOT IMPLEMENTED: delete character %d",val)
}

func (buffer *TermBuffer) eraseLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s erase line %d",buffer.Desc(),val) }
    switch val {
        case 0:
            for c:=buffer.cursor.x; c<=buffer.max.x; c++ {
                buffer.buffer[buffer.cursor.y][c] = rune(' ')    
            }
        default:
            log.Warning("NOT IMPLEMENTED: erase line %d",val)
    }
}


func (buffer *TermBuffer) insertLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s insert line %d",buffer.Desc(),val) }
    log.Warning("NOT IMPLEMENTED: insert line %d",val)
        
}

func (buffer *TermBuffer) linePositionAbsolute(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s line position absolute %d",buffer.Desc(),val) }
    //FIXME: probably need to implement data mode and others..?
    buffer.cursor.x = 1
    buffer.cursor.y = val
}


func (buffer *TermBuffer) cursorCharacterAbsolute(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s cursor character absolute %d",buffer.Desc(),val) }
    buffer.cursor.x = val
}



func (buffer *TermBuffer) setMode(val string) {
            switch val {
                case "?1049":
                    buffer.saveBuffer()
                    buffer.saveCursor()
                default:
//                    if DEBUG_TERMBUFFER { log.Debug("%s ignore set mode '%s'",buffer.Desc(),lookupMode(val)) }
            }
}

func (buffer *TermBuffer) resetMode(val string) {
            switch val {
                case "?1049":
                    buffer.restoreBuffer()
                    buffer.restoreCursor()
                default:
//                    if DEBUG_TERMBUFFER { log.Debug("%s ignore reset mode '%s'",buffer.Desc(),lookupMode(val)) }

            }
}


func (buffer *TermBuffer) ProcessSequence(seq *ansi.S) {
    // lock mutex?

    sequence, ok := lookupSequence(seq.Code)
    if !ok {
        return
        //unlock mutex tho?
    }

    switch sequence {
        
        case ansi.Table[ansi.ED]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.erasePage(val)
            
        case ansi.Table[ansi.CUP]:
            var x,y uint
            fmt.Sscanf(seq.Params[0],"%d",&y)
            fmt.Sscanf(seq.Params[1],"%d",&x)
            buffer.setCursor(x,y)
            
        case ansi.Table[ansi.CUU]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.cursorUp(val)
            
        case ansi.Table[ansi.CUF]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.cursorRight(val)
            

        case ansi.Table[ansi.EL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.eraseLine(val)

        case ansi.Table[ansi.IL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.insertLine(val)
        
        case ansi.Table[ansi.DL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.deleteLine(val)


        case ansi.Table[ansi.DCH]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.deleteCharacter(val)
        

        case ansi.Table[ansi.VPA]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.linePositionAbsolute(val)


        case ansi.Table[ansi.CHA]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.cursorCharacterAbsolute(val)
            
        case ansi.Table[ansi.RI]:
            buffer.reverseLineFeed()
            
        case ansi.Table[ansi.SM]:
            var val string = seq.Params[0]
            buffer.setMode(val)

        case ansi.Table[ansi.RM]:
            var val string = seq.Params[0]
            buffer.resetMode(val)

        case ansi.Table[ansi.SGR]:
            break
            
        
        case Table[DECSTBM]:
            var f,t uint
            fmt.Sscanf(seq.Params[0],"%d",&f)
            fmt.Sscanf(seq.Params[1],"%d",&t)
            buffer.setScrollRegion(f,t)
            
        default:
//            if true && DEBUG_TERMBUFFER { log.Debug("%s unhandled sequence 0x%x",buffer.Desc(),seq.Code) }
            if DEBUG_TERMBUFFER { log.Debug("%s unhandled sequence %s '%s'",buffer.Desc(),sequence.Desc,sequence.Name) }
//            if false && DEBUG_TERMBUFFER { log.Debug("%s unhandled sequence %s(%s) %s",buffer.Desc(),sequence.Name,strings.Join(seq.Params,","),sequence.Desc) }
        
    }
    
//    buffer.trace()

}


//func (buffer *TermBuffer) trace() {
//    os.Stdout.Write( []byte(buffer.Dump()) )
//    os.Stdout.Write( []byte("\n") )
//    os.Stdout.Sync()
//}



func (buffer *TermBuffer) Desc() string {
    alt := ""
    if buffer.altBuffer != nil || buffer.altCursor != nil {
        alt =  " alt"
    }
    scr := ""
    if buffer.scroll.from != 1 || buffer.scroll.to != buffer.max.y {
        scr = fmt.Sprintf(" %d-%d",buffer.scroll.from,buffer.scroll.to)
    }
    return fmt.Sprintf("termbuffer[%2dx%-2d %2d,%-2d%s%s]",buffer.cols,buffer.rows,buffer.cursor.x,buffer.cursor.y,alt,scr)
}

func (buffer *TermBuffer) Dump() string {
    ret := ""
    ret += "+"
    for c:=uint(1); c<=buffer.max.x; c++ { ret += "-" }
    ret += "+\n"
    for r:=uint(1); r<=buffer.max.y; r++ {
        ret += "|"
        for c:=uint(1); c<=buffer.max.x; c++ {
            if c == buffer.cursor.x && r == buffer.cursor.y { ret += "\033[7m" }
            ret += fmt.Sprintf("%c",buffer.buffer[r][c])
            if c == buffer.cursor.x && r == buffer.cursor.y { ret += "\033[27m" }
        }
        ret += "| "
        if r % 10 == 0 {
            ret += fmt.Sprintf("%01d",(r/10)%10)
        } else {
            ret += " "
        }
        ret += fmt.Sprintf("%01d\n",r%10)
    }
    ret += "+"
    for c:=uint(1); c<=buffer.max.x; c++ { ret += "-" }
    ret += "+\n "
    for c:=uint(1); c<=buffer.max.x; c++ { 
        if c % 10 == 0 { ret += fmt.Sprintf("%01d",(c/10)%10) 
        } else { ret += " " }
    }
    ret += "\n "
    for c:=uint(1); c<=buffer.max.x; c++ { ret += fmt.Sprintf("%01d",c%10) }
    ret += "\n"
    ret += fmt.Sprintf("cursor %d,%d\n",buffer.cursor.x,buffer.cursor.y)
    return ret
}





