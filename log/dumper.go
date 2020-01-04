package log

import (
	"fmt"
)

func Dump(in []byte, count, offset int) string {
	off := offset % (4 * 4)
	ret := ""
	left, right := "", ""
	//first line spaces for offset
	for i := 0; i < off; i++ {
		left += "  "
		right += " "
		if (i+1)%4 == 0 {
			left += " "
		} else if (i+1)%(4*4) != 0 {
			left += " "
		}
	}
	for i, s := range in {

		if count != 0 && i >= count {
			break
		}

		if i > 0 && (i+off)%(4*4) == 0 {
			ret += left + "    " + right + "\n"
			left, right = "", ""
		}
		left += fmt.Sprintf("%02x", s)
		if s >= 0x20 && s <= 0x7f {
			right += fmt.Sprintf("%c", s)
		} else {
			right += "."
		}
		if (i+off+1)%(2*4) == 0 {
			left += "  "
			right += " "
		} else if (i+off+1)%4 == 0 {
			left += " "
		} else if (i+off+1)%(4*4) != 0 {
			left += ":"
		}

	}
	//fill up remaining space
	for i := len(left); i < len("00:00:00:00 00:00:00:00  00:00:00:00 00:00:00:00  "); i++ {
		left += " "
	}
	for i := len(right); i < len("XXXXXXXX XXXXXXXX"); i++ {
		right += " "
	}
	ret += left + "    " + right
	return ret
}
