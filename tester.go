package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"sync"
	"time"

	facade "./facade"
	gfx "./gfx"
	log "./log"
	//    proto "./facade/proto"
)

type Tester struct {
	Terminal bool

	mode  facade.Mode
	debug bool

	font       *gfx.Font
	vert, frag *gfx.Shader

	shaderConfig *facade.ShaderConfig
	gridConfig   *facade.GridConfig

	lineBuffer *facade.LineBuffer
	termBuffer *facade.TermBuffer
	wordBuffer *facade.WordBuffer

	fontService    *gfx.FontService
	programService *gfx.ProgramService

	mutex     *sync.Mutex
	directory string

	prevFrame gfx.ClockFrame

	image *image.RGBA

	refreshChan chan bool
}

func NewTester(directory string) *Tester {
	ret := &Tester{directory: directory}
	ret.mutex = &sync.Mutex{}
	ret.refreshChan = make(chan bool, 1)
	return ret
}

func (tester *Tester) switchShader(shaderName string, shaderType gfx.ShaderType) error {
	var err error

	shaderName = "grid" + "/" + strings.ToLower(shaderName)
	name := shaderName // + "." + string(shaderType)

	log.Debug("tester load shader %s", name)
	err = tester.programService.LoadShader(shaderName, shaderType)
	if err != nil {
		log.Error("tester fail load shader %s: %s", name, err)
		//return log.NewError("tester fail load shader %s: %s", name, err)
	}

	var shader *gfx.Shader
	shader, err = tester.programService.GetShader(name, shaderType)
	if err != nil {
		log.PANIC("tester fail get shader %s: %s", name, err)
		return log.NewError("tester fail get shader %s: %s", name, err)
	}

	log.Debug("tester switch to shader %s", name)
	switch shaderType {
	case gfx.VertType:
		tester.vert = shader
	case gfx.FragType:
		tester.frag = shader
	}
	return nil
}

func (tester *Tester) switchFont(name string) error {
	var err error

	if name != tester.font.GetName() {

		log.Debug("tester load font %s", name)
		err = tester.fontService.LoadFont(name)
		if err != nil {
			log.Debug("tester fail load font %s: %s", name, err)
			return log.NewError("tester fail load font %s: %s", name, err)
		}

		var font *gfx.Font
		font, err = tester.fontService.GetFont(name)
		if err != nil {
			log.Debug("tester fail get font %s: %s", name, err)
			return log.NewError("tester fail get font %s: %s", name, err)
		}

		log.Debug("tester switch to font %s", name)
		tester.font = font

	}

	return nil
}

func (tester *Tester) Init() error {
	var err error

	log.Debug("init tester[%s]", tester.directory)

	if strings.HasPrefix(tester.directory, "~/") {
		tester.directory = os.Getenv("HOME") + tester.directory[1:]
	}

	tester.shaderConfig = &facade.ShaderDefaults
	tester.shaderConfig.SetFrag = true
	tester.shaderConfig.SetVert = true

	tester.gridConfig = &facade.GridDefaults
	tester.gridConfig.SetWidth = true
	tester.gridConfig.SetHeight = true

	tester.fontService = gfx.NewFontService(tester.directory+"/font", facade.FontAsset)
	tester.programService = gfx.NewProgramService(tester.directory+"/shader", facade.ShaderAsset)

	tester.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width), uint(facade.GridDefaults.Height))
	tester.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height), uint(facade.LineDefaults.Buffer), tester.refreshChan)
	tester.wordBuffer = facade.NewWordBuffer(tester.refreshChan)

	err = tester.switchFont(facade.FontDefaults.Name)
	if err != nil {
		log.PANIC("tester missing default font: %s", err)
	}

	err = tester.switchShader(facade.ShaderDefaults.Vert, gfx.VertType)
	if err != nil {
		log.PANIC("tester missing default vert shader: %s", err)
	}

	err = tester.switchShader(facade.ShaderDefaults.Frag, gfx.FragType)
	if err != nil {
		log.PANIC("tester missing default frag shader: %s", err)
	}

	gfx.WorldClock().Reset()
	return nil
}

func (tester *Tester) Configure(config *facade.Config) error {
	var err error
	if config == nil {
		return nil
	}
	log.Debug("tester config %s", config.Desc())

	if config.GetSetDebug() {
		tester.debug = config.GetDebug()
	} else {
		tester.debug = false
	}

	if cfg := config.GetFont(); cfg != nil {

		if cfg.GetSetName() && cfg.GetName() != tester.font.GetName() {
			err = tester.switchFont(cfg.GetName())
			if err != nil {
				log.Error("tester fail switch font: %s", err)
			}
		}
	}

	var shader *facade.ShaderConfig = nil
	var grid *facade.GridConfig = nil

	if terminal := config.GetTerminal(); terminal != nil {

		if terminal.GetShader() != nil {
			shader = terminal.GetShader()
		}

		if terminal.GetGrid() != nil {
			grid = terminal.GetGrid()
		}

	}

	if lines := config.GetLines(); lines != nil {

		if lines.GetShader() != nil {
			shader = lines.GetShader()
		}

		if lines.GetGrid() != nil {
			grid = lines.GetGrid()
		}

		if lines.GetSetBuffer() {
			tester.lineBuffer.Resize(uint(tester.gridConfig.GetHeight()), uint(lines.GetBuffer()))
		}

		if lines.GetSetSpeed() {
			tester.lineBuffer.SetSpeed(float32(lines.GetSpeed()))
		}

		if lines.GetSetFixed() {
			tester.lineBuffer.Fixed = lines.GetFixed()
		}

		if lines.GetSetDrop() {
			tester.lineBuffer.Drop = lines.GetDrop()
		}

		if lines.GetSetSmooth() {
			tester.lineBuffer.Smooth = lines.GetSmooth()
		}

	}

	if shader != nil {

		if shader.GetSetVert() {
			err = tester.switchShader(shader.GetVert(), gfx.VertType)
			if err != nil {
				log.Error("tester fail switch shader: %s", err)
			} else {
				tester.shaderConfig.Vert = shader.GetVert()
			}
		}

		if shader.GetSetFrag() {
			err = tester.switchShader(shader.GetFrag(), gfx.FragType)
			if err != nil {
				log.Error("tester fail switch shader: %s", err)
			} else {
				tester.shaderConfig.Frag = shader.GetFrag()
			}
		}

	}

	if grid != nil {

		if grid.GetSetWidth() {
			tester.gridConfig.Width = grid.GetWidth()
		}

		if grid.GetSetHeight() {
			tester.gridConfig.Height = grid.GetHeight()
		}

		if grid.GetSetWidth() || grid.GetSetHeight() {
			tester.termBuffer.Resize(uint(tester.gridConfig.Width), uint(tester.gridConfig.Height))
		}

		if grid.GetSetHeight() {
			tester.lineBuffer.Resize(uint(tester.gridConfig.Height), uint(tester.lineBuffer.GetBuffer()))
		}

		if grid.GetSetFill() {

			if err := tester.render(grid.GetFill()); err != nil {
				log.Error("fail render '%': %s", grid.GetFill(), err)
			}

		}

	}

	if config.GetSetMode() {

		if tester.mode != config.GetMode() {
			log.Info("switch mode %s", config.GetMode())
			tester.mode = config.GetMode()
		}

	}

	return nil

}

func (tester *Tester) render(fill string) error {
	if tester.font != nil {
		log.Debug("tester render '%s' with %s", fill, tester.font.Desc())
	}
	var err error

	if fill == "" { //render out glyphmap

		tester.image, err = tester.font.RenderMap(tester.debug)

	} else { // render out given string

		tester.image, err = tester.font.RenderText(fill, tester.debug)
	}

	if tester.image == nil {
		return log.NewError("fail render '%s' with %s: %s", fill, tester.font.Desc(), err)
	}
	log.Debug("tester rendered '%s' with %s", fill, tester.font.Desc())
	return nil

}

//rem, should not need this, can do directly?
func (tester *Tester) ProcessConf(confChan chan facade.Config) {
	select {
	case conf := <-confChan:
		tester.Configure(&conf)

	default:
		//nop
	}
}

func (tester *Tester) ProcessTextSeqs(bufChan chan facade.TextSeq) error {

	for {
		item := <-bufChan
		text, seq := item.Text, item.Seq
		if text != nil && len(text) > 0 {
			tester.lineBuffer.ProcessRunes(text)
			tester.termBuffer.ProcessRunes(text)
			tester.wordBuffer.ProcessRunes(text)
		}
		if seq != nil {
			tester.lineBuffer.ProcessSequence(seq)
			tester.termBuffer.ProcessSequence(seq)
			tester.wordBuffer.ProcessSequence(seq)
		}
	}
	return nil
}

func (tester *Tester) ProcessRawConfs(rawChan chan facade.Config, confChan chan facade.Config) error {
	for {
		rawConf := <-rawChan

		//        tester.mutex.Lock()
		//        // prep some stuff i guess?
		//        tester.mutex.Unlock()

		confChan <- rawConf

	}
	return nil
}

func (tester *Tester) InfoMode() string {
	mode := ""
	switch tester.mode {
	case facade.Mode_TERM:
		mode = "term " + tester.shaderConfig.Desc() + " " + tester.gridConfig.Desc()
	case facade.Mode_LINES:
		mode = "line " + tester.shaderConfig.Desc() + " " + tester.gridConfig.Desc()
	case facade.Mode_TAGS:
		mode = "tags " + tester.shaderConfig.Desc()
	case facade.Mode_WORDS:
		mode = "words " + tester.shaderConfig.Desc()
	}
	dbg := ""
	if tester.debug {
		dbg = " DEBUG"
	}
	return fmt.Sprintf("%s\n  %s%s", mode, tester.font.Desc(), dbg)

}

func (tester *Tester) Test(confChan chan facade.Config) error {
	const FRAME_RATE = 60.
	tester.prevFrame = gfx.WorldClock().Frame()

	for {

		gfx.WorldClock().Tick()
		verboseFrame := gfx.WorldClock().VerboseFrame()

		tester.mutex.Lock()

		tester.ProcessConf(confChan)
		tester.programService.CheckRefresh()

		if DEBUG_PERIODIC && verboseFrame {

			tester.printDebug()
			tester.prevFrame = gfx.WorldClock().Frame()
		}

		tester.mutex.Unlock()

		if tester.image != nil {

			outPath := fmt.Sprintf("./%s.png",tester.font.GetName())
			log.Info("render image to %s", outPath)
			outFile, err := os.Create(outPath)
			if err != nil {
				log.Error("fail to create file %s: %s", outPath, err)
			}
			writer := bufio.NewWriter(outFile)
			if err := png.Encode(writer, tester.image); err != nil {
				log.Error("fail to encode rendered image: %s", err)
				return log.NewError("fail to encode rendered image: %s", err)
			}
			writer.Flush()
			outFile.Close()
			tester.image = nil

		}

		time.Sleep(time.Duration(int64(time.Second / FRAME_RATE)))

	}

	return nil
}

func (tester *Tester) Info() string {
	ret := ""

	ret += InfoAuthor()
	ret += InfoVersion()
	ret += "\n\n"

	ret += gfx.WorldClock().Info(tester.prevFrame)
	ret += "\n  " + tester.InfoMode()
	ret += "\n  " + tester.lineBuffer.Desc()
	ret += "\n  " + tester.termBuffer.Desc()
	ret += "\n  " + tester.fontService.Desc()
	ret += "\n  " + tester.programService.Desc()
	ret += "\n\n"

	return ret
}

func (tester *Tester) ProcessQueries(queryChan chan (chan string)) error {

	log.Debug("tester start process info queries")

	for {

		chn := <-queryChan
		info := tester.Info()

		select {
		case chn <- info:
			continue

		case <-time.After(1000. * time.Millisecond):
			continue
		}

	}

}

func (tester *Tester) printDebug() {

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Debug("")
	}

	if DEBUG_MEMORY {
		log.Debug("memory usage %s", MemUsage())
	}

	if DEBUG_DIAG {
		log.Debug("%s    %s", gfx.WorldClock().Info(tester.prevFrame), InfoDiag())
	}

	if DEBUG_CLOCK {
		log.Info("%s", gfx.WorldClock().Info(tester.prevFrame))
	}

	if DEBUG_MODE {
		log.Debug("  %s", tester.InfoMode())
		log.Debug("  %s", tester.lineBuffer.Desc())
		log.Debug("  %s", tester.termBuffer.Desc())
		log.Debug("  %s", tester.wordBuffer.Desc())
	}

	if DEBUG_FONT {
		log.Info("  %s", tester.fontService.Desc())
		if tester.font != nil {
			log.Debug("  %s", tester.font.Desc())
		}
	}

	if DEBUG_BUFFER && log.DebugLogging() {
		tester.dumpBuffer()
	}

	if DEBUG_MEMORY || DEBUG_DIAG || DEBUG_CLOCK || DEBUG_MODE || DEBUG_FONT {
		log.Debug("")
	}

}

func (tester *Tester) dumpBuffer() {
	if !DEBUG_BUFFER {
		return
	}
	if tester.mode == facade.Mode_TERM {
		os.Stdout.Write([]byte(tester.termBuffer.Dump()))
	} else if tester.mode == facade.Mode_LINES {
		os.Stdout.Write([]byte(tester.lineBuffer.Dump(uint(tester.gridConfig.GetWidth()))))
	} else if tester.mode == facade.Mode_TAGS {
		os.Stdout.Write([]byte(tester.wordBuffer.Dump()))
	} else if tester.mode == facade.Mode_WORDS {
		os.Stdout.Write([]byte(tester.wordBuffer.Dump()))
	}
	os.Stdout.Write([]byte("\n"))
	os.Stdout.Sync()
}

//func SaveRGBA(img *image.RGBA,outPath string)  {
//
//    if img == nil {
//        log.Error("no image to save at "+outPath)
//        return
//    }
//
//    outFile, err := os.Create(outPath)
//    if err != nil {
//        log.PANIC("fail to create file %s: %s",outPath,err)
//    }
//    defer outFile.Close()
//
//    writer := bufio.NewWriter(outFile)
//    if err := png.Encode(writer, img); err != nil {
//        log.PANIC("fail to encode image to %s: %s",err,err)
//    }
//
//    writer.Flush()
//    log.Info("wrote image to %s",outPath)
//
//}

//func (tester *Tester) testCharMap() (*image.RGBA,error) {
//
//    ret, err := tester.font.RenderMapRGBA()
//    if err != nil {
//        log.Error("fail to render glyphmap for %s: %s",tester.font.Desc(),err)
//        return nil,err
//    }
//    return ret, nil
//}
//
//
//func (tester *Tester) testTextTex(str string) (*image.RGBA,error) {
//
//
//    ret,err := tester.font.RenderTextRGBA(str)
//
//    if err != nil {
//        log.Error("fail to generate texture for '%s': %s",str,err)
//        return nil,err
//    }
//    return ret, nil
//}

//func (tester *Tester) Test(str string, confChan chan facade.Config, textChan chan facade.RawText) error {
//    var err error
//    switch tester.mode {
//        case facade.GRID:
//            test,_ := tester.testCharMap()
//            SaveRGBA(test,fmt.Sprintf("%s/test/map-%s.png",tester.directory,tester.name))
//
//        case facade.LINES:
//            test,_ := tester.testTextTex(str)
//            SaveRGBA(test,fmt.Sprintf("%s/test/text-%s-%s.png",tester.directory,tester.name,str))
//
//        default:
//            err = tester.testAnsi(confChan,textChan)
//
//    }
//    return err
//}

//func SaveRGBA(img *image.RGBA,outPath string)  {
//
//    if img == nil {
//        log.Error("no image to save at "+outPath)
//        return
//    }
//
//    outFile, err := os.Create(outPath)
//    if err != nil {
//        log.PANIC("fail to create file %s: %s",outPath,err)
//    }
//    defer outFile.Close()
//
//    writer := bufio.NewWriter(outFile)
//    if err := png.Encode(writer, img); err != nil {
//        log.PANIC("fail to encode image to %s: %s",err,err)
//    }
//
//    writer.Flush()
//    log.Info("wrote image to %s",outPath)
//
//}

//func (tester *Tester) testAnsi(rawConfs chan facade.Config, rawTexts chan facade.RawText) error {
//    term := gfx.NewTermBuffer(20,8)
//    for {
//        select {
//            case txt := <- rawTexts:
//    //                        log.Debug("recv %d byte text",len(text))
//
//    //                		os.Stdout.Write([]byte(text))
//                text := gfx.NewText( string(txt) )
//                term.WriteText( text )
//    //                        os.Stdout.Write( []byte("\n") )
//    //                        os.Stdout.Write( []byte(ansi.Dump()) )
//
//
//            case conf := <- rawConfs:
//                log.Debug("recv conf %s",conf.Desc())
//                if grid,ok := conf.Grid(); ok {
//                    var w,h uint = 0,0
//                    w,_ = grid.Width()
//                    h,_ = grid.Height()
//                    if w!=0 && h!= 0 {
//                        term.Resize(w,h)
//                    }
//                }
//
//            case <- time.After( 1 * time.Second ):
//                log.Debug(term.Desc() )
//
//            default:
//                //nop
//        }
//
//    //            for {
//    //                time.Sleep( time.Duration( int64(time.Second)) )
//    //            }
//    }
//    return nil
//}
