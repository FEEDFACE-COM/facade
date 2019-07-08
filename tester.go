

package main

import (
    "strings"
    "os"
    "sync"
    "bufio"
    "image/png"
    "image"
    "time"
    log "./log"
    gfx "./gfx"
    facade "./facade"
//    proto "./facade/proto"
)

type Tester struct {

    Terminal bool
    
    Mode facade.Mode
    debug bool

    font *gfx.Font; 
    vert,frag *gfx.Shader
    
    
    lineBuffer *facade.LineBuffer
    termBuffer *facade.TermBuffer
    
    fontService *gfx.FontService
    programService *gfx.ProgramService

    mutex *sync.Mutex
    directory string
    
    prevFrame gfx.ClockFrame
    
    image *image.RGBA
    
    
    refreshChan chan bool
        
    
}



func NewTester(directory string) *Tester { 
    ret := &Tester{directory: directory}
    ret.mutex = &sync.Mutex{}
    ret.refreshChan = make( chan bool, 1 )
    return ret    
}


func (tester *Tester) switchShader(shaderName string, shaderType gfx.ShaderType) error {
    var err error
    shaderName = strings.ToLower(tester.Mode.String()) + "/" + strings.ToLower(shaderName)
    name := shaderName + "." + string(shaderType)
    
    
    log.Debug("tester load shader %s",name)
    err = tester.programService.LoadShader(shaderName,shaderType)
    if err != nil {
        log.Error("tester fail load shader %s: %s",name,err)
        return log.NewError("tester fail load shader %s: %s",name,err)
    }

    var shader *gfx.Shader
    shader,err = tester.programService.GetShader( name, shaderType )
    if err != nil {
        log.PANIC("tester fail get shader %s: %s",name,err)
        return log.NewError("tester fail get shader %s: %s",name,err)
    }

    log.Debug("tester switch to shader %s",name)
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

        log.Debug("tester load font %s",name )
        err = tester.fontService.LoadFont( name )
        if err != nil {
            log.Debug("tester fail load font %s: %s",name,err)
            return log.NewError("tester fail load font %s: %s",name,err)
        }
        
        var font *gfx.Font
        font,err = tester.fontService.GetFont( name )
        if err != nil {
            log.Debug("tester fail get font %s: %s",name,err)
            return log.NewError("tester fail get font %s: %s",name,err)
        }


        log.Debug("tester switch to font %s",name)
        tester.font = font

    }        
    
    return nil
}



func (tester *Tester) Init() error {
    var err error    

    log.Debug("init tester[%s]",tester.directory,)

    if strings.HasPrefix(tester.directory, "~/") {
        tester.directory = os.Getenv("HOME") + tester.directory[1:]
    }

    tester.fontService = gfx.NewFontService(tester.directory+"/font")
    tester.programService = gfx.NewProgramService(tester.directory+"/shader")

    tester.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width),uint(facade.GridDefaults.Height)) 
    tester.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height),uint(facade.GridDefaults.Buffer),tester.refreshChan) 


    err = tester.switchFont(facade.FontDefaults.Name)
    if err != nil {
        log.PANIC("tester missing default font: %s",err)     
    }

    err = tester.switchShader(facade.GridDefaults.Vert, gfx.VertType)
    if err != nil {
        log.PANIC("tester missing default vert shader: %s",err)    
    }

    err = tester.switchShader(facade.GridDefaults.Frag, gfx.FragType)
    if err != nil {
        log.PANIC("tester missing default frag shader: %s",err)    
    }
    
    
    gfx.WorldClock().Reset()
    return nil   
}




func (tester *Tester) Configure(config *facade.Config) error {
    var err error
    if config == nil { return nil }
    log.Debug("tester config %s",config.Desc())



	if config.GetSetDebug() {
		tester.debug = config.GetDebug()
	} else {
		tester.debug = false	
	}


    if cfg:=config.GetFont(); cfg!=nil {

        if cfg.GetSetName() && cfg.GetName() != tester.font.GetName() {
            err = tester.switchFont( cfg.GetName() )
            if err != nil {
                log.Error("tester fail switch font: %s",err)     
            }
        }
    }


    if grid := config.GetGrid(); grid != nil {
            
        if grid.GetSetVert() && grid.GetVert() != facade.GridDefaults.Vert {
            err = tester.switchShader( grid.GetVert(), gfx.VertType )
            if err != nil {
                log.Error("tester fail switch shader: %s",err)     
            }
        }
        
        if grid.GetSetFrag() && grid.GetFrag() != facade.GridDefaults.Frag {
            err = tester.switchShader( grid.GetFrag(), gfx.FragType )
            if err != nil {
                log.Error("tester fail switch shader: %s",err)     
            }
        }
        

        if grid.GetSetTerminal() { tester.Terminal = grid.GetTerminal() } 

        if ( grid.GetSetWidth() && grid.GetWidth() != tester.termBuffer.GetWidth())    ||
           ( grid.GetSetHeight() && grid.GetHeight() != tester.termBuffer.GetHeight() ) {
        
            tester.termBuffer.Resize(uint(grid.GetWidth()),uint(grid.GetHeight()))
        }

        if ( grid.GetSetBuffer() && grid.GetBuffer() != tester.lineBuffer.GetBuffer() )   ||
           ( grid.GetSetHeight() && grid.GetHeight() != tester.lineBuffer.GetHeight() ) { 

            tester.lineBuffer.Resize(uint(grid.GetHeight()),uint(grid.GetBuffer()))
        }
        
	    if grid.GetSetSpeed()    { tester.lineBuffer.SetSpeed( float32(grid.GetSpeed() ) ) }
    	if grid.GetSetAdaptive() { tester.lineBuffer.Adaptive = grid.GetAdaptive() }
        if grid.GetSetDrop()     { tester.lineBuffer.Drop = grid.GetDrop() }
        if grid.GetSetSmooth()   { tester.lineBuffer.Smooth = grid.GetSmooth() }
        
        if grid.GetSetFill() {
            
            if err:=tester.render( grid.GetFill() );  err!=nil {
                log.Error("fail render '%': %s",grid.GetFill(),err)
            }

        }
        
	}

    return nil

}


func (tester *Tester) render(fill string) error {
    if tester.font != nil {
        log.Debug("tester render '%s' with %s",fill,tester.font.Desc())
    }
    var err error

    if fill == "" { //render out glyphmap
        
        tester.image,err = tester.font.RenderMap(tester.debug)
            
    } else { // render out given string
        
        tester.image,err = tester.font.RenderText( fill, tester.debug )
    }
    
    if tester.image == nil {
        return log.NewError("fail render '%s' with %s: %s",fill,tester.font.Desc(),err)
    }
    log.Debug("tester rendered '%s' with %s",fill,tester.font.Desc())
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
        item := <- bufChan    
        text, seq := item.Text, item.Seq
        if text != nil && len(text) > 0 {
            tester.lineBuffer.ProcessRunes( text )
            tester.termBuffer.ProcessRunes( text )    
        }
        if seq != nil {
            tester.lineBuffer.ProcessSequence( seq )
            tester.termBuffer.ProcessSequence( seq )
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
    ret := strings.ToLower(tester.Mode.String())    
    ret += "["
    if tester.Terminal { 
        ret += "TT "
    }
    if tester.debug {
        ret += "DEBUG "	
    }
    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
    
}





func (tester *Tester) Test(confChan chan facade.Config) error {
    const FRAME_RATE = 60.
    gfx.WorldClock().Tick()
    tester.prevFrame = gfx.WorldClock().Frame()



    for {
        tester.mutex.Lock()

        tester.ProcessConf(confChan)

        if gfx.WorldClock().VerboseFrame() {

            log.Debug("")
            log.Info( "%s", gfx.WorldClock().Info(tester.prevFrame) )
//            log.Info("  %s", MemUsage() )
            log.Info("  %s", tester.InfoMode() ) 
            log.Info("  %s", tester.lineBuffer.Desc() )
            log.Info("  %s", tester.termBuffer.Desc() )
            log.Info("  %s",tester.fontService.Desc())
            if tester.font != nil { log.Info("  %s",tester.font.Desc()) }
            log.Info("  %s",tester.programService.Desc())
            if tester.vert != nil { log.Info("  %s",tester.vert.Desc()) }
            if tester.frag != nil { log.Info("  %s",tester.frag.Desc()) }
            if DEBUG_BUFFER && tester.Mode == facade.Mode_GRID {
                log.Info("")
                if tester.Terminal { 
                    log.Info(tester.termBuffer.Dump() ) 
                } else { 
                    log.Info(tester.lineBuffer.Dump( uint(tester.termBuffer.GetWidth())) ) 
                }
            }
            log.Debug("")
        }
            

//            if DEBUG_BUFFER && tester.Mode == facade.Mode_GRID {
//                if tester.Terminal {
//                    os.Stdout.Write( []byte( tester.termBuffer.Dump() ) )
//                } else {
//                    os.Stdout.Write( []byte( tester.lineBuffer.Dump( tester.GetWidth ) ) ) 
//                }
//                os.Stdout.Write( []byte( "\n" ) )
//                os.Stdout.Sync()
//            }
//            
//
//
//
//        }

        tester.mutex.Unlock()


        if tester.image != nil {
    
            outPath := "./font.png"
            log.Info("write rendered image to %s",outPath)
            outFile, err := os.Create(outPath)
            if err != nil {
                log.Error("fail to create file %s: %s",outPath,err)
            }
            writer := bufio.NewWriter(outFile)
            if err := png.Encode(writer, tester.image); err != nil {
                log.Error("fail to encode rendered image: %s",err)
                return log.NewError("fail to encode rendered image: %s",err)
            }
            writer.Flush()
            outFile.Close()
            tester.image = nil
            
        }

        
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        gfx.WorldClock().Tick()


        
    }
    

    
    
    return nil
}

func (tester *Tester) Info() string { 
    ret := ""
    
    ret += InfoVersion()
    ret += InfoAssets( tester.programService.GetAvailableNames(), tester.fontService.GetAvailableNames() )
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



func (tester *Tester) ProcessQueries(queryChan chan (chan string) ) error {

    log.Debug("tester start process info queries")

    for {
    
        chn := <- queryChan
        info := tester.Info()
        
        select {
            case chn <- info:
                continue
            
            case <-time.After(1000. * time.Millisecond):
                continue
        }
        
        
    }
    
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




