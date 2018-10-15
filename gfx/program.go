
// +build linux,arm


package gfx

import (
    "strings"
    "fmt"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)


type Program struct {
    Name string

    Program uint32
    vertexShader *Shader
    fragmentShader *Shader
    debugFlag float32
}



type Refresh struct{ program *Program; shader *Shader }
var refreshChan chan Refresh = make( chan Refresh )


func RefreshPrograms() {
    var err error
    select {
        
        case refresh := <-refreshChan:
        
            log.Debug("refresh %s %s",refresh.program.Desc(),refresh.shader.Desc())
            err = refresh.program.ReloadShader( refresh.shader )
            if err != nil {
                log.Debug("fail refresh %s %s",refresh.program.Desc(),refresh.shader.Desc())
            }
        
        default:
            //nop
        
    }
    
}


func (program *Program) ReloadShader(shader *Shader) error {
    var err error

    err = shader.CompileShader()
    if err != nil { return err }

//    err = program.vertexShader.CompileShader()
//    if err != nil { return err }
//    err = program.fragmentShader.CompileShader()
//    if err != nil { return err }

//    switch shader.Type {
//        case VERTEX_SHADER:
//            program.vertexShader = shader
//            
//        case FRAGMENT_SHADER:    
//            program.fragmentShader = shader
//    }
    
    err = program.LinkProgram()
    if err != nil { return err }

    log.Debug("reload %s %s",program.Desc(),shader.Desc())
    return nil    
}

func GetProgram(name string) *Program {
    ret := &Program{Name: name}
    return ret
}


func (program *Program) UseProgram(debug bool) { 
    gl.UseProgram(program.Program) 
    program.debugFlag = float32( 0 )
    if debug {
        program.debugFlag = float32( 1 )
    }
    program.Uniform1f(DEBUGFLAG, program.debugFlag)
}




func (program *Program) LoadShaders(vertName, fragName string) error {
    var err error
    program.vertexShader, err = GetShader(vertName,VERTEX_SHADER,program)
    if err != nil { return log.NewError("fail to get shader: %s",err) }
    err = program.vertexShader.CompileShader()
    if err != nil { return log.NewError("fail to compile shader: %s",err) }
    
    
    program.fragmentShader, err = GetShader(fragName,FRAGMENT_SHADER,program)
    if err != nil { return log.NewError("fail to get shader: %s",err) }
    err = program.fragmentShader.CompileShader()
    if err != nil { return log.NewError("fail to compile shader: %s",err) }
    
    return nil   
}



func (program *Program) LinkProgram() error {
    var ret error

    if program.vertexShader == nil || program.fragmentShader == nil {
        return log.NewError("missing shader in %s",program.Name)
    }

    if program.Program > 0 {
        gl.DeleteProgram(program.Program)
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

        src := ""
        if strings.Contains(logs,"vertex shader")   { src += "\n" + program.vertexShader.Source }
        if strings.Contains(logs,"fragment shader") { src += "\n" + program.fragmentShader.Source }

		log.Error("fail link program %s/%s: %s%s", program.vertexShader.Name,program.fragmentShader.Name, logs,src)
		ret = log.NewError("fail link program %s/%s",program.vertexShader.Name,program.fragmentShader.Name)
	}

	gl.DetachShader(program.Program, program.vertexShader.Shader)
	gl.DetachShader(program.Program, program.fragmentShader.Shader)


	return ret
    
    
}

func (program *Program) Close() {
    if program.Program > 0 {
        gl.DeleteProgram(program.Program)
        program.Program = 0
    }
}


func (program *Program) VertexAttribPointer(name AttribName, size int32, stride int32, offset int) error {
    ret := gl.GetAttribLocation(program.Program, gl.Str(string(name)+"\x00")) 
    if ret < 0 {
        return log.NewError("no pointer for attribute '%s' by program %s",name,program.Name)
    }
    gl.EnableVertexAttribArray( uint32(ret) )
    gl.VertexAttribPointer( uint32(ret), size, gl.FLOAT, false, stride, gl.PtrOffset(offset) )
    return nil
}


func (program *Program) uniformLocation(name UniformName) (int32,error) {
    ret := gl.GetUniformLocation(program.Program, gl.Str(string(name)+"\x00") )
    if ret <= 0 { 
        return -1,log.NewError("no location for uniform '%s' by program %s",name,program.Name)
    }
    return ret,nil;
}

func (program *Program) UniformMatrix4fv(name UniformName, count int32, value *float32) (int32,error) {
    ret,err := program.uniformLocation(name)
    if err != nil {
        return -1,err
    }
    gl.UniformMatrix4fv(ret, count, false, value)
    return ret,nil
}


func (program *Program) Uniform2f(name UniformName, value0, value1 float32) (int32,error) {
    ret,err := program.uniformLocation(name)
    if err != nil {
        return -1,err
    }
    gl.Uniform2f(ret, value0, value1)
    return ret,nil
}

func (program *Program) Uniform2fv(name UniformName, count int32, value *float32) (int32,error) {
    ret,err := program.uniformLocation(name)
    if err != nil {
        return -1,err
    }
    gl.Uniform2fv(ret,count, value)
    return ret,nil
}

func (program *Program) Uniform1f(name UniformName, value float32) (int32,error) {
    ret,err := program.uniformLocation(name)
    if err != nil {
        return -1,err
    }
    gl.Uniform1f(ret, value)
    return ret,nil
}

func (program *Program) Uniform1i(name UniformName, value int32) (int32,error) {
    ret,err := program.uniformLocation(name)
    if err != nil {
        return -1,err
    }
    gl.Uniform1i(ret, value)
    return ret,nil
}

func (program *Program) Desc() string {
    tmp := ""
//    if program.vertexShader != nil { tmp += " " + program.vertexShader.Name }
//    if program.fragmentShader != nil { tmp += "/" + program.fragmentShader.Name }
    return fmt.Sprintf("program[%s%s]",program.Name,tmp)
}


