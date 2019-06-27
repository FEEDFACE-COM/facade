
package gfx

import (
    "fmt"
    "strings"
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




func (service *ShaderService) LoadShader(name string, typ ShaderType) error {
    name = strings.ToLower(name)
    

    var err error
    var data []byte = []byte{}

    {
        var path = ""
        path,err = service.getFilePathForName(name,typ)
    
        if err == nil {
            
            if DEBUG_SHADERSERVICE { log.Debug("%s read shader %s.%s from %s",service.Desc(),name,typ,path) }
            data, err = ioutil.ReadFile(path)
            if err != nil {
                return log.NewError("fail to read shader %s.%s from %s: %s",name,typ,path,err)
            }
            
        } else  { // no file found, try embedded

            if DEBUG_SHADERSERVICE { log.Debug("%s no file for shader %s%s: %s",service.Desc(),name,typ,err ) }

            encoded := ""
            switch typ {
                case VertType:
                    encoded = VertexShaderAsset[name]
                case FragType: 
                    encoded = FragmentShaderAsset[name]
            }

            if encoded == "" {
                return log.NewError("no data for embedded shader %s.%s",name,typ)    
            }

            if DEBUG_SHADERSERVICE { log.Debug("%s decode embedded shader %s.%s",service.Desc(),name,typ) }
            data,err = base64.StdEncoding.DecodeString( encoded )
            if err != nil {
                return log.NewError("fail to decode embedded shader %s.%s: %s",name,typ,err)    
            }
        }    
    }
    
    if len(data) <= 0 {
        if DEBUG_SHADERSERVICE { log.Debug("%s no data for shader %s.%s",service.Desc(),name,typ) }
        return log.NewError("no data for shader %s.%s",name,typ)
        
    }
    
    shader := service.shaders[typ][name]
    if  shader == nil {
        shader = NewShader(name,typ)    
    }
    
    err = shader.loadData( data )
    if err != nil {
        log.Debug("%s fail load shader %s.%s data: %s",service.Desc(),name,typ,err)
        return log.NewError("fail to load shader %s.%s data",name,typ)
    }


    if DEBUG_FONTSERVICE { log.Debug("%s add shader %s.%s",service.Desc(),name,typ) }
    service.shaders[typ][name] = shader
    
    

    return nil
}

func (service *ShaderService) getFilePathForName(name string, typ ShaderType) (string,error) {
    
    return (service.directory + "/" + name + "." + string(typ)), nil
    
}



func (service *ShaderService) Desc() string {
    ret := "shaderservice["
    ret += fmt.Sprintf("%d verts ",len(service.shaders[VertType]))
    ret += fmt.Sprintf("%d frags",len(service.shaders[FragType]))
    ret += "]"
    return ret
}


