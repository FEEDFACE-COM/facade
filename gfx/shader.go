
// +build linux,arm


package gfx

import (
    "fmt"
    "os"
    "path"
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
    DOWNWARD   UniformName = "downwardFlag"
    SCROLLER   UniformName = "scroller"
    TIMER      UniformName = "timer"
)    

type AttribName string
const (
    VERTEX     AttribName = "vertex"    
    COLOR      AttribName = "color"    
    TEXCOORD   AttribName = "texCoord"
    TILECOORD  AttribName = "tileCoord"
)


type ShaderType string
const (
    VERTEX_SHADER   ShaderType = "vert"
    FRAGMENT_SHADER ShaderType = "frag"
)

var shaderPath string
func SetShaderDirectory(path string) { shaderPath = path.Clean(path) }


func (shader *Shader) LoadShader(filePath string) error {
    var err error
    shader.Source, err = loadShaderFile(filePath)       
    if err != nil { return err }
    return nil
}

func loadShaderFile(filePath string) (string, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", log.NewError("fail read shader file: %s",err)
    }
    log.Debug("read shader file %s",filePath)
    return string(data), nil
}




func GetShader(name string, shaderType ShaderType, program *Program) (*Shader,error) {
    var ret *Shader = nil
    filePath := fmt.Sprintf("%s/%s.%s",shaderPath,name,string(shaderType))
    src, err := loadShaderFile(filePath)
    if err == nil {
        
        ret = NewShader(name, src, shaderType)
        go WatchShaderFile(filePath, program, ret)
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

func WatchShaderFile(path string, program *Program, shader *Shader) {

        

        info,err := os.Stat(path)
        if err != nil {
            log.Debug("fail stat %s",shader.Desc())
            return
        }
        last := info.ModTime()
        for {
            
            time.Sleep( time.Duration( int64(time.Second)) )
            info,err = os.Stat(path)
            if err != nil { continue }
            if info.ModTime().After( last ) {
                log.Debug("alert %s %s",program.Desc(),shader.Desc())
                shader.LoadShader(path)
                refreshChan <-Refresh{program: program, shader: shader}
                last = info.ModTime()
            }
            
        }


        
        
        
        
}

func (shader *Shader) Desc() string {
    return fmt.Sprintf("shader[%s.%s]",shader.Name,shader.Type)
}

func NewShader(name string, source string, shaderType ShaderType) *Shader {
    ret := &Shader{Name: name, Source: source, Type: shaderType}
    return ret    
}

func (shader *Shader) Close() {
    if shader.Shader > 0 {
        log.Debug("delete %s",shader.Desc())    
    	gl.DeleteShader(shader.Shader)
    }   
}

func (shader *Shader) CompileShader() error {
    if shader.Shader <= 0 {
        switch shader.Type {
            case VERTEX_SHADER:    shader.Shader = gl.CreateShader(gl.VERTEX_SHADER)
            case FRAGMENT_SHADER:  shader.Shader = gl.CreateShader(gl.FRAGMENT_SHADER)
        }
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



