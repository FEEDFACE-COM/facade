
// +build !linux !arm

package main

import (
    "errors"
    facade "./facade"
)

const RENDERER_AVAILABLE = false

type Renderer struct {}
func NewRenderer() *Renderer { return &Renderer{} }

func (renderer *Renderer) Init(config *facade.Config) error { return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Configure(config *facade.Config) error { return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Render(confChan chan facade.Config, textChan chan string) error { return errors.New("RENDERER NOT AVAILABLE") } 
func (renderer *Renderer) QueueTexts(textChan chan facade.RawText) error { return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) QueueConfs(confChan chan facade.Config) error { return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessConf(rawChan chan facade.Config, confChan chan facade.Config) error { return errors.New("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessText(rawChan chan facade.RawText, textChan chan string) error { return errors.New("RENDERER NOT AVAILABLE") }


