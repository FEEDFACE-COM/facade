
// +build linux,arm


package gfx

import (
    "strings"
    "fmt"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type UniformName string
const (
    PROJECTION UniformName = "projection"
    MODEL      UniformName = "model"
    VIEW       UniformName = "view"
    TEXTURE    UniformName = "texture"
)    

type AttribName string
const (
    VERTEX     AttribName = "vertex"    
    COLOR      AttribName = "color"    
    TEXCOORD   AttribName = "texcoord"
)


type Program struct {
    Name string

    Program uint32
    vertexShader *Shader
    fragmentShader *Shader
}



func NewProgram(name string) *Program {
    ret := &Program{Name: name}
    return ret
}

func (program *Program) UseProgram() { gl.UseProgram(program.Program) }


func (program *Program) LoadVertexShader(vertName string) error {
    var err error
    
    // look up in maps
    if VertexShader[vertName] != "" {
        program.vertexShader = NewShader(vertName, VertexShader[vertName], gl.VERTEX_SHADER)
        log.Debug("found vertex shader %s in map:\n%s",program.vertexShader.Name,program.vertexShader.ShaderSource)
    }

    if program.vertexShader == nil {
        return log.NewError("unknown vertex shader in %s",vertName)
    }

    err = program.vertexShader.CompileShader()
    if err != nil {
        return log.NewError("fail compile vertex shader %s: %s",vertName,err)
    }

    return nil
}    
    
    

func (program *Program) LoadFragmentShader(fragName string) error {
    var err error
    
    if FragmentShader[fragName] != "" {
        program.fragmentShader = NewShader(fragName, FragmentShader[fragName], gl.FRAGMENT_SHADER)
        log.Debug("found fragment shader %s in map:\n%s",program.fragmentShader.Name,program.fragmentShader.ShaderSource)
    }


    if program.fragmentShader == nil {
        return log.NewError("unknown fragment shader in %s",fragName)
    }
    
    err = program.fragmentShader.CompileShader()
    if err != nil {
        return log.NewError("fail compile fragment shader %s: %s",fragName,err)
    }
    
    
    return nil

}

func (program *Program) LoadShaders(vertName, fragName string) error {
    var err error
    err = program.LoadVertexShader(vertName)
    if err != nil { return err }
    err = program.LoadFragmentShader(fragName)
    if err != nil { return err }
    return nil
}


func (program *Program) LinkProgram() error {

    //todo: cleanup if already present?

    if program.vertexShader == nil || program.fragmentShader == nil {
        return log.NewError("missing shader in %s",program.Name)
    }

    

	program.Program = gl.CreateProgram()
	gl.AttachShader(program.Program, program.vertexShader.Shader)
	gl.AttachShader(program.Program, program.fragmentShader.Shader)
	gl.LinkProgram(program.Program)

	var status int32
	gl.GetProgramiv(program.Program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.Program, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.Program, logLength, nil, gl.Str(logs))

		log.Error("fail link %s %s: %v", program.vertexShader.Name, program.fragmentShader.Name, logs)
		return log.NewError("fail link %s,%s:",program.vertexShader.Name, program.fragmentShader.Name)
	}

	gl.DeleteShader(program.vertexShader.Shader)
	gl.DeleteShader(program.fragmentShader.Shader)

	return nil
    
    
}


func (program *Program) VertexAttribPointer(name AttribName, size int32, stride int32, offset int) uint32 {
    ret := uint32( gl.GetAttribLocation(program.Program, gl.Str(string(name)+"\x00")) )
    gl.EnableVertexAttribArray(ret)
    gl.VertexAttribPointer(ret, size, gl.FLOAT, false, stride, gl.PtrOffset(offset) )
    return ret
}


func (program *Program) UniformMatrix4fv(name UniformName, count int32, value *float32) int32 {
    ret := gl.GetUniformLocation(program.Program, gl.Str(string(name)+"\x00") )
    gl.UniformMatrix4fv(ret, count, false, value)
    return ret
}

func (program *Program) Uniform1i(name UniformName, value int32) int32 {
    ret := gl.GetUniformLocation(program.Program, gl.Str(string(name)+"\x00") )
    gl.Uniform1i(ret, value)
    return ret
}

func (program *Program) Desc() string {
    tmp := ""
    if program.vertexShader != nil { tmp += " " + program.vertexShader.Name }
    if program.fragmentShader != nil { tmp += " " + program.fragmentShader.Name }
    return fmt.Sprintf("program[%s%s]",program.Name,tmp)
}

