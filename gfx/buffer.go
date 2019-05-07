
// +build linux,arm

package gfx




type Buffer interface {

    WriteRaw
    WriteText(txt *Text)
    WriteString(str string)
    
}

