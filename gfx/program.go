
// +build linux,arm


package gfx

import (
    "strings"
    "fmt"
    "io/ioutil"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)


type Program struct {
    Name string

    Program uint32
    vertexShader *Shader
    fragmentShader *Shader
    debugFlag float32
    broken bool
}



func NewProgram(name string) *Program {
    ret := &Program{Name: name}
    return ret
}

func (program *Program) UseProgram(debug bool) { 
    gl.UseProgram(program.Program) 
    program.debugFlag = float32(0)
    if debug {
        program.debugFlag = float32(1)
    }
        program.Uniform1f(DEBUGFLAG, program.debugFlag)
}



func (program *Program) loadShaderFile(shaderName string, shaderType uint32) (string, error) {
    var data []byte
    var err error
    filePath := "/home/folkert/src/gfx/facade/shader/" + shaderName
    switch (shaderType) {
        case gl.VERTEX_SHADER: filePath += ".vert"
        case gl.FRAGMENT_SHADER: filePath += ".frag"
        default: return "", log.NewError("unknown shadertype %d for file %s",shaderType,shaderName)    
    }
    data, err = ioutil.ReadFile(filePath)
    if err != nil {
        log.Error("fail read shader file %s: %s",filePath,err)
        return "", log.NewError("fail read shader file: %s",err)
    }
//    log.Debug("read shader file %s",filePath)
    return string(data), nil    
}

func (program *Program) LoadVertexShader(vertName string) error {
    var err error
    var src string    
    
    // try from file
    src,err = program.loadShaderFile(vertName,gl.VERTEX_SHADER)
    if err == nil { //success
        program.vertexShader = NewShader(vertName, src, gl.VERTEX_SHADER)    
//        log.Debug("load vertex shader %s from file",program.vertexShader.Name)
    } else if VertexShader[vertName] != "" {
        program.vertexShader = NewShader(vertName, VertexShader[vertName], gl.VERTEX_SHADER)
//        log.Debug("load vertex shader %s from map",program.vertexShader.Name)
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
    var src string    
    
    // try from file
    src,err = program.loadShaderFile(fragName,gl.FRAGMENT_SHADER)
    if err == nil { //success
        program.fragmentShader = NewShader(fragName, src, gl.FRAGMENT_SHADER)    
//        log.Debug("load fragment shader %s from file",program.fragmentShader.Name)
    } else if FragmentShader[fragName] != "" {
        program.fragmentShader = NewShader(fragName, FragmentShader[fragName], gl.FRAGMENT_SHADER)
//        log.Debug("load fragment shader %s from map",program.fragmentShader.Name)
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
	var status int32

    //todo: cleanup if already present?

    if program.vertexShader == nil || program.fragmentShader == nil {
        program.broken = true
        return log.NewError("missing shader in %s",program.Name)
    }

    

	program.Program = gl.CreateProgram()
	gl.AttachShader(program.Program, program.vertexShader.Shader)
	gl.AttachShader(program.Program, program.fragmentShader.Shader)
	gl.LinkProgram(program.Program)

	gl.GetProgramiv(program.Program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.Program, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.Program, logLength, nil, gl.Str(logs))

        src := ""
        if strings.Contains(logs,"vertex shader")   { src += "\n" + program.vertexShader.Source }
        if strings.Contains(logs,"fragment shader") { src += "\n" + program.fragmentShader.Source }


        program.broken = true
		log.Error("fail link program %s/%s: %s%s", program.vertexShader.Name,program.fragmentShader.Name, logs,src)
		return log.NewError("fail link program %s/%s",program.vertexShader.Name,program.fragmentShader.Name)
	}

	gl.DeleteShader(program.vertexShader.Shader)
	gl.DeleteShader(program.fragmentShader.Shader)

	return nil
    
    
}


func (program *Program) VertexAttribPointer(name AttribName, size int32, stride int32, offset int) {
    ret := gl.GetAttribLocation(program.Program, gl.Str(string(name)+"\x00")) 
    if ret < 0 {
//        log.Debug("no vertexattrib %s: %d",name,ret)
        return;
    }
    gl.EnableVertexAttribArray( uint32(ret) )
    gl.VertexAttribPointer( uint32(ret), size, gl.FLOAT, false, stride, gl.PtrOffset(offset) )
    return
}


func (program *Program) uniformLocation(name UniformName) int32 {
    ret := gl.GetUniformLocation(program.Program, gl.Str(string(name)+"\x00") )
    if ret < 0 { 
        log.Debug("fail get uniform location %s program %s: %d",name,program.Name,ret)
    }
    return ret;
}

func (program *Program) UniformMatrix4fv(name UniformName, count int32, value *float32) int32 {
    if program.broken { return -1 }
    ret := program.uniformLocation(name)
    if ret <= 0 {
        return ret;
    }
    gl.UniformMatrix4fv(ret, count, false, value)
    return ret
}


func (program *Program) Uniform2f(name UniformName, value0, value1 float32) int32 {
    if program.broken { return -1 }
    ret := program.uniformLocation(name)
    if ret <= 0 {
        return ret;
    }
    gl.Uniform2f(ret, value0, value1)
    return ret
}

func (program *Program) Uniform2fv(name UniformName, count int32, value *float32) int32 {
    if program.broken { return -1 }
    ret := program.uniformLocation(name)
    if ret <= 0 {
        return ret;
    }
    gl.Uniform2fv(ret,count, value)
    return ret
}

func (program *Program) Uniform1f(name UniformName, value float32) int32 {
    if program.broken { return -1 }
    ret := program.uniformLocation(name)
    if ret <= 0 {
        return ret;
    }
    gl.Uniform1f(ret, value)
    return ret
}

func (program *Program) Uniform1i(name UniformName, value int32) int32 {
    if program.broken { return -1 }
    ret := program.uniformLocation(name)
    if ret <= 0 {
        return ret;
    }
    gl.Uniform1i(ret, value)
    return ret
}

func (program *Program) Desc() string {
    tmp := ""
    if program.vertexShader != nil { tmp += " " + program.vertexShader.Name }
    if program.fragmentShader != nil { tmp += "/" + program.fragmentShader.Name }
    return fmt.Sprintf("program[%s%s]",program.Name,tmp)
}


