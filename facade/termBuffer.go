

package facade

import(
    "fmt"
//    "os"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_ANSI = true


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type pos struct {
    x,y uint
}

type TermBuffer struct {
    cols uint
    rows  uint
    buf [][]rune
    cursor pos    
}



func makeBuffer(cols,rows uint) [][]rune {
    ret := make([][]rune, rows)
    for r:=uint(0); r<rows; r++ { 
        ret[r] = makeRow(cols)
    }
    return ret
}

func makeRow(cols uint) []rune {
    ret := make([]rune,cols)
    for c:=uint(0); c<cols; c++ { 
        ret[c] = rune(' ')
    }
    return ret
}




func NewTermBuffer(cols, rows uint) *TermBuffer {
    ret := &TermBuffer{cols:cols, rows:rows}
    ret.buf = makeBuffer(cols,rows)
    ret.cursor = pos{0,0}

    return ret
}

func (buffer *TermBuffer) Desc() string {
    return fmt.Sprintf("termbuffer[%dx%d]",buffer.cols,buffer.rows)
}




func (buffer *TermBuffer) Resize(cols, rows uint) {
    log.Debug("resize %dx%d %s",cols,rows,buffer.Desc())
    buf := makeBuffer(cols,rows)
    for r:=uint(0); r<rows && r<buffer.rows; r++ {
        for c:=uint(0); c<cols && c<buffer.cols; c++ {
            if r < buffer.rows && c < buffer.cols {
                buf[r][c] = buffer.buf[r][c]
            }
        }
    }
    buffer.cols, buffer.rows = cols,rows
    buffer.cursor.x %= cols
    buffer.cursor.y %= rows
    buffer.buf = buf 
}


func (buffer *TermBuffer) Dump() string {
    ret := ""
//    ret += fmt.Sprintf("cursor %02d,%02d\n",buffer.cursor.x,buffer.cursor.y)
    ret += "+"
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += "+\n"
    for r:=uint(0); r<buffer.rows; r++ {
        ret += "|"
        for c:=uint(0); c<buffer.cols; c++ {
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
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += "+\n "
    for c:=0; c<int(buffer.cols); c++ { 
        if c % 10 == 0 { ret += fmt.Sprintf("%01d",(c/10)%10) 
        } else { ret += " " }
    }
    ret += "\n "
    for c:=0; c<int(buffer.cols); c++ { ret += fmt.Sprintf("%01d",c%10) }
    ret += "\n"
    return ret
}


func (buffer *TermBuffer) LineForRow(row int) string {
    r := uint(row)
    if r >= buffer.rows {
        log.Error("line for row %d > rows %d",row,buffer.rows)
        return ""
    }
    
    ret := ""
    for c:=uint(0); c<buffer.cols;c++ {
        chr := buffer.buf[r][c]
        if chr < ' ' || chr >= 0x7f {
            chr = ' '
        }
        ret += fmt.Sprintf("%c",chr)
    }
    
    return ret
    
}




func (buffer *TermBuffer) Scroll() {
    for r:=uint(0); r<buffer.rows-1; r++ {
        buffer.buf[r] = buffer.buf[r+1]        
    }
    buffer.buf[buffer.rows-1] = makeRow(buffer.cols)
}


func (buffer *TermBuffer) processRunes(runes []rune) {
    cur := buffer.cursor
    rows,cols := buffer.rows,buffer.cols
    buf := buffer.buf

    cnt := 0

    
    for _,run := range(runes) {
        
        switch (run) {
            case '\n':
                if DEBUG_ANSI { log.Debug("LF") }
                cur.x = 0
                cur.y += 1
                if cur.y >= rows-1 {  // last row
                    cur.y = rows-1
                    buffer.Scroll()
                } else {
                    //new empty last row
                    buf[ cur.y ] = makeRow(cols)
                }

            
            
            case '\t':
                if DEBUG_ANSI { log.Debug("TAB") }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {


                    buf[cur.y][cur.x] = rune(' ')
                    cnt += 1

                    cur.x += 1

                    if int(cur.x) % TABWIDTH == 0 { //hit tab stop
                        break
                    }



                    if cur.x >= cols {
                        cur.x = 0
                        cur.y += 1
                        if cur.y == rows {
                            cur.y = rows-1
                            buffer.Scroll()
                        }
                    }
                    
                }
            
            case '\r':
                if DEBUG_ANSI { log.Debug("CR") }
            
            case '\a':
                if DEBUG_ANSI { log.Debug("BEL") }
            
            case '\b':
                if DEBUG_ANSI { log.Debug("BS") }
                cur.x -= 1
                if cur.x <= 0 { cur.x = 0 }
            
            default:
                log.Debug("rune %c",run)
                buf[cur.y][cur.x] = run
                cnt += 1
                
                cur.x += 1
                if cur.x >= cols {
                    cur.x = 0
                    cur.y += 1
                    if cur.y == rows {
                        cur.y = rows-1
                        buffer.Scroll()
                    }
                }

        }
        
    }

    cur.x %= cols
    cur.y %= rows
    
    
    buffer.cursor = cur
    
    if cnt > 0 {
        if DEBUG_ANSI { log.Debug("print %d runes.",cnt) }
    }
}





func (buffer *TermBuffer) ProcessBytes(raw []byte) {
    var err error
    var seq *ansi.S

    var ptr []byte = raw
    var rem []byte = raw
    
    var tmp []rune = []rune{}
    var str string = ""
    
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
                str = str + seq.String()
                if DEBUG_ANSI { log.Debug("plain %s",string(s)) }

            case "C0":
                if DEBUG_ANSI { log.Debug("ansi C0 byte.") }
            case "C1":
                if DEBUG_ANSI { log.Debug("ansi C1 byte.") }

            case "CSI", "IF":
            
                buffer.processRunes(tmp)
                tmp = []rune{}

                params := ""
                for _,v := range(seq.Params) { 
                    params += string(v) + ", "
                }
                sequence, ok := ansi.Table[seq.Code]
                if !ok {
                    log.Error("ansi %s 0x%x not in table",seq.Type,seq.Code)    
                } else {
                    if DEBUG_ANSI { log.Debug("ansi %s %s: %s(%s)",seq.Type,sequence.Desc,sequence.Name,params) }
                    buffer.processSequence(sequence)
                }
            
            default:
                log.Error("ansi unknown sequence type %s",seq.Type)
        }

        ptr = rem
        
    }
    
    buffer.processRunes(tmp)
    
}





func (buffer *TermBuffer) processSequence(seq *ansi.Sequence) {

    cur := buffer.cursor
    cols,rows := buffer.cols,buffer.rows
    buf := buffer.buf
    
    switch seq {
        
        case ansi.Table[ansi.ED]: 
            if DEBUG_ANSI { log.Debug("ED") }
            cur = pos{0,0}
            for r:=uint(0); r<rows; r++ {
                buf[r] = makeRow(cols)
            }
        
    }

    cur.x %= cols
    cur.y %= rows
    
    buffer.cursor = cur
    
}





