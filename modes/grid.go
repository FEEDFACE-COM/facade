
package modes

import(
    "fmt"
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
    
    white *gfx.Texture
}



func (grid *Grid) Render(camera *gfx.Camera, debug bool) {
    gl.ClearColor(0.5,0.5,0.5,1.0)
    
    
    gl.ActiveTexture(gl.TEXTURE0)
    grid.program.UseProgram()
    grid.object.BindBuffer()
    camera.Uniform(grid.program)
    
    grid.program.VertexAttribPointer(gfx.VERTEX,3,5*4,0)
    grid.program.VertexAttribPointer(gfx.TEXCOORD,2,5*4,3*4)
    
    
    grid.texture.BindTexture()
    gl.DrawArrays(gl.TRIANGLES, 0, 2*3)
    
    
    if true {
        gl.LineWidth(3.0)
        grid.white.BindTexture()
        gl.DrawArrays(gl.LINE_STRIP, 0, 2*3)        
    }
    
    
}







func (grid *Grid) Queue(text string) {
    newText := gfx.NewText(text)
    grid.buffer.Queue( newText )
    grid.generateData()
    log.Debug("queued text: %s",text)
}














func (grid *Grid) generateData() {
    grid.data = []float32{}
    
    
    grid.data = append(grid.data, gfx.QuadVertices(1.,1.)...)
    grid.object.BufferData(len(grid.data)*4,grid.data)
    
}

func (grid *Grid) Init(camera *gfx.Camera) {
    var err error
    log.Debug("create %s",grid.Desc())


    err = grid.texture.LoadFile("/home/folkert/src/gfx/facade/asset/test.png")
    if err != nil {
        log.Error("fail load grid file")
    }
    grid.white = gfx.WhiteColor()
    grid.object.Init()
    
    grid.generateData()

    err = grid.program.LoadShaders("grid","grid")
    if err != nil { log.Error("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { log.Error("fail link grid program: %v",err) }

    

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
    if config == nil { config = conf.NewGridConfig() }
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


