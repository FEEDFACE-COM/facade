

package facade

import(
    "strings"
    log "../log"
    "github.com/pborman/ansi"
)

const DEBUG_ANSI = false
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
    var err error
    var seq *ansi.S

    var ptr []byte = raw
    var rem []byte = raw
    
    var tmp []byte = []byte{}
    
//    if DEBUG_ANSI { log.Debug("process raw %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }
    
    
    for ptr != nil && len(ptr) > 0 {


        rem,seq,err = ansi.Decode(ptr)
        if err != nil {
            
            switch err {
            
                case ansi.LoneEscape:
                    log.Debug("ansi lone escape: %s",log.Dump(ptr,0,0)) 
                    sendBytes(tmp, bufChan)
                    return ptr, log.NewError("ansi lone escape")    
                    
                case ansi.UnknownEscape:
                    log.Warning("ansi unknown sequence 0x%x",seq.Code)    
            
                default:
                    log.Error("ansi fail decode: %s\n%s",err,log.Dump(ptr,0,0)) 
                    sendBytes(tmp, bufChan)
                    return rem, log.NewError("ansi fail decode")    
                
            }
            
        }

        
        switch seq.Type {
    
            case "":  // no ansi sequence
                s := seq.String()
                if DEBUG_ANSI_DUMP { log.Debug("ansi text %d byte:\n%s",len(s),log.Dump([]byte(s),len(s),0) ) 
                } else if DEBUG_ANSI { log.Debug("ansi text %d byte",len(s)) }
                tmp = append(tmp, []byte(s) ... )

            case "C0":
                sendBytes(tmp, bufChan)
                tmp = []byte{}
                s, ok := ansi.Table[seq.Code]
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
                    tmp = append(tmp, ptr[0] )
                } else {
                    sendBytes(tmp, bufChan)
                    tmp = []byte{}
                    s, ok := ansi.Table[seq.Code]
                    if ok {
                        if DEBUG_ANSI { log.Debug("ansi C1 %s %s(%s)",s.Desc,s.Name,strings.Join(seq.Params,",")) }
                        sendSequence(seq, bufChan)
                    } else {
                        log.Warning("ansi unknown C1 sequence 0x%x",seq.Code)
                    }
                }
            case "CSI", "ICF", "CS":
                sendBytes(tmp, bufChan)
                tmp = []byte{}
                s, ok := ansi.Table[seq.Code]
                if ok {
                    if DEBUG_ANSI { log.Debug("ansi %s %s(%s)",s.Desc,s.Name,strings.Join(seq.Params,",")) }
                    sendSequence(seq, bufChan)
                } else {
                    log.Warning("ansi unknown %s sequence 0x%x:\n%s",seq.Type,seq.Code,log.Dump(ptr,len(ptr)-len(rem),0))
                }

            case "ESC":

                if len(rem) < 3 { // no full sequence, return ptr to pick up more
                    sendBytes(tmp, bufChan)
                    tmp = []byte{}
                    if DEBUG_ANSI { log.Debug("ansi short sequence, return ptr to pick up more") }
                    return ptr, log.NewError("ansi short escape sequence")
                }

                switch seq.Code {
                    case "\033(", "\033=":
                        if DEBUG_ANSI { log.Debug("ansi skip ESC sequence 0x%0x plus one byte",seq.Code) }
                        rem = rem[1:]

                    
                    default:
                        if DEBUG_ANSI { log.Debug("ansi unexpected sequence 0x%x, ptr %s",seq.Code,log.Dump(ptr,16,0)) }
                        log.Warning("ansi unexpected sequence 0x%x",seq.Code)
                    
                }
                                
            
            default:
                log.Warning("ansi unknown sequence type %s",seq.Type)
        }

        ptr = rem
        
    }
    sendBytes(tmp, bufChan)
    

    return []byte{}, nil 
}






func ansiModeName(val string) string {

// https://www.real-world-systems.com/docs/ANSIcode.html
// https://ttssh2.osdn.jp/manual/en/usage/tips/vim.html
// https://chromium.googlesource.com/apps/libapps/+/a5fb83c190aa9d74f4a9bca233dac6be2664e9e9/hterm/doc/ControlSequences.md


    switch val {

        case "1"     : return "GUARDED AREA TRANSFER MODE"
        case "2"     : return "KEYBOARD ACTION MODE"
        case "3"     : return "CONTROL REPRESENTATION MODE"
        case "4"     : return "INSERTION REPLACEMENT MODE"
        case "5"     : return "STATUS REPORT TRANSFER MODE"
        case "6"     : return "ERASURE MODE"
        case "7"     : return "LINE EDITING MODE"
        case "8"     : return "BI-DIRECTIONAL SUPPORT MODE"
        case "9"     : return "DEVICE COMPONENT SELECT MODE"
        case "10"    : return "CHARACTER EDITING MODE"
        case "11"    : return "POSITIONING UNIT MODE"
        case "12"    : return "SEND/RECEIVE MODE"
        case "13"    : return "FORMAT EFFECTOR ACTION MODE"
        case "14"    : return "FORMAT EFFECTOR TRANSFER MODE"
        case "15"    : return "MULTIPLE AREA TRANSFER MODE"
        case "16"    : return "TRANSFER TERMINATION MODE"
        case "17"    : return "SELECTED AREA TRANSFER MODE"
        case "18"    : return "TABULATION STOP MODE"
        case "21"    : return "GRAPHIC RENDITION COMBINATION"
        case "22"    : return "ZERO DEFAULT MODE"

        case "?1"    : return "Application Cursor Keys"
        case "?12"   : return "Start Blinking Cursor"
        case "?25"   : return "Show Cursor"
        case "?2004" : return "Bracketed Paste Mode"
        case "?1049" : return "Use Alternate Screen Buffer / Save cursor as in DECSC"
        
        default:       return "????????"
    }
}







