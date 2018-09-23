
// +build linux,arm


package gfx

import (
    "strings"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Shader struct {
    Name string
    Source string
    Type uint32
    Shader uint32
}





func NewShader(name string, source string, shaderType uint32) *Shader {
    ret := &Shader{Name: name, Source: source, Type: shaderType}
    return ret    
}




func (shader *Shader) CompileShader() error {
    shader.Shader = gl.CreateShader(shader.Type)
    
    sources, free := gl.Strs(shader.Source+"\x00")
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
        log.Error("fail compile shader %s: %s",shader.Name,logs)
        return log.NewError("fail compile shader %s: %s",shader.Name,logs)
    }
    
    
    return nil
}



