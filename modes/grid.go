
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
    config conf.GridConfig    

    buffer *gfx.Buffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    black *gfx.Texture
    white *gfx.Texture
}



func (grid *Grid) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {
    gl.ClearColor(0,0,0,1)
//    gl.ClearColor(0.5,0.5,0.5,1.)
    gl.ActiveTexture(gl.TEXTURE0)
    

    
    grid.program.UseProgram()
    grid.object.BindBuffer()
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.config.Width),float32(grid.config.Height), }
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
    
    count := int32(grid.config.Width*grid.config.Height)

    if true {    
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, 0, count*(2*3)  )
    }

    if debug {
        gl.LineWidth(3.0)
        grid.white.BindTexture()
        for i:=0; i<int(count); i++ {
            gl.DrawArrays(gl.LINE_STRIP, int32(i * (2*3)), int32(2*3) )        
        }
    }
    
    
}






func gridVertices(size gfx.Size, glyphSize gfx.Size, tileCoord gfx.Coord, texOffset gfx.Point, maxSize gfx.Size) []float32 {
    
    w, h := size.W, size.H
    x, y := float32(tileCoord.X), float32(tileCoord.Y)
    
    tw, th := 1./float32(gfx.GlyphCols)  , 1./float32(gfx.GlyphRows)
    if true {
        tw = glyphSize.W / ( maxSize.W * float32(gfx.GlyphCols) )
    }
    
    offx, offy := texOffset.X, texOffset.Y
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



func (grid *Grid) generateData(font *gfx.Font) {
    grid.data = []float32{}
    tmp := ""
    w,h := int(grid.config.Width), int(grid.config.Height)
    for r:=0; r<h; r++ {
        y:= -1 * (r-h/2)
        
        line  := grid.buffer.Tail(uint(r))
        for c:=0; c<w; c++ {
            x:= c-w/2 + (1-w%2)
            
            chr := byte('#')
            if line != nil && int(c) < len(line.Text) {
                chr = line.Text[c]
            }    

            tileCoord := gfx.Coord{X: x, Y:y}
            glyphCoord := getGlyphCoord( byte(chr) )
            glyphSize := font.Size[glyphCoord.X][glyphCoord.Y]


            maxSize := font.MaxSize()

            size := gfx.Size{
                W: glyphSize.W / glyphSize.H,
                H: glyphSize.H / glyphSize.H,
            }
            
            texOffset := gfx.Point{
                X: float32(glyphCoord.X) / (gfx.GlyphCols),
                Y: float32(glyphCoord.Y) / (gfx.GlyphRows),
            }

            grid.data = append(grid.data, gridVertices(size,glyphSize,tileCoord,texOffset,maxSize)... )

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



    grid.texture.Init()
    grid.RenderMap(font)
    grid.texture.TexImage()
    
    grid.black = gfx.BlackColor()
    grid.white = gfx.WhiteColor()
    grid.object.Init()
    
    grid.generateData(font)

    err = grid.program.LoadShaders("grid","grid")
    if err != nil { log.Error("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { log.Error("fail link grid program: %v",err) }

    

}


func (grid *Grid) RenderMap(font *gfx.Font) error {

    rgba, err := font.RenderMapRGBA()
    if err != nil {
        log.Error("fail render font map: %s",err)
        return log.NewError("fail render font map: %s",err)
    }
    err = grid.texture.LoadRGBA(rgba)
    if err != nil {
        log.Error("fail load font map: %s",err)
        return log.NewError("fail to load font map: %s",err)
    }
    return nil
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
    
    if config.Width != grid.config.Width || config.Height != grid.config.Height {
        grid.config = *config
        grid.buffer.Resize(grid.config.Height)    
    }


    if true {
        grid.RenderMap(font)
        grid.texture.TexImage()
    }

    grid.generateData(font)

    log.Debug("configured grid: %s",config.Desc())
}

func NewGrid(config *conf.GridConfig) *Grid {
    if config == nil { 
        config = conf.NewGridConfig() 
    }
    ret := &Grid{config: *config}
    ret.buffer = gfx.NewBuffer(config.Height)
    ret.program = gfx.NewProgram("grid")
    ret.object = gfx.NewObject("grid")
    ret.texture = gfx.NewTexture("grid")
    return ret
}

func (grid *Grid) Desc() string {
    ret := grid.config.Desc()
    if grid.buffer.Tail(0) != nil {
        ret += " '" + (*grid.buffer.Tail(0)).Desc() + "'"
    }
    return ret
}

func (grid *Grid) Dump() string {
    return grid.buffer.Dump()
}


