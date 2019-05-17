

package facade

import(
    "fmt"
//    "os"
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



func (buffer *TermBuffer) scroll() {
    if DEBUG_TERMBUFFER { log.Debug("scroll %s",buffer.Desc() ) }
    for r:=uint(1); r<buffer.max.y; r++ {
        buffer.buf[r] = buffer.buf[r+1]
    }
    buffer.buf[ buffer.max.y ] = makeRow(buffer.max.x)
}


func (buffer *TermBuffer) ProcessRunes(runes []rune) {
    cur := buffer.cursor
    buf := buffer.buf
    max := buffer.max

    cnt := 0

    
    for _,run := range(runes) {
        
        switch (run) {
            
            case '\n':
                if DEBUG_TERMBUFFER { log.Debug("LF %d,%d",cur.x,cur.y) }
                cur.x = 1
                cur.y += 1
                if cur.y > max.y {  // scroll last row
                    cur.y = max.y
                    buffer.scroll()
                } else { //new empty last row
                    buf[ cur.y ] = makeRow(max.x)  //for ps ax output??
                }

            
            
            case '\t':
//                if DEBUG_TERMBUFFER { log.Debug("TAB") }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {

                    buf[cur.y][cur.x] = rune(' ')
                    cnt += 1
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
                        cur.y = max.y
                        buffer.scroll()
                }
            }
                
            
            case '\r':
//                if DEBUG_TERMBUFFER { log.Debug("CR") }
                cur.x = 1
            
            case '\a':
                if DEBUG_TERMBUFFER { log.Debug("BEL") }
            
            case '\b':
//                if DEBUG_TERMBUFFER { log.Debug("BS") }
                cur.x -= 1
                if cur.x <= 1 { cur.x = 1 }

            
            default:
                buf[cur.y][cur.x] = run
                cnt += 1
                
                cur.x += 1
                if cur.x > max.x {
                    cur.x = 1
                    cur.y += 1
                    if cur.y > max.y {
                        cur.y = max.y
                        buffer.scroll()
                    } 
                }

        }
        
    }

    
    buffer.cursor = cur
    
    if cnt > 0 {
        if DEBUG_TERMBUFFER { log.Debug("print %d runes %s",cnt,buffer.Desc()) }
    }
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
            fmt.Sscanf(seq.Params[0],"%d",&x)
            fmt.Sscanf(seq.Params[1],"%d",&y)
            buffer.setCursor(x,y)

        case ansi.Table[ansi.EL]:
            var val uint
            fmt.Sscanf(seq.Params[0],"%d",&val)
            buffer.eraseLine(val)

        case ansi.Table[ansi.SGR]:
            break
            
        default:
            if DEBUG_TERMBUFFER {             
                tmp := ""
                for _,v := range(seq.Params) { 
                    tmp += string(v) + ", "
                }
                log.Debug("sequence unhandled: %s %s(%s)",sequence.Desc,sequence.Name,tmp)
            }
        
    }

}






func (buffer *TermBuffer) Desc() string {
    return fmt.Sprintf("termbuffer[%dx%d]",buffer.cols,buffer.rows)
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
    ret += fmt.Sprintf("cursor %2d,%2d\n",buffer.cursor.x,buffer.cursor.y)
    return ret
}





