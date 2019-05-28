

package facade

import(
    "strings"
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_ANSI = true
const DEBUG_ANSI_DUMP = false

type Line []rune

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
    var decodeErr error
    var seq *ansi.S

    var ptr []byte = raw
    var rem []byte = raw
    var txt []byte = []byte{} // keep track of non-sequence bytes that might be multibyte characters
    
//    if DEBUG_ANSI { log.Debug("process raw %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }
    
    
    for ptr != nil && len(ptr) > 0 {


        rem,seq,decodeErr = ansi.Decode(ptr)
            
        switch decodeErr {
        
            case nil:
                break
        
            case ansi.LoneEscape:
                log.Debug("ansi lone escape: %s",log.Dump(ptr,0,0)) 
                sendBytes(txt, bufChan)
                return ptr, log.NewError("ansi lone escape")    
                
            case ansi.UnknownEscape:
                log.Warning("ansi unknown sequence 0x%x",seq.Code)    
                //handle below
        
            case ansi.NoST:
                log.Warning("ansi missing terminator for sequence 0x%x",seq.Code)
                //look for terminating BEL (xterm) or ST (ansi) 
    
                var tmp []byte
                for tmp = ptr; len(tmp) > 0; tmp = tmp[1:] {
                    if tmp[0] == 0x07 {  // BEL terminator (xterm)
                        ptr = tmp[1:]
                        break 
                    }    
                    if tmp[0] == 0x9c { // C1 terminator (ansi)
                        ptr = tmp[1:]
                        break
                    }
                    if len(tmp) > 1 && tmp[0] == 0x1b && tmp[1] == 0x5c { // ESC terminator (ansi)
                        ptr = tmp[2:]
                        break
                    }
                }
                
                if len(tmp) > 0 {
                    if DEBUG_ANSI { log.Debug("ansi found missing terminator") }
                    continue
                } else { // did not find terminator, return and wait for more bytes
                    sendBytes(txt, bufChan)
                    if DEBUG_ANSI { log.Debug("ansi missing terminator, return ptr to pick up more") }
                    return ptr, log.NewError("ansi missing terminator")
                }
                
                
        
            default:
                log.Error("ansi fail decode: %s\n%s",decodeErr,log.Dump(ptr,0,0)) 
                sendBytes(txt, bufChan)
                return rem, log.NewError("ansi fail decode")    
            
        }
            

        
        switch seq.Type {
    
            case "":  // no sequence
                s := seq.String()
                if DEBUG_ANSI_DUMP { log.Debug("ansi text %d byte:\n%s",len(s),log.Dump([]byte(s),len(s),0) ) 
                } else if DEBUG_ANSI { log.Debug("ansi text %d byte",len(s)) }
                txt = append(txt, []byte(s) ... )

            case "C0":
                sendBytes(txt, bufChan)
                txt = []byte{}
                s, ok := lookupSequence(seq.Code)
                if ok {
                    if DEBUG_ANSI { log.Debug("ansi C0 %s %s",s.Desc,s.Name) }
                    sendSequence(seq, bufChan)
                } else {
                    log.Warning("ansi unknown C0 sequence 0x%x",seq.Code)
                }

            case "C1":
                // The C1 control set has both a two byte and a single byte representation.  The
                // two byte representation is an Escape followed by a byte in the range of 0x40
                // to 0x5f.  They may also be specified by a single byte in the range of 0x80 -
                // 0x9f. 
                if ptr[0] >= 0x80 && ptr[0] <= 0x9f {
                    if DEBUG_ANSI { log.Debug("ansi skip probable UTF8 byte 0x%02x",ptr[0]) }
                    txt = append(txt, ptr[0] )
                } else {
                    sendBytes(txt, bufChan)
                    txt = []byte{}
                    s, ok := lookupSequence(seq.Code)
                    if ok {
                        if DEBUG_ANSI { log.Debug("ansi C1 %s %s(%s)",s.Desc,s.Name,strings.Join(seq.Params,",")) }
                        sendSequence(seq, bufChan)
                    } else {
                        log.Warning("ansi unknown C1 sequence 0x%x",seq.Code)
                    }
                }
            case "CSI", "ICF":
                sendBytes(txt, bufChan)
                txt = []byte{}
                s, ok := lookupSequence(seq.Code)
                if ok {
                    if DEBUG_ANSI { log.Debug("ansi %s sequence 0x%x %s '%s'",seq.Type,seq.Code,s.Name,s.Desc) }
                    sendSequence(seq, bufChan)
                } else {
                    log.Warning("ansi unknown %s sequence 0x%x:\n%s",seq.Type,seq.Code,log.Dump(ptr,len(ptr)-len(rem),0))
                }


            case "ESC":

                if len(rem) < 3 { // no full sequence, return ptr to pick up more
                    sendBytes(txt, bufChan)
                    txt = []byte{}
                    if DEBUG_ANSI { log.Debug("ansi short escape sequence, return ptr to pick up more") }
                    return ptr, log.NewError("ansi short escape sequence")
                }

                switch seq.Code {
                    case "\033(":
                        if DEBUG_ANSI { log.Debug("ansi skip escape sequence 0x%0x plus one byte",seq.Code) }
                        rem = rem[1:]


//                    case "\033=":
//                        if DEBUG_ANSI { log.Debug("ansi skip escape sequence 0x%0x plus one byte",seq.Code) }
//                        rem = rem[1:]

                    
                    case "\033]":
                        if DEBUG_ANSI { log.Debug("ansi skip escape sequence 0x%0x plus one byte",seq.Code) }
                        rem = rem[1:]
                        
                    
                    default:
                        if DEBUG_ANSI { log.Debug("ansi unexpected escape sequence 0x%x, ptr %s",seq.Code,log.Dump(ptr,16,0)) }
                        log.Warning("ansi unexpected escape sequence 0x%x",seq.Code)
                    
                }
                                
            
            default:
                log.Warning("ansi unknown sequence type %s",seq.Type)
        }

        ptr = rem
        
    }
    sendBytes(txt, bufChan)
    

    return []byte{}, nil 
}




