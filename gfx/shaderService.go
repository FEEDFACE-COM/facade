
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
    ret.shaders = make( map[ShaderType](map[string]*Shader )
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




func (service *FontService) LoadShader(name string, typ ShaderType) error {
    name = strings.ToLower(name)
    

    var err error
    var data []byte = []byte{}

    {

        path,err := service.getFilePathForName(name)
    
        if err == nil {
            
            if DEBUG_SHADERSERVICE { log.Debug("%s read shader %s.%s from %s",service.Desc(),name,typ,path) }
            data, err = ioutil.ReadFile(path)
            if err != nil {
                return log.NewError("fail to read shader %s.%s from %s: %s",name,typ,path,err)
            }
            
        } else  { // no file found, try embedded

            if DEBUG_SHADERSERVICE { log.Debug("%s no file for shader %s%s: %s",service.Desc(),name,typ,err ) }

            encoded = ""
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


}