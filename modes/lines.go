
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

    buffer Buffer 
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
        item := lines.buffer.Item(i)
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

func (line *Line) Desc() string { return line.Text }

//func (line *Line) Close() { log.Debug("close line[%s]",line.Text) }
//
//func (line *Line) Bind(program uint32)  { 
//    line.Texture.Bind()
//    line.Quad.Bind()
//    line.Quad.VertexAttribPointer(program)
//}

func (lines *Lines) Queue(text string, font *gfx.Font) {
    newLine := NewLine(text,font)
    log.Debug("queue %s",text)
    lines.buffer.Queue( newLine )
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
    ret.buffer = NewBuffer(config.LineCount)
    return ret
}

func (lines *Lines) Desc() string {
    ret := fmt.Sprintf("lines[%d]",lines.lineCount)
    item  := lines.buffer.Item(0)
    if item != nil {
        ret += " '" + (*item).Desc() + "'"
    }
//    ret += "\n>> "
//    for i:=uint(0);i<lines.lineCount;i++ {
//        item := lines.buffer.Item(i)
//        if item != nil { ret += (*item).Desc() }
//        ret += ","
//    }
//    ret += "\n<< "
//    for i:=uint(lines.lineCount);i>0;i-- {
//        item := lines.buffer.Item(i-1)
//        if item != nil { ret += (*item).Desc() }
//        ret += ","
//    }
//    ret += "\n"
    return ret
}

func (lines *Lines) Dump() string {
    return lines.buffer.Dump()
}

func NewLine(text string, font *gfx.Font) Line {
    ret := Line{Text: text}
    
    ret.Texture = gfx.NewTexture()
	ret.Texture.LoadFile("/home/folkert/src/gfx/facade/asset/FEEDFACE.COM.white.png")
	ret.Texture.GenTexture()
    
    
    return ret
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


func (lines *Lines) Render(camera *gfx.Camera) {

    gl.ClearColor(0x00,0x80,0x80,1.0)

    lines.model = mgl32.Ident4()

    gl.UseProgram(lines.program)
    gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])

    camera.Uniform(lines.program)
    gl.ActiveTexture(gl.TEXTURE0)




    line  := lines.buffer.Item(0)
    if line != nil { 
        line.Texture.Bind()
        gl.BindBuffer(gl.ARRAY_BUFFER,lines.object) 
        gl.DrawArrays(gl.TRIANGLES, 0, 2*3)
    }



}




