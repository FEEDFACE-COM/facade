package gfx

import (
	//    "image"
	log "../log"
)

//type TextTexture interface {
//
//    Close()
//    LoadRGBA(*image.RGBA)
//
//}
//

type Text struct {
	Text    string
	Texture *Texture
}

func (text *Text) Bytes() []byte {
	return []byte(text.Text)
}

func (text *Text) Close() {
	if text.Texture != nil {
		text.Texture.Close()
	}
}

func (text *Text) Desc() string { return text.Text }

func NewText(text string) *Text {
	ret := &Text{Text: text}
	return ret
}

const MAX_LENGTH = 20

func (text *Text) RenderTexture(font *Font) error {

	if text.Texture != nil {
		//REM, cleanup old, then rerender!
		log.Error("TODO: CLEANUP!!")
		return nil
	}

	txt := text.Text
	if len(txt) > MAX_LENGTH {
		txt = text.Text[:MAX_LENGTH]
	}

	text.Texture = NewTexture("text")
	if text.Text == "" {
		text.Texture.LoadEmpty()
	} else {
		rgba, err := font.RenderText(txt, false)
		if err != nil {
			return log.NewError("fail render '%s': %s", txt, err)
		} else {
			text.Texture.LoadRGBA(rgba)
		}
	}
	text.Texture.TexImage()
	return nil
}


