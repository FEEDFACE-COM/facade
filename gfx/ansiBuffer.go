

package gfx

import(
    "fmt"
//    "os"
    log "../log"
    "github.com/pborman/ansi"



)

const DEBUG_ANSI = true


/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) */ 

type AnsiBuffer struct {
    cols uint
    rows  uint
    buf [][]rune
    i,j uint
    
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




func NewAnsiBuffer(cols, rows uint) *AnsiBuffer {
    ret := &AnsiBuffer{cols:cols, rows:rows}
    ret.buf = makeBuffer(cols,rows)
    return ret
}

func (buffer *AnsiBuffer) Desc() string {
    return fmt.Sprintf("ansi[%dx%d]",buffer.cols,buffer.rows)
}




func (buffer *AnsiBuffer) Resize(cols, rows uint) {
    buf := makeBuffer(cols,rows)
    for r:=uint(0); r<rows && r<buffer.rows; r++ {
        for c:=uint(0); c<cols && c<buffer.cols; c++ {
            if r < buffer.rows && c < buffer.cols {
                buf[r][c] = buffer.buf[r][c]
            }
        }
    }
    buffer.cols, buffer.rows = cols,rows
    buffer.buf = buf 
}


func (buffer *AnsiBuffer) Dump() string {
    ret := ""

    ret += "+ "
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += " +\n"
    for r:=0; r<int(buffer.rows); r++ {
        ret += "| "
        for c:=0; c<int(buffer.cols); c++ {
            ret += fmt.Sprintf("%c",buffer.buf[r][c])
        }
        ret += " |\n"
    }
    ret += "+ "
    for c:=0; c<int(buffer.cols); c++ { ret += "-" }
    ret += " +\n"
    return ret
}


func (buffer *AnsiBuffer) LineForRow(row int) string {
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




func (buffer *AnsiBuffer) writeString(text string) {
    i,j := buffer.i, buffer.j
    rows,cols := buffer.rows,buffer.cols
    buf := buffer.buf

    cnt := 0

    var runes []rune = []rune(text)
    for _,run := range(runes) {
        
        switch (run) {
            case '\n':
                if DEBUG_ANSI { log.Debug("LF") }
                i = 0
                j += 1
                if j >= rows-1 {  // last row
                    j = rows-1
                    //shift all rows up one
                    for r:=uint(0); r<rows-1; r++ {
                        buf[r] = buf[r+1]        
                    }
                    //new empty last row
                    buf[j] = makeRow(cols)


                } else {
                    //new empty last row
                    buf[j] = makeRow(cols)
                }

            
            
            case '\t':
                if DEBUG_ANSI { log.Debug("TAB") }

                TABWIDTH := 8
                for c:=0; c<TABWIDTH ; c++ {


                    buf[j][i] = rune(' ')
                    cnt += 1

                    i += 1

                    if int(i) % TABWIDTH == 0 { //hit tab stop
                        break
                    }



                    if i >= cols {
                        i = 0
                        j += 1
                        if j >= rows {
                             j = rows-1 
                        }
                    }
                    
                }
            
            case '\r':
                if DEBUG_ANSI { log.Debug("CR") }
            
            case '\a':
                if DEBUG_ANSI { log.Debug("BEL") }
            
            case '\b':
                if DEBUG_ANSI { log.Debug("BS") }
                i -= 1
                if i <= 0 { i = 0 }
            
            default:
                if run>=' ' && run<='~' {
//                    log.Debug("0x%04x '%c'",run,run)
                }
                buf[j][i] = run
                cnt += 1
                i += 1
                if i >= cols {
                    j += 1
                    if j >= rows { j = rows-1 }
                }
                i %= cols

        }
        
    }

    i %= cols
    j %= rows
    
    
    buffer.i,buffer.j = i,j
    
    if cnt > 0 {
        if DEBUG_ANSI { log.Debug("print %d runes.",cnt) }
    }
}

func (buffer *AnsiBuffer) Write(raw []byte) {
    var err error 
//    if DEBUG_ANSI { log.Debug("write %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }

    var ptr []byte = raw
    var rem []byte = raw
    off := 0
    var seq *ansi.S
    chr := []byte{}
    str := string(chr)


    for rem != nil {
        rem,seq,err = ansi.Decode( ptr )
        if err != nil {
            if DEBUG_ANSI { log.Debug("fail ansi decode: %s",err) }
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

//            if DEBUG_ANSI { log.Debug("c1 byte:\n%s",log.Dump(s,0,off)) }
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
            if ok {
                if DEBUG_ANSI { log.Debug("ansi %s %-32s: %s( %s)",seq.Type,sequence.Desc,sequence.Name,tmp) }
            } else {    
                if DEBUG_ANSI { log.Debug("ansi %s 0x%x not in table",seq.Type,seq.Code) }
            }
            off += l
            ptr = rem
            
////            foo := ansi.SGR
//            foo := seq.Code
//            bar := ansi.Table[foo]
//            log.Debug(" %s %s %s",bar.Name,bar.Desc,bar.Type)

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







