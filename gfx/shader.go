
// +build linux,arm


package gfx

import (
    "fmt"
    "strings"
    "io/ioutil"
    "time"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Shader struct {
    Name string
    Source string
    Type ShaderType
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


type ShaderType string
const (
    VERTEX_SHADER   ShaderType = "vert"
    FRAGMENT_SHADER ShaderType = "frag"
)

var shaderDirectory string
func SetShaderDirectory(directory string) { shaderDirectory = directory }

func loadShader(shaderName string, shaderType ShaderType) (string, error) {
    var data []byte
    var err error
    filePath := fmt.Sprintf("%s/%s.%s",shaderDirectory,shaderName,string(shaderType))
    data, err = ioutil.ReadFile(filePath)
    if err != nil {
        log.Error("fail read shader file %s: %s",filePath,err)
        return "", log.NewError("fail read shader file: %s",err)
    }
    log.Debug("read shader file %s",filePath)
    
    go func(){
        time.Sleep( time.Duration( int64(time.Second*2)) )
        log.Debug("still here watching %s",filePath)
    }()
    
    return string(data), nil    
}

func GetShader(name string, shaderType ShaderType) (*Shader,error) {
    var ret *Shader = nil
    src, err := loadShader(name, shaderType)
    if err == nil {
        
        ret = NewShader(name, src, shaderType)
        //TODO: register callback
        return ret,nil
        
    } else {
        src := ""
        switch (shaderType) {
            case VERTEX_SHADER:   src = VertexShader[name]
            case FRAGMENT_SHADER: src = FragmentShader[name]
        }
        
        if src == "" {
            return nil, log.NewError("fail get shader %s from file and map",name)
        }
        
        ret := NewShader(name, src, shaderType)
        return ret,nil
    }
    
    
}




func NewShader(name string, source string, shaderType ShaderType) *Shader {
    ret := &Shader{Name: name, Source: source, Type: shaderType}
    return ret    
}



func (shader *Shader) CompileShader() error {
    switch shader.Type {
        case VERTEX_SHADER:    shader.Shader = gl.CreateShader(gl.VERTEX_SHADER)
        case FRAGMENT_SHADER:  shader.Shader = gl.CreateShader(gl.FRAGMENT_SHADER)
    }
        
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



