

package main

import (
//    "fmt"
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

    Width, Height uint
    Speed float32
    Terminal bool
    Buffer uint
    
    Mode facade.Mode

    font *gfx.Font; 
    
    
    lineBuffer *facade.LineBuffer
    termBuffer *facade.TermBuffer

    mutex *sync.Mutex
    directory string
    
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
    gfx.SetFontDirectory(tester.directory+"/font")
    
    
    
    
//    var err error
    
    //setup things 
//	tester.state = facade.Defaults
//    tester.state.ApplyConfig(config)
    
//	fontConfig := gfx.FontDefaults.Config()
//	if cfg,ok := config.Font(); ok {
//		fontConfig.ApplyConfig( &cfg )	
//	}
//	tester.font,err = gfx.GetFont( fontConfig )
//    if err != nil {
//        log.PANIC("no default font: %s",err)    
//    }
//	tester.font.Init()
//	
//	


//    grid := proto.GridConfig{}
//	gridConfig := facade.GridDefaults.Config()
//    if cfg,ok := config.Grid(); ok {
//        gridConfig.ApplyConfig(&cfg)
//    }
//    

    tester.Width  = uint(facade.GridDefaults.Width) 
    tester.Height = uint(facade.GridDefaults.Height)
    tester.Buffer = uint(facade.GridDefaults.Buffer)
    tester.Terminal = facade.GridDefaults.Terminal

    
    if grid := config.GetGrid(); grid!=nil {
        if grid.GetCheckWidth() { tester.Width = uint(grid.GetWidth()) }
        if grid.GetCheckHeight() { tester.Height = uint(grid.GetHeight()) }
        if grid.GetCheckBuffer() { tester.Buffer = uint(grid.GetBuffer()) }
        if grid.GetCheckTerminal() { tester.Terminal = grid.GetTerminal() }
    }

    tester.termBuffer = facade.NewTermBuffer(tester.Width,tester.Height) 
    tester.lineBuffer = facade.NewLineBuffer(tester.Height,tester.Buffer,tester.refreshChan) 
    
    

    
    gfx.ClockReset()
    return nil   
}



func (tester *Tester) Desc() string {

    tmp := facade.GridConfig{
        CheckWidth: true,  Width: uint64(tester.Width),
        CheckHeight: true, Height: uint64(tester.Height),
        CheckBuffer: true, Buffer: uint64(tester.Buffer),
        CheckTerminal: true, Terminal: tester.Terminal,
    }
    
    return "tester[" + tmp.Desc() + "]"
}

func (tester *Tester) Configure(config *facade.Config) error {
    
    if config == nil { return nil }
    log.Debug("%s configure %s",tester.Desc(),config.Desc())


    if grid := config.GetGrid(); grid != nil {
        
        resize := false

        if grid.GetCheckWidth() { resize = true;  tester.Width = uint(grid.GetWidth()) } 
        if grid.GetCheckHeight() { resize = true; tester.Height = uint(grid.GetHeight()) } 
        if grid.GetCheckBuffer() { resize = true; tester.Buffer = uint(grid.GetBuffer()) } 
        if grid.GetCheckTerminal() { tester.Terminal = grid.GetTerminal() } 

        if resize {
            tester.termBuffer.Resize(tester.Width,tester.Height) 
            tester.lineBuffer.Resize(tester.Height,tester.Buffer) 
        }

	}
    
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



func (tester *Tester) ProcessBufferItems(bufChan chan facade.BufferItem) error {

    for {
        item := <- bufChan    
//        if DEBUG_MESSAGES { log.Debug("buffer %s",item.Desc()) }
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
        if DEBUG_MESSAGES { log.Debug("process raw conf: %s",rawConf.Desc()) }

//        tester.mutex.Lock()
//        // prep some stuff i guess?
//        tester.mutex.Unlock()
        
        confChan <- rawConf

    }
    return nil
}
















func (tester *Tester) Test(confChan chan facade.Config) error {
    const FRAME_RATE = 60.
    gfx.ClockTick()
    var prev gfx.Clock = *gfx.NewClock()

//    log.Debug("test %s",tester.state.Desc())
    for {
        tester.mutex.Lock()

        tester.ProcessConf(confChan)

        if gfx.ClockVerboseFrame() {

            if DEBUG_BUFFER && tester.Mode == facade.Mode_GRID {
                if tester.Terminal {
                    os.Stdout.Write( []byte( tester.termBuffer.Dump() ) )
                } else {
                    os.Stdout.Write( []byte( tester.lineBuffer.Dump( tester.Width ) ) ) 
                }
                os.Stdout.Write( []byte( "\n" ) )
                os.Stdout.Sync()
            }
            

            if DEBUG_CLOCK { 
                log.Debug("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(prev)) 
                prev = *gfx.NewClock() 
            }
            
            if DEBUG_MODE && tester.Mode == facade.Mode_GRID {
                log.Debug("mode[grid] %s",tester.Desc() )
            }
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




