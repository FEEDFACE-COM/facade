
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



func (grid *Grid) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {
//    gl.ClearColor(0,0,0,1)
    gl.ClearColor(0.5,0.5,0.5,1.)
    gl.ActiveTexture(gl.TEXTURE0)
    

    
    grid.program.UseProgram()
    grid.object.BindBuffer()
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.width),float32(grid.height), }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
//    tileSize := mgl32.Vec2{ 34./70., 70./70. }
    tileSize := mgl32.Vec2{ font.MaxSize().W/font.MaxSize().H, font.MaxSize().H/font.MaxSize().H }
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
func gridVertices(size gfx.Size, glyphSize gfx.Size, tileCoord gfx.Coord, texOffset gfx.Point) []float32 {
    
    w, h := size.W, size.H
    x, y := float32(tileCoord.X), float32(tileCoord.Y)
//    tw, th := /*size.W * */float32(1./(gfx.GlyphCols)),  /*size.H * */float32(1./(gfx.GlyphRows))

    r := glyphSize.W / glyphSize.H
    r = float32(1.)
    tw, th := 1./float32(gfx.GlyphCols) * r , 1./float32(gfx.GlyphRows)
    
    offx, offy := texOffset.X, texOffset.Y
    left := float32(0)
    return []float32{
            //vertex                       //texcoords        // coordinates
        -w/2 + left,  h/2, 0,                0 +offx,  0 + offy,      x, y,    
        -w/2 + left, -h/2, 0,                0 +offx, th + offy,      x, y,    
         w/2 + left, -h/2, 0,               tw +offx, th + offy,      x, y,    
         w/2 + left, -h/2, 0,               tw +offx, th + offy,      x, y,    
         w/2 + left,  h/2, 0,               tw +offx,  0 + offy,      x, y,    
        -w/2 + left,  h/2, 0,                0 +offx,  0 + offy,      x, y,    
        
    }
    
}



func (grid *Grid) generateData(font *gfx.Font) {
    grid.data = []float32{}
    tmp := ""
    w,h := int(grid.width), int(grid.height)
    for r:=0; r<h; r++ {
        y:= -1 * (r-h/2)
        
        line  := grid.buffer.Tail(uint(r))
//        totalWidth := float32(0)
        for c:=0; c<w; c++ {
            x:= c-w/2 + (1-w%2)
            
            chr := byte('#')
            if line != nil && int(c) < len(line.Text) {
                chr = line.Text[c]
            }    
            tileCoord := gfx.Coord{X: x, Y:y}


            
            glyphCoord := getGlyphCoord( byte(chr) )
            
            glyphSize := font.Size[glyphCoord.X][glyphCoord.Y]


            

            size := gfx.Size{
//                W: 34./70. ,
//                H: 70./70. ,
                W: font.MaxSize().W / font.MaxSize().H ,
                H: font.MaxSize().H / font.MaxSize().H ,
            }

            
            texOffset := gfx.Point{
                X: float32(glyphCoord.X) / (gfx.GlyphCols),
                Y: float32(glyphCoord.Y) / (gfx.GlyphRows),
            }
            grid.data = append(grid.data, gridVertices(size,glyphSize,tileCoord,texOffset)... )
//            totalWidth += float32(font.Size[r][c].W)
//            tmp += fmt.Sprintf("%+d/%+d %.0fx%0.f    ",x,y,float32(glyphSize.W),float32(glyphSize.H))
            tmp += fmt.Sprintf("%+d/%+d %.0fx%0.f    ",x,y,float32(glyphSize.W),float32(glyphSize.H))
        } 
        tmp += "\n"
    }
    log.Debug(tmp)
    grid.object.BufferData(len(grid.data)*4,grid.data)
}



func getGlyphCoord(chr byte) gfx.Coord {
    cols := byte(gfx.GlyphCols)

    col := chr % cols
    row := chr / cols
    return gfx.Coord{
        X: int(col),
        Y: int(row),
    }
}


func (grid *Grid) Init(camera *gfx.Camera, font *gfx.Font) {
    var err error
    log.Debug("create %s",grid.Desc())


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
    
    grid.generateData(font)

    err = grid.program.LoadShaders("grid","grid")
    if err != nil { log.Error("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { log.Error("fail link grid program: %v",err) }

    

}




func (grid *Grid) Queue(text string, font *gfx.Font) {
    newText := gfx.NewText(text)
    grid.buffer.Queue( newText )
    grid.generateData(font)
    log.Debug("queued text: %s",text)
}







func (grid *Grid) Configure(config *conf.GridConfig, font *gfx.Font) {
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

    grid.generateData(font)

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


