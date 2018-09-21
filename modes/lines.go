
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
    texture *gfx.Texture
    camera *gfx.Camera
    program uint32
    model mgl32.Mat4
    modelUniform int32
    quad *gfx.Quad
}


type lineItem struct {
    text string
}

func (item *lineItem) Desc() string { return item.text }


func (lines *Lines) Queue(text string) {
    newItem := lineItem{text: text}
    lines.buffer.Queue( &newItem )
}

func (lines *Lines) Configure(config *conf.LineConfig) {
    if config == nil {
        return
    }
    log.Debug("configure line: %s",config.Desc())
    if config.LineCount != lines.lineCount {
        lines.lineCount = config.LineCount
        lines.buffer.Resize(config.LineCount)
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
    ret += "\n>> "
    for i:=uint(0);i<lines.lineCount;i++ {
        item := lines.buffer.Item(i)
        if item != nil { ret += (*item).Desc() }
        ret += ","
    }
    ret += "\n<< "
    for i:=uint(lines.lineCount);i>0;i-- {
        item := lines.buffer.Item(i-1)
        if item != nil { ret += (*item).Desc() }
        ret += ","
    }
    ret += "\n"
    return ret
}

func (lines *Lines) Dump() string {
    return lines.buffer.Dump()
}


func (lines *Lines) Init(camera *gfx.Camera) {
    var err error

	lines.texture = gfx.NewTexture()
	lines.texture.LoadFile("/home/folkert/src/gfx/facade/asset/FEEDFACE.COM.white.png")
	lines.texture.GenTexture()

    fragment := gfx.NewShader("identity",gfx.IDENTITY_FRAGMENT,gl.FRAGMENT_SHADER)
    vertex := gfx.NewShader("identity",gfx.IDENTITY_VERTEX,gl.VERTEX_SHADER)
    
    err = fragment.Compile()
    if err != nil {
        log.Error("fail compile fragment: %v",err)
    }

    err = vertex.Compile()
    if err != nil {
        log.Error("fail compile vertex: %v",err)
    }

    lines.program, err = gfx.NewProgram(&vertex,&fragment)

    if err != nil {
        log.Error("fail new program: %v",err)    
    }

	gl.UseProgram(lines.program)

	camera.Uniform(lines.program)

	lines.model = mgl32.Ident4()
	lines.modelUniform = gl.GetUniformLocation(lines.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])

    lines.texture.Uniform(lines.program)

    lines.quad = gfx.NewQuad(lines.texture.Size.Width,lines.texture.Size.Height)
    lines.quad.VertexAttribPointer(lines.program)

    
}


func (lines *Lines) Render(camera *gfx.Camera) {

    gl.ClearColor(0x00,0x80,0x80,1.0)

    lines.model = mgl32.Ident4()

    gl.UseProgram(lines.program)
    gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])




    camera.Uniform(lines.program)

    gl.ActiveTexture(gl.TEXTURE0)
    lines.texture.Bind()

    gl.DrawArrays(gl.TRIANGLES, 0, 2*3)


}


var quad = []float32 { 
// x, y, z, u, v    

    
    -1.0,  1.0, 0.0, 0.0, 0.0,
    -1.0, -1.0, 0.0, 0.0, 1.0,
     1.0, -1.0, 0.0, 1.0, 1.0,
     
    -1.0,  1.0, 0.0, 0.0, 0.0,
     1.0, -1.0, 0.0, 1.0, 1.0,
     1.0,  1.0, 0.0, 1.0, 0.0,    
    
}



