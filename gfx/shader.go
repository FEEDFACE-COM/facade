
// +build linux,arm


package gfx

import (
    "strings"
    "errors"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Shader struct {
    Name string
    ShaderSource string
    ShaderType uint32
    Shader uint32
}




func NewShader(name string, source string, shaderType uint32) *Shader {
    ret := &Shader{Name: name, ShaderSource: source, ShaderType: shaderType}
    return ret    
}



//func VertexAttribPointer(program uint32, name AttribName, size int32, stride int32, offset int) uint32 {
//    ret := uint32( gl.GetAttribLocation(program, gl.Str(string(name)+"\x00")) )
//    gl.EnableVertexAttribArray(ret)
//    gl.VertexAttribPointer(ret, size, gl.FLOAT, false, stride, gl.PtrOffset(offset) )
//    return ret
//}


//func UniformMatrix4fv(program uint32, name UniformName, count int32, value *float32) int32 {
//    ret := gl.GetUniformLocation(program, gl.Str(string(name)+"\x00") )
//    gl.UniformMatrix4fv(ret, count, false, value)
//    return ret
//}
//
//func Uniform1i(program uint32, name UniformName, value int32) int32 {
//    ret := gl.GetUniformLocation(program, gl.Str(string(name)+"\x00") )
//    gl.Uniform1i(ret, value)
//    return ret
//}

func (shader *Shader) CompileShader() error {
    log.Debug("shader compile %s",shader.Name)
    shader.Shader = gl.CreateShader(shader.ShaderType)
    
    sources, free := gl.Strs(shader.ShaderSource+"\x00")
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
        return errors.New("fail compile shader")
    }
    
    
    return nil
}




var IDENTITY_VERTEX = `
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertex;
attribute vec2 texcoord;

varying vec2 fragcoord;

void main() {
    fragcoord = texcoord;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
` 



var IDENTITY_FRAGMENT = `
uniform sampler2D texture;

varying vec2 fragcoord;

void main() {
    vec4 tex = texture2D(texture,fragcoord);
    gl_FragColor = tex;
}
`

