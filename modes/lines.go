
// +build linux,arm

package modes

import(
    "fmt"
    "strings"
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
//    program  uint32
    object uint32
    vertexShader gfx.Shader
    fragmentShader gfx.Shader
    program uint32
    projection mgl32.Mat4
    projectionUniform int32
    camera mgl32.Mat4
    cameraUniform int32
    model mgl32.Mat4
    modelUniform int32
    textureUniform int32
    vertAttrib uint32
    texCoordAttrib uint32
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


func (lines *Lines) newProgram() {

    lines.fragmentShader = gfx.NewShader("identity",gfx.IDENTITY_FRAGMENT,gl.FRAGMENT_SHADER)
    lines.vertexShader = gfx.NewShader("identity",gfx.IDENTITY_VERTEX,gl.VERTEX_SHADER)
    
    
    lines.fragmentShader.Compile()
    lines.vertexShader.Compile()

    lines.program = gl.CreateProgram()
 
    gl.AttachShader( lines.program, lines.vertexShader.Shader )
    gl.AttachShader( lines.program, lines.fragmentShader.Shader )
    gl.LinkProgram( lines.program)
    
	var status int32
	gl.GetProgramiv(lines.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(lines.program, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(lines.program, logLength, nil, gl.Str(logs))

		log.Error("fail link program: %v", logs)
	}
    
        
}

func (lines *Lines) Init() {


    lines.newProgram()

    windowWidth := 1280
    windowHeight := 960

    
	lines.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0)
	lines.projectionUniform = gl.GetUniformLocation(lines.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(lines.projectionUniform, 1, false, &lines.projection[0])

	lines.camera = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	lines.cameraUniform = gl.GetUniformLocation(lines.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(lines.cameraUniform, 1, false, &lines.camera[0])

	lines.model = mgl32.Ident4()
	lines.modelUniform = gl.GetUniformLocation(lines.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])

	lines.textureUniform = gl.GetUniformLocation(lines.program, gl.Str("tex\x00"))
	gl.Uniform1i(lines.textureUniform, 0)



    gl.GenBuffers(1,&lines.object)
    gl.BindBuffer(gl.ARRAY_BUFFER, lines.object)
    gl.BufferData(gl.ARRAY_BUFFER, len(quad)*4, gl.Ptr(quad), gl.STATIC_DRAW)


	lines.vertAttrib = uint32(gl.GetAttribLocation(lines.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(lines.vertAttrib) 
	gl.VertexAttribPointer(lines.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	lines.texCoordAttrib = uint32(gl.GetAttribLocation(lines.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(lines.texCoordAttrib)
	gl.VertexAttribPointer(lines.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))



    
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
    

    lines.texture = gfx.NewTexture()
    lines.texture.LoadFile("/home/folkert/src/gfx/facade/asset/test.png")
    lines.texture.GenTexture()
    
}


func (lines *Lines) Render() {

    gl.ClearColor(0x00,0x80,0x80,1.0)

	gl.UseProgram(lines.program)
    gl.UniformMatrix4fv(lines.modelUniform, 1, false, &lines.model[0])

    gl.ActiveTexture(gl.TEXTURE0)
    (*lines.texture).Bind()
    gl.BindBuffer(gl.ARRAY_BUFFER, lines.object)
    gl.DrawArrays(gl.TRIANGLES, 0, 2*3 )
    

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



