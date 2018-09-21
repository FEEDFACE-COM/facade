
// +build linux,arm


package gfx

import (
    "strings"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Shader struct {
    Name string
    ShaderSource string
    ShaderType uint32
    Shader uint32
}


func NewShader(name string, source string, shaderType uint32) Shader {
    ret := Shader{Name: name, ShaderSource: source, ShaderType: shaderType}
    return ret    
}

func (shader *Shader) Compile() {
    log.Debug("shader compile %s",shader.Name)
    shader.Shader = gl.CreateShader(shader.ShaderType)
    
    sources, free := gl.Strs(shader.ShaderSource)
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
    }
    
    
}


var IDENTITY_VERTEX = `
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
attribute vec3 vert;
attribute vec2 vertTexCoord;
varying vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"



var IDENTITY_FRAGMENT = `
uniform sampler2D tex;
varying vec2 fragTexCoord;
void main() {
    gl_FragColor = texture2D(tex,fragTexCoord);
}
` + "\x00"

