
package modes

import(
    "fmt"
	"github.com/go-gl/mathgl/mgl32"    
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)

type Grid struct {
    width uint
    height uint
    

    buffer *gfx.Buffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    black *gfx.Texture
}



func (grid *Grid) Render(camera *gfx.Camera, debug bool) {
    gl.ClearColor(1.0,1.0,1.0,1.0)
    gl.ActiveTexture(gl.TEXTURE0)
    

    
    grid.program.UseProgram()
    grid.object.BindBuffer()
    

    model := mgl32.Ident4()
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    camera.Uniform(grid.program)
    grid.program.Uniform1i(gfx.TEXTURE,0)
    
    grid.program.VertexAttribPointer(gfx.VERTEX,3,5*4,0)
    grid.program.VertexAttribPointer(gfx.TEXCOORD,2,5*4,3*4)
    

    if true {    
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, 0, (2*3)*int32(grid.height*grid.width)  )
    }
    
    
    if true {
        gl.LineWidth(3.0)
        grid.black.BindTexture()
        gl.DrawArrays(gl.LINES, 0, (2*3)*int32(grid.height*grid.width) )        
    }
    
    
}










func gridVertices(x,y,w,h float32) []float32 {
    return []float32{
        
        -w/2+x,  h/2+y, 0,            0,  0,
        -w/2+x, -h/2+y, 0,            0,  1,
         w/2+x, -h/2+y, 0,            1,  1,
         w/2+x, -h/2+y, 0,            1,  1,
         w/2+x,  h/2+y, 0,            1,  0,
        -w/2+x,  h/2+y, 0,            0,  0,
        
    }
    
}





func (grid *Grid) generateData() {
    grid.data = []float32{}
    
    var w,h float32 = 1.0,1.0
    
    for y := uint(0); y<grid.height; y++ {
        for x:=uint(0); x<grid.width; x++ {
            
            var d float32
            if grid.height % 2 == 0 {
                d = 0.5
            } else {
                d = -0.5 
            }
            
            xx := float32(x) - w * float32(grid.width)/2.  + d
            yy := float32(y) - h * float32(grid.height)/2. + d
            grid.data = append(grid.data, gridVertices( xx,yy, w,h )... )
        }
        
    }
    
    grid.data = append(grid.data, gfx.QuadVertices(w,h)...)
    grid.object.BufferData(len(grid.data)*4,grid.data)
    
}




func (grid *Grid) Init(camera *gfx.Camera, font *gfx.Font) {
    var err error
    log.Debug("create %s",grid.Desc())


//    err = grid.texture.LoadFile("/home/folkert/src/gfx/facade/asset/test.png")

    rgba, err := font.RenderMapRGBA()
    if err != nil {
        log.Error("fail render font map: %s",err)
    }
    err = grid.texture.LoadRGBA(rgba)
    if err != nil {
        log.Error("fail load font map: %s",err)
    }


    
    grid.texture.TexImage2D()
    
    grid.black = gfx.BlackColor()
    grid.object.Init()
    
    grid.generateData()

    err = grid.program.LoadShaders("grid","grid")
    if err != nil { log.Error("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { log.Error("fail link grid program: %v",err) }

    

}




func (grid *Grid) Queue(text string) {
    newText := gfx.NewText(text)
    grid.buffer.Queue( newText )
    grid.generateData()
    log.Debug("queued text: %s",text)
}







func (grid *Grid) Configure(config *conf.GridConfig) {
    if config == nil {
        return
    }
    
    if config.Width != grid.width {
        grid.width = config.Width    
    }
    
    if config.Height != grid.height {
        grid.height = config.Height
        grid.buffer.Resize(config.Height)    
    }

    log.Debug("configured grid: %s",config.Desc())
}

func NewGrid(config *conf.GridConfig) *Grid {
    if config == nil { 
        config = conf.NewGridConfig() 
    }
    ret := &Grid{width: config.Width, height: config.Height}
    ret.buffer = gfx.NewBuffer(config.Height)
    ret.program = gfx.NewProgram("grid")
    ret.object = gfx.NewObject("grid")
    ret.texture = gfx.NewTexture("grid")
    return ret
}

func (grid *Grid) Desc() string {
    ret := fmt.Sprintf("grid[%dx%d]",grid.width,grid.height)
    if grid.buffer.Tail(0) != nil {
        ret += " '" + (*grid.buffer.Tail(0)).Desc() + "'"
    }
    return ret
}

func (grid *Grid) Dump() string {
    return grid.buffer.Dump()
}


