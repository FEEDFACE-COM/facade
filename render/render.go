
// +build !darwin,!amd64


package render

import (
    log "../log"
)

const RENDERER_AVAILABLE = false

type Renderer struct {}
func NewRenderer() *Renderer { return &Renderer{} }
func (renderer *Renderer) Init() error { log.Fatal("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Start() error { log.Fatal("RENDERER NOT AVAILABLE") } 

