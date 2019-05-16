

package gfx


//import(
//    "strings"
//    log "../log"
//    "github.com/pborman/ansi"
//)
//
//
//type Buffer interface {
//    
//    ProcessRunes(runes []rune)
//    ProcessSequence(seq *ansi.Sequence)
//
//}
//
//
//func ProcessBytes(buffer Buffer, raw []byte) {
//    
//    var err error
//    var seq *ansi.S
//
//    var ptr []byte = raw
//    var rem []byte = raw
//    
//    for rem != nil {
//
//        //fixme, need to handle shorter strings somehow
//        if len(ptr) >= 3 {
//            switch string(ptr[:3])  {
//                case "\033(B":
//                    log.Debug("SETUSG0 (skip 3 byte)")
//                    ptr = ptr[3:]
//                    continue
//            }
//            
//        }
//
//        rem,seq,err = ansi.Decode(ptr)
//        if err != nil {
//            log.Error("fail ansi decode: %s\n%s",err,log.Dump(ptr,0,0)) 
//            break    
//        }
//
//        if seq == nil {
//            log.Error("ansi sequence nil")
//            break
//        }
//        
//        switch seq.Type {
//    
//            case "":  // no sequence
//                s := []rune( seq.String() )
//                if DEBUG_ANSI { log.Debug("plain %s",string(s)) }
//                buffer.ProcessRunes(s)
//        
//            case "C0":
//                if DEBUG_ANSI { log.Debug("ansi C0 byte.") }
//
//            case "C1":
//                if DEBUG_ANSI { log.Debug("ansi C1 byte.") }
//            
//            case "CSI", "IF":
//                params := ""
//                for _,v := range(seq.Params) { 
//                    params += string(v) + ", "
//                }
//                params = strings.TrimSpace(params)
//                sequence, ok := ansi.Table[seq.Code]
//                if !ok {
//                    log.Error("ansi %s 0x%x not in table",seq.Type,seq.Code)    
//                } else {
//                    if DEBUG_ANSI { log.Debug("ansi %s %s: %s(%s)",seq.Type,sequence.Desc,sequence.Name,params) }
//                    buffer.ProcessSequence(sequence)
//                }
//            
//            default:
//                log.Error("ansi unknown sequence type %s",seq.Type)
//        }
//
//        ptr = rem
//        
//    }
//
//}
