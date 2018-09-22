
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
    camera *gfx.Camera

    program uint32
    model mgl32.Mat4
    modelUniform int32


    object uint32
    verts []float32

    vertAttrib uint32
    texCoordAttrib uint32

}





func (lines *Lines) setVBO() {
    lines.verts = []float32{}
//    lines.verts = append( lines.verts, gfx.QuadVertices(2.0,1.0)...)
    for i:=uint(0);i<lines.lineCount;i++ {
        item := lines.buffer.Tail(i)
        var w,h float32
        if item == nil {
            w,h = 0.0,0.0
        } else {
            w = item.Texture.Size.Width / item.Texture.Size.Height
            h = item.Texture.Size.Height / item.Texture.Size.Height
        }
        lines.verts = append(lines.verts, gfx.QuadVertices(w,h)...   )
    }
    gl.BindBuffer(gl.ARRAY_BUFFER, lines.object)
    gl.BufferData(gl.ARRAY_BUFFER, len(lines.verts)*4, gl.Ptr(lines.verts), gl.STATIC_DRAW)
}



func (lines *Lines) Queue(text string, font *gfx.Font) {
    log.Debug("queue text: %s",text)
    newText := gfx.NewText(text)
    newText.RenderTexture(font)
    lines.buffer.Queue( newText )
    lines.setVBO()
}

func (lines *Lines) Configure(config *conf.LineConfig) {
    if config == nil {
        return
    }
    log.Debug("configure line: %s",config.Desc())
    if config.LineCount != lines.lineCount {
        lines.lineCount = config.LineCount
        lines.buffer.Resize(config.LineCount)
        lines.setVBO()
    }
    

}

func NewLines(config *conf.LineConfig) *Lines {
    if config == nil {
        config = conf.NewLineConfig()
    }
    ret := &Lines{lineCount: config.LineCount}
    ret.buffer = gfx.NewBuffer(config.LineCount)
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
//    return lines.buffer.Dump()
    return lines.dumpVBO()
}


func (lines *Lines) Init(camera *gfx.Camera) {
    var err error

    log.Debug("create vbo[%d]",lines.lineCount)
    gl.GenBuffers(1,&lines.object)

    lines.setVBO()

    fragment := gfx.NewShader("identity",gfx.IDENTITY_FRAGMENT,gl.FRAGMENT_SHADER)
    vertex := gfx.NewShader("identity",gfx.IDENTITY_VERTEX,gl.VERTEX_SHADER)
    
    lines.program, err = gfx.NewProgram(&vertex,&fragment)
    if err != nil {
        log.Error("fail new program: %v",err)    
        return
    }

	gl.UseProgram(lines.program)

	camera.Uniform(lines.program)

	lines.model = mgl32.Ident4()
	lines.modelUniform = gl.GetUniformLocation(lines.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])


	lines.vertAttrib = uint32(gl.GetAttribLocation(lines.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(lines.vertAttrib) 
	gl.VertexAttribPointer(lines.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	lines.texCoordAttrib = uint32(gl.GetAttribLocation(lines.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(lines.texCoordAttrib)
	gl.VertexAttribPointer(lines.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))


    
}


func (lines *Lines) dumpVBO() string {
    ret := ""
    t:=-1
    for i:=0;i<len(lines.verts); i+=5 {
        if (i)%(5 * 6)  == 0 {
            ret += "\n"
            t += 1    
        }
        ret += fmt.Sprintf("  #%02d x%5.1f y%5.1f z%5.1f  u%5.1f v%5.1f\n", t,
            lines.verts[i], lines.verts[i+1], lines.verts[i+2], lines.verts[i+3], lines.verts[i+4])
            
    }
    return ret    
}

func (lines *Lines) Render(camera *gfx.Camera) {

    gl.ClearColor(0.23,0.23,0.23,1.0)

    c := float32(lines.lineCount)  
    z := 1./c
//    d := 1.5/c
    lines.model = mgl32.Ident4()
    lines.model = lines.model.Mul4( mgl32.Scale3D(z,z,z) )
    lines.model = lines.model.Mul4( mgl32.Translate3D(0.0,c/2.+0.5,0.0) )
    

    gl.UseProgram(lines.program)
//    gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])

    camera.Uniform(lines.program)
    gl.ActiveTexture(gl.TEXTURE0)


const DRAW_TEXT = true
const DRAW_BOX = false

//    log.Debug(lines.dumpVBO())

    tmp := "dump VBO:\n"
    for i:=uint(0);i<lines.lineCount;i++ {
        line  := lines.buffer.Tail(i)
//        idx := lines.buffer.Index(i)    
        lines.model = lines.model.Mul4( mgl32.Translate3D(0.0,-1.0,0.0) )

        if DRAW_TEXT {
            if line != nil { 
                tmp += line.Text + " "
    //            log.Debug("got tex %5.0fx%5.0f",line.Texture.Size.Width,line.Texture.Size.Height)
                gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])
                line.Texture.Bind()
                gl.BindBuffer(gl.ARRAY_BUFFER,lines.object) 
                gl.DrawArrays(gl.TRIANGLES, 0, 2*3)
            } else { tmp += "XXX " }

        }
        if DRAW_BOX {
            if line != nil {
                gl.LineWidth(3.0)
                gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])
                gl.BindBuffer(gl.ARRAY_BUFFER,lines.object) 
                gl.DrawArrays(gl.LINE_STRIP, 6, 2*3)
            }
        }    
        
    }
//    log.Debug("draw: %s",tmp)
    



}




