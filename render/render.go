
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

func (renderer *Renderer) Init(config *conf.Config) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Configure(config *conf.Config) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Render(confChan chan conf.Config, textChan chan string) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") } 
func (renderer *Renderer) ReadText(textChan chan conf.RawText) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ReadConf(confChan chan conf.Config) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessConf(rawChan chan conf.Config, confChan chan conf.Config) error { log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessText(rawChan chan conf.RawText, textChan chan string) error {log.PANIC("RENDERER NOT AVAILABLE"); return errors.New("RENDERER NOT AVAILABLE") }


