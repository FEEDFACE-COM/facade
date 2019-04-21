
// +build linux,arm

package gfx

type Buffer interface {

    WriteText(txt *Text)
    WriteString(str string)
    
}

