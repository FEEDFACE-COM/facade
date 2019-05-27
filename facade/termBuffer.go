

package facade

import(
    "fmt"
//    "os"
//    "strings"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_TERMBUFFER = true


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type pos struct {
    x,y uint
}

type TermBuffer struct {
    cols uint    // runes per line
    rows  uint   // lines on screen
    buffer [][]rune  // cols+1 x rows+1
    max pos      // max row / column
    cursor pos    
    
    altbuffer *[][]rune
    altcursor *pos
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
    
    //throw away alternate buffer/cursor
    buffer.altbuffer = nil
    buffer.altcursor = nil
    
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
    buf := buffer.buffer
    max := buffer.max


    if DEBUG_TERMBUFFER { log.Debug("%s process %d runes",buffer.Desc(),len(runes)) }

//    tmp := []rune{}
    
    for _,run := range(runes) {
        
        switch (run) {
            
            case '\n':
//                if DEBUG_TERMBUFFER { log.Debug("%s linefeed",buffer.Desc()) }
                buffer.cursor.x = 1
                buffer.cursor.y += 1
                if buffer.cursor.y > max.y {  // scroll last row
                    if DEBUG_TERMBUFFER { log.Debug("%s linefeed scroll",buffer.Desc()) }
                    buffer.lineFeed()
                    buffer.buffer[buffer.max.y] = makeRow(buffer.max.x)
                    buffer.cursor.y = max.y
//                } else { //new empty last row
//                    buf[ buffer.cursor.y ] = makeRow(max.x)  //for ps ax output??
                }

            
            
            case '\t':
//                if DEBUG_TERMBUFFER { log.Debug("%s tabulator",buffer.Desc()) }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {

                    if buffer.cursor.x > max.x {
                        buffer.cursor.x = 1
                        buffer.cursor.y += 1
                    }
                    if buffer.cursor.y > max.y {
                        if DEBUG_TERMBUFFER { log.Debug("%s tabulator scroll",buffer.Desc()) }
                        buffer.lineFeed()
                        buffer.cursor.x = 1
                        buffer.cursor.y = max.y
                    }

                    buf[buffer.cursor.y][buffer.cursor.x] = rune(' ')
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
                if buffer.cursor.x > max.x {
                    buffer.cursor.x = 1
                    buffer.cursor.y += 1
                }
                if buffer.cursor.y > max.y {
                    if DEBUG_TERMBUFFER { log.Debug("%s rune scroll",buffer.Desc()) }
                    buffer.lineFeed()
                    buffer.cursor.x = 1
                    buffer.cursor.y = max.y
                } 
//                if DEBUG_TERMBUFFER { log.Debug("rune %c %d,%d",run,cur.x,cur.y) }
                buf[buffer.cursor.y][buffer.cursor.x] = run
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
    buffer.altbuffer = &alt
}

func (buffer *TermBuffer) restoreBuffer() {
    if buffer.altbuffer == nil {
        log.Warning("%s cannot restore nil buffer",buffer.Desc())
        return
    }
    // rem check for same size!!
    if DEBUG_TERMBUFFER { log.Debug("%s restore buffer",buffer.Desc()) }

    var alt [][]rune = *(buffer.altbuffer)
    for r:=uint(0); r<=buffer.rows; r++ {
        for c:=uint(0); c<=buffer.cols; c++ {
            buffer.buffer[r][c] = alt[r][c]
        }
    }
    buffer.altbuffer = nil
}


func (buffer *TermBuffer) saveCursor() {
    if DEBUG_TERMBUFFER { log.Debug("%s save cursor",buffer.Desc()) }
    alt := pos{ buffer.cursor.x, buffer.cursor.y }
    buffer.altcursor = &alt
}


func (buffer *TermBuffer) restoreCursor() {
    if buffer.altcursor == nil {
        log.Warning("%s cannot restore nil cursor",buffer.Desc())
        return
    }
    if DEBUG_TERMBUFFER { log.Debug("%s restore cursor",buffer.Desc()) }
    buffer.cursor = pos{buffer.altcursor.x, buffer.altcursor.y }
    buffer.altcursor = nil
}





func (buffer *TermBuffer) lineFeed() {
    if DEBUG_TERMBUFFER { log.Debug("%s linefeed",buffer.Desc()) }
    for r:=uint(1); r<buffer.max.y; r++ {
        buffer.buffer[r] = buffer.buffer[r+1]
    }
    buffer.buffer[ buffer.max.y ] = makeRow(buffer.max.x)
}

func (buffer *TermBuffer) reverseLineFeed() {
    if DEBUG_TERMBUFFER { log.Debug("%s reverse linefeed",buffer.Desc()) }
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
                    if DEBUG_TERMBUFFER { log.Debug("%s ignore set mode '%s'",buffer.Desc(),ansiModeName(val)) }
            }
}

func (buffer *TermBuffer) resetMode(val string) {
            switch val {
                case "?1049":
                    buffer.restoreBuffer()
                    buffer.restoreCursor()
                default:
                    if DEBUG_TERMBUFFER { log.Debug("%s ignore reset mode '%s'",buffer.Desc(),ansiModeName(val)) }

            }
}


func (buffer *TermBuffer) ProcessSequence(seq *ansi.S) {
    // lock mutex?

    sequence, ok := ansi.Table[seq.Code]
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
    if buffer.altbuffer != nil || buffer.altcursor != nil {
        alt =  " alt"
    }
    return fmt.Sprintf("termbuffer[%2dx%-2d %2d,%-2d%s]",buffer.cols,buffer.rows,buffer.cursor.x,buffer.cursor.y,alt)
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





