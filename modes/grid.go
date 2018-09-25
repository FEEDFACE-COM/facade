
package modes

import(
    "fmt"
	"github.com/go-gl/mathgl/mgl32"    
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)

type Coord struct {
    x int
    y int    
}

type Size struct {
    w float32
    h float32  
}

type Point struct {
    x float32
    y float32    
}

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
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.width),float32(grid.height), }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
    tileSize := mgl32.Vec2{1.0,1.0}
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );
    

    model := mgl32.Ident4()
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    camera.Uniform(grid.program)
    grid.program.Uniform1i(gfx.TEXTURE,0)
    
    grid.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+2)*4,  0*4)
    grid.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+2)*4, (3)*4)
    grid.program.VertexAttribPointer(gfx.TILECOORD, 2, (3+2+2)*4, (3+2)*4)
    
    count := (2*3)*int32(grid.height*grid.width)

    if true {    
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, 0, count  )
    }

    if debug {
        gl.LineWidth(3.0)
        grid.white.BindTexture()
        for i:=0; i<int(grid.height*grid.width); i++ {
            gl.DrawArrays(gl.LINE_STRIP, int32(i * (2*3)), int32(2*3) )        
        }
    }
    
    
}






//func gridVertices(size,tileCoord,texCoord,texOffset struct{x,y float32}) []float32 {
func gridVertices(size Size, tileCoord Coord, texOffset Point) []float32 {
    
    w, h := size.w, size.h
    x, y := float32(tileCoord.x), float32(tileCoord.y)
    tw, th := float32(1./(gfx.GlyphCols)),  float32(1./(gfx.GlyphRows))
    
    offx, offy := texOffset.x, texOffset.y
    return []float32{
            //vertex                       //texcoords        // coordinates
        -w/2,  h/2, 0,                0 +offx,  0 + offy,      x, y,    
        -w/2, -h/2, 0,                0 +offx, th + offy,      x, y,    
         w/2, -h/2, 0,               tw +offx, th + offy,      x, y,    
         w/2, -h/2, 0,               tw +offx, th + offy,      x, y,    
         w/2,  h/2, 0,               tw +offx,  0 + offy,      x, y,    
        -w/2,  h/2, 0,                0 +offx,  0 + offy,      x, y,    
        
    }
    
}



func (grid *Grid) generateData() {
    grid.data = []float32{}
    
    w,h := int(grid.width), int(grid.height)
    for r:=0; r<h; r++ {
        y:= -1 * (r-h/2)
        
        line  := grid.buffer.Tail(uint(r))
        
        for c:=0; c<w; c++ {
            x:= c-w/2 + (1-w%2)
            
            chr := byte('#')
            if line != nil && int(c) < len(line.Text) {
                chr = line.Text[c]
            }    
            tileCoord := Coord{x: x, y:y}

            size := Size{w: 1., h: 1. }
            
            glyphCoord := getGlyphCoord( byte(chr) )
            texOffset := Point{
                x: float32(glyphCoord.x) / (gfx.GlyphCols),
                y: float32(glyphCoord.y) / (gfx.GlyphRows),
            }
            grid.data = append(grid.data, gridVertices(size,tileCoord,texOffset)... )
        } 
    }
            
    grid.object.BufferData(len(grid.data)*4,grid.data)
}



func getGlyphCoord(chr byte) Coord {
    cols := byte(gfx.GlyphCols)

    col := chr % cols
    row := chr / cols
    return Coord{
        x: int(col),
        y: int(row),
    }
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


