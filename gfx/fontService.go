package gfx

import (
	log "../log"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"strings"
)

const DEBUG_FONTSERVICE = false

var ForegroundColor = image.White
var BackgroundColor = image.Transparent
var DebugColor = image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255})

var Extensions = []string{".ttf", ".ttc"}

type FontService struct {
	fonts   map[string]*Font
	scratch *image.RGBA
	asset   map[string]string

	directory string
}

func NewFontService(directory string, asset map[string]string) *FontService {
	ret := &FontService{directory: directory}
	ret.fonts = make(map[string]*Font)
	ret.scratch = image.NewRGBA(image.Rect(0, 0, FontScratchSize, FontScratchSize))
	ret.asset = asset
	return ret
}

func (service *FontService) GetFont(fontName string) (*Font, error) {
	fontName = strings.ToLower(fontName)

	if service.fonts[fontName] == nil {
		var err error
		if DEBUG_FONTSERVICE {
			log.Debug("%s loading font %s", service.Desc(), fontName)
		}
		err = service.LoadFont(fontName)
		if err != nil {
			log.Error("%s fail get font %s: %s", service.Desc(), fontName, err)
		}
	}

	if service.fonts[fontName] == nil {
		return nil, log.NewError("no font named %s", fontName)
	}
	return service.fonts[fontName], nil

}

func (service *FontService) LoadFont(name string) error {
	name = strings.ToLower(name)

	var err error
	var data []byte = []byte{}

	if font := service.fonts[name]; font != nil {
		return nil
		//        return log.NewError("font %s already loaded: %s",name,font.Desc())
	}

	{

		path, err := service.getFilePathForName(name)

		if err == nil {

			if DEBUG_FONTSERVICE {
				log.Debug("%s read font %s from %s", service.Desc(), name, path)
			}
			data, err = ioutil.ReadFile(path)
			if err != nil {
				return log.NewError("fail to read font from %s: %s", path, err)
			}

		} else { // no file found, try asset

			if DEBUG_FONTSERVICE {
				log.Debug("%s no file for font %s: %s", service.Desc(), name, err)
			}

			encoded := service.asset[name]
			if encoded == "" {
				return log.NewError("no asset data for font %s", name)
			}

			if DEBUG_FONTSERVICE {
				log.Debug("%s decode font %s", service.Desc(), name)
			}
			data, err = base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return log.NewError("fail to decode font %s: %s", name, err)
			}
		}
	}

	if len(data) <= 0 {
		if DEBUG_FONTSERVICE {
			log.Debug("%s no data for font %s", service.Desc(), name)
		}
		return log.NewError("no data for font %s", name)

	}

	font := NewFont(name, service.scratch)
	err = font.loadData(data)
	if err != nil {
		if DEBUG_FONTSERVICE {
			log.Debug("%s fail load font %s data: %s", service.Desc(), name, err)
		}
		return log.NewError("fail to load font %s data", name)
	}

	if DEBUG_FONTSERVICE {
		log.Debug("%s add font %s", service.Desc(), name)
	}
	service.fonts[name] = font

	return nil

}

func (service *FontService) getFilePathForName(fontName string) (string, error) {
	var err error
	files, err := ioutil.ReadDir(service.directory)
	if err != nil {
		return "", log.NewError("fail list fonts: %s", err)
	}
	for _, f := range files {
		for _, ext := range Extensions {
			if strings.ToLower(f.Name()) == strings.ToLower(fontName+ext) {
				if DEBUG_FONTSERVICE {
					log.Debug("%s found font file %s", service.Desc(), f.Name())
				}
				return service.directory + "/" + f.Name(), nil
			}
		}
	}
	return "", log.NewError("no file for font %s", fontName)
}

func (service *FontService) Desc() string {
	ret := "fontservice["
	ret += fmt.Sprintf("%d", len(service.fonts))
	ret += "]"
	return ret
}
