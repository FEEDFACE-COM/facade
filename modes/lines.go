
// +build linux,arm

package modes

import(
    "fmt"
//    "strings"
	"github.com/go-gl/mathgl/mgl32"    
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)



type Lines struct {
    lineCount uint

    buffer *gfx.Buffer 
    program *gfx.Program


    object *gfx.Object
    data []float32

    white *gfx.Texture

}


func (lines *Lines) Reform() {
    lines.data = []float32{}
    
    for i:=uint(0);i<lines.lineCount;i++ {

        item := lines.buffer.Tail(i)
        var w,h float32
        if item == nil {
            w,h = 0.0,0.0
        } else {
            w = item.Texture.Size.Width / item.Texture.Size.Height
            h = item.Texture.Size.Height / item.Texture.Size.Height
        }
        lines.data = append(lines.data, gfx.QuadVertices(w,h)...   )
    }
    lines.object.BufferData(len(lines.data)*4,lines.data)

}




func (lines *Lines) Queue(text string, font *gfx.Font) {
    log.Debug("queue text: %s",text)
    newText := gfx.NewText(text)
    newText.RenderTexture(font)
    lines.buffer.Queue( newText )
    lines.Reform()
}



func (lines *Lines) Configure(config *conf.LinesConfig) {
    if config == nil {
        return
    }
    log.Debug("configure line: %s",config.Desc())
    if config.LineCount != lines.lineCount {
        lines.lineCount = config.LineCount
        lines.buffer.Resize(config.LineCount)
        lines.Reform()
    }
    

}



func (lines *Lines) Init(camera *gfx.Camera, font *gfx.Font) {
    var err error
    log.Debug("create %s",lines.Desc())


    lines.white = gfx.WhiteColor()
    lines.object.Init()

    lines.Reform()


    err = lines.program.LoadShaders("ident","ident")
    if err != nil { log.Error("fail load lines shaders: %s",err) }
    err = lines.program.LinkProgram(); 
    if err != nil { log.Error("fail link lines program: %v",err) }


}


func (lines *Lines) Render(camera *gfx.Camera, font *gfx.Font, debug bool) {
    

    gl.ClearColor(0.0,0.0,0.0,1.0)

    model := mgl32.Ident4()

    gl.ActiveTexture(gl.TEXTURE0)
    lines.program.UseProgram()
    lines.object.BindBuffer()
    camera.Uniform(lines.program)

    lines.program.VertexAttribPointer(gfx.VERTEX,3,5*4,0)
    lines.program.VertexAttribPointer(gfx.TEXCOORD,2,5*4,3*4)
    modelUniform := lines.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    var d float32 
    if lines.lineCount % 2 == 0 {
        d = float32( int(lines.lineCount/2) ) + 0.5
    } else {
        d = float32( int(lines.lineCount/2) ) + 1.0
    }


    model = model.Mul4( mgl32.Translate3D(0.0,-d,0.0) )
    
    for i:=uint(0);i<lines.lineCount;i++ {
        line  := lines.buffer.Tail(i)

        model = model.Mul4( mgl32.Translate3D(0.0,1.0,0.0) )
        gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
        
        if true && line != nil {
            line.Texture.BindTexture()
            gl.DrawArrays(gl.TRIANGLES, int32(i* 2*3), 2*3)

        }
        if true && line != nil {
            gl.LineWidth(3.0)
            lines.white.BindTexture()
            gl.DrawArrays(gl.LINE_STRIP, int32(i* 2*3), 2*3)
        }    
        
    }
    



}






func NewLines(config *conf.LinesConfig) *Lines {
    if config == nil {
        config = conf.NewLinesConfig()
    }
    ret := &Lines{lineCount: config.LineCount}
    ret.buffer = gfx.NewBuffer(config.LineCount)
    ret.program = gfx.NewProgram("lines")
    ret.object = gfx.NewObject("lines")
    return ret
}

func (lines *Lines) Desc() string {
    ret := fmt.Sprintf("lines[%d]",lines.lineCount)
    item  := lines.buffer.Item(0)
    if item != nil {
        ret += " '" + (*item).Desc() + "'"
    }
    return ret
}

func (lines *Lines) Dump() string {
    return lines.buffer.Dump()   
//    ret := ""; t:=-1
//    for i:=0;i<len(lines.data); i+=5 {
//        if i%(5 * 6)  == 0 { ret += "\n"; t += 1 }
//        ret += fmt.Sprintf("  #%02d pos%v tex%v\n",t,lines.data[i:i+3],lines.data[i+3:i+5])
//    }
//    return ret    
}




