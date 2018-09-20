
// +build linux,arm

package modes

import(
    "fmt"
//    "strings"
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)

type Lines struct {
//    buffer Buffer 
//    texture *gfx.Texture
//    program  uint32
    vertexShader gfx.Shader
    fragmentShader gfx.Shader
}





func (lines *Lines) Queue(text string) {
//    newLine := gfx.NewLine(text)
    
//    lines.buffer.Queue(newLine)
}

func (lines *Lines) Configure(config *conf.LineConfig) {
    log.Debug("configure line: %s",config.Desc())
    
//    lines.texture = gfx.NewTexture()
//    lines.texture.LoadFile("/home/folkert/src/gfx/facade/asset/test.png")
    
//    lines.texture.GenTexture()


    
    
}

func NewLines() Lines {
    ret := Lines{}
//    ret.buffer = NewBuffer(1)
//    ret.Configure(conf.NewLineConfig())
    return ret
}

func (lines *Lines) Desc() string {
    ret := fmt.Sprintf("lines[]")
//    ret := fmt.Sprintf("lines[%d]",lines.buffer.count)
//    ret += lines.buffer.Describe() + "\n"
    return ret
}


func (lines *Lines) Setup() {
    
        _foo := gl.FRAGMENT_SHADER
//    _foo := "xx"      
        
    
    log.Debug("SETUP WITH SHADER %v",_foo)
    
//    lines.fragmentShader = gfx.NewShader("identity",gfx.FRAGMENT_IDENTITY,gl.FRAGMENT_SHADER)
//    lines.vertexShader = gfx.NewShader("identity",gfx.VERTEX_IDENTITY,gl.VERTEX_SHADER)
//    
//    
//    lines.fragmentShader.Compile()
//    lines.vertexShader.Compile()
//    
//    lines.program = gl.CreateProgram()
//    
//    gl.AttachShader( lines.program, lines.vertexShader.Shader )
//    gl.AttachShader( lines.program, lines.fragmentShader.Shader )
//    gl.LinkProgram( lines.program)
//    
//    
//    //check
//
//	var status int32
//	gl.GetProgramiv(lines.program, gl.LINK_STATUS, &status)
//	if status == gl.FALSE {
//		var logLength int32
//		gl.GetProgramiv(lines.program, gl.INFO_LOG_LENGTH, &logLength)
//
//		logs := strings.Repeat("\x00", int(logLength+1))
//		gl.GetProgramInfoLog(lines.program, logLength, nil, gl.Str(logs))
//
//		log.Error("fail link program: %s",logs)
//	}
//    

    
//	lines.projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0)
//	lines.projectionUniform := gl.GetUniformLocation(lines.program, gl.Str("projection\x00"))
//	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
//
//	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
//	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
//	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
//
//	model := mgl32.Ident4()
//	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
//	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
//
//	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
//	gl.Uniform1i(textureUniform, 0)

    
    
}


func (lines *Lines) Render() {

//    gl.ActiveTexture(gl.TEXTURE0)
//    gl.BindTexture(gl.TEXTURE_2D, lines.texture)
//    
//    gl.DrawArrays(GL.TRIANGLES, 0, 2*2 )
    

}



