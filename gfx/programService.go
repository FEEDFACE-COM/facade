//go:build (linux && arm) || DARWIN_GUI
// +build RENDERER

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

const DEBUG_PROGRAMSERVICE = false

type ProgramService struct {
	asset map[string]string

	shaders   map[string]*Shader
	programs  map[string]*Program
	directory string

	refresh []*Shader
}

func NewProgramService(directory string, asset map[string]string) *ProgramService {
	ret := &ProgramService{directory: directory}
	ret.shaders = make(map[string]*Shader)
	ret.programs = make(map[string]*Program)
	ret.refresh = []*Shader{}
	ret.asset = asset
	return ret
}

func (service *ProgramService) GetProgram(name, mode string) *Program {

	if service.programs[name] == nil {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s creating program %s", service.Desc(), name)
		}
		service.programs[name] = NewProgram(name, mode, service)

	}
	return service.programs[name]
}

func (service *ProgramService) GetShader(shaderName string, shaderType ShaderType) (*Shader, error) {
	shaderName = strings.ToLower(shaderName)
	indexName := shaderName + "." + string(shaderType)

	if service.shaders[indexName] == nil {
		var err error
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s loading shader %s", service.Desc(), indexName)
		}
		err = service.LoadShader(shaderName, shaderType)
		if err != nil {
			log.Error("%s fail get shader %s: %s", service.Desc(), shaderName, err)
		}
	}

	if service.shaders[indexName] == nil {
		return nil, log.NewError("no shader named %s", indexName)
	}
	return service.shaders[indexName], nil
}

func watchShaderFile(filePath string, shader *Shader, service *ProgramService) {
	info, err := os.Stat(filePath)
	if err != nil {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s fail stat %s", service.Desc(), shader.Desc())
		}
		return
	}

	log.Info("%s watch %s", service.Desc(), filePath)

	last := info.ModTime()
	for {

		time.Sleep(time.Duration(int64(time.Second)))
		info, err = os.Stat(filePath)
		if err != nil {
			continue
		}
		if info.ModTime().After(last) { // modified

			service.shaderFileChanged(shader, filePath)
			last = info.ModTime()

		}

	}
}

func (service *ProgramService) CheckRefresh() {
	var err error

	for _, shdr := range service.refresh {

		err = shdr.CompileShader()
		if err != nil {
			log.Error("%s fail compile %s: %s", service.Desc(), shdr.Desc(), err)
			continue
		}

		for _, prog := range service.programs {
			if prog.HasShader(shdr) {
				if DEBUG_PROGRAMSERVICE {
					log.Debug("%s refresh %s", service.Desc(), prog.Desc())
				}
				err = prog.Relink()
				if err != nil {
					log.Error("%s fail refresh %s: %s", service.Desc(), prog.Desc(), err)
				}
			}
		}
	}

	service.refresh = []*Shader{}

}

func (service *ProgramService) shaderFileChanged(shader *Shader, filePath string) {

	log.Notice("%s reload %s", service.Desc(), filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.NewError("fail read shader %s from %s: %s", shader.IndexName(), filePath, err)
		return
	}

	shader.LoadSource(string(data))
	service.refresh = append(service.refresh, shader)
}

func (service *ProgramService) LoadShader(shaderName string, shaderType ShaderType) error {
	shaderName = strings.ToLower(shaderName)
	indexName := shaderName + "." + string(shaderType)

	var err error
	var data []byte = []byte{}

	if service.shaders[indexName] != nil {
		return log.NewError("refuse load shader %s already have %s", indexName, service.shaders[indexName].Desc())
	}

	shader := NewShader(shaderName, shaderType)

	{
		var filePath = ""
		filePath, err = service.getFilePathForName(shaderName, shaderType)

		if err == nil { // file found, try reading

			if DEBUG_PROGRAMSERVICE {
				log.Debug("%s read shader %s from %s", service.Desc(), shader.IndexName(), filePath)
			}
			data, err = ioutil.ReadFile(filePath)
			if err != nil {
				return log.NewError("fail read shader %s from %s: %s", shaderName, filePath, err)
			} else {

				go watchShaderFile(filePath, shader, service)

			}

		} else { // no file found, lookup embedded

			if DEBUG_PROGRAMSERVICE {
				log.Debug("%s %s", service.Desc(), err)
			}

			encoded := service.asset[indexName]
			if encoded == "" {
				return log.NewError("no asset data for shader %s", indexName)
			}

			if DEBUG_PROGRAMSERVICE {
				log.Debug("%s decode embedded shader %s", service.Desc(), shaderName)
			}
			data, err = base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return log.NewError("fail decode embedded shader %s: %s", shaderName, err)
			}
		}
	}

	if len(data) <= 0 {
		if DEBUG_PROGRAMSERVICE {
			log.Debug("%s no data for shader %s", service.Desc(), shaderName)
		}
		return log.NewError("no data for shader %s", shaderName)

	}

	shader.LoadSource(string(data))

	if DEBUG_PROGRAMSERVICE {
		log.Debug("%s add shader %s", service.Desc(), shader.IndexName())
	}
	service.shaders[shader.IndexName()] = shader

	return nil
}

func (service *ProgramService) getFilePathForName(shaderName string, shaderType ShaderType) (string, error) {

	ret := service.directory + "/" + shaderName + "." + string(shaderType)
	_, err := os.Stat(ret)
	if os.IsNotExist(err) {
		return "", log.NewError("no file for shader %s", shaderName+"."+string(shaderType))
	} else if err != nil {
		return "", log.NewError("fail stat file %s: %s", shaderName+"."+string(shaderType), err)
	}
	return ret, nil

}

func (service *ProgramService) GetAvailableNames() []string {

	var ret []string
	for n, _ := range service.asset {
		ret = append(ret, fmt.Sprintf("%s", n))
	}
	sort.Strings(ret)
	return ret
}

func (service *ProgramService) Desc() string {
	ret := "programservice["
	ret += fmt.Sprintf("%d", len(service.shaders))
	ret += "]"
	return ret
}
