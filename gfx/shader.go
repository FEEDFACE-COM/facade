
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

const DEBUG_SHADER = false

type Shader struct {
    Name string
    Source string
    Type ShaderType
    Shader uint32
}


type UniformName string
const (
    DEBUGFLAG   UniformName = "debugFlag"
    PROJECTION  UniformName = "projection"
    MODEL       UniformName = "model"
    VIEW        UniformName = "view"
    SCREENRATIO UniformName = "ratio"
    TEXTURE     UniformName = "texture"
    TILECOUNT   UniformName = "tileCount"
    TILESIZE    UniformName = "tileSize"
    TILEOFFSET  UniformName = "tileOffset"
    CURSORPOS   UniformName = "cursorPos"
    SCROLLER    UniformName = "scroller"
    DOWNWARD    UniformName = "downward"
    CLOCKNOW    UniformName = "now"
)    

type AttribName string
const (
    VERTEX     AttribName = "vertex"    
    COLOR      AttribName = "color"    
    TEXCOORD   AttribName = "texCoord"
    TILECOORD  AttribName = "tileCoord"
    GRIDCOORD  AttribName = "gridCoord"
)


type ShaderType string
const (
    VertType ShaderType = "vert"
    FragType ShaderType = "frag"
)

var shaderPath string
func SetShaderDirectory(directory string) { shaderPath = path.Clean(directory) }

func (shader *Shader) loadSource(src string) error {
//    var err error

    shader.Source = src
    if DEBUG_FONTSERVICE { log.Debug("%s shader setup",shader.Desc()) }
    return nil



}


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
//    log.Debug("read %s",filePath)
    return string(data), nil
}


func shaderFilePath(shaderName string, shaderType ShaderType) string {
    return path.Clean( fmt.Sprintf("%s/%s.%s",shaderPath,shaderName,string(shaderType)) )    
}



func GetShader(shaderPrefix string, shaderName string, shaderType ShaderType, program *Program) (*Shader,error) {
    var ret *Shader = nil
    filePath := shaderFilePath(shaderPrefix+shaderName,shaderType)
    src, err := loadShaderFile(filePath)
    if err == nil {
        
        ret = NewShader(shaderName, shaderType)
        ret.loadSource( src )
        go WatchShaderFile(filePath, program, ret)
        return ret,nil
        
    } else {

        src = ShaderAsset[shaderPrefix+shaderName+"."+string(shaderType)]
        
        if src == "" {
            return nil, log.NewError("fail get %s shader %s from file and map",string(shaderType),shaderName)
        }
        
        ret := NewShader(shaderPrefix+shaderName, shaderType)
        ret.loadSource( src )
        return ret,nil
    }
    
    
}



func WatchShaderFile(filePath string, program *Program, shader *Shader) {
        // REM, should verify that we're not alreay watching this file..
        if DEBUG_SHADER { log.Debug("watch %s",filePath) }
        

        info,err := os.Stat(filePath)
        if err != nil {
            log.Debug("fail stat %s",shader.Desc())
            return
        }
        last := info.ModTime()
        for {
            
            time.Sleep( time.Duration( int64(time.Second)) )
            info,err = os.Stat(filePath)
            if err != nil { continue }
            if info.ModTime().After( last ) {
                if ! program.HasShader(shader) {
                    break
                }
                shader.Source, err = loadShaderFile(filePath)       
//                log.Debug("alert %s %s",program.Desc(),shader.Desc())
                refreshChan <-Refresh{program: program, shader: shader}
                last = info.ModTime()
            }
            
        }
    if DEBUG_SHADER { log.Debug("unwatch %s",filePath) }
}



func (shader *Shader) Desc() string {
    return fmt.Sprintf("%s[%s]",shader.Type,shader.Name)
}

func NewShader(name string, shaderType ShaderType) *Shader {
    ret := &Shader{Name: name, Type: shaderType}
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
            case VertType:    shader.Shader = gl.CreateShader(gl.VERTEX_SHADER)
            case FragType:  shader.Shader = gl.CreateShader(gl.FRAGMENT_SHADER)
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



