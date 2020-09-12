// +build !linux !arm

package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
)

const RENDERER_AVAILABLE = false

type Renderer struct{mode facade.Mode}

func NewRenderer(string, chan bool) *Renderer { return &Renderer{} }

func (renderer *Renderer) Init() error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Configure(config *facade.Config) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Render(chan facade.Config) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) ProcessTextSeqs(chan facade.TextSeq) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) ProcessQueries(chan (chan string)) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Finish() error {
	return log.NewError("RENDERER NOT AVAILABLE")
}
