
// +build !linux !arm

package main

import (
    log "./log"
    facade "./facade"
)

const RENDERER_AVAILABLE = false

type Renderer struct {}
func NewRenderer(_ string) *Renderer { return &Renderer{} }

func (renderer *Renderer) Init(config *facade.Config) error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Configure(config *facade.Config) error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Render(confChan chan facade.Config, textChan chan string) error { return log.NewError("RENDERER NOT AVAILABLE") } 
func (renderer *Renderer) QueueTexts(textChan chan facade.RawText) error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) QueueConfs(confChan chan facade.Config) error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessRawConfs(rawChan chan facade.Config, confChan chan facade.Config) error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) ProcessRawTexts(rawChan chan facade.RawText, textChan chan string) error { return log.NewError("RENDERER NOT AVAILABLE") }


