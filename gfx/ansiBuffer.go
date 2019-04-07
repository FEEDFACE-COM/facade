

package gfx

import(
    "fmt"
//    "os"
    log "../log"
    "github.com/pborman/ansi"



)


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
        ret[r] = make([]rune, cols)
        for c:=uint(0); c<cols; c++ { 
            ret[r][c] = ' '
        }
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

func (buffer *AnsiBuffer) Consume(chr rune) {
    i,j := buffer.i, buffer.j
    rows,cols := buffer.rows,buffer.cols
    buf := buffer.buf
    switch (chr) {
        
        case '\n':
            j += 1
            i = 0
        
        case '\t':
        case '\r':
        
        default:
            buf[j][i] = chr
            i += 1
            if i == cols {
                i = 0
                j += 1
            }
    }
    
    i %= cols
    j %= rows
    
    
    buffer.i,buffer.j = i,j
}

func (buffer *AnsiBuffer) WriteText(text string) {
    i,j := buffer.i, buffer.j
    rows,cols := buffer.rows,buffer.cols
//    buf := buffer.buf

    for _,c := range(text) {
        log.Debug("write %c",c) 
    }

    i %= cols
    j %= rows
    
    
    buffer.i,buffer.j = i,j
}

func (buffer *AnsiBuffer) Write(raw []byte) {
    var err error 
//    strip, err := ansi.Strip(raw)
//    if err != nil {
//        log.Error("fail ansi strip: %s",err)
//        return
//    }
//    log.Debug("ansi strip: " + string(strip) )

    log.Debug("write %d byte:\n%s",len(raw),log.Dump(raw,0,0))

    var ptr []byte = raw
    var rem []byte = raw
    off := 0
    var seq *ansi.S

    for rem != nil {
        rem,seq,err = ansi.Decode( ptr )
        if err != nil || seq == nil {
            log.Debug("fail ansi decode: %s",err)
            break
        }
    
        if seq.Type == "" {
                s := []byte(seq.String())
                l := len(s)
                log.Debug("text %d byte:\n%s",l,log.Dump(s,0,off))
                buffer.WriteText(string(s))
//                os.Stdout.Write( s )
                off += l
                ptr = rem

        } else if seq.Type == "C1" {
            
            s := ptr[0:2]
            log.Debug("c1:\n%s",log.Dump(s,0,off))
            buffer.WriteText(string(s))
//            os.Stdout.Write(s)
            off += 2
            ptr = ptr[2:]   
                
        } else {
            l := len(ptr) - len(rem)
            s := ptr[:l]
            str := ""
            str += fmt.Sprintf(" %x ",seq.Code)
            
//            name := ansi.Name(
//            str += ansi.S

            str += " "
            for _,v := range(seq.Params) {
                str += string(v)
            }
            log.Debug("ansi %s %s:\n%s",seq.Type,str,log.Dump(s,0,off))
            off += l
            ptr = rem
        }
    }
    
  
}







