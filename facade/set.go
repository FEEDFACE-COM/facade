// +build linux,arm

package facade

import (
//    "fmt"
    "strings"
    
    gfx "../gfx"
    log "../log"
    
//	gl "github.com/FEEDFACE-COM/piglet/gles2"
//	"github.com/go-gl/mathgl/mgl32"
)    


const DEBUG_SET = true

type Set struct {

    vert, frag string

    buffer *SetBuffer

	program *gfx.Program
    
	refreshChan chan bool
    
}



func (set *Set) ScheduleRefresh() {

	select {
	case set.refreshChan <- true:
	default:
	}

}

func (set *Set) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
		case refresh := <-set.refreshChan:
			if refresh {
				ret = true
			}

		default:
			return ret
		}
	}
	return ret
}


func NewSet(setBuffer *SetBuffer) *Set {
    ret := &Set{}

    ret.vert = ShaderDefaults.GetVert()
    ret.frag = ShaderDefaults.GetFrag()
    
    ret.refreshChan = make(chan bool, 1)
    ret.buffer = setBuffer
    return ret
}

func (set *Set) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	if set.checkRefresh() {
		if DEBUG_GRID {
			log.Debug("%s refresh", set.Desc())
		}
	}


	set.program.UseProgram(debug)


}


func (set *Set) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", set.Desc())


	set.program = programService.GetProgram("set", "set/")
	set.program.Link(set.vert, set.frag)

	set.ScheduleRefresh()
        
}


func (set *Set) Configure(config *TagConfig, camera *gfx.Camera, font *gfx.Font) {
    
    var shader *ShaderConfig = nil
    
    shader = config.GetShader()
    
	{
		changed := false
		vert, frag := set.vert, set.frag

        if shader != nil {
                            
            if shader.GetSetVert() {
                changed = true
                set.vert = shader.GetVert()
    		}
    		
            if shader.GetSetFrag() {
                changed = true
                set.frag = shader.GetFrag()
            }
        }
            
		if changed {
			err := set.program.Link(set.vert, set.frag)
			//			err := grid.LoadShaders()
			if err != nil {
				set.vert = vert
				set.frag = frag
			}
		}
	}

	if config.GetSetFill() {
		fillStr := set.fill( config.GetFill() ) 
		if set.buffer != nil {
    		set.buffer.Fill(fillStr)
		}
	}

	set.ScheduleRefresh()
    
    
    
}


func (set *Set) fill(name string) []string{
    switch name {
        case "alpha":
            return strings.Split(`
alpha
bravo
charlie
delta
echo
foxtrott
golf
hotel
india
juliet
kilo
lima
mike
november
oscar
papa
quebec
romeo
sierra
tango
uniform
victor
whiskey
xray
yankee
zulu
`           ,"\n")[1:]    
    }
    return []string{}
}

func (set *Set) Desc() string {
    return set.Config().Desc()
}

func (set *Set) Config() *TagConfig {
    return &TagConfig{}
}

func (set *Set) ShaderConfig() *ShaderConfig {
    ret := &ShaderConfig{
        SetVert: true, Vert: set.vert,
        SetFrag: true, Frag: set.frag,
    }
    return ret
}



