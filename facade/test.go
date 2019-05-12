
package facade

import(
    "../gfx"
    log "../log"    
)

type Test struct {
    termBuffer   *gfx.TermBuffer
    lineBuffer *gfx.LineBuffer
    
    state TestState
    refreshChan chan bool
    
    
}

func (test *Test) Width() uint { return test.state.Width }
func (test *Test) Height() uint { return test.state.Height }
func (test *Test) Buffer() BufferName { return test.state.Buffer }
func (test *Test) BufLen() uint { return test.state.BufLen }

func NewTest(config *TestConfig, termBuffer *gfx.TermBuffer, lineBuffer *gfx.LineBuffer) *Test {

    ret := &Test{}
    ret.state = TestDefaults
    ret.state.ApplyConfig(config)
    ret.refreshChan = make( chan bool, 1 )
    ret.termBuffer = termBuffer
    ret.lineBuffer = lineBuffer
    ret.termBuffer.Resize(ret.state.Width,ret.state.Height)   
    ret.lineBuffer.Resize(ret.state.Height)
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

    if height,ok := config.Height(); ok && height != 0 { 
	    test.state.Height = height 
        test.lineBuffer.Resize(test.state.Height + test.state.BufLen)
        test.termBuffer.Resize(test.state.Width,test.state.Height)   
    }
    
    if buflen,ok := config.BufLen(); ok && buflen != 0 {
        test.state.BufLen = buflen
        test.lineBuffer.Resize(test.state.Height + test.state.BufLen)
    }
    
    if buffer,ok := config.Buffer(); ok && buffer != test.state.Buffer {
        test.state.Buffer = buffer
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
