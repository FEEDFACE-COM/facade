
package gfx

import (
    log "../log"
)

type Line struct {
    Text string
}

func NewLine(text string) *Line {
    log.Debug("+line(%s)",text)   
    return &Line{Text: text}
}


func (line *Line) Close() {
    log.Debug("~line(%s)",line.Text)   
}


