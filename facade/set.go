// +build linux,arm

package facade

import (
//    "fmt"
    "strings"
    "hash/crc32"
    
    gfx "../gfx"
    log "../log"
    
    
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"
)    


const DEBUG_SET = true



type TexItem struct {
    texture *gfx.Texture
    item *SetItem
}


type Set struct {

    vert, frag string


    
    
    texItem map[string] *TexItem
        
    

    buffer *SetBuffer

	program *gfx.Program
 	object  *gfx.Object


	data []float32
	tags []string

	refreshChan chan bool
    
}


const (
    TAGMAX      gfx.UniformName = "tagMax"
    TAGINDEX    gfx.UniformName = "tagIndex"
	TAGCOUNT    gfx.UniformName = "tagCount"
	TAGWIDTH    gfx.UniformName = "tagWidth"
    TAGFADER    gfx.UniformName = "tagFader"
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


func NewSet(setBuffer *SetBuffer) *Set {
    ret := &Set{}
    
    ret.vert = ShaderDefaults.GetVert()
    ret.frag = ShaderDefaults.GetFrag()
    
    ret.buffer = setBuffer
    ret.texItem = make( map[string] *TexItem, ret.buffer.Max())
    
    ret.refreshChan = make(chan bool, 1)
    return ret
}

func (set *Set) generateData(font *gfx.Font) {


    old := set.texItem
    
    set.texItem = make( map[string] *TexItem, set.buffer.Max())
    
    
    bufferItems := set.buffer.Items(set.buffer.Max())


    for _,item := range bufferItems {
        
        tag := item.tag

        if len(set.texItem) >= set.buffer.Max() {
            log.Error("%s stop render %d/%d reached", set.Desc(), len(set.texItem),set.buffer.Max())
            break
        }

        if old[tag] != nil {   //reuse existing textures

            set.texItem[tag] = old[tag]
            delete(old, tag)

        } else {               //create new texture
            
            
            rgba, err := font.RenderText(tag, false)
            if err != nil {
                log.Error("%s fail render '%s': %s", set.Desc(), tag, err)
                continue
            } 

            texture := gfx.NewTexture(tag)
            texture.Init()
            
            err = texture.LoadRGBA(rgba)
            if err != nil {
                log.Error("%s fail load rgba '%s': %s", set.Desc(), tag, err)
                texture.Close()
                continue
            }
            
            err = texture.TexImage()
            if err != nil {
                log.Error("%s fail teximage '%s': %s", set.Desc(), tag, err)
                texture.Close()
                continue
            }

            set.texItem[tag] = &TexItem{}
            set.texItem[tag].item = item
            set.texItem[tag].texture = texture
            
            if DEBUG_SET {
                log.Debug("%s prepped %s %.1f",set.Desc(),set.texItem[tag].texture.Desc(),set.texItem[tag].item.timer.Fader())
            }
            
        }
    }
    
    // remove old textures
    for _,item := range old {
        texture := item.texture
        idx := item.item.index
        if DEBUG_SET {
            log.Debug("%s drop #%d: %s",set.Desc(),idx,texture.Desc())
        }
        texture.Close()
    }
    
    

    //setup vertex + bind order arrays

    set.data = []float32{}
    set.tags = []string{}

    idx := 0
    for tag,item := range set.texItem  {

        texture := item.texture

        w := float32( texture.Size.Width / texture.Size.Height )
        h := float32( 1. )

//        w := texture.Size.Width / 2.
//        h := texture.Size.Height/ 2.
//        
//        w = w/h
//        h = h/h

/* 

     A          D
 -w/2,h/2____w/2,h/2
     |          |
     |          |
 -w/2,-h/2___w/2,-h/2
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
        //   x,   y,    z,             tx, ty,    index,  fader
            -w/2.,  +h/2.,  0.0,       0., 0.,   /*   idx,    fdr, */  // A
            -w/2.,  -h/2.,  0.0,       0., 1.,   /*   idx,    fdr, */  // B
            +w/2.,  -h/2.,  0.0,       1., 1.,   /*   idx,    fdr, */  // C
                                                 /*                */
            +w/2.,  -h/2.,  0.0,       1., 1.,   /*   idx,    fdr, */  // C
            +w/2.,  +h/2.,  0.0,       1., 0.,   /*   idx,    fdr, */  // D
            -w/2.,  +h/2.,  0.0,       0., 0.,   /*   idx,    fdr, */  // A
        }
        set.data = append(set.data, data...)
        set.tags = append(set.tags, tag)
//        if DEBUG_SET {
//            log.Debug("%s append #%d '%s' %s",set.Desc(),idx,tag,texture.Desc())
//        }
        idx += 1
        
    }
    
    set.object.BufferData(len(set.data) * 4, set.data)
//    if DEBUG_SET {
//        log.Debug("%s generated %d tags %d float",set.Desc(),len(set.tags),len(set.data))
//    }
    
    
}


func (set *Set) autoScale(camera *gfx.Camera) float32 {

    scaleHeight := float32(1.) / float32(set.buffer.Max())
    return scaleHeight * 2.

}


func (set *Set) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {

	if set.checkRefresh() {
		if DEBUG_SET {
			log.Debug("%s refresh", set.Desc())
		}
		set.generateData(font)
	}

    gl.ActiveTexture(gl.TEXTURE0)
    
	set.program.UseProgram(debug)
	set.object.BindBuffer()

    tagMax := float32( set.buffer.Max() )
    set.program.Uniform1fv(TAGMAX, 1, &tagMax)

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
	set.program.VertexAttribPointer(gfx.TEXCOORD, 2, (3+2)*4, (3)*4)

    count := int32(1)
	offset := int32(0)
	
	for _,tag := range set.tags {

        texture := set.texItem[tag].texture
        item := set.texItem[tag].item
    	

        texture.BindTexture()
        texture.Uniform(set.program)

        var fader float32
        fader = item.timer.Fader()
    	set.program.Uniform1fv(TAGFADER, 1, &fader)
    	
    	var index float32;
    	crc := crc32.Checksum( []byte(tag) , crc32.IEEETable)
        index = float32( uint32(crc) % uint32(set.buffer.Max()) )    	
//    	index = float32(item.index)
    	set.program.Uniform1fv(TAGINDEX, 1, &index)

        tagCount := float32( item.count )
        set.program.Uniform1fv(TAGCOUNT, 1, &tagCount)


//        if DEBUG_SET && verbose {
//        	log.Debug("%s has hash %08x index #%.0f",tag,crc,index)
//        }	

        var width float32;
        width = float32(texture.Size.Width / texture.Size.Height)        
        set.program.Uniform1fv(TAGWIDTH, 1, &width)

        if DEBUG_SET && verbose {
            log.Debug("%s render #%.0f f%.1f  %s",set.Desc(),index,fader,texture.Desc())
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


func (set *Set) Configure(config *TagConfig, camera *gfx.Camera, font *gfx.Font) {

    log.Debug("%s configure %s", set.Desc(), config.Desc())
    
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
			//			err := set.LoadShaders()
			if err != nil {
				set.vert = vert
				set.frag = frag
			}
		}
	}
	
    if config.GetSetDuration() {
        set.buffer.SetDuration( float32(config.GetDuration()) )
    }
	
	if config.GetSetSlot() {
    	set.buffer.Resize( int(config.GetSlot()) )
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
    ret := &TagConfig{
        SetDuration: true, Duration: float64(set.buffer.Duration()),
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



