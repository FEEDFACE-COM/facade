

package facade

import(
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_BUFFER = true

type BufferItem struct {
    Text []rune
    Seq  *ansi.S
}


func (item *BufferItem) Desc() string {
    ret := ""
    if item.Text != nil && len(item.Text) > 0 {
        ret += "text " + string(item.Text)
    } else if item.Seq != nil {
        ret += "ansi " + string(item.Seq.Code)
    } else {
        ret += "nil"
    }
    return ret
}



func sendBytes(raw []byte, bufChan chan BufferItem) {
    var item = BufferItem{}
    str := string(raw)
    item.Text = []rune(str)
    bufChan <- item
}

func sendSequence(seq *ansi.S, bufChan chan BufferItem) {
    var item = BufferItem{}
    item.Seq = seq
    bufChan <- item    
}


// process raw bytes, 
// split into runes and ansi sequences, 
// send to channel, 
// return leftover bytes

func ProcessRaw(raw []byte, bufChan chan BufferItem) ([]byte, error) {
    var err error
    var seq *ansi.S

    var ptr []byte = raw
    var rem []byte = raw
    
    var tmp []byte = []byte{}
    
    if DEBUG_BUFFER { log.Debug("process raw %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }
    
    for ptr != nil && len(ptr) > 0 {


        rem,seq,err = ansi.Decode(ptr)
        if err != nil {
            
            switch err {
            
                case ansi.LoneEscape:
                    log.Warning("ansi lone escape: %s",log.Dump(ptr,0,0)) 
                    sendBytes(tmp, bufChan)
                    return ptr, log.NewError("ansi lone escape")    
                    
                case ansi.UnknownEscape:
                    log.Warning("ansi unknown %d byte sequence: 0x%x",len(seq.Code),seq.Code)    
            
                default:
                    log.Error("fail ansi decode: %s\n%s",err,log.Dump(ptr,0,0)) 
                    return rem, log.NewError("fail ansi decode")    
                
            }
            
        }

        
        switch seq.Type {
    
            case "":  // no ansi sequence
                s := seq.String()
                if DEBUG_BUFFER { log.Debug("process text %d byte:\n%s",len(s),log.Dump([]byte(s),len(s),0) ) }
                tmp = append(tmp, []byte(s) ... )

            case "C0":
                if DEBUG_BUFFER { log.Debug("process ansi C0 byte: 0x%02x",ptr[0]) }

            case "C1":
                // The C1 control set has both a two byte and a single byte representation.  The
                // two byte representation is an Escape followed by a byte in the range of 0x40
                // to 0x5f.  They may also be specified by a single byte in the range of 0x80 -
                // 0x9f. 
                if DEBUG_BUFFER { log.Debug("process ansi C1 byte: 0x%02x",ptr[0]) }
                if ptr[0] >= 0x80 && ptr[0] <= 0x9f {
                    tmp = append(tmp, ptr[0] )
                }
            case "CSI", "IF":
                sendBytes(tmp, bufChan)
                tmp = []byte{}
                s, ok := ansi.Table[seq.Code]
                if ok {
                    if DEBUG_BUFFER { log.Debug("process ansi: %s %s",seq.Type,s.Name) }
                    sendSequence(seq, bufChan)
                } else {
                    log.Error("ansi %s 0x%x not in table",seq.Type,seq.Code)
                }

            case "ESC":
                log.Debug("got escape, here's 16 byte: %s",log.Dump(ptr,16,0))
                log.Debug("       and remainn 16 byte: %s",log.Dump(rem,16,0))
                switch seq.Code {
                    case "\033(":
                        if len(rem) >= 1 { //no full sequence, return ptr to pick up more
                            rem = rem[1:]
                            log.Debug("vt100 0x%x skip %d byte: 0x%x",len(seq.Code),seq.Code)
                        } else {
                            sendBytes(tmp, bufChan)
                            return ptr, log.NewError("ansi short escape sequence")
                        }
                    
                    default:
                        log.Warning("vt100 unknown %d byte:  0x%x ",len(seq.Code),seq.Code)
                    
                }
                                
            
            default:
                log.Error("ansi unknown sequence type %s",seq.Type)
        }

        ptr = rem
        
    }
    sendBytes(tmp, bufChan)
    

    return []byte{}, nil 
}
