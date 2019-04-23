

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
)

type Tester struct {
    state facade.State

    font *gfx.Font; 
    
    ringBuffer *gfx.RingBuffer
    termBuffer *gfx.TermBuffer

    mutex *sync.Mutex
    directory string
}

func NewTester(directory string) *Tester { 
    ret := &Tester{directory: directory} 
    ret.mutex = &sync.Mutex{}
    ret.ringBuffer = gfx.NewRingBuffer(10) //FIXME
    ret.termBuffer = gfx.NewTermBuffer(10,10) //FIXME
    return ret    
}



func (tester *Tester) Init(config *facade.Config) error {
    log.Debug("init tester[%s] %s",tester.directory,config.Desc())
    if strings.HasPrefix(tester.directory, "~/") {
        tester.directory = os.Getenv("HOME") + tester.directory[1:]
    }
    gfx.SetFontDirectory(tester.directory+"/font")
    return nil   
}



func (tester *Tester) Configure(config *facade.Config) error {
    
    if config == nil { log.Error("tester config nil") ;return nil }
    if len(*config) <= 0 { return nil }
    
    log.Debug("tester config %s",config.Desc())

    if tmp,ok := config.Font(); ok {
		newFont, err := gfx.GetFont(&tmp)
		if err != nil {
			log.Error("fail to get font %s",tmp.Desc())
		} else {
			newFont.Init()
			tester.font = newFont
		}
	}
    
    return nil

}







func (tester *Tester) ProcessText(textChan chan string) {
    select {
        case txt := <-textChan:
            text := gfx.NewText( txt )
            tester.ringBuffer.WriteText( text )
            tester.termBuffer.WriteText( text )
        	
        default:
            //nop    
    }
}



func (tester *Tester) ProcessConf(confChan chan facade.Config) {
    select {
        case conf := <-confChan:
            tester.Configure(&conf)
        
        default:
            //nop    
    }
}


func (tester *Tester) ProcessRawTexts(rawChan chan facade.RawText, textChan chan string) error {

    for {
        rawText := <-rawChan
        if DEBUG_MESSAGES { log.Debug("process raw text: %s",rawText) }
        text := rawText.Sanitize()
        
//        tester.mutex.Lock()
//        tester.buffer.Queue( gfx.NewText(text) )
//        tester.mutex.Unlock()
        
        textChan <- text
        
    }
    return nil    
    
}





func (tester *Tester) ProcessRawConfs(rawChan chan facade.Config, confChan chan facade.Config) error {
    for {
        rawConf := <-rawChan
        if DEBUG_MESSAGES { log.Debug("process raw: %s",rawConf.Desc()) }
        conf := rawConf.Sanitize()

//        tester.mutex.Lock()
//        // prep some stuff i guess?
//        tester.mutex.Unlock()
        
        confChan <- conf

    }
    return nil
}














func (tester *Tester) Test(confChan chan facade.Config, textChan chan string) error {
    const FRAME_RATE = 60.
    gfx.ClockTick()
    var prev gfx.Clock = *gfx.NewClock()

    log.Debug("test %s",tester.state.Desc())
    for {
        verbose := gfx.ClockDebug()
        tester.mutex.Lock()

        tester.ProcessText(textChan)
        tester.ProcessConf(confChan)

        if verbose { 
            tester.PrintDebug(prev); 
            prev = *gfx.NewClock() 
        }

        tester.mutex.Unlock()
        
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
        gfx.ClockTick()
        
        
    }
    return nil
}



func (tester *Tester) PrintDebug(prev gfx.Clock) {

    if DEBUG_CLOCK { log.Debug("%s    %4.1ffps",gfx.ClockDesc(),gfx.ClockDelta(prev)) }


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




