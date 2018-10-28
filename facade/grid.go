
// +build linux,arm

package facade

import(
    "fmt"
//    "math"
    gfx "../gfx"
    log "../log"
//    math "../math32"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)


type Grid struct {
    config GridConfig    

    buffer *gfx.Buffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    scroller *gfx.Scroller
    
    refreshChan chan bool
}



const DEBUG_GRID = false

func (grid *Grid) MarkDirty() {
    
}

func (grid *Grid) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {
    
    select {
    
        case refresh := <- grid.refreshChan:
            if refresh {
//                log.Debug("refresh %s",grid.Desc())
                grid.generateData(font)
            }

        default:
            break    
    }
    
    gl.ActiveTexture(gl.TEXTURE0)
    
    grid.program.UseProgram(debug)
    grid.object.BindBuffer()
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.config.Width),float32(grid.config.Height), }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
    tileSize := mgl32.Vec2{ font.MaxSize().W/font.MaxSize().H, font.MaxSize().H/font.MaxSize().H }
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );

    grid.program.Uniform1f(gfx.CLOCKNOW, gfx.NOW() )
        
    grid.scroller.Uniform(grid.program, grid.config.Downward)
    camera.Uniform(grid.program)
    grid.program.Uniform1i(gfx.TEXTURE,0)

    { 
        dw := float32(0.0); 
        if grid.config.Downward { dw = 1.0 }
        grid.program.Uniform1f(gfx.DOWNWARD, dw)
    }
    
    
    scale := float32( 1.0 )
    const autoScale = true //scale model to fit screen
    if autoScale {

        fontRatio := font.Ratio()
        screenRatio := camera.Ratio()
        
        ratio := screenRatio / fontRatio
        
        scaleWidth :=  ratio * 2. / float32(grid.config.Width) 
        scaleHeight :=      2. / float32(grid.config.Height - 1) //minus one line below
        
        if scaleWidth < scaleHeight { 
            scale = scaleWidth
        } else { 
            scale = scaleHeight
        }
                
    }

    trans := float32(-0.5)
    if ( grid.config.Downward ) {
        trans *= -1.
    } 

    model := mgl32.Ident4()
    model = model.Mul4( mgl32.Scale3D(scale,scale,scale) )
    model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    grid.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+2)*4, (0)*4 )
    grid.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+2)*4, (3)*4 )
    grid.program.VertexAttribPointer(gfx.TILECOORD, 2, (3+2+2)*4, (3+2)*4 )
    
    
    count := int32(grid.config.Width*grid.config.Height)

    if !debug {    
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, 0, count*(2*3)  )
    }

    if debug {
        gl.LineWidth(3.0)
        
        w,h := int(grid.config.Width), int(grid.config.Height)
        for r:=0; r<h; r++ {
            var line *gfx.Text
            if grid.config.Downward {
                line  = grid.buffer.Tail(uint(r))
            } else {
                line  = grid.buffer.Head(uint(r))
            }        
            for c:=0; c<w; c++ {
                if line != nil && int(c) < len(line.Text) /**/ && line.Text[c] != byte(' ') /**/ {
                    gl.DrawArrays(gl.LINE_STRIP, int32((r*w+c) * (2*3)), int32(2*3) )
                }
            }
        }
        
        
    }
//    if verbose { log.Debug( grid.scroller.Desc() ) }
}



func (grid *Grid) Fill(font *gfx.Font) {
    
    fill := grid.config.Fill
    
    switch fill {
    
        case "title1": 
            for _,line := range []string{
                "                    ",
                " _  _   _  _   _   _",
                "|_ |_| /  |_| | \\ |_", 
                "|  | | \\_ | | |_/ |_",
                "                    ",
                "     by FEEDFACE.COM",
                "                    ",
            } {
                grid.Queue(line,font)    
            }
            
        case "title2": 
            for _,line := range []string{
                "              ",
                "F A C A D E   ",
                "              ",
                "            by",
                "  FEEDFACE.COM",
                "              ",
            } {
                grid.Queue(line,font)    
            }
            
        case "title3": 
            for _,line := range []string{
                "F A C A D E",
            } {
                grid.Queue(line,font)    
            }
        
        
        case "coord":
            w,h := int(grid.config.Width), int(grid.config.Height)
            for r:=0; r<h; r++ {
                line := ""
                for c:=0; c<w; c++ {
                    d := "."
                    if c % 5 == 0 { d = fmt.Sprintf("%d",r%10) }
                    if r % 5 == 0 { d = fmt.Sprintf("%d",c%10) }
                    if c % 5 == 0 && r % 5 == 0 { d = "#" }
    
                    line += fmt.Sprintf("%s",d)        
                }
                grid.Queue(line,font)
            }
            
            
        case "alpha":
            w,h := int(grid.config.Width), int(grid.config.Height)
            alpha := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`~!@#$^&*()-_=+[{]}|;:',<.>/?"
            s := 0
            for r:=0; r<h; r++ {
                line := alpha[s%len(alpha):(s+w)%len(alpha)]
                grid.Queue(line,font)
                s += 1
            }

    }    
}



func gridVertices(size gfx.Size, glyphSize gfx.Size, tileCoord gfx.Coord, texOffset gfx.Point, maxSize gfx.Size) []float32 {
    
    w, h := size.W, size.H
    x, y := float32(tileCoord.X), float32(tileCoord.Y)
    ox, oy := texOffset.X, texOffset.Y

    th := 1./float32(gfx.GlyphRows)
    tw := glyphSize.W / ( maxSize.W * float32(gfx.GlyphCols) )
    
    return []float32{
        //vertex                     //texcoords        // tile coords
        -w/2,  h/2, 0,                 0+ox,  0+oy,      x, y,
        -w/2, -h/2, 0,                 0+ox, th+oy,      x, y,
         w/2, -h/2, 0,                tw+ox, th+oy,      x, y,
         w/2, -h/2, 0,                tw+ox, th+oy,      x, y,
         w/2,  h/2, 0,                tw+ox,  0+oy,      x, y,
        -w/2,  h/2, 0,                 0+ox,  0+oy,      x, y,
        
    }
    
}

func min(a,b float32) float32 { 
    if a<b { return a } 
    return b 
}


const DEBUG_DATA = false

func (grid *Grid) generateData(font *gfx.Font) {
    grid.data = []float32{}
    tmp := fmt.Sprintf("generate %s %s",grid.Desc(),font.Desc())
    w,h := int(grid.config.Width), int(grid.config.Height)
    for r:=0; r<h; r++ {
        tmp += "\n"
        y:= -1 * (r-h/2)

        var line *gfx.Text
        if grid.config.Downward {
            line  = grid.buffer.Tail(uint(r))
        } else {
            line  = grid.buffer.Head(uint(r))
        }        
        
        for c:=0; c<w; c++ {
            x:= c-w/2 + (1-w%2)
            
            chr := byte(' ')
            if DEBUG_GRID { chr = byte('#') }
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

            tmp += fmt.Sprintf("%+d/%+d    ",x,y)
//            tmp += fmt.Sprintf(" %.0fx%0.f    ",float32(glyphSize.W),float32(glyphSize.H))
        } 
    }
    if DEBUG_DATA { log.Debug(tmp) }
    grid.object.BufferData( len(grid.data)*4,grid.data )
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

    grid.config = grid.autoWidth(grid.config,camera,font)
    

    log.Debug("init %s",grid.Desc())
    grid.texture.Init()
    grid.RenderMap(font)
    grid.texture.TexImage()
    

    grid.object.Init()
    grid.LoadShaders()
        
    select { case grid.refreshChan <- true: ; default: ; }

}

func (grid *Grid) LoadShaders() {
    var err error
    err = grid.program.GetCompileShaders("grid/",grid.config.Vert,grid.config.Frag)
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
    grid.scroller.Once()
    select { case grid.refreshChan <- true: ; default: ; }
    
}



func (grid *Grid) autoWidth(config GridConfig, camera *gfx.Camera, font *gfx.Font) GridConfig {
    ret := config
    if config.Width == 0 {
        w := camera.Ratio() / font.Ratio() * float32(config.Height-1)  //minus one for line below waiting to scroll in
        if config.Height == 1 { w = 5. }
        ret.Width = uint(w)
        log.Debug("autowidth %s -> %s",config.Desc(),ret.Desc())
    } 
    if ret.Width == 0 { ret.Width = 5 }    
    return ret
}



func (grid *Grid) Configure(config *GridConfig, camera *gfx.Camera, font *gfx.Font) {
    if config == nil { return }

    autoConf := grid.autoWidth(*config,camera,font)
    config = &autoConf

    if *config == grid.config { return }
    

    
    log.Debug("config %s -> %s",grid.Desc(),config.Desc())
    old := grid.config
    grid.config = *config

    if config.Scroll != old.Scroll || config.Speed != old.Speed {
        grid.scroller.SetScrollSpeed(config.Scroll, float32(config.Speed))
    }
        
    if config.Height != old.Height {
        grid.buffer.Resize(grid.config.Height)    
        log.Debug(grid.buffer.Dump())
    }

    if true { // REM, optimize, only rerender if font changed?
        grid.RenderMap(font)
        grid.texture.TexImage()
    }
    
    if config.Vert != old.Vert || config.Frag != old.Frag {
        grid.LoadShaders()    
    }


    if config.Fill != "" {
        grid.Fill(font)
    }

    select { case grid.refreshChan <- true: ; default: ; }

}

func NewGrid(config *GridConfig) *Grid {
    if config == nil { 
        config = NewGridConfig() 
    }
    ret := &Grid{config: *config}
    ret.refreshChan = make( chan bool, 1 )
    ret.buffer = gfx.NewBuffer(config.Height)
    ret.program = gfx.GetProgram("grid")
    ret.object = gfx.NewObject("grid")
    ret.texture = gfx.NewTexture("grid")
    ret.scroller = gfx.NewScroller(config.Scroll,float32(config.Speed))
    return ret
}

func (grid *Grid) Desc() string {
    tmp := ""
    if grid.scroller != nil { tmp = " " + grid.scroller.Desc() }
//    tmp += " " + grid.program.Desc()
    return grid.config.Desc() + tmp
}

func (grid *Grid) Dump() string {
    return grid.buffer.Dump()
}


