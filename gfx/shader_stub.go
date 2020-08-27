// +build !linux,!arm

package gfx

import (
	log "../log"
	"fmt"
)

const DEBUG_SHADER = false

type Shader struct {
	Name   string
	Source string
	Type   ShaderType
	Shader uint32
}

type UniformName string

const (
	DEBUGFLAG   UniformName = "debugFlag"
	PROJECTION  UniformName = "projection"
	MODEL       UniformName = "model"
	VIEW        UniformName = "view"
	SCREENRATIO UniformName = "screenRatio"
	FONTRATIO   UniformName = "fontRatio"
	TEXTURE     UniformName = "texture"
	CLOCKNOW    UniformName = "now"
)

type AttribName string

const (
	VERTEX   AttribName = "vertex"
	COLOR    AttribName = "color"
	TEXCOORD AttribName = "texCoord"
)

type ShaderType string

const (
	VertType ShaderType = "vert"
	FragType ShaderType = "frag"
)

func (shader *Shader) LoadSource(src string) error {
	//    var err error

	shader.Source = src
	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s shader setup", shader.Desc())
	}
	return nil

}

func (shader *Shader) IndexName() string { return shader.Name + "." + string(shader.Type) }

func (shader *Shader) Desc() string {
	return fmt.Sprintf("%s[%s]", shader.Type, shader.Name)
}

func NewShader(name string, shaderType ShaderType) *Shader {
	ret := &Shader{Name: name, Type: shaderType}
	return ret
}

func (shader *Shader) Close() {
	if shader.Shader > 0 {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("delete %s", shader.Desc())
		}
	}
}

func (shader *Shader) CompileShader() error {
	return log.NewError("NO COMPILE")
}
