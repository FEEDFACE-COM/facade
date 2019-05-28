

package facade

import(
    "github.com/pborman/ansi"
)





// https://www.aivosto.com/articles/control-characters.html
// https://en.wikipedia.org/wiki/ANSI_escape_code
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html

func lookupMode(val string) string {

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


func lookupSequence(code ansi.Name) (*ansi.Sequence,bool) {
    var ret *ansi.Sequence
    var ok bool
    ret,ok = ansi.Table[code]
    if ok {
        return ret,true    
    }
    ret, ok = xtermTable[code]
    if ok {
        return ret,true    
    }
    return &ansi.Sequence{},false
}
    
const (
    DECSTBM = ansi.Name("\033[r") // Set Top and Bottom Margins
)



var xtermTable = map[ansi.Name]*ansi.Sequence {

// https://vt100.net/docs/vt510-rm/DECSTBM.html
    DECSTBM: &ansi.Sequence{     
                Name:      "DECSTBM",
                Desc:      "Set Top and Bottom Margins",
                Type:      ansi.CSI,
                Notation:  "Pn1;Pn2",
                NParam:    2,
                Defaults:  []string{"1","-1"},
                Code:      []byte{'r'},
            },
}

