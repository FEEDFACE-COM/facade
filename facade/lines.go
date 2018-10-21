
// +build linux,arm

package facade

import(
//    "fmt"
//    "strings"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)



type Lines struct {
    config LinesConfig

    buffer *gfx.Buffer 
    program *gfx.Program


    object *gfx.Object
    data []float32

    white *gfx.Texture

}


func (lines *Lines) generateData() {
    lines.data = []float32{}
    
    for i:=uint(0);i<lines.config.Height;i++ {

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
    newText := gfx.NewText(text)
    newText.RenderTexture(font)
    lines.buffer.Queue( newText )
    lines.generateData()
//    log.Debug("queued text: %s",text)
}



func (lines *Lines) Configure(config *LinesConfig) {
    if config == nil { return }
    if *config == lines.config { return}
    
    old := lines.config
    lines.config = *config
    log.Debug("config %s",config.Desc())
    
    if config.Height != old.Height {
        lines.buffer.Resize(config.Height)
        lines.generateData()
    }

}



func (lines *Lines) Init(camera *gfx.Camera, font *gfx.Font) {
    var err error
    log.Debug("create %s",lines.Desc())


    lines.white = gfx.WhiteColor()
    lines.object.Init()

    lines.generateData()


    err = lines.program.GetCompileShaders("","identity","identity")
    if err != nil { log.Error("fail load lines shaders: %s",err) }
    err = lines.program.LinkProgram(); 
    if err != nil { log.Error("fail link lines program: %v",err) }


}


func (lines *Lines) Render(camera *gfx.Camera, debug, verbose bool) {
    

    gl.ClearColor(0.0,0.0,0.0,1.0)

    model := mgl32.Ident4()

    gl.ActiveTexture(gl.TEXTURE0)
    lines.program.UseProgram(debug)
    lines.object.BindBuffer()
    camera.Uniform(lines.program)

    lines.program.VertexAttribPointer(gfx.VERTEX,3,5*4,0)
    lines.program.VertexAttribPointer(gfx.TEXCOORD,2,5*4,3*4)
    modelUniform,_ := lines.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    var d float32 
    if lines.config.Height % 2 == 0 {
        d = float32( int(lines.config.Height/2) ) + 0.5
    } else {
        d = float32( int(lines.config.Height/2) ) + 1.0
    }


    model = model.Mul4( mgl32.Translate3D(0.0,-d,0.0) )
    
    for i:=uint(0);i<lines.config.Height;i++ {
        line  := lines.buffer.Tail(i)

        model = model.Mul4( mgl32.Translate3D(0.0,1.0,0.0) )
        gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
        
        if true && line != nil { //textures
            line.Texture.BindTexture()
            gl.DrawArrays(gl.TRIANGLES, int32(i* 2*3), 2*3)

        }
        if false && line != nil { //lines
            gl.LineWidth(3.0)
            lines.white.BindTexture()
            gl.DrawArrays(gl.LINE_STRIP, int32(i* 2*3), 2*3)
        }    
        
    }
    



}






func NewLines(config *LinesConfig) *Lines {
    if config == nil {
        config = NewLinesConfig()
    }
    ret := &Lines{config: *config}
    ret.buffer = gfx.NewBuffer(config.Height)
    ret.program = gfx.GetProgram("lines")
    ret.object = gfx.NewObject("lines")
    return ret
}

func (lines *Lines) Desc() string { return lines.config.Desc() }

func (lines *Lines) Dump() string {
    ret := lines.config.Desc()
    if lines.buffer.Tail(0) != nil {
        ret += " '" + (*lines.buffer.Tail(0)).Desc() + "'"
    }
    return lines.buffer.Dump()   
}




