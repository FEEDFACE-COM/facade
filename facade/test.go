
package facade

import(
    "../gfx"
    log "../log"    
)

type Test struct {
    ringBuffer *gfx.RingBuffer
    termBuffer   *gfx.TermBuffer
    textBuffer *gfx.TextBuffer
    
    state TestState
    refreshChan chan bool
    
    
}

func (test *Test) Width() uint { return test.state.Width }

func NewTest(config *TestConfig, ringBuffer *gfx.RingBuffer, termBuffer *gfx.TermBuffer, textBuffer *gfx.TextBuffer) *Test {

    ret := &Test{}
    ret.state = TestDefaults
    ret.state.ApplyConfig(config)
    ret.refreshChan = make( chan bool, 1 )
    ret.ringBuffer = ringBuffer
    ret.termBuffer = termBuffer
    ret.textBuffer = textBuffer
    ret.ringBuffer.Resize(ret.state.Height) 
    ret.termBuffer.Resize(ret.state.Width,ret.state.Height)   
    ret.textBuffer.Resize(ret.state.Height)
    return ret

}


func (test *Test) Init(font *gfx.Font) {

    log.Debug("init %s",test.Desc())
    
}


func (test *Test) Configure(config *TestConfig, font *gfx.Font) {
    if config == nil { return }
    log.Debug("config %s",config.Desc())

    if width,ok := config.Width(); ok && width != 0 && width != test.state.Width { 
	    test.state.Width = width 
        test.termBuffer.Resize(test.state.Width,test.state.Height)   
	} 

    if height,ok := config.Height(); ok && height != 0 && height != test.state.Height { 
	    test.state.Height = height 
        test.ringBuffer.Resize(test.state.Height) 
        test.termBuffer.Resize(test.state.Width,test.state.Height)   
    }
    
}


func (test *Test) Desc() string { return test.state.Desc() }

func (test *Test) ScheduleRefresh() {

    select { case test.refreshChan <- true: ; default: ; }
	
}


func (test *Test) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
			case refresh := <- test.refreshChan:
				if refresh {
					ret = true
				}

			default:
				return ret
		}
	}
	return ret
}
