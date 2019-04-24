
package facade

import(
    "../gfx"
    log "../log"    
)

type Test struct {
    ringBuffer *gfx.RingBuffer
    termBuffer   *gfx.TermBuffer
    
    state TestState
    
}


func NewTest(config *TestConfig, ringBuffer *gfx.RingBuffer, termBuffer *gfx.TermBuffer) *Test {

    ret := &Test{}
    ret.state = TestDefaults
    ret.ringBuffer = ringBuffer
    ret.termBuffer = termBuffer
    ret.state.ApplyConfig(config)
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

