// +build linux,arm

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"fmt"
	"regexp"
	"strings"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
)

type Program struct {
	name string
	mode string

	program        uint32
	vertexShader   *Shader
	fragmentShader *Shader
	debugFlag      float32

	programService *ProgramService
}

func (program *Program) HasShader(shader *Shader) bool {
	if shader.Type == VertType && shader == program.vertexShader {
		return true
	}
	if shader.Type == FragType && shader == program.fragmentShader {
		return true
	}
	return false
}

func NewProgram(name, mode string, programService *ProgramService) *Program {
	ret := &Program{name: name, mode: mode, programService: programService}
	return ret
}

func (program *Program) SetDebug(debug bool) {
	program.debugFlag = float32(0)
	if debug {
		program.debugFlag = float32(1)
	}
	program.Uniform1fv(DEBUGFLAG, 1, &program.debugFlag)
}

func (program *Program) UseProgram(debug bool) {
	gl.UseProgram(program.program)
	program.SetDebug(debug)
}

func (program *Program) Relink() error {
	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s relink", program.Desc())
	}

	err := program.attachShadersAndLinkProgram()
	if err != nil {
		return log.NewError("fail relink: %s", err)
	}

	return nil
}

func (program *Program) Link(vertName, fragName string) error {

	var err error
	var vert, frag *Shader = program.vertexShader, program.fragmentShader

	err = program.getAndCompileShaders(vertName, fragName)
	if err != nil {
		program.vertexShader = vert
		program.fragmentShader = frag
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s fail compile shaders %s,%s: %s", program.Desc(), vertName, fragName, err)
		}
		return log.NewError("fail compile: %s", err)
	}

	err = program.attachShadersAndLinkProgram()
	if err != nil {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s fail link %s,%s: %s", program.Desc(), vertName, fragName, err)
		}
		program.vertexShader = vert
		program.fragmentShader = frag
		err2 := program.attachShadersAndLinkProgram()
		if err2 != nil {
			if DEBUG_PROGRAMSERVICE {
				log.Debug("%s fail reset to previous shaders: %s", program.Desc(), err2)
			}
		} else {
			if DEBUG_PROGRAMSERVICE {
				log.Debug("%s reset to previous shaders", program.Desc())
			}
		}
		return log.NewError("fail link: %s", err)
	}

	return nil

}

func (program *Program) getAndCompileShaders(vertName, fragName string) error {

	var err error
	var mode = program.mode

	var vert, frag = program.vertexShader, program.fragmentShader

	if vert == nil || vert.Name != vertName {
		program.vertexShader, err = program.programService.GetShader(mode+vertName, VertType)
		if err != nil {
			return log.NewError("fail get vert shader: %s", err)
		}
	}
	err = program.vertexShader.CompileShader()
	if err != nil {
		return log.NewError("fail compile vert shader: %s", err)
	}
	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s compiled vert shader %s", program.Desc(), program.vertexShader.Name)
	}

	if frag == nil || frag.Name != fragName {
		program.fragmentShader, err = program.programService.GetShader(mode+fragName, FragType)
		if err != nil {
			return log.NewError("fail get frag shader: %s", err)
		}
	}
	err = program.fragmentShader.CompileShader()
	if err != nil {
		return log.NewError("fail compile frag shader: %s", err)
	}
	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s compiled frag shader %s", program.Desc(), program.fragmentShader.Name)
	}

	return nil
}

func (program *Program) attachShadersAndLinkProgram() error {
	var ret error

	if program.vertexShader == nil || program.fragmentShader == nil {
		return log.NewError("%s fail link missing shader", program.Desc())
	}

	if program.program > 0 {
		gl.DeleteProgram(program.program)
	}
	program.program = gl.CreateProgram()

	gl.AttachShader(program.program, program.vertexShader.Shader)
	gl.AttachShader(program.program, program.fragmentShader.Shader)
	gl.LinkProgram(program.program)

	var status int32
	gl.GetProgramiv(program.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		error, source := program.getLinkError()
		log.Error("%s %s%s", program.Desc(), error, source)
		ret = log.NewError("fail link: %s", error)
	} else {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s attached and linked", program.Desc())
		}
	}

	gl.DetachShader(program.program, program.vertexShader.Shader)
	gl.DetachShader(program.program, program.fragmentShader.Shader)

	return ret

}

func (program *Program) getLinkError() (string, string) {

	var logLength int32
	gl.GetProgramiv(program.program, gl.INFO_LOG_LENGTH, &logLength)

	logs := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program.program, logLength, nil, gl.Str(logs))

	source := ""

	var shaderType string
	var lineNumber int

	{
		re := regexp.MustCompile(`shader, line ([0-9]+)\)`)
		str := re.FindStringSubmatch(logs)
		if len(str) > 1 {
			fmt.Sscanf(str[1], "%d", &lineNumber)
		}
	}
	{
		re := regexp.MustCompile(`(vertex|fragment) shader`)
		str := re.FindStringSubmatch(logs)
		if len(str) > 1 {
			shaderType = str[1]
		}
	}

	if shaderType == "vertex" {
		if lineNumber <= 0 {
			source += "\n" + program.vertexShader.Source
		} else {
			for i, str := range strings.Split(program.vertexShader.Source, "\n") {
				if i > (lineNumber-5) && (i < lineNumber+5) {
					source += "\n" + fmt.Sprintf("%4d  ", i+1) + str
				}
			}
		}

	} else if shaderType == "fragment" {
		if lineNumber <= 0 {
			source += "\n" + program.fragmentShader.Source
		} else {
			for i, str := range strings.Split(program.fragmentShader.Source, "\n") {
				if i > (lineNumber-5) && (i < lineNumber+5) {
					source += "\n" + fmt.Sprintf("%4d  ", i) + str
				}
			}
		}
	}

	return logs, source

}

func (program *Program) Close() {
	if program.program > 0 {
		gl.DeleteProgram(program.program)
		program.program = 0
	}
}

func (program *Program) VertexAttribPointer(name AttribName, size int32, stride int32, offset int) error {
	ret := gl.GetAttribLocation(program.program, gl.Str(string(name)+"\x00"))
	if ret < 0 {
		return log.NewError("no pointer for attribute '%s' by program %s", name, program.name)
	}
	gl.EnableVertexAttribArray(uint32(ret))
	gl.VertexAttribPointer(uint32(ret), size, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	return nil
}

func (program *Program) uniformLocation(name UniformName) (int32, error) {
	ret := gl.GetUniformLocation(program.program, gl.Str(string(name)+"\x00"))
	if ret < 0 {
		return -1, log.NewError("no location for uniform '%s' by program %s", name, program.name)
	}
	return ret, nil
}

func (program *Program) UniformMatrix4fv(name UniformName, count int32, value *float32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.UniformMatrix4fv(ret, count, false, value)
	return ret, nil
}

func (program *Program) Uniform2f(name UniformName, value0, value1 float32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.Uniform2f(ret, value0, value1)
	return ret, nil
}

func (program *Program) Uniform2fv(name UniformName, count int32, value *float32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.Uniform2fv(ret, count, value)
	return ret, nil
}

func (program *Program) Uniform1fv(name UniformName, count int32, value *float32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.Uniform1fv(ret, count, value)
	return ret, nil
}

func (program *Program) Uniform1f(name UniformName, value float32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.Uniform1f(ret, value)
	return ret, nil
}

func (program *Program) Uniform1i(name UniformName, value int32) (int32, error) {
	ret, err := program.uniformLocation(name)
	if err != nil {
		return -1, err
	}
	gl.Uniform1i(ret, value)
	return ret, nil
}

func (program *Program) Desc() string {
	tmp := program.name + " "
	if program.vertexShader != nil {
		tmp += program.vertexShader.Name
	}
	if program.fragmentShader != nil {
		tmp += "," + program.fragmentShader.Name
	}
	return fmt.Sprintf("prog[%s]", tmp)
}
