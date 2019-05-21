

package facade

import(
    "fmt"
//    "os"
    "strings"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_TERMBUFFER = true


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type pos struct {
    x,y uint
}

type TermBuffer struct {
    cols uint    // lines on screen
    rows  uint   // runes per line
    max pos      // max row / column
    buf [][]rune  // cols+1 x rows+1
    cursor pos    
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
    ret.buf = makeBuffer(cols,rows)
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
        
            buffer.buf[r+1][c+1] = line[c]
            
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
            buf[r][c] = buffer.buf[r][c]
        }
    }
    buffer.cols, buffer.rows = cols,rows
    buffer.buf = buf
    buffer.max = max
    buffer.cursor = pos{1,1}
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
    return buffer.buf[idx+1][1:]
}







func (buffer *TermBuffer) ProcessRunes(runes []rune) {
    buf := buffer.buf
    max := buffer.max


    if DEBUG_TERMBUFFER { log.Debug("%s process %d runes",buffer.Desc(),len(runes)) }

    
    for _,run := range(runes) {
        
        switch (run) {
            
            case '\n':
//                if DEBUG_TERMBUFFER { log.Debug("linefeed %d,%d",cur.x,cur.y) }
                buffer.cursor.x = 1
                buffer.cursor.y += 1
                if buffer.cursor.y > max.y {  // scroll last row
                    if DEBUG_TERMBUFFER { log.Debug("%s linefeed",buffer.Desc()) }
                    buffer.lineFeed()
                    buffer.cursor.y = max.y
//                } else { //new empty last row
//                    buf[ buffer.cursor.y ] = makeRow(max.x)  //for ps ax output??
                }

            
            
            case '\t':
//                if DEBUG_TERMBUFFER { log.Debug("tabulator %d,%d",cur.x,cur.y) }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {

                    if buffer.cursor.x > max.x {
                        buffer.cursor.x = 1
                        buffer.cursor.y += 1
                    }
                    if buffer.cursor.y > max.y {
                        if DEBUG_TERMBUFFER { log.Debug("%s linefeed for tabulator",buffer.Desc()) }
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
                if DEBUG_TERMBUFFER { log.Debug("%s carriage return",buffer.Desc()) }
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
                    if DEBUG_TERMBUFFER { log.Debug("%s linefeed for rune",buffer.Desc()) }
                    buffer.lineFeed()
                    buffer.cursor.x = 1
                    buffer.cursor.y = max.y
                } 
//                if DEBUG_TERMBUFFER { log.Debug("rune %c %d,%d",run,cur.x,cur.y) }
                buf[buffer.cursor.y][buffer.cursor.x] = run
                buffer.cursor.x += 1
                

        }
        
    }
    
}

func (buffer *TermBuffer) lineFeed() {
    for r:=uint(1); r<buffer.max.y; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[ buffer.max.y ] = makeRow(buffer.max.x)
}

func (buffer *TermBuffer) reverseLineFeed() {
    for r:=uint(buffer.max.y); r>1; r-- {
        buffer.buf[r] = buffer.buf[r-1]
    }
    buffer.buf[ 1 ] = makeRow(buffer.max.x)
}





func (buffer *TermBuffer) clear() {
    if DEBUG_TERMBUFFER { log.Debug("%s clear",buffer.Desc()) }
    buffer.buf = makeBuffer(buffer.cols,buffer.rows)
}

func (buffer *TermBuffer) erasePage(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s erase page(%d)",buffer.Desc(),val) }
    switch val {
        case 2:
            buffer.buf = makeBuffer(buffer.cols,buffer.rows)
            buffer.cursor = pos{1,1}
        default:
            log.Warning("NOT IMPLEMENTED: erase page(%d)",val)
    }
}


func (buffer *TermBuffer) setCursor(x,y uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s set cursor(%d,%d)",buffer.Desc(),x,y) }
    buffer.cursor = pos{x,y}
}

func (buffer *TermBuffer) eraseLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s erase line(%d)",buffer.Desc(),val) }
    switch val {
        case 0:
            for c:=buffer.cursor.x; c<=buffer.max.x; c++ {
                buffer.buf[buffer.cursor.y][c] = rune(' ')    
            }
        default:
            log.Warning("NOT IMPLEMENTED: erase line(%d)",val)
    }
}


func (buffer *TermBuffer) insertLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s insert line(%d)",buffer.Desc(),val) }
    switch val {
//        case 1:
//            for r:=uint(buffer.max.y); r>buffer.cursor.y; r-- {
//                buffer.buf[r] = buffer.buf[r-1]
//            }
//            buffer.buf[ buffer.cursor.y ] = makeRow(buffer.max.x)
        default:
            log.Warning("NOT IMPLEMENTED: insert line(%d)",val)
    }
        
}

func (buffer *TermBuffer) linePositionAbsolute(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s line position absolute(%d)",buffer.Desc(),val) }
    //FIXME: probably need to implement data mode and others..?
    buffer.cursor.x = 1
    buffer.cursor.y = val
}


func (buffer *TermBuffer) cursorCharacterAbsolute(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("%s cursor character absolute(%d)",buffer.Desc(),val) }
    buffer.cursor.x = val
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

        case ansi.Table[ansi.EL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.eraseLine(val)

        case ansi.Table[ansi.IL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.insertLine(val)
        

        case ansi.Table[ansi.VPA]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.linePositionAbsolute(val)


        case ansi.Table[ansi.CHA]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.cursorCharacterAbsolute(val)

        case ansi.Table[ansi.SGR]:
            break
            
        case ansi.Table[ansi.RI]:
            buffer.reverseLineFeed()
            
        case ansi.Table[ansi.SM]:
            var val int
            fmt.Sscanf(seq.Params[0],"%d",&val)
            if DEBUG_TERMBUFFER { log.Debug("%s set mode(%d) IGNORED",buffer.Desc(),val) }


        case ansi.Table[ansi.RM]:
            var val int
            fmt.Sscanf(seq.Params[0],"%d",&val)
            if DEBUG_TERMBUFFER { log.Debug("%s reset mode(%d) IGNORED",buffer.Desc(),val) }

            
        default:
            if DEBUG_TERMBUFFER { log.Debug("%s unhandled sequence %s(%s) %s",buffer.Desc(),sequence.Name,strings.Join(seq.Params,","),sequence.Desc) }
        
    }

}






func (buffer *TermBuffer) Desc() string {
    return fmt.Sprintf("termbuffer[%dx%d %d,%d]",buffer.cols,buffer.rows,buffer.cursor.x,buffer.cursor.y)
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
            ret += fmt.Sprintf("%c",buffer.buf[r][c])
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





