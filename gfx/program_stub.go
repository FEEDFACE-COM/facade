// +build !linux !arm

package gfx

import (
	"fmt"
)

type Program struct {
	name string
	mode string

	programService *ProgramService
}

func (program *Program) HasShader(shader *Shader) bool { return false }
func (program *Program) Desc() string                  { return fmt.Sprintf("prog[%s]", program.name) }
func NewProgram(name, mode string, programService *ProgramService) *Program {
	return &Program{name: name, mode: mode}
}
func (program *Program) Relink() error { return nil }
