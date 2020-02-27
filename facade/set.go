// +build linux,arm

package facade

import (
    "fmt"
    "strings"
    gfx "../gfx"
    log "../log"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"
)    


const DEBUG_SET = false



type Set struct {

    vert, frag string
    
    textures map[string] *gfx.Texture
    words []Word
        
    wordBuffer *WordBuffer

	program *gfx.Program
 	object  *gfx.Object
	data []float32

	refreshChan chan bool
    
}


const (
    WORDMAX   gfx.UniformName = "wordMax"
    WORDINDEX gfx.UniformName = "wordIndex"
	WORDCOUNT gfx.UniformName = "wordCount"
	WORDWIDTH gfx.UniformName = "wordWidth"
    WORDTIMER gfx.UniformName = "wordTimer"
    WORDFADER gfx.UniformName = "wordFader"
)

const (
)


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


func NewSet(buffer *WordBuffer) *Set {
    ret := &Set{
        wordBuffer: buffer,
    }
    
    ret.vert = ShaderDefaults.GetVert()
    ret.frag = ShaderDefaults.GetFrag()
    
    ret.words = []Word{}
    ret.textures = make( map[string] *gfx.Texture, ret.wordBuffer.SlotCount() )

    
    ret.refreshChan = make(chan bool, 1)
    return ret
}

func (set *Set) generateData(font *gfx.Font) {
    old := set.textures
    
    set.textures = make( map[string] *gfx.Texture, set.wordBuffer.SlotCount())
    set.words = set.wordBuffer.Words()
        
        
    //generate textures
    for _,word := range set.words {
        
        text := word.text
        if len(text) <= 0 {
            continue
        }

        if old[text] != nil {   //reuse existing textures

            set.textures[text] = old[text]
//            if DEBUG_SET {
//                log.Debug("%s texture reused: %s",set.Desc(),old[text].Desc())
//            }

        } else {               //create new texture
            
            rgba, err := font.RenderText(text, false)
            if err != nil {
                log.Error("%s texture fail render '%s': %s", set.Desc(), text, err)
                continue
            } 

            texture := gfx.NewTexture(text)
            texture.Init()
            
            err = texture.LoadRGBA(rgba)
            if err != nil {
                log.Error("%s texture fail load rgba '%s': %s", set.Desc(), text, err)
                texture.Close()
                continue
            }
            
            err = texture.TexImage()
            if err != nil {
                log.Error("%s texture fail teximage '%s': %s", set.Desc(), text, err)
                texture.Close()
                continue
            }

            set.textures[text] = texture
            
            if DEBUG_SET {
                log.Debug("%s texture prepped: %s",set.Desc(),set.textures[text].Desc())
            }
            
        }
    }
    
    // remove unused textures
    for text,texture := range old {
        _,ok := set.textures[text]
        if !ok {
            if DEBUG_SET {
                log.Debug("%s texture close: %s",set.Desc(),texture.Desc())
            }
            texture.Close()
        }
    }
    
    

    //setup vertex + bind order arrays

    set.data = []float32{}

    for _,word := range set.words  {

        text := word.text
        if len(text) <= 0 {
            continue
        }
        
        texture,ok := set.textures[text]
        if !ok {
            log.Debug("%s texture generate miss: %s",set.Desc(),text)
            continue
        }

        w := float32( texture.Size.Width / texture.Size.Height )
        h := float32( 1. )

/* 

     A          D
 -w/2,h/2____w/2,h/2
     |          |
     |          |
 -w/2,-h/2___w/2,-h/2
     B          C


     A          D
    0,0________1,0
     |          |
     |          |
    0,1________1,1
     B          C


  A
  |\
  |_\
  B C
 
  A_D
  \ |
   \|
    C

*/

        data := []float32{
        //   x,   y,    z,             tx, ty,   
            -w/2.,  +h/2.,  0.0,       0., 0.,    // A
            -w/2.,  -h/2.,  0.0,       0., 1.,    // B
            +w/2.,  -h/2.,  0.0,       1., 1.,    // C

            +w/2.,  -h/2.,  0.0,       1., 1.,    // C
            +w/2.,  +h/2.,  0.0,       1., 0.,    // D
            -w/2.,  +h/2.,  0.0,       0., 0.,    // A
        }
        set.data = append(set.data, data...)
        if DEBUG_SET {
            log.Debug("%s data generate '%s' %s",set.Desc(),text,texture.Desc())
        }
        
    }
    
    set.object.BufferData(len(set.data) * 4, set.data)
//    if DEBUG_SET {
//        log.Debug("%s generated %d words (%d floats)",set.Desc(),cnt,len(set.data))
//    }
    
    
}


func (set *Set) autoScale(camera *gfx.Camera) float32 {

    scaleHeight := float32(1.) / float32(set.wordBuffer.SlotCount())
    return scaleHeight * 2.

}


func (set *Set) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	if set.checkRefresh() {
//		if DEBUG_SET {
//			log.Debug("%s refresh", set.Desc())
//		}
		set.generateData(font)
	}

    gl.ActiveTexture(gl.TEXTURE0)
    
	set.program.UseProgram(debug)
	set.object.BindBuffer()

    tagMax := float32( set.wordBuffer.SlotCount() )
    set.program.Uniform1fv(WORDMAX, 1, &tagMax)

    set.program.Uniform1f(gfx.SCREENRATIO, camera.Ratio())
    
	clocknow := float32(gfx.Now())
	set.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow)
	
	camera.Uniform(set.program)
	
	scale := float32(1.0)
	scale = set.autoScale(camera)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
	//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
	set.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])
	

	set.program.VertexAttribPointer(gfx.VERTEX,   3, (3+2)*4, (0)*4)
	set.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2)*4, (0+3)*4)

    count := int32(1)
	offset := int32(0)
	
	for _,word := range set.words {

        text := word.text
        
        if len(text) <= 0 {
            continue
        }
        
        texture,ok := set.textures[text]
        if !ok {
            log.Debug("%s texture render miss: %s",set.Desc(),text)
            continue
        }


        texture.BindTexture()
        texture.Uniform(set.program)

        var timer float32
        timer = word.timer.Edge(gfx.Now())
    	set.program.Uniform1fv(WORDTIMER, 1, &timer)

        var fader float32
        fader = word.timer.Fader()
    	set.program.Uniform1fv(WORDFADER, 1, &fader)
    	
    	var index float32;
    	index = float32(word.index)
    	set.program.Uniform1fv(WORDINDEX, 1, &index)

        tagCount := float32( word.count )
        set.program.Uniform1fv(WORDCOUNT, 1, &tagCount)


        var width float32;
        width = float32(texture.Size.Width / texture.Size.Height)        
        set.program.Uniform1fv(WORDWIDTH, 1, &width)

        if DEBUG_SET && verbose {
            log.Debug("%s render #%d f%.1f  %s",set.Desc(),word.index,fader,texture.Desc())
        }
        
        if !debug || debug {
            set.program.SetDebug(false)
            gl.DrawArrays(gl.TRIANGLES, int32(offset*(2*3)), int32(count*2*3))
            set.program.SetDebug(debug)
        }
        
        if debug {
            gl.LineWidth(3.0)
            gl.BindTexture(gl.TEXTURE_2D, 0)
            gl.DrawArrays(gl.LINE_STRIP, int32(offset*(2*3)), int32(count*2*3))
        }
        offset += 1
    }
    


}




func (set *Set) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", set.Desc())

	set.object = gfx.NewObject("set")
	set.object.Init()

	set.program = programService.GetProgram("set", "set/")
	set.program.Link(set.vert, set.frag)

	set.ScheduleRefresh()
        
}


func (set *Set) Configure(words *WordConfig, tags *TagConfig, camera *gfx.Camera, font *gfx.Font) {
    var shader *ShaderConfig = nil
	var config *SetConfig = nil

	if tags != nil {
        log.Debug("%s configure %s", set.Desc(), tags.Desc())
        shader = tags.GetShader()
        config = tags.GetSet()
    } else if words != nil {
        log.Debug("%s configure %s", set.Desc(), words.Desc())
        shader = words.GetShader()
        config = words.GetSet()
	} else {
		log.Debug("%s cannot configure", set.Desc())
		return
	}

    
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
			//			err := set.LoadShaders()
			if err != nil {
				set.vert = vert
				set.frag = frag
			}
		}
	}
	
    if config.GetSetDuration() {
        set.wordBuffer.SetDuration( float32(config.GetDuration()) )
    }
	
	if config.GetSetSlot() {
    	set.wordBuffer.Resize( int(config.GetSlot()) )
    }
    
    if config.GetShuffle() {
        set.wordBuffer.SetShuffle( config.GetShuffle() )
    }

	if config.GetSetFill() {
		fillStr := set.fill( config.GetFill() ) 
		if set.wordBuffer != nil {
    		set.wordBuffer.Fill(fillStr)
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
    ret := "set"
    ret += "["
    ret += fmt.Sprintf("%d/%d",set.wordBuffer.WordCount(),set.wordBuffer.SlotCount())
    ret += fmt.Sprintf(" %.1f",set.wordBuffer.Duration())
    ret += "]"
    return ret 
    
}

func (set *Set) Config() *SetConfig {
    ret := &SetConfig{
        SetDuration: true, Duration: float64(set.wordBuffer.Duration()),
        SetSlot: true, Slot: uint64(set.wordBuffer.SlotCount()),
        SetShuffle: true, Shuffle: bool(set.wordBuffer.Shuffle()),
    }
    return ret
    
}

func (set *Set) ShaderConfig() *ShaderConfig {
    ret := &ShaderConfig{
        SetVert: true, Vert: set.vert,
        SetFrag: true, Frag: set.frag,
    }
    return ret
}



