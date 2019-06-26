

package main

import (
    "fmt"
    "strings"
    "os"
    "sync"
//    "bufio"
//    "image/png"
//    "image"
    "time"
    log "./log"
    gfx "./gfx"
    facade "./facade"
//    proto "./facade/proto"
)

type Tester struct {

    Terminal bool
    
    Mode facade.Mode

    font *gfx.Font; 
    
    
    lineBuffer *facade.LineBuffer
    termBuffer *facade.TermBuffer
    
    fontService *gfx.FontService

    mutex *sync.Mutex
    directory string
    
    prevClock gfx.Clock
    
    refreshChan chan bool
        
    
}

func NewTester(directory string) *Tester { 
    ret := &Tester{directory: directory}
    ret.mutex = &sync.Mutex{}
    ret.refreshChan = make( chan bool, 1 )
    return ret    
}




func (tester *Tester) Init(config *facade.Config) error {
    log.Debug("init tester[%s] %s",tester.directory,config.Desc())

    if strings.HasPrefix(tester.directory, "~/") {
        tester.directory = os.Getenv("HOME") + tester.directory[1:]
    }

    tester.fontService = gfx.NewFontService(tester.directory+"/font")

    
    var err error    
    {
    	var name = facade.DEFAULT_FONT

        if cfg := config.GetFont(); cfg!=nil {
            if cfg.GetSetName() {
                name = cfg.GetName()
            }
        }
        
        log.Debug("load font %s",name)
        err = tester.fontService.LoadFont(name)
        if err != nil {
            log.PANIC("tester fail load default font %s: %s",name,err)
        }

        tester.font,err = tester.fontService.GetFont( name )
        if err != nil {
            log.PANIC("tester fail get default font %s: %s",name,err)
        }


    }

    
    
    
    if grid := config.GetGrid(); grid!=nil {
        
        if grid.GetSetTerminal() { tester.Terminal = grid.GetTerminal() }
        
    }

    tester.termBuffer = facade.NewTermBuffer(uint(facade.GridDefaults.Width),uint(facade.GridDefaults.Height)) 
    tester.lineBuffer = facade.NewLineBuffer(uint(facade.GridDefaults.Height),uint(facade.GridDefaults.Buffer),tester.refreshChan) 
    
    if grid := config.GetGrid(); grid!=nil {
	    if grid.GetSetSpeed() { 
            tester.lineBuffer.SetSpeed( float32(grid.GetSpeed() ) )
    	}
    	
    	if grid.GetSetAdaptive() {
            tester.lineBuffer.Adaptive = grid.GetAdaptive()
        }
        if grid.GetSetDrop() {
            tester.lineBuffer.Drop = grid.GetDrop()
        }

        if grid.GetSetSmooth() {
            tester.lineBuffer.Smooth = grid.GetSmooth()
        }
    }
    
    gfx.ClockReset()
    return nil   
}



func (tester *Tester) GridConfig() *facade.GridConfig {

    return &facade.GridConfig{
        SetWidth: true,  Width: uint64(tester.termBuffer.GetWidth()),
        SetHeight: true, Height: uint64(tester.termBuffer.GetHeight()),
        SetBuffer: true, Buffer: uint64(tester.lineBuffer.GetBuffer()),
        SetTerminal: true, Terminal: tester.Terminal,
    }
}

func (tester *Tester) Configure(config *facade.Config) error {
    var err error
    
    if config == nil { return nil }
    log.Debug("tester config %s",config.Desc())


    if cfg := config.GetFont(); cfg!=nil {
    
        if cfg.GetSetName() {
            name := cfg.GetName()
            if name != tester.font.GetName() {
        
                log.Debug("tester add font %s",name )
                err = tester.fontService.LoadFont( name )
                if err != nil {
                    log.Debug("tester fail add font %s: %s",name,err)
                }
                
                var font *gfx.Font
                font,err = tester.fontService.GetFont( name )
                if err == nil {
                    log.Debug("tester switch to font %s",name)
                    tester.font = font
                } else {
                    log.Debug("tester fail get font %s: %s",name,err)
                }
                
                
            }        
        
            
        }
        
    }


    if grid := config.GetGrid(); grid != nil {
        

        if grid.GetSetTerminal() { tester.Terminal = grid.GetTerminal() } 

        if ( grid.GetSetWidth() && grid.GetWidth() != tester.termBuffer.GetWidth())    ||
           ( grid.GetSetHeight() && grid.GetHeight() != tester.termBuffer.GetHeight() ) {
        
            tester.termBuffer.Resize(uint(grid.GetWidth()),uint(grid.GetHeight()))
        }

        if ( grid.GetSetBuffer() && grid.GetBuffer() != tester.lineBuffer.GetBuffer() )   ||
           ( grid.GetSetHeight() && grid.GetHeight() != tester.lineBuffer.GetHeight() ) { 

            tester.lineBuffer.Resize(uint(grid.GetHeight()),uint(grid.GetBuffer()))
        }
        
	    if grid.GetSetSpeed() { 
            tester.lineBuffer.SetSpeed( float32(grid.GetSpeed() ) )
    	}
    	
    	if grid.GetSetAdaptive() {
            tester.lineBuffer.Adaptive = grid.GetAdaptive()
        }
        if grid.GetSetDrop() {
            tester.lineBuffer.Drop = grid.GetDrop()
        }

        if grid.GetSetSmooth() {
            tester.lineBuffer.Smooth = grid.GetSmooth()
        }
        
        if grid.GetSetFill() {
            tester.render( grid.GetFill() )
        }
        
        

	}



    
    return nil

}


func (tester *Tester) render(fill string) {
    if tester.font != nil {
        log.Debug("tester render %s with %s",fill,tester.font.Desc())
    }
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
//    if tester.font != nil {
//        ret += " " + tester.font.Desc()
//    }
    return ret
    
}

func (tester *Tester) InfoClock() string {
    return fmt.Sprintf("%s    %4.1f",gfx.ClockDesc(),gfx.ClockDelta(tester.prevClock)  )
}





func (tester *Tester) Test(confChan chan facade.Config) error {
    const FRAME_RATE = 60.
    gfx.ClockTick()
    tester.prevClock = *gfx.NewClock()

    for {
        tester.mutex.Lock()

        tester.ProcessConf(confChan)

        if gfx.ClockVerboseFrame() {


            if DEBUG_CLOCK||DEBUG_MODE||DEBUG_BUFFER {
                log.Debug("")
            }
        
        
            if DEBUG_CLOCK { log.Info( "%s", tester.InfoClock() ) }
            
            if DEBUG_DIAG { log.Info("  %s", MemUsage() ) }
                
            if DEBUG_MODE { 
                    log.Info("  %s", tester.InfoMode() ) 
                    log.Info("  %s", tester.lineBuffer.Desc() )
                    log.Info("  %s", tester.termBuffer.Desc() )
            }
            
            if DEBUG_FONT {
                log.Info("  %s",tester.fontService.Desc())
                if tester.font != nil {
                    log.Info("  %s",tester.font.Desc())
                }
            }    
                        
            if DEBUG_BUFFER && tester.Mode == facade.Mode_GRID {
                log.Info("")
                if tester.Terminal { 
                    log.Info(tester.termBuffer.Dump() )
                } else {
                    log.Info(tester.lineBuffer.Dump( uint(tester.termBuffer.GetWidth())) )    
                }
                
            }
        
//            if DEBUG_BUFFER &&  log.DebugLogging() { renderer.dumpBuffer() }
         
            if DEBUG_CLOCK||DEBUG_MODE||DEBUG_BUFFER {
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

        }
        tester.mutex.Unlock()
        
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        gfx.ClockTick()
        
        
    }
    return nil
}



















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




