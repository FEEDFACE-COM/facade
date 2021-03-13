// +build darwin,amd64

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"fmt"
    gl "github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

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

func (shader *Shader) LoadSource(src string) {
	shader.Source = src
	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s source loaded", shader.Desc())
	}
}

func (shader *Shader) IndexName() string {
	return fmt.Sprintf("%s.%s", shader.Name, string(shader.Type))
}

func (shader *Shader) Desc() string {
	return fmt.Sprintf("shader[%s.%s]", shader.Name, string(shader.Type))
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
		gl.DeleteShader(shader.Shader)
	}
}

func (shader *Shader) CompileShader() error {
	if shader.Shader <= 0 {
		switch shader.Type {
		case VertType:
			shader.Shader = gl.CreateShader(gl.VERTEX_SHADER)
		case FragType:
			shader.Shader = gl.CreateShader(gl.FRAGMENT_SHADER)
		}
	}

	sources, free := gl.Strs(shader.Source + "\x00")
	gl.ShaderSource(shader.Shader, 1, sources, nil)
	free()
	gl.CompileShader(shader.Shader)

	//check
	var status int32
	gl.GetShaderiv(shader.Shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader.Shader, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader.Shader, logLength, nil, gl.Str(logs))
		log.Error("%s fail compile shader: %s", shader.Desc(), logs)
		return log.NewError("fail compile shader %s: %s", shader.Name, logs)
	}

	return nil
}
