// +build linux,arm

package facade

import (
    "fmt"
    "strings"
    "unicode/utf8"
    gfx "../gfx"
    log "../log"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"
)    


const DEBUG_SET = true



type Set struct {

    vert, frag string
    
//    textures map[string] *gfx.Texture
    words []Word
    widths []float32
        
    wordBuffer *WordBuffer

	texture *gfx.Texture
	program *gfx.Program
 	object  *gfx.Object
	data []float32

	refreshChan chan bool
    
}


const (
    WORDCOUNT  gfx.UniformName = "wordCount"
    WORDINDEX  gfx.UniformName = "wordIndex"
	WORDVALUE  gfx.UniformName = "wordValue"
	WORDWIDTH  gfx.UniformName = "wordWidth"
    WORDTIMER  gfx.UniformName = "wordTimer"
    WORDFADER  gfx.UniformName = "wordFader"

    CHARCOUNT  gfx.UniformName = "charCount"
)


const (
	CHARINDEX gfx.AttribName = "charIndex"
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
    ret.widths = []float32{}
//    ret.textures = make( map[string] *gfx.Texture, ret.wordBuffer.SlotCount() )

    
    ret.refreshChan = make(chan bool, 1)
    return ret
}

func (set *Set) generateData(font *gfx.Font) {
//    old := set.textures
//    
//    set.textures = make( map[string] *gfx.Texture, set.wordBuffer.SlotCount())
    set.words = set.wordBuffer.Words()
    set.widths = make( []float32, len(set.words) )
        
        
//    //generate textures
//    for _,word := range set.words {
//        
//        text := word.text
//        if len(text) <= 0 {
//            continue
//        }
//
//        if old[text] != nil {   //reuse existing textures
//
//            set.textures[text] = old[text]
//            if DEBUG_SET {
//                log.Debug("%s texture reused: %s",set.Desc(),old[text].Desc())
//            }
//
//        } else {               //create new texture
//            
//            rgba, err := font.RenderText(text, false)
//            if err != nil {
//                log.Error("%s texture fail render '%s': %s", set.Desc(), text, err)
//                continue
//            } 
//
//            texture := gfx.NewTexture(text)
//            texture.Init()
//            
//            err = texture.LoadRGBA(rgba)
//            if err != nil {
//                log.Error("%s texture fail load rgba '%s': %s", set.Desc(), text, err)
//                texture.Close()
//                continue
//            }
//            
//            err = texture.TexImage()
//            if err != nil {
//                log.Error("%s texture fail teximage '%s': %s", set.Desc(), text, err)
//                texture.Close()
//                continue
//            }
//
//            set.textures[text] = texture
//            
//            if DEBUG_SET {
//                log.Debug("%s texture prepped: %s",set.Desc(),set.textures[text].Desc())
//            }
//            
//        }
//    }
//    
//    // remove unused textures
//    for text,texture := range old {
//        _,ok := set.textures[text]
//        if !ok {
//            if DEBUG_SET {
//                log.Debug("%s texture close: %s",set.Desc(),texture.Desc())
//            }
//            texture.Close()
//        }
//    }
    
    

    //setup vertex + bind order arrays
    set.data = []float32{}
    
    charCount := 0    
    for i,word := range set.words  {
        
        width := float32(0.)
        for _,run := range word.text {
            glyphCoord := getGlyphCoord(run)
            glyphSize := font.Size(glyphCoord.X, glyphCoord.Y)
            width += glyphSize.W/glyphSize.H
            charCount += 1
        }
        set.widths[i] = width
        set.data = append(set.data, set.vertices(word,font,width)... )
    }

    set.object.BufferData(len(set.data) * 4, set.data)
    if DEBUG_SET {
        log.Debug("%s generated words:%d chars:%d floats:%d",set.Desc(),len(set.words),charCount,len(set.data))
    }

}

func (set *Set) vertices(
    word Word, 
    font *gfx.Font, 
    totalWidth float32,
) []float32 {

    var ret = []float32{}

    charIndex := 0
    offset := -totalWidth/2.
    for _,run := range word.text {
        
        glyphCoord := getGlyphCoord(run)
        glyphSize := font.Size(glyphCoord.X, glyphCoord.Y)
        maxGlyphSize := font.MaxSize()

        texOffset := gfx.Point{
            X: float32(glyphCoord.X) / (gfx.GlyphMapCols),
            Y: float32(glyphCoord.Y) / (gfx.GlyphMapRows),
        }



        ox, oy := texOffset.X, texOffset.Y
        th := 1. / float32(gfx.GlyphMapRows)
        tw := glyphSize.W / (maxGlyphSize.W * float32(gfx.GlyphMapCols))


        w := glyphSize.W / glyphSize.H
        h := float32(1.)


        dx := offset + w/2.

        idx := float32(charIndex)

        data := []float32{
        //   x,   y,    z,             tx, ty,   
           dx -w/2.,  +h/2.,  0.0,    0. + ox, 0. + oy,  idx,      // A
           dx -w/2.,  -h/2.,  0.0,    0. + ox, th + oy,  idx,      // B
           dx +w/2.,  -h/2.,  0.0,    tw + ox, th + oy,  idx,      // C

           dx +w/2.,  -h/2.,  0.0,    tw + ox, th + oy,  idx,      // C
           dx +w/2.,  +h/2.,  0.0,    tw + ox, 0. + oy,  idx,      // D
           dx -w/2.,  +h/2.,  0.0,    0. + ox, 0. + oy,  idx,      // A
        }
        ret = append(ret, data...)

        offset += w
        charIndex += 1
    }

    if DEBUG_SET {
        log.Debug("%s data generate '%s'",set.Desc(),word.text)
    }

    return ret
}

//        text := word.text
//        if len(text) <= 0 {
//            continue
//        }
        
//        texture,ok := set.textures[text]
//        if !ok {
//            log.Debug("%s texture generate miss: %s",set.Desc(),text)
//            continue
//        }
//
//        w := float32( texture.Size.Width / texture.Size.Height )
//        h := float32( 1. )
//
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
		set.renderMap(font)
	}

    gl.ActiveTexture(gl.TEXTURE0)
    
	set.program.UseProgram(debug)
	set.object.BindBuffer()

    wordCount := float32( set.wordBuffer.SlotCount() )
    set.program.Uniform1fv(WORDCOUNT, 1, &wordCount)

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
	

	set.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+1)*4, (0)*4)
	set.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+1)*4, (0+3)*4)
	set.program.VertexAttribPointer(    CHARINDEX, 1, (3+2+1)*4, (0+3+2)*4)

    count := int32(1)
	offset := int32(0)
	
	for i,word := range set.words {

        text := word.text
        
        if len(text) <= 0 {
            continue
        }


    	set.texture.Uniform(set.program)

        
//        texture,ok := set.textures[text]
//        if !ok {
//            log.Debug("%s texture render miss: %s",set.Desc(),text)
//            continue
//        }
//
//
//        texture.BindTexture()
//        texture.Uniform(set.program)


        count = int32( utf8.RuneCountInString(word.text) )

        var timer float32
        timer = word.timer.Edge(gfx.Now())
    	set.program.Uniform1fv(WORDTIMER, 1, &timer)

        var fader float32
        fader = word.timer.Fader()
    	set.program.Uniform1fv(WORDFADER, 1, &fader)
    	
    	var index float32;
    	index = float32(word.index)
    	set.program.Uniform1fv(WORDINDEX, 1, &index)

        wordValue := float32( word.count )
        set.program.Uniform1fv(WORDVALUE, 1, &wordValue)


        var width float32;
        width = set.widths[i]
        set.program.Uniform1fv(WORDWIDTH, 1, &width)


        var charCount float32;
        charCount = float32( count )
        set.program.Uniform1fv(CHARCOUNT, 1, &charCount)

        if DEBUG_SET && verbose {
            log.Debug("%s render #%d width:%.1f fader:%.1f",set.Desc(),index,width,fader)
        }
                
        if !debug || debug {
            set.program.SetDebug(false)
        	set.texture.BindTexture()
            gl.DrawArrays(gl.TRIANGLES, int32(offset*(2*3)), int32(count*2*3))
            set.program.SetDebug(debug)
        }
        
        if debug {
            gl.LineWidth(3.0)
            gl.BindTexture(gl.TEXTURE_2D, 0)
            gl.DrawArrays(gl.LINE_STRIP, int32(offset*(2*3)), int32(count*2*3))
        }
        offset += count
    }
    


}




func (set *Set) Init(programService *gfx.ProgramService, font *gfx.Font) {
	log.Debug("%s init", set.Desc())

	set.object = gfx.NewObject("set")
	set.object.Init()

    set.texture = gfx.NewTexture("set")
    set.texture.Init()

	set.program = programService.GetProgram("set", "set/")
	set.program.Link(set.vert, set.frag)

    set.renderMap(font)

	set.ScheduleRefresh()
        
}

func (set *Set) renderMap(font *gfx.Font) error {

	if DEBUG_SET {
		log.Debug("%s render texture map %s", set.Desc(), font.Desc())
	}

	rgba, err := font.RenderMap(false)
	if err != nil {
		log.Error("%s fail render font map: %s", set.Desc(),err)
		return log.NewError("fail render font map: %s", err)
	}
	err = set.texture.LoadRGBA(rgba)
	if err != nil {
		log.Error("fail load font map: %s", err)
		return log.NewError("fail to load font map: %s", err)
	}
	set.texture.TexImage()

	return nil
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
        case "nato":
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



