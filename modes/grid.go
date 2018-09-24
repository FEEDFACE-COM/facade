
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
    white *gfx.Texture
}



func (grid *Grid) Render(camera *gfx.Camera, debug, verbose bool) {
    gl.ClearColor(0,0,0,1)
    gl.ActiveTexture(gl.TEXTURE0)
    

    
    grid.program.UseProgram()
    grid.object.BindBuffer()
    
    
    
    gridSize := mgl32.Vec2{ 
        float32(grid.width)/2.  - 0.5 ,
        float32(grid.height)/2. - 0.5 ,
    }
    grid.program.Uniform2fv(gfx.GRIDSIZE, 1, &gridSize[0] );
    
    tileSize := mgl32.Vec2{
        1.0,
        1.0,    
    }
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );

    model := mgl32.Ident4()
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    camera.Uniform(grid.program)
    grid.program.Uniform1i(gfx.TEXTURE,0)
    
    grid.program.VertexAttribPointer(gfx.VERTEX,3,(3+2+2+2)*4,0)
    grid.program.VertexAttribPointer(gfx.GRIDCOORD,2,(3+2+2+2)*4,(3)*4)
    grid.program.VertexAttribPointer(gfx.TEXCOORD,2,(3+2+2+2)*4,(3+2)*4)
    grid.program.VertexAttribPointer(gfx.TEXOFFSET,2,(3+2+2+2)*4,(3+2+2)*4)
    
    count := (2*3)*int32(grid.height*grid.width)

    if true {    
        if (verbose) { log.Debug("draw triangles") }
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, 0, count  )
    }
    

    if debug {
        if (verbose) { log.Debug("draw lines") }
        gl.LineWidth(3.0)
        grid.white.BindTexture()
        for i:=0; i<int(grid.height*grid.width); i++ {
            gl.DrawArrays(gl.LINE_STRIP, int32(i * (2*3)), int32(2*3) )        
        }
    }
    
    
}










func gridVertices(size,gridcoord,tilesize,offset mgl32.Vec2) []float32 {
    w, h := size[0], size[1]
    x, y := gridcoord[0], gridcoord[1]
    tw, th := tilesize[0], tilesize[1]
    offx,offy := offset[0],offset[1]
    return []float32{
        -w/2,  h/2, 0,            x, y,        0,  0,       offx, offy,
        -w/2, -h/2, 0,            x, y,        0, th,       offx, offy,
         w/2, -h/2, 0,            x, y,       tw, th,       offx, offy,
         w/2, -h/2, 0,            x, y,       tw, th,       offx, offy,
         w/2,  h/2, 0,            x, y,       tw,  0,       offx, offy,
        -w/2,  h/2, 0,            x, y,        0,  0,       offx, offy,
        
    }
    
}


func getOffset(chr byte) mgl32.Vec2 {
    cols := byte(gfx.GlyphCols)
//    rows:= byte(gfx.GlyphRows)

    col := chr % cols
    row := chr / cols
    return mgl32.Vec2{
        float32(col),
        float32(row),
        
    }
}


func (grid *Grid) generateData() {
    log.Debug("grid generate data") 
    grid.data = []float32{}
    
    
    for y := uint(0); y<grid.height; y++ {
        line  := grid.buffer.Tail(y)

        for x:=uint(0); x<grid.width; x++ {
            chr := byte(' ')
            if line != nil && int(x) < len(line.Text) {
                chr = line.Text[x]
            }    

            size := mgl32.Vec2{1, 1}
            gridcoord := mgl32.Vec2{float32(x),float32(y)}
            tilesize := mgl32.Vec2{  1, 1 }
            offset := getOffset( byte(chr) )
            grid.data = append(grid.data, gridVertices(size,gridcoord,tilesize,offset)... )

        }
            
        
    }
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
    grid.white = gfx.WhiteColor()
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

    grid.generateData()

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


