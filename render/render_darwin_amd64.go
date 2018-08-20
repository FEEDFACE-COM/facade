
package render

import (
    "fmt"
    "time"
    "sync"
    log "../log"
    conf "../conf"
    grid "../grid"
    font "../font"
)





const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

type Renderer struct {
    mode conf.Mode
    grid *grid.Grid
    font *font.Font
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    return ret
}

const DEBUG_CLOCK  = false
const DEBUG_CONF   = true
const DEBUG_TEXT   = true

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init() error {
    var err error
    log.Debug("initializing renderer")
    if err != nil {
        log.Fatal("could not initialize renderer: %s",err)    
    }
    
    config := conf.NewConfig(conf.DEFAULT)
    renderer.Configure(config)
    return err
}


func (renderer *Renderer) Configure(config *conf.Config) {
    renderer.mode = config.Mode
    switch (config.Mode) {
        case conf.GRID:
            renderer.grid = grid.NewGrid()
    }
    if renderer.font == nil {
        renderer.font = font.NewFont()    
    }
}



func (renderer *Renderer) Render() error {

    InitClock()

    var now  Clock = Clock{frame: 0}
    var prev Clock = Clock{frame: 0}
    


    log.Debug("renderer start")
    for {
        now.Tick()

        //draw
        // TBD
        renderer.mutex.Lock()
        


        switch renderer.mode {
            case conf.GRID:
                renderer.grid.Render()
                
        }


        if now.frame % DEBUG_FRAMES == 0 {
            if DEBUG_CLOCK   {
                fps := float64(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Describe(),fps)
                prev = now
            }

            if DEBUG_CONF {
                str := fmt.Sprintf("mode[%s]",string(renderer.mode))
                switch renderer.mode {
                    case conf.GRID:
                        str += " " + renderer.grid.Describe()
                    }
                str += " " + renderer.font.Describe()
                log.Debug(str)
            }
            
            if DEBUG_TEXT {
                str := ""
                switch renderer.mode {
                    case conf.GRID:
                        str = renderer.grid.Buffer.Debug(grid.PageDown)
                }    
                if str != "" {
                    log.Debug(str)
                }
            }
                
            
        }
        renderer.mutex.Unlock()
        
        // wait for next frame
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}


func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    for {
        text := <-textChan
//        log.Debug("read text: %s",text)
        renderer.mutex.Lock()
        renderer.grid.Queue( string(text) )
        renderer.mutex.Unlock()
    }
    return nil
    
}


func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    for {
        config := <-confChan
//        log.Debug("read config: %s",config.Describe())    
        renderer.mutex.Lock()
        renderer.Configure(&config)
        renderer.mutex.Unlock()
    }
    return nil
}

