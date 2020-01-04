// +build !linux
// +build !arm

package main

import (
	facade "./facade"
	log "./log"
)

const RENDERER_AVAILABLE = false

type Renderer struct{}

func NewRenderer(_ string) *Renderer { return &Renderer{} }

func (renderer *Renderer) Init() error { return log.NewError("RENDERER NOT AVAILABLE") }
func (renderer *Renderer) Configure(config *facade.Config) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}
func (renderer *Renderer) Render(confChan chan facade.Config) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}
func (renderer *Renderer) ProcessTextSeqs(bufChan chan facade.TextSeq) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}
func (renderer *Renderer) ProcessQueries(queryChan chan (chan string)) error {
	return log.NewError("RENDERER NOT AVAILABLE")
}
