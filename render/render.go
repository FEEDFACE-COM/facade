
// +build !linux !arm

package render

import (
    "errors"
    log "../log"
    conf "../conf"
)

const RENDERER_AVAILABLE = false

type Renderer struct {}

func NewRenderer() *Renderer { 
    return &Renderer{} 
}

func (renderer *Renderer) Init(config *conf.Config) error {
    log.PANIC("RENDERER NOT AVAILABLE")  
    return errors.New("RENDERER NOT AVAILABLE") 
}

func (renderer *Renderer) Configure(config *conf.Config) error {
    log.PANIC("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Render(confChan chan conf.Config, textChan chan conf.Text) error { 
    log.PANIC("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
} 

func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    log.PANIC("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    log.PANIC("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}