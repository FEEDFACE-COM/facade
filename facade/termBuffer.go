

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
    if DEBUG_TERMBUFFER { log.Debug("resize %dx%d %s",cols,rows,buffer.Desc()) }
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
    cur := buffer.cursor
    buf := buffer.buf
    max := buffer.max


    if DEBUG_TERMBUFFER { log.Debug("process %d runes %s",len(runes),buffer.Desc()) }

    
    for _,run := range(runes) {
        
        switch (run) {
            
            case '\n':
                if DEBUG_TERMBUFFER { log.Debug("linefeed %d,%d",cur.x,cur.y) }
                cur.x = 1
                cur.y += 1
                if cur.y > max.y {  // scroll last row
                    if DEBUG_TERMBUFFER { log.Debug("scroll for linefeed.") }
                    buffer.scrollLine()
                    cur.y = max.y
//                } else { //new empty last row
//                    buf[ cur.y ] = makeRow(max.x)  //for ps ax output??
                }

            
            
            case '\t':
                if DEBUG_TERMBUFFER { log.Debug("tabulator %d,%d",cur.x,cur.y) }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {

                    buf[cur.y][cur.x] = rune(' ')
                    cur.x += 1

                    if int(cur.x) % TABWIDTH == 1 { //hit tab stop
                        break
                    }

                    if cur.x > max.x {
                        break
                    }
                }
                if cur.x > max.x {
                    cur.x = 1
                    cur.y += 1
                    if cur.y > max.y {
                        if DEBUG_TERMBUFFER { log.Debug("scroll for tabulator.") }
                        buffer.scrollLine()
                        cur.y = max.y
                }
            }
                
            
            case '\r':
                if DEBUG_TERMBUFFER { log.Debug("carriage return %d,%d",cur.x,cur.y) }
                cur.x = 1
            
            case '\a':
                if DEBUG_TERMBUFFER { log.Debug("bell.") }
            
            case '\b':
                if DEBUG_TERMBUFFER { log.Debug("backspace %d,%d",cur.x,cur.y) }
                cur.x -= 1
                if cur.x <= 1 { cur.x = 1 }

            
            default:
                if DEBUG_TERMBUFFER { log.Debug("rune %c %d,%d",run,cur.x,cur.y) }
                buf[cur.y][cur.x] = run
                cur.x += 1
                
                if cur.x > max.x {
                    cur.x = 1
                    cur.y += 1
                    if cur.y > max.y {
                        if DEBUG_TERMBUFFER { log.Debug("scroll for rune.") }
                        buffer.scrollLine()
                        cur.y = max.y
                    } 
                }

        }
        
    }

    
    buffer.cursor = cur
    
}

func (buffer *TermBuffer) scrollLine() {
    for r:=uint(1); r<buffer.max.y; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[ buffer.max.y ] = makeRow(buffer.max.x)
}





func (buffer *TermBuffer) clear() {
    if DEBUG_TERMBUFFER { log.Debug("clear %s",buffer.Desc()) }
    buffer.buf = makeBuffer(buffer.cols,buffer.rows)
}


func (buffer *TermBuffer) setCursor(x,y uint) {
    if DEBUG_TERMBUFFER { log.Debug("set cursor %d,%d %s",x,y,buffer.Desc()) }
    buffer.cursor = pos{x,y}
}

func (buffer *TermBuffer) eraseLine(val uint) {
    if DEBUG_TERMBUFFER { log.Debug("erase line(%d) %s",val,buffer.Desc()) }
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
    if DEBUG_TERMBUFFER { log.Debug("insert line(%d) %s",val,buffer.Desc()) }
    switch val {
        case 1:
            for r:=uint(buffer.max.y); r>buffer.cursor.y; r-- {
                buffer.buf[r] = buffer.buf[r-1]
            }
            buffer.buf[ buffer.cursor.y ] = makeRow(buffer.max.x)
        default:
            log.Warning("NOT IMPLEMENTED: insert line(%d)",val)
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
            buffer.clear()
            
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
        

        case ansi.Table[ansi.SGR]:
            break
            
        case ansi.Table[ansi.SM]:
            if DEBUG_TERMBUFFER { log.Debug("ignore set mode request") }


        case ansi.Table[ansi.RM]:
            if DEBUG_TERMBUFFER { log.Debug("ignore reset mode request") }

            
        default:
            if DEBUG_TERMBUFFER { log.Debug("sequence unhandled: %s %s(%s)",sequence.Desc,sequence.Name,strings.Join(seq.Params,",")) }
        
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





