
// +build linux,arm

package gfx

import(
//    log "../log"    
)

type Text struct {
    Text string
    Texture *Texture
}

func (text *Text) Close() {
    if text.Texture != nil { text.Texture.Close() }
}

func (text *Text) Desc() string { return text.Text }


func NewText(text string) *Text {
    ret := &Text{Text: text}
    return ret
}


func (text *Text) RenderTexture(font *Font) {

    if text.Texture != nil { 
        //REM, cleanup old, rerender!
        return; 
    }
    
    text.Texture = NewTexture()
    if text.Text == "" {
        text.Texture.LoadEmpty()
    } else {
        rgba, err := font.RenderTextRGBA(text.Text)
        if err != nil {
            
        } else {
            text.Texture.LoadRGBA(rgba)    
        }
    }
//	text.Texture.LoadFile("/home/folkert/src/gfx/facade/asset/FEEDFACE.COM.white.png")
    text.Texture.GenTexture()
        
}
