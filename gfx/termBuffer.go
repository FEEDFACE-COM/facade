

package gfx

import(
    "fmt"
//    "os"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_ANSI = false


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type pos struct {
    x,y uint
}

type TermBuffer struct {
    cols uint
    rows  uint
    buf [][]rune
    cursor pos
//    i,j uint
    
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
    ret += "+ "
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += " +\n"
    for r:=uint(0); r<buffer.rows; r++ {
        ret += "| "
        for c:=uint(0); c<buffer.cols; c++ {
            if c == buffer.cursor.x && r == buffer.cursor.y { ret += "\033[7m" }
            ret += fmt.Sprintf("%c",buffer.buf[r][c])
            if c == buffer.cursor.x && r == buffer.cursor.y { ret += "\033[27m" }
        }
        ret += fmt.Sprintf(" | %02d\n",r)
    }
    ret += "+ "
//    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += " +\n  "
    for c:=0; c<int(buffer.cols); c++ { 
        if c % 10 == 0 { ret += fmt.Sprintf("%01d",(c/10)%10) 
        } else { ret += " " }
    }
    ret += "\n  "
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


func (buffer *TermBuffer) writeString(text string) {
    cur := buffer.cursor
    rows,cols := buffer.rows,buffer.cols
    buf := buffer.buf

    cnt := 0

    var runes []rune = []rune(text)
    
    
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


func (buffer *TermBuffer) WriteBytes(raw []byte) {
    var err error 
//    if DEBUG_ANSI { log.Debug("write %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }

    var ptr []byte = raw
    var rem []byte = raw
    off := 0
    var seq *ansi.S
    chr := []byte{}
    str := string(chr)


    for rem != nil {
        
        if len(ptr) >= 3 { 
        
            switch ( string(ptr[:3]) ) {
             
                case "\033(B":
                    log.Debug("skip 3 byte setusg0")
                    ptr = ptr[3:]
                    continue   
                
            }
            
        }
        
        
        
        rem,seq,err = ansi.Decode( ptr )
        if err != nil {
            if DEBUG_ANSI { 
                log.Debug("fail ansi decode: %s\n%s",err,log.Dump(ptr,0,0)) 
            }
            break
        }
        if seq == nil {
//            log.Debug("ansi decode empty seq")
            break
        }
    
        if seq.Type == "" {
                s := []byte(seq.String())
                l := len(s)
                chr = append(chr, s ...)
                off += l
                ptr = rem

        } else if seq.Type == "C1" {

            s := ptr[0:1]
            chr = append(chr,s ...)

            if DEBUG_ANSI { log.Debug("c1 byte:\n%s",log.Dump(s,0,off)) }
            off += 1
            ptr = ptr[1:]
                
        } else if seq.Type == "C0" {
            
            if DEBUG_ANSI { log.Debug("C0 byte.") }
            
        } else if seq.Type == "CSI" || seq.Type == "IF" {
            
            //flush
            str = string(chr)
            buffer.writeString(str)    
            chr,str = []byte{}, string(chr)
            
            l := len(ptr) - len(rem)
            tmp := ""
            for _,v := range(seq.Params) {
                tmp += string(v) + ", "
            }
            
            sequence,ok := ansi.Table[seq.Code]
            if !ok {
                if DEBUG_ANSI { log.Debug("ansi %s 0x%x not in table",seq.Type,seq.Code) }
                off += l
                ptr = rem
                continue
            }
                
            if DEBUG_ANSI { log.Debug("ansi %s %s: %s( %s)",seq.Type,sequence.Desc,sequence.Name,tmp) }
            buffer.handleSequence( sequence )

            off += l
            ptr = rem

        } else {

            log.Error("unknown sequence type %s",seq.Type)  
            ptr = rem

        }
        
        
        
    }
    //flush
    str = string(chr)
    buffer.writeString(str)    
    chr,str = []byte{}, string(chr)
    
    
  
}



func (buffer *TermBuffer) handleSequence(seq *ansi.Sequence) {

    cur := buffer.cursor
    rows,cols := buffer.rows,buffer.cols
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





