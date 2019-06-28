
package gfx

import (
    "fmt"
    "strings"
    "os"
    "io/ioutil"
    "encoding/base64"
    "sort"
    log "../log"       
    
)

const DEBUG_SHADERSERVICE = true


type ShaderService struct {

    shaders map[string]*Shader
    directory string    

}




func NewShaderService(directory string) *ShaderService {
    ret := &ShaderService{directory: directory}
    ret.shaders = make( map[string]*Shader )
    return ret
}


func (service *ShaderService) GetShader(shaderName string, shaderType ShaderType) (*Shader,error) {
    shaderName = strings.ToLower(shaderName)
    indexName := shaderName + "." + string(shaderType)

    if service.shaders[indexName] == nil {
        var err error
        if DEBUG_SHADERSERVICE { log.Debug("%s shader %s not loaded",service.Desc(),indexName) }
        err = service.LoadShader(shaderName,shaderType)
        if err != nil {
            log.Error("%s fail get shader %s: %s",service.Desc(),shaderName,err)
        }
    }

    if service.shaders[indexName] == nil {
        return nil, log.NewError("no %s shader named %s",shaderType,shaderName)
    }
    return service.shaders[indexName], nil
}




func (service *ShaderService) LoadShader(shaderName string, shaderType ShaderType) error {
    shaderName = strings.ToLower(shaderName)
    indexName := shaderName + "." + string(shaderType)

    var err error
    var data []byte = []byte{}
    
    if service.shaders[indexName] != nil {
        return log.NewError("refuse load %s shader %s already have %s",shaderType,shaderName,service.shaders[indexName].Desc())     
    }
    
    {
        var path = ""
        path,err = service.getFilePathForName(shaderName,shaderType)
    
        if err == nil { // file found, try reading
            
            if DEBUG_SHADERSERVICE { log.Debug("%s read %s shader %s from %s",service.Desc(),string(shaderType),shaderName,path) }
            data, err = ioutil.ReadFile(path)
            if err != nil {
                return log.NewError("fail read shader %s from %s: %s",shaderName,path,err)
            }
            
        } else  { // no file found, lookup embedded

            if DEBUG_SHADERSERVICE { log.Debug("%s %s",service.Desc(),err ) }

            encoded := ShaderAsset[shaderName]
            if encoded == "" {
                return log.NewError("no asset data for shader %s",shaderName)    
            }

            if DEBUG_SHADERSERVICE { log.Debug("%s decode embedded shader %s",service.Desc(),shaderName) }
            data,err = base64.StdEncoding.DecodeString( encoded )
            if err != nil {
                return log.NewError("fail decode embedded shader %s: %s",shaderName,err)    
            }
        }    
    }
    
    if len(data) <= 0 {
        if DEBUG_SHADERSERVICE { log.Debug("%s no data for shader %s",service.Desc(),shaderName) }
        return log.NewError("no data for shader %s",shaderName)
        
    }
    
    shader := NewShader(shaderName,shaderType)    
    
    err = shader.loadSource( string(data) )
    if err != nil {
        log.Debug("%s fail load %s shader %s data: %s",service.Desc(),shaderType,shaderName,err)
        return log.NewError("fail load %s shader %s data",shaderType,shaderName)
    }


    if DEBUG_SHADERSERVICE { log.Debug("%s add shader %s",service.Desc(),indexName) }
    service.shaders[indexName] = shader
    

    return nil
}

func (service *ShaderService) getFilePathForName(shaderName string, shaderType ShaderType) (string,error) {

    ret := service.directory + "/" + shaderName + "." + string(shaderType)
    _,err := os.Stat(ret)
    if os.IsNotExist(err) {
        return "", log.NewError("no file for shader %s",shaderName+"."+string(shaderType) )    
    } else if err != nil {
        return "", log.NewError("fail stat file %s: %s",shaderName+"."+string(shaderType),err )
    }
    return ret, nil

    
    
}

func (service *ShaderService) GetAvailableNames() []string {

    var ret = GetShaderAssetNames()

//    files, err := ioutil.ReadDir(service.directory)
//    if err != nil {
//        return ret
//    }
//    var Extensions []string = []string{ "."+string(VertType), "."+string(FragType) }
//
//    for _, f := range files {
//        for _, ext := range Extensions {
//            name := strings.ToLower( f.Name() )
//            if strings.HasSuffix( name, ext ) {
//                name = strings.TrimSuffix( name, ext )
//                if ShaderAsset[name] == "" {
//
//                    ret = append(ret, name)
//
//                }                
//            }
//        }
//    }
    return ret
}


func GetShaderAssetNames() []string {
    var ret []string
    for n, _ := range ShaderAsset {
        ret = append(ret,fmt.Sprintf("%s",n)) 
    }
    sort.Strings(ret)
    return ret
}



func (service *ShaderService) Desc() string {
    ret := "shaderservice["
    ret += fmt.Sprintf("%d",len(service.shaders))
//    ret += fmt.Sprintf("%d",len(service.shaders[FragType]))
    ret += "]"
    return ret
}


