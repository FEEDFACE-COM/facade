
// +build linux,arm


package gfx

import (
    "strings"
    "io/ioutil"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Shader struct {
    Name string
    Source string
    Type uint32
    Shader uint32
}


type UniformName string
const (
    DEBUGFLAG  UniformName = "debugFlag"
    PROJECTION UniformName = "projection"
    MODEL      UniformName = "model"
    VIEW       UniformName = "view"
    TEXTURE    UniformName = "texture"
    TILECOUNT  UniformName = "tileCount"
    TILESIZE   UniformName = "tileSize"
    SCROLLER   UniformName = "scroller"
    TIMER      UniformName = "timer"
)    

type AttribName string
const (
    VERTEX     AttribName = "vertex"    
    COLOR      AttribName = "color"    
    TEXCOORD   AttribName = "texcoord"
    TILECOORD  AttribName = "tileCoord"
    TOTALWIDTH AttribName = "totalWidth"
)

var shaderDirectory string

func SetShaderDirectory(directory string) { shaderDirectory = directory }

func loadShaderFile(shaderName string, shaderType uint32) (string, error) {
    var data []byte
    var err error
    ext := ""
    switch (shaderType) {
        case gl.VERTEX_SHADER:   ext = ".vert"
        case gl.FRAGMENT_SHADER: ext = ".frag"
    }
    filePath := shaderDirectory + shaderName + ext
    data, err = ioutil.ReadFile(filePath)
    if err != nil {
        log.Error("fail read shader file %s: %s",filePath,err)
        return "", log.NewError("fail read shader file: %s",err)
    }
    log.Debug("read shader file %s",filePath)
    return string(data), nil    
}

func GetShader(name string, shaderType uint32) (*Shader,error) {
    var ret *Shader = nil
    src, err := loadShaderFile(name, shaderType)
    if err == nil {
        
        ret = NewShader(name, src, shaderType)
        //TODO: register callback
        return ret,nil
        
    } else {
        src := ""
        switch (shaderType) {
            case gl.VERTEX_SHADER:   src = VertexShader[name]
            case gl.FRAGMENT_SHADER: src = FragmentShader[name]
        }
        
        if src == "" {
            return nil, log.NewError("fail get shader %s from file and map",name)
        }
        
        ret := NewShader(name, src, shaderType)
        return ret,nil
    }
    
    
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



