
// +build !linux
// +build !arm


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

func (renderer *Renderer) Init() error {
    log.Fatal("RENDERER NOT AVAILABLE")  
    return errors.New("RENDERER NOT AVAILABLE") 
}

func (renderer *Renderer) Configure(config *conf.Config) error {
    log.Fatal("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Render() error { 
    log.Fatal("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
} 

func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    log.Fatal("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    log.Fatal("RENDERER NOT AVAILABLE")
    return errors.New("RENDERER NOT AVAILABLE")
}