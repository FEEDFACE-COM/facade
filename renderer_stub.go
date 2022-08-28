//go:build !RENDERER
// +build !RENDERER

package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
)

const RENDERER_AVAILABLE = false

type Renderer struct{}

func NewRenderer(string, chan Tick) *Renderer { return &Renderer{} }

func (renderer *Renderer) Init() error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Configure(*facade.Config) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}

func (renderer *Renderer) Render(chan facade.Config, bool) error {
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
