// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/log"
	"fmt"
	//    "os"
	//    "strings"
	"github.com/pborman/ansi"
)

const DEBUG_TERMBUFFER = false
const DEBUG_TERMBUFFER_DUMP = false

/* An array of rows ( ie arrays of cols ( ie multibyte characters ( ie runes ) ) ) */

type pos struct {
	x, y uint
}

type region struct {
	top, bot uint
}

type TermBuffer struct {
	cols   uint     // runes per line
	rows   uint     // lines on screen
	buffer [][]rune // cols+1 x rows+1
	max    pos      // max row / column
	cursor pos

	scroll region

	altBuffer *[][]rune
	altCursor *pos
}

func makeBuffer(cols, rows uint) [][]rune {
	ret := make([][]rune, rows+1)
	for r := uint(0); r <= rows; r++ {
		ret[r] = makeRow(cols)
	}
	return ret
}

func makeRow(cols uint) []rune {
	ret := make([]rune, cols+1)
	for c := uint(0); c <= cols; c++ {
		ret[c] = rune(' ')
	}
	return ret
}

func NewTermBuffer(cols, rows uint) *TermBuffer {
	if rows == 0 {
		rows = 1
	}
	if cols == 0 {
		cols = 1
	}
	ret := &TermBuffer{cols: cols, rows: rows}
	ret.max = pos{cols, rows}
	ret.buffer = makeBuffer(cols, rows)
	ret.cursor = pos{1, 1}
	ret.scroll = region{1, rows}
	if DEBUG_TERMBUFFER {
		log.Debug("%s created", ret.Desc())
	}
	return ret
}

func (buffer *TermBuffer) Fill(fill []string) {

	// lock lock lock

	rows := uint(len(fill))
	if DEBUG_TERMBUFFER {
		log.Debug("%s fill %d lines", buffer.Desc(), rows)
	}

	for r := uint(0); r < rows && r < buffer.rows; r++ {
		line := Line(fill[r])
		cols := uint(len(line))
		for c := uint(0); c < cols && c < buffer.cols; c++ {
			buffer.buffer[r+1][c+1] = line[c]
		}
		buffer.cursor = pos{1, r + 1}
	}

}

func (buffer *TermBuffer) Resize(cols, rows uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s resize %dx%d", buffer.Desc(), cols, rows)
	}

	if rows == 0 {
		rows = 1
	}
	if cols == 0 {
		cols = 1
	}

	max := pos{cols, rows}
	buf := makeBuffer(cols, rows)
	for r := uint(1); r < max.y && r < buffer.max.y; r++ {
		for c := uint(1); c < max.x && c < buffer.max.x; c++ {
			buf[r][c] = buffer.buffer[r][c]
		}
	}
	buffer.cols, buffer.rows = cols, rows
	buffer.buffer = buf
	buffer.max = max
	buffer.cursor = pos{1, 1}
	buffer.scroll = region{1, rows}

	//throw away alternate buffer/cursor
	buffer.altBuffer = nil
	buffer.altCursor = nil

}

func (buffer *TermBuffer) GetCursor() (uint, uint) {
	return buffer.cursor.x - 1, buffer.cursor.y - 1
}

func (buffer *TermBuffer) GetLine(idx uint) Line {
	// REM probably should lock mutex?
	if idx == buffer.rows {
		return Line{}
	} else if idx >= buffer.rows {
		log.Error("%s no line %d", buffer.Desc(), idx)
		return Line{}
	}
	return buffer.buffer[idx+1][1:]
}

func (buffer *TermBuffer) ProcessRunes(runes []rune) {

	if DEBUG_TERMBUFFER_DUMP {
		log.Debug("%s process %d runes:\n%s", buffer.Desc(), len(runes), log.Dump([]byte(string(runes)), 0, 0))
	} else if DEBUG_TERMBUFFER {
		log.Debug("%s process %d runes", buffer.Desc(), len(runes))
	}

	for _, run := range runes {

		switch run {

		case '\n':
			if DEBUG_TERMBUFFER {
				log.Debug("%s linefeed", buffer.Desc())
			}
			buffer.cursor.y += 1
			// NEWLINE does no CARRIAGE RETURN
			if buffer.shouldScroll() {
				buffer.scrollLine()
				buffer.cursor.y = buffer.max.y
				buffer.cursor.x = 1
			}

		case '\t':
			if DEBUG_TERMBUFFER {
				log.Debug("%s tabulator", buffer.Desc())
			}

			for c := 0; c < TabWidth; c++ {

				// ?TWEAK - checking and updating before writing fixes 'man foo' at width 64
				if buffer.cursor.x > buffer.max.x {
					buffer.cursor.x = 1
					buffer.cursor.y += 1
				}

				if buffer.shouldScroll() {
					buffer.scrollLine()
					buffer.cursor.x = 1
					buffer.cursor.y = buffer.max.y
				}

				buffer.buffer[buffer.cursor.y][buffer.cursor.x] = rune(' ')
				buffer.cursor.x += 1

				if int(buffer.cursor.x)%TabWidth == 1 { //hit tab stop
					break
				}

				//                    if buffer.cursor.x > max.x {
				//                        break
				//                    }
			}
			//            }

		case '\r':
			if DEBUG_TERMBUFFER {
				log.Debug("%s carriage return", buffer.Desc())
			}
			buffer.cursor.x = 1

		case '\a':
			if DEBUG_TERMBUFFER {
				log.Debug("%s bell.", buffer.Desc())
			}

		case '\b':
			if DEBUG_TERMBUFFER {
				log.Debug("%s backspace", buffer.Desc())
			}
			buffer.cursor.x -= 1
			if buffer.cursor.x <= 1 {
				buffer.cursor.x = 1
			}

		default:
			//if DEBUG_TERMBUFFER { log.Debug("%s rune %c",buffer.Desc(),run) }

			// ?TWEAK - checking and updating before writing fixes 'man foo' at width 64
			if buffer.cursor.x > buffer.max.x {
				buffer.cursor.x = 1
				buffer.cursor.y += 1
			}

			if buffer.shouldScroll() {
				buffer.scrollLine()
				buffer.cursor.x = 1
				buffer.cursor.y = buffer.max.y
			}

			buffer.buffer[buffer.cursor.y][buffer.cursor.x] = run
			buffer.cursor.x += 1

		}

	}

}

func (buffer *TermBuffer) saveBuffer() {
	if DEBUG_TERMBUFFER {
		log.Debug("%s save buffer", buffer.Desc())
	}
	alt := makeBuffer(buffer.cols, buffer.rows)
	for r := uint(0); r <= buffer.rows; r++ {
		for c := uint(0); c <= buffer.cols; c++ {
			alt[r][c] = buffer.buffer[r][c]
		}
	}
	buffer.altBuffer = &alt
}

func (buffer *TermBuffer) restoreBuffer() {
	if buffer.altBuffer == nil {
		log.Warning("%s cannot restore nil buffer", buffer.Desc())
		return
	}
	// rem check for same size!!
	if DEBUG_TERMBUFFER {
		log.Debug("%s restore buffer", buffer.Desc())
	}

	var alt [][]rune = *(buffer.altBuffer)
	for r := uint(0); r <= buffer.rows; r++ {
		for c := uint(0); c <= buffer.cols; c++ {
			buffer.buffer[r][c] = alt[r][c]
		}
	}
	buffer.altBuffer = nil
}

func (buffer *TermBuffer) saveCursor() {
	if DEBUG_TERMBUFFER {
		log.Debug("%s save cursor", buffer.Desc())
	}
	alt := pos{buffer.cursor.x, buffer.cursor.y}
	buffer.altCursor = &alt
}

func (buffer *TermBuffer) restoreCursor() {
	if buffer.altCursor == nil {
		log.Warning("%s cannot restore nil cursor", buffer.Desc())
		return
	}
	if DEBUG_TERMBUFFER {
		log.Debug("%s restore cursor", buffer.Desc())
	}
	buffer.cursor = pos{buffer.altCursor.x, buffer.altCursor.y}
	buffer.altCursor = nil
}

func (buffer *TermBuffer) shouldScroll() bool {
	return buffer.cursor.y > buffer.scroll.bot
}

func (buffer *TermBuffer) scrollLine() {
	if DEBUG_TERMBUFFER {
		log.Debug("%s scroll one line", buffer.Desc())
	}
	for r := uint(buffer.scroll.top); r < buffer.scroll.bot; r++ {
		buffer.buffer[r] = buffer.buffer[r+1]
	}
	buffer.buffer[buffer.scroll.bot] = makeRow(buffer.max.x)
}

func (buffer *TermBuffer) scrollLineReverse() {
	if DEBUG_TERMBUFFER {
		log.Debug("%s scroll one line reverse", buffer.Desc())
	}
	for r := uint(buffer.scroll.bot); r > buffer.scroll.top; r-- {
		buffer.buffer[r] = buffer.buffer[r-1]
	}
	buffer.buffer[buffer.scroll.top] = makeRow(buffer.max.x)
}

func (buffer *TermBuffer) clear() {
	if DEBUG_TERMBUFFER {
		log.Debug("%s clear", buffer.Desc())
	}
	buffer.buffer = makeBuffer(buffer.cols, buffer.rows)
}

func (buffer *TermBuffer) erasePage(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s erase in page %d", buffer.Desc(), val)
	}
	switch val {
	case 0:
		for c := buffer.cursor.x; c <= buffer.max.x; c++ {
			buffer.buffer[buffer.cursor.y][c] = rune(' ')
		}
		for r := buffer.cursor.y; r <= buffer.max.y; r++ {
			for c := uint(1); c <= buffer.max.x; c++ {
				buffer.buffer[r][c] = rune(' ')
			}
		}
	case 2:
		buffer.buffer = makeBuffer(buffer.cols, buffer.rows)
		buffer.cursor = pos{1, 1}
	default:
		log.Warning("erase page %d not implemented!!", val)
	}
}

func (buffer *TermBuffer) setScrollRegion(top, bottom uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s set scroll region %d-%d", buffer.Desc(), top, bottom)
	}
	if top <= 0 || top > buffer.max.y {
		top = 1
	}
	if bottom <= 0 || bottom > buffer.max.y {
		bottom = buffer.max.y
	}
	if top < bottom {
		buffer.scroll = region{top, bottom}
	}

}

func (buffer *TermBuffer) scrollUp(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s scroll up %d", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		buffer.scrollLine()
	}
}

func (buffer *TermBuffer) scrollDown(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s scroll down %d", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		buffer.scrollLineReverse()
	}
}

func (buffer *TermBuffer) setCursor(x, y uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s set cursor %d,%d", buffer.Desc(), x, y)
	}
	buffer.cursor = pos{x, y}
}

func (buffer *TermBuffer) cursorUp(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s cursor up %d", buffer.Desc(), val)
	}
	if int(buffer.cursor.y)-int(val) >= 1 {
		buffer.cursor.y = buffer.cursor.y - val
	}

}

func (buffer *TermBuffer) cursorRight(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s cursor right %d", buffer.Desc(), val)
	}
	if buffer.cursor.x+val <= buffer.max.x {
		buffer.cursor.x = buffer.cursor.x + val
	}

}

func (buffer *TermBuffer) deleteLine(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s delete %d lines", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		for r := uint(buffer.cursor.y); r < buffer.scroll.bot; r++ {
			buffer.buffer[r] = buffer.buffer[r+1]
		}
		buffer.buffer[buffer.scroll.bot] = makeRow(buffer.max.x)
	}
}

func (buffer *TermBuffer) deleteCharacter(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s delete %d chars", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		for c := uint(buffer.cursor.x); c < buffer.max.x; c++ {
			buffer.buffer[buffer.cursor.y][c] = buffer.buffer[buffer.cursor.y][c+1]
		}
		buffer.buffer[buffer.cursor.y][buffer.max.x] = rune(' ')
	}
}

func (buffer *TermBuffer) eraseCharacter(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s erase %d characters", buffer.Desc(), cnt)
	}
	for c := uint(buffer.cursor.x); c < buffer.cursor.x+cnt && c <= buffer.max.x; c++ {
		buffer.buffer[buffer.cursor.y][c] = ' '
	}
}

func (buffer *TermBuffer) insertCharacter(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s insert %d chars", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		for c := uint(buffer.max.x - 1); c > buffer.cursor.x; c-- {
			buffer.buffer[buffer.cursor.y][c] = buffer.buffer[buffer.cursor.y][c-1]
		}
		buffer.buffer[buffer.cursor.y][buffer.cursor.x] = rune(' ')
	}
}

func (buffer *TermBuffer) eraseLine(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s erase line %d", buffer.Desc(), val)
	}
	switch val {
	case 0:
		for c := buffer.cursor.x; c <= buffer.max.x; c++ {
			buffer.buffer[buffer.cursor.y][c] = rune(' ')
		}
	case 1:
		for c := uint(1); c <= buffer.cursor.x && c <= buffer.max.x; c++ {
			buffer.buffer[buffer.cursor.y][c] = rune(' ')
		}
	case 2:
		for c := uint(1); c <= buffer.max.x; c++ {
			buffer.buffer[buffer.cursor.y][c] = rune(' ')
		}
	default:
		log.Warning("erase line %d not implemented!!", val)
	}
}

func (buffer *TermBuffer) insertLine(cnt uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s insert %d lines", buffer.Desc(), cnt)
	}
	for i := uint(0); i < cnt; i++ {
		for r := uint(buffer.scroll.bot); r > buffer.cursor.y; r-- {
			buffer.buffer[r] = buffer.buffer[r-1]
		}
		buffer.buffer[buffer.cursor.y] = makeRow(buffer.max.x)
	}
}

func (buffer *TermBuffer) linePositionAbsolute(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s line position absolute %d", buffer.Desc(), val)
	}
	//    buffer.cursor.x = 1  // TWEAK? setting x breaks top(1)
	buffer.cursor.y = val
}

func (buffer *TermBuffer) cursorCharacterAbsolute(val uint) {
	if DEBUG_TERMBUFFER {
		log.Debug("%s cursor character absolute %d", buffer.Desc(), val)
	}
	buffer.cursor.x = val
}

func (buffer *TermBuffer) setMode(val string) {
	switch val {
	case "?1049":
		buffer.saveBuffer()
		buffer.saveCursor()
	default:
		if DEBUG_TERMBUFFER {
			log.Debug("%s ignore set mode '%s'", buffer.Desc(), lookupMode(val))
		}
	}
}

func (buffer *TermBuffer) resetMode(val string) {
	switch val {
	case "?1049":
		buffer.restoreBuffer()
		buffer.restoreCursor()
	default:
		if DEBUG_TERMBUFFER {
			log.Debug("%s ignore reset mode '%s'", buffer.Desc(), lookupMode(val))
		}

	}
}

func (buffer *TermBuffer) ProcessSequence(seq *ansi.S) {
	// lock mutex?

	sequence, ok := lookupSequence(seq.Code)
	if !ok {
		return
		//unlock mutex tho?
	}

	switch sequence {

	case ansi.Table[ansi.ECH]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.eraseCharacter(val)

	case ansi.Table[ansi.ED]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.erasePage(val)

	case ansi.Table[ansi.CUP]:
		var val pos
		fmt.Sscanf(seq.Params[0], "%d", &val.y)
		fmt.Sscanf(seq.Params[1], "%d", &val.x)
		buffer.setCursor(val.x, val.y)

	case ansi.Table[ansi.CUU]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.cursorUp(val)

	case ansi.Table[ansi.CUF]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.cursorRight(val)

	case ansi.Table[ansi.EL]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.eraseLine(val)

	case ansi.Table[ansi.IL]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.insertLine(val)

	case ansi.Table[ansi.DL]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.deleteLine(val)

	case ansi.Table[ansi.ICH]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.insertCharacter(val)

	case ansi.Table[ansi.DCH]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.deleteCharacter(val)

	case ansi.Table[ansi.VPA]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.linePositionAbsolute(val)

	case ansi.Table[ansi.CHA]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.cursorCharacterAbsolute(val)

	case ansi.Table[ansi.RI]:
		buffer.scrollLineReverse()

	case ansi.Table[ansi.SM]:
		var val string = seq.Params[0]
		buffer.setMode(val)

	case ansi.Table[ansi.RM]:
		var val string = seq.Params[0]
		buffer.resetMode(val)

	case ansi.Table[ansi.SU]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.scrollUp(val)

	case ansi.Table[ansi.SD]:
		var val uint
		fmt.Sscanf(seq.Params[0], "%d", &val)
		buffer.scrollDown(val)

	case xtermTable[DECSTBM]:
		var val region
		if len(seq.Params) > 0 {
			fmt.Sscanf(seq.Params[0], "%d", &val.top)
		}
		if len(seq.Params) > 1 {
			fmt.Sscanf(seq.Params[1], "%d", &val.bot) // REM, crashes here, with exec w3m feedface.com
		}
		buffer.setScrollRegion(val.top, val.bot)

	case ansi.Table[ansi.SGR]:
		break // no support for color / weight / decorations

	default:
		log.Warning("%s unhandled sequence %s '%s'", buffer.Desc(), sequence.Name, sequence.Desc)

	}

}

func (buffer *TermBuffer) Desc() string {
	alt := ""
	if buffer.altBuffer != nil || buffer.altCursor != nil {
		alt = " alt"
	}
	scr := ""
	if buffer.scroll.top != 1 || buffer.scroll.bot != buffer.max.y {
		scr = fmt.Sprintf(" %d-%d", buffer.scroll.top, buffer.scroll.bot)
	}
	return fmt.Sprintf("termbuffer[%2dx%-2d %2d,%-2d%s%s]", buffer.cols, buffer.rows, buffer.cursor.x, buffer.cursor.y, alt, scr)
}

func (buffer *TermBuffer) Dump() string {
	ret := ""
	ret += "+"
	for c := uint(1); c <= buffer.max.x; c++ {
		ret += "-"
	}
	ret += "+\n"
	for r := uint(1); r <= buffer.max.y; r++ {
		ret += "|"
		for c := uint(1); c <= buffer.max.x; c++ {
			if c == buffer.cursor.x && r == buffer.cursor.y {
				ret += "\033[7m"
			}
			ret += fmt.Sprintf("%c", buffer.buffer[r][c])
			if c == buffer.cursor.x && r == buffer.cursor.y {
				ret += "\033[27m"
			}
		}
		ret += "| "
		if r%10 == 0 {
			ret += fmt.Sprintf("%01d", (r/10)%10)
		} else {
			ret += " "
		}
		ret += fmt.Sprintf("%01d\n", r%10)
	}
	ret += "+"
	for c := uint(1); c <= buffer.max.x; c++ {
		ret += "-"
	}
	ret += "+\n "
	for c := uint(1); c <= buffer.max.x; c++ {
		if c%10 == 0 {
			ret += fmt.Sprintf("%01d", (c/10)%10)
		} else {
			ret += " "
		}
	}
	ret += "\n "
	for c := uint(1); c <= buffer.max.x; c++ {
		ret += fmt.Sprintf("%01d", c%10)
	}
	ret += "\n"
	return ret
}

func (buffer *TermBuffer) GetWidth() uint64  { return uint64(buffer.cols) }
func (buffer *TermBuffer) GetHeight() uint64 { return uint64(buffer.rows) }
