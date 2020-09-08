package facade

import (
	log "../log"
	"fmt"
	"github.com/pborman/ansi"
	"strings"
)

const DEBUG_ANSI = false
const DEBUG_ANSI_DUMP = false

type Line []rune

type TextSeq struct {
	Text []rune
	Seq  *ansi.S
}

const TabWidth = 8

func (item *TextSeq) Desc() string {
	ret := ""
	if item.Seq != nil {
		ret += fmt.Sprintf("ansi 0x%x", item.Seq.Code)
	} else if len(item.Text) > 0 {
		ret += fmt.Sprintf("text %d byte", len(item.Text))
	} else {
		ret += "empty"
	}
	return ret
}

func sendBytes(raw []byte, bufChan chan TextSeq) {
	var item = TextSeq{}
	str := string(raw)
	item.Text = []rune(str)
	bufChan <- item
}

func sendSequence(seq *ansi.S, bufChan chan TextSeq) {
	var item = TextSeq{}
	item.Seq = seq
	bufChan <- item
}

// process raw bytes,
// split into runes and ansi sequences,
// send to channel,
// return leftover bytes

func ProcessRaw(raw []byte, bufChan chan TextSeq) ([]byte, error) {
	var decodeErr error
	var seq *ansi.S

	var ptr []byte = raw
	var rem []byte = raw
	var txt []byte = []byte{} // keep track of non-sequence bytes that might be multibyte characters

	//    if DEBUG_ANSI { log.Debug("process raw %d byte:\n%s",len(raw),log.Dump(raw,0,0)) }

	for ptr != nil && len(ptr) > 0 {

		rem, seq, decodeErr = ansi.Decode(ptr)

		switch decodeErr {

		case nil:
			break

		case ansi.UnknownEscape:
			break // handle below

		case ansi.LoneEscape:
			if DEBUG_ANSI {
				log.Debug("[ansi] lone escape: %s", log.Dump(ptr, 0, 0))
			}
			sendBytes(txt, bufChan)
			return ptr, log.NewError("ansi lone escape")

		case ansi.NoST:
			if DEBUG_ANSI {
				log.Debug("[ansi] missing terminator for sequence 0x%x", seq.Code)
			}

			//look for terminating BEL (xterm) or ST (ansi)
			var tmp []byte
			for tmp = ptr; len(tmp) > 0; tmp = tmp[1:] {
				if tmp[0] == 0x07 { // BEL terminator (xterm)
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
				if DEBUG_ANSI {
					log.Debug("[ansi] found missing terminator")
				}
				continue
			} else { // did not find terminator, return and wait for more bytes
				sendBytes(txt, bufChan)
				if DEBUG_ANSI {
					log.Debug("[ansi] missing terminator, return ptr to pick up more")
				}
				return ptr, log.NewError("ansi missing terminator")
			}

		default:
			if DEBUG_ANSI {
				log.Debug("[ansi] fail decode: %s\n%s", decodeErr, log.Dump(ptr, 0, 0))
			}
			sendBytes(txt, bufChan)
			return ptr, log.NewError("ansi fail decode") //TWEAK? fix 'man man' escape sequence split error

		}

		switch seq.Type {

		case "": // no sequence
			s := seq.String()
			if DEBUG_ANSI_DUMP {
				log.Debug("[ansi] text %d byte:\n%s", len(s), log.Dump([]byte(s), len(s), 0))
			} else if DEBUG_ANSI {
				log.Debug("[ansi] text %d byte", len(s))
			}
			txt = append(txt, []byte(s)...)

		case "C0":
			sendBytes(txt, bufChan)
			txt = []byte{}
			s, ok := lookupSequence(seq.Code)
			if ok {
				if DEBUG_ANSI {
					log.Debug("[ansi] C0 %s %s", s.Desc, s.Name)
				}
				sendSequence(seq, bufChan)
			} else {
				if DEBUG_ANSI {
					log.Debug("[ansi] unknown C0 sequence 0x%x", seq.Code)
				}
			}

		case "C1":
			// The C1 control set has both a two byte and a single byte representation.  The
			// two byte representation is an Escape followed by a byte in the range of 0x40
			// to 0x5f.  They may also be specified by a single byte in the range of 0x80 -
			// 0x9f.
			if ptr[0] >= 0x80 && ptr[0] <= 0x9f {
				if DEBUG_ANSI {
					log.Debug("[ansi] skip probable UTF8 byte 0x%02x", ptr[0])
				}
				txt = append(txt, ptr[0])
			} else {
				sendBytes(txt, bufChan)
				txt = []byte{}
				s, ok := lookupSequence(seq.Code)
				if ok {
					if DEBUG_ANSI {
						log.Debug("[ansi] C1 %s %s(%s)", s.Desc, s.Name, strings.Join(seq.Params, ","))
					}
					sendSequence(seq, bufChan)
				} else {
					if DEBUG_ANSI {
						log.Debug("[ansi] unknown C1 sequence 0x%x", seq.Code)
					}
				}
			}
		case "CSI", "ICF":
			sendBytes(txt, bufChan)
			txt = []byte{}
			s, ok := lookupSequence(seq.Code)
			if ok {
				if DEBUG_ANSI {
					log.Debug("[ansi] %s sequence 0x%x %s '%s'", seq.Type, seq.Code, s.Name, s.Desc)
				}
				sendSequence(seq, bufChan)
			} else {
				if DEBUG_ANSI {
					log.Debug("[ansi] unknown %s sequence 0x%x:\n%s", seq.Type, seq.Code, log.Dump(ptr, len(ptr)-len(rem), 0))
				}
			}

		case "ESC":

			if len(rem) < 2 { // no full sequence, return ptr to pick up more
				sendBytes(txt, bufChan)
				txt = []byte{}
				if DEBUG_ANSI {
					log.Debug("[ansi] short escape sequence, return ptr to pick up more")
				}
				return ptr, log.NewError("ansi short escape sequence")
			}

			s, ok := lookupSequence(seq.Code)
			if ok {

				if DEBUG_ANSI {
					log.Debug("[ansi] ESC sequence 0x%x %s '%s'", seq.Code, s.Name, s.Desc)
				}

			} else {

				switch seq.Code {
				// OpenBSD manpages have Set Character Set sequences
				// https://vt100.net/docs/vt510-rm/SCS.html
				// remove those, including the intermediate byte
				case "\033(", "\033)", "\033*", "\033+", "\033-", "\033.", "\033/":
					if DEBUG_ANSI {
						log.Debug("[ansi] skip SCS sequence 0x%0x plus one byte", seq.Code)
					}
					rem = rem[1:]

				default:
					if DEBUG_ANSI {
						log.Debug("[ansi] unexpected escape sequence 0x%x, ptr %s", seq.Code, log.Dump(ptr, 16, 0))
					}

				}

			}

		default:
			if DEBUG_ANSI {
				log.Debug("[ansi] unknown sequence type %s", seq.Type)
			}
		}

		ptr = rem

	}
	sendBytes(txt, bufChan)

	return []byte{}, nil
}
