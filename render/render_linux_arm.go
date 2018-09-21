
// +build linux,arm

package render

import (
//    "fmt"
    "os"
    "time"
    "strings"
    "sync"
    log "../log"
    conf "../conf"
    modes "../modes"
    gfx "../gfx"
    "src.feedface.com/gfx/piglet"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
    "image"
    "image/draw"
)





const RENDERER_AVAILABLE = true

const FRAME_RATE = 60.0

const BUFFER_SIZE = 40

type Renderer struct {
    size struct{width int32; height int32}

    mode conf.Mode
    grid *modes.Grid
    lines *modes.Lines
    font *gfx.Font

    now Clock
    buffer *Buffer
    mutex *sync.Mutex
}

func NewRenderer() *Renderer {
    ret := &Renderer{}
    ret.mutex = &sync.Mutex{}
    ret.buffer = &Buffer{}
    ret.buffer.Resize(BUFFER_SIZE)
    return ret
}

const DEBUG_CLOCK  = false
const DEBUG_MODE   = true
const DEBUG_BUFFER = true
 

const DEBUG_FRAMES = 90

func (renderer *Renderer) Init(config *conf.Config) error {
    var err error
    log.Debug("initializing renderer")
    
    err = piglet.CreateContext()
    if err != nil {
        log.PANIC("fail to initialize renderer: %s",err)    
    }
    
    w,h := piglet.GetDisplaySize()
    renderer.size = struct{width int32; height int32} {int32(w),int32(h)}
    log.Info("got display %dx%d",renderer.size.width,renderer.size.height)
    

    piglet.MakeCurrent()

    err = gl.InitWithProcAddrFunc( piglet.GetProcAddress )
    if err != nil {
        log.PANIC("fail to init GLES: %s",err)    
    }
    

    log.Debug("got renderer %s %s", gl.GoStr(gl.GetString((gl.VENDOR))),gl.GoStr(gl.GetString((gl.RENDERER))));
    log.Debug("got version %s %s", gl.GoStr(gl.GetString((gl.VERSION))),gl.GoStr(gl.GetString((gl.SHADING_LANGUAGE_VERSION))));


    //setup things    
    renderer.mode = config.Mode
    renderer.grid = modes.NewGrid(config.Grid)
    renderer.lines = modes.NewLines(config.Line)
    renderer.font = gfx.NewFont(config.Font)


    InitClock()
    renderer.now = Clock{frame: 0}

    return err
}


func (renderer *Renderer) Configure(config *conf.Config) error {
    
    if config == nil {
        return nil
    }
    
    log.Debug("configure %s",config.Desc())
    
    if renderer.mode != config.Mode {
        log.Debug("switch mode to %s",string(config.Mode))
    }
    
    renderer.mode = config.Mode
    renderer.font.Configure(config.Font,conf.DIRECTORY)
    renderer.lines.Configure(config.Line)
    renderer.grid.Configure(config.Grid)
    return nil
}



func (renderer *Renderer) Render(confChan chan conf.Config, textChan chan conf.Text) error {

    
    var now *Clock = &renderer.now
    var prev Clock = *now
    


    log.Debug("renderer start")
//    gl.ClearColor(0x0, 0x0, 0x0, 1.0)
    gl.ClearColor(0xFE,0xED,0xFA,0xCE)
    gl.Viewport(0, 0, renderer.size.width,renderer.size.height)


    renderer.lines.Init()


//////////
    var vertexShader = `

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
attribute vec3 vert;
attribute vec2 vertTexCoord;
varying vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

    var fragmentShader = `
uniform sampler2D tex;
varying vec2 fragTexCoord;
void main() {
    gl_FragColor = texture2D(tex,fragTexCoord);
}
` + "\x00"

	program, err := newProgram(vertexShader, fragmentShader)

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(renderer.size.width)/float32(renderer.size.height), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)


  // Load the texture
	texture, err := newTexture("/home/folkert/src/gfx/facade/asset/test.png")
	if err != nil {
		log.PANIC("fail to load texture: %v",err)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib) 
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

//    gl.Enable(gl.CULL_FACE)
//    gl.CullFace(gl.BACK)


//////////


    for {
        now.Tick()
        
        renderer.mutex.Lock()
        piglet.MakeCurrent()
        
        
        select {
            case config := <-confChan:
                log.Debug("conf: %s",config.Desc())
                renderer.Configure(&config)
            default:
        }
        
        select {
            case text := <-textChan:
                log.Debug("read: %s",text)
                renderer.buffer.Queue(renderer.now.Time(),text)
                renderer.lines.Queue( string(text) )
                renderer.grid.Queue( string(text) )
            default:
        }

        
        
        gl.BindFramebuffer(gl.FRAMEBUFFER,0)
        gl.Clear( gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )


//////////



        angle := renderer.now.cycle
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
//
//		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])


		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)








//////////



        switch renderer.mode {
            case conf.GRID:
                //renderer.grid.Render()
            case conf.LINE:
                //renderer.lines.Render()
        }
        
        


        if now.frame % DEBUG_FRAMES == 0 {
            if DEBUG_CLOCK   {
                fps := float32(now.frame - prev.frame) / (now.time - prev.time)
                log.Debug("frame %05d %s    %4.1ffps",now.frame,now.Desc(),fps)
                prev = *now
            }
            
//            if DEBUG_BUFFER {
//                log.Debug("%5.1f %5.1f %s",renderer.now.Time(),now.Time(),renderer.buffer.Dump())    
//            }
            
            if DEBUG_MODE {
                switch renderer.mode { 
                    case conf.LINE:
                        log.Debug( renderer.lines.Desc() + " " +renderer.font.Desc() )
                        if DEBUG_BUFFER {
                            log.Debug( renderer.lines.Dump() )    
                        }
                    case conf.GRID:
                        log.Debug( renderer.grid.Desc() + " " +renderer.font.Desc() )
                        if DEBUG_BUFFER {
                            log.Debug( renderer.grid.Dump() )    
                        }
                }
            }
            
        }
        piglet.SwapBuffers()
        renderer.mutex.Unlock()
        
        // wait for next frame
        // FIXME, maybe dont wait as long??
        time.Sleep( time.Duration( int64(time.Second / FRAME_RATE) ) )
    }
    return nil
}


func (renderer *Renderer) ReadText(textChan chan conf.Text) error {
    for {
        text := <-textChan
        log.Debug("read: %s",text)
        renderer.mutex.Lock()
        renderer.buffer.Queue(renderer.now.Time(),text)
//        renderer.lines.Queue(string(text))
//        renderer.grid.Queue( string(text) )
//        renderer.lines.Queue( string(text) )
        renderer.mutex.Unlock()
    }
    return nil
    
}


func (renderer *Renderer) ReadConf(confChan chan conf.Config) error {
    for {
        config := <-confChan
        log.Debug("conf: %s",config.Desc())    
        renderer.mutex.Lock()
//        renderer.Configure(&config)
        renderer.mutex.Unlock()
    }
    return nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
    log.Debug("new program...")
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logs))

		log.Error("failed to link program: %v", logs)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
    log.Debug("compile shader...")
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logs))

		log.Error("failed to compile %v: %v", source, logs)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
    log.Debug("new texture...")
	imgFile, err := os.Open(file)
	if err != nil {
		log.PANIC("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.PANIC("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	 1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0,  1.0, 0.0, 1.0,
	 1.0, -1.0, -1.0, 1.0, 0.0,
	 1.0, -1.0,  1.0, 1.0, 1.0,
	-1.0, -1.0,  1.0, 0.0, 1.0,

    // Top
    -1.0, 1.0, -1.0, 0.0, 0.0,
    -1.0, 1.0,  1.0, 0.0, 1.0,
     1.0, 1.0, -1.0, 1.0, 0.0,
     1.0, 1.0, -1.0, 1.0, 0.0,
    -1.0, 1.0,  1.0, 0.0, 1.0,
     1.0, 1.0,  1.0, 1.0, 1.0,

//	// Front
    -1.0, -1.0, 1.0, 1.0, 0.0,
     1.0, -1.0, 1.0, 0.0, 0.0,
    -1.0,  1.0, 1.0, 1.0, 1.0,
     1.0, -1.0, 1.0, 0.0, 0.0,
     1.0,  1.0, 1.0, 0.0, 1.0,
    -1.0,  1.0, 1.0, 1.0, 1.0,

//	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

    // Right
    1.0, -1.0,  1.0, 1.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0,  1.0, -1.0, 0.0, 0.0,
    1.0, -1.0,  1.0, 1.0, 1.0,
    1.0,  1.0, -1.0, 0.0, 0.0,
    1.0,  1.0,  1.0, 0.0, 1.0,
}

