
package gfx

import (
    "fmt"
    "strings"
    "os"
    "io/ioutil"
    "encoding/base64"
    log "../log"       
    
)

const DEBUG_SHADERSERVICE = true


type ShaderService struct {

    shaders map[ShaderType](map[string]*Shader)
    
    directory string    
    
}


func NewShaderService(directory string) *ShaderService {
    ret := &ShaderService{directory: directory}
    ret.shaders = make( map[ShaderType](map[string]*Shader) )
    ret.shaders[VertType] = make(map[string]*Shader)
    ret.shaders[FragType] = make(map[string]*Shader)
    return ret
}


func (service *ShaderService) GetShader(name string, typ ShaderType) (*Shader,error) {
    name = strings.ToLower(name)
    ret := service.shaders[typ][name]
    if ret == nil {
        return nil, log.NewError("no %s shader for name %s",typ,name)
    }
    return ret, nil
}




func (service *ShaderService) LoadShader(shaderName string, shaderType ShaderType) error {
    shaderName = strings.ToLower(shaderName)
    name := shaderName + "." + string(shaderType)

    var err error
    var data []byte = []byte{}
    

    {
        var path = ""
        path,err = service.getFilePathForName(shaderName,shaderType)
    
        if err == nil { // file found, try reading
            
            if DEBUG_SHADERSERVICE { log.Debug("%s read shader %s from %s",service.Desc(),name,path) }
            data, err = ioutil.ReadFile(path)
            if err != nil {
                return log.NewError("fail to read shader %s from %s: %s",name,path,err)
            }
            
        } else  { // no file found, lookup embedded

            if DEBUG_SHADERSERVICE { log.Debug("%s %s",service.Desc(),err ) }

            encoded := ShaderAsset[name]
            if encoded == "" {
                return log.NewError("no asset data for shader %s",name)    
            }

            if DEBUG_SHADERSERVICE { log.Debug("%s decode embedded shader %s",service.Desc(),name) }
            data,err = base64.StdEncoding.DecodeString( encoded )
            if err != nil {
                return log.NewError("fail to decode embedded shader %s: %s",name,err)    
            }
        }    
    }
    
    if len(data) <= 0 {
        if DEBUG_SHADERSERVICE { log.Debug("%s no data for shader %s",service.Desc(),name) }
        return log.NewError("no data for shader %s",name)
        
    }
    
    shader := service.shaders[shaderType][name]
    if  shader == nil {
        shader = NewShader(name,shaderType)    
    }
    
    err = shader.loadData( data )
    if err != nil {
        log.Debug("%s fail load shader %s data: %s",service.Desc(),name,err)
        return log.NewError("fail to load shader %s data",name)
    }


    if service.shaders[shaderType][name] != shader {
        if DEBUG_SHADERSERVICE { log.Debug("%s add shader %s",service.Desc(),name) }
        service.shaders[shaderType][name] = shader
    }
    
    

    return nil
}

func (service *ShaderService) getFilePathForName(name string, typ ShaderType) (string,error) {

    ret := service.directory + "/" + name + "." + string(typ)
    _,err := os.Stat(ret)
    if os.IsNotExist(err) {
        return "", log.NewError("no file for shader %s",name+"."+string(typ) )    
    } else if err != nil {
        return "", log.NewError("fail stat file %s: %s",name+"."+string(typ),err )
    }
    return ret, nil

    
    
}



func (service *ShaderService) Desc() string {
    ret := "shaderservice["
    ret += fmt.Sprintf("%d verts ",len(service.shaders[VertType]))
    ret += fmt.Sprintf("%d frags",len(service.shaders[FragType]))
    ret += "]"
    return ret
}


