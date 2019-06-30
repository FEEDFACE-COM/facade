
// +build linux,arm

package facade

import(
    "fmt"
    "strings"
//    "math"
    gfx "../gfx"
    log "../log"
//    math "../math32"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)

const DEBUG_GRID = true


type Grid struct {

    width, height uint

    downward bool
    
    buffer uint
    terminal bool


    vert,frag string

    lineBuffer   *LineBuffer
    termBuffer   *TermBuffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    
        
    refreshChan chan bool
}







func (grid *Grid) ScheduleRefresh() {

    select { case grid.refreshChan <- true: ; default: ; }
	
}


func (grid *Grid) checkRefresh() bool {
	ret := false
	for { //read all messages from channel
		select {
			case refresh := <- grid.refreshChan:
				if refresh {
					ret = true
				}

			default:
				return ret
		}
	}
	return ret
}

func (grid *Grid) Render(camera *gfx.Camera, font *gfx.Font, debug, verbose bool) {


    // REM maybe also if grid.checkReconfig then grid.Configure??

    
    if grid.checkRefresh() {

	    grid.GenerateData(font)
	    
	}
    
    
    gl.ActiveTexture(gl.TEXTURE0)
    
    grid.program.UseProgram(debug)
    grid.object.BindBuffer()
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.width), float32(grid.height) }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
    tileSize := mgl32.Vec2{ font.MaxSize().W/font.MaxSize().H, font.MaxSize().H/font.MaxSize().H }
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );
    
    tileOffset := mgl32.Vec2{-1., 0.0}
    if grid.width % 2 == 0 { //even columns
        tileOffset[0] = 0.5
    }
    if grid.height % 2 == 0 { //even rows
        tileOffset[1] = -0.5
    }
    if grid.downward && ! grid.terminal {
        tileOffset[1] += 1.    
    }
    grid.program.Uniform2fv(gfx.TILEOFFSET, 1, &tileOffset[0] );


    cursorPos := mgl32.Vec2{-1., -1.}
    if grid.terminal {
        x,y := grid.termBuffer.GetCursor()
        cursorPos[0] = float32(x)
        cursorPos[1] = float32(y)
    }
    grid.program.Uniform2fv(gfx.CURSORPOS, 1, &cursorPos[0] );

    clocknow := float32( gfx.NOW() )
    grid.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow )

    scroller := float32(0.0)
    if ! grid.terminal {
        scroller = float32( grid.lineBuffer.GetScroller() )
        if grid.downward {
            scroller *= -1.
        }
    }
    grid.program.Uniform1f(gfx.SCROLLER,scroller)

    camera.Uniform(grid.program)
    
    grid.texture.Uniform(grid.program)

    
    
    scale := float32( 1.0 )
    scale = grid.autoScale(camera,font)


    model := mgl32.Ident4()
    model = model.Mul4( mgl32.Scale3D(scale,scale,scale) )
//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    grid.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+2+2)*4, (0)*4 )
	grid.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+2+2)*4, (3)*4 )
    grid.program.VertexAttribPointer(gfx.TILECOORD, 2, (3+2+2+2)*4, (3+2)*4 )
    grid.program.VertexAttribPointer(gfx.GRIDCOORD, 2, (3+2+2+2)*4, (3+2+2)*4 )
    
    
    count := int32(grid.width*(grid.height+1))
	offset := int32(0)


    if !debug || debug {    
	    grid.program.SetDebug(false)
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, int32(offset * 2*3), (count)*(2*3)  )
	    grid.program.SetDebug(debug)
    }

    if debug {
        gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
        off := offset
        // REM, use single gl.DrawElements call instead (create indice array before)
        for r:=0; r<int(grid.height+1); r++ {
            for c:=0; c<int(grid.width); c++ {
                gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(1*2*3))
                off += int32(1)          
            }
//    	   gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(grid.Width*2*3) )    
//    	   off += int32(grid.Width)
        }	  
        
    }
}


//func (grid *Grid) Height() uint { return grid.Height }







func (grid *Grid) fill(name string) []string {
    
    switch name {
    
    	//todo: cheeck widht, switch different titles
    	//also, clear!

    	
        case "title":
            return strings.Split(`
 _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _      
|_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|
|   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |
`,           
            "\n")[1:]

    	
    	
        case "title2": 
            return strings.Split(`
 _  _   _  _   _   _
|_ |_| /  |_| | \ |_
|  | | \_ | | |_/ |_
                    
     by FEEDFACE.COM
`,           
            "\n")[1:]

            
            
        case "title3": 
            return strings.Split(`
              
F A C A D E   
              
            by
  FEEDFACE.COM
              
`,           
             "\n")[1:]

            
        case "title4": 
            return []string{
                "F A C A D E",
            }
        
        
        
        case "grid":
            ret := []string{}
            w,h := int(grid.width), int(grid.height)
            for r:=0; r<h; r++ {
                tmp := ""
                for c:=0; c<w; c++ {
                    d := "."
                    if c % 5 == 0 { d = fmt.Sprintf("%d",r%10) }
                    if r % 5 == 0 { d = fmt.Sprintf("%d",c%10) }
                    if c % 5 == 0 && r % 5 == 0 { d = "#" }
                    tmp += fmt.Sprintf("%s",d)
                }
                ret = append(ret, tmp )
            }
            return ret
            
            
        case "alpha":
            ret := []string{}
            w,h := int(grid.width), int(grid.height)
            alpha := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`~!@#$^&*()-_=+[{]}|;:',<.>/?"
            s := 0
            for r:=0; r<h; r++ {
                tmp := alpha[ s%len(alpha) : min(s+w,len(alpha)-1) ]
                ret = append(ret, tmp )
                s += 1
            }
            return ret
            
            
        case "clear":
            ret := []string{}
        	h := int(grid.height)
        	for r:=0; r<h; r++ {
	        	ret = append(ret, "" )
	        }
	        return ret

    }    

    return []string{}

}


func min(a,b int) int { if a < b { return a; }; return b; }


func gridVertices(
    tileSize     gfx.Size,  // dimensions of the tile
    glyphSize    gfx.Size,  // dimensions of the glyph
    maxGlyphSize gfx.Size,  // max dimension of glyph in font
    gridCoord    gfx.Coord, // row/col coord of tile
    tilePos      gfx.Point, // position of tile in grid 
    texOffset    gfx.Point, // glyph offset in texture
) []float32 {
    
    col,row := float32(gridCoord.X), float32(gridCoord.Y)
    
    w, h := tileSize.W, tileSize.H
    x, y := tilePos.X, tilePos.Y
    ox, oy := texOffset.X, texOffset.Y

    th := 1./float32(gfx.GlyphMapRows)
    tw := glyphSize.W / ( maxGlyphSize.W * float32(gfx.GlyphMapCols) )

    return []float32{
        //vertex            //texcoords        // tile coords     // grid coords
        -w/2,  h/2, 0,        0+ox,  0+oy,      x, y,             col,row,
        -w/2, -h/2, 0,        0+ox, th+oy,      x, y,             col,row,
         w/2, -h/2, 0,       tw+ox, th+oy,      x, y,             col,row,
         w/2, -h/2, 0,       tw+ox, th+oy,      x, y,             col,row,
         w/2,  h/2, 0,       tw+ox,  0+oy,      x, y,             col,row,
        -w/2,  h/2, 0,        0+ox,  0+oy,      x, y,             col,row,
        
    }
    
}




func (grid *Grid) GenerateData(font *gfx.Font) {
    grid.data = []float32{}
    
//    if DEBUG_GRID { log.Debug("%s generate vector data %s",grid.Desc(),font.Desc()) }
    
    width, height := int(grid.width), int(grid.height)
    
    for r:=0; r<=height; r++ {
        y := -1 * (r - height/2)
        
        row := r
        if grid.downward && ! grid.terminal {
            row = height - r    
        }
        
        line := Line("")
        if grid.terminal {
            line = grid.termBuffer.GetLine( uint(row) )    
        } else {
            line = grid.lineBuffer.GetLine( uint(row) )
        }
        
        for col:=0; col<width; col++ {
            x := col - width/2 + (width%2)
            run := rune(' ')
            if col < len(line) {
                run = line[col]
            }
            
            gridCoord := gfx.Coord{col,row}
            tilePos := gfx.Point{ float32(x), float32(y) }
            glyphCoord := getGlyphCoord( run )
            glyphSize := font.Size[glyphCoord.X][glyphCoord.Y]
            maxGlyphSize := font.MaxSize() 

            tileSize := gfx.Size{
                W: glyphSize.W / glyphSize.H,
                H: glyphSize.H / glyphSize.H,
            }
            
            texOffset := gfx.Point{
                X: float32(glyphCoord.X) / (gfx.GlyphMapCols),
                Y: float32(glyphCoord.Y) / (gfx.GlyphMapRows),
            }

            grid.data = append(grid.data, gridVertices(tileSize,glyphSize,maxGlyphSize,gridCoord,tilePos,texOffset)... )

            
        }

        
    }
    grid.object.BufferData( len(grid.data)*4,grid.data )
    
}



func getGlyphCoord(run rune) gfx.Coord {
    if run <= 0x20 || run >= 0x80 {
        return gfx.Coord{X: 0, Y: 0}    
    }
    chr := byte(run)
    
    cols := byte(gfx.GlyphMapCols)

    col := chr % cols
    row := chr / cols
    return gfx.Coord{
        X: int(col),
        Y: int(row),
    }
}


func (grid *Grid) Init(shaderService *gfx.ShaderService, font *gfx.Font) {
    log.Debug("%s init",grid.Desc())


    grid.program = gfx.NewProgram("grid",shaderService)
    grid.object = gfx.NewObject("grid")
    grid.texture = gfx.NewTexture("grid")


    grid.texture.Init()
    grid.RenderMap(font)
//    grid.texture.TexImage()
    

    grid.object.Init()
    grid.LoadShaders()

//	grid.empty.RenderTexture(font)
	grid.ScheduleRefresh()        

}

func (grid *Grid) LoadShaders() error {
	var err error
	if DEBUG_GRID { log.Debug("%s load shaders",grid.Desc()) }
    err = grid.program.GetCompileShaders("grid/",grid.vert,grid.frag)
    if err != nil { return log.NewError("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { return log.NewError("fail link grid program: %v",err) }
    return nil
}


func (grid *Grid) RenderMap(font *gfx.Font) error {

//    if DEBUG_GRID { log.Debug("%s  render map with %s",grid.Desc(),font.Desc()) }

    rgba, err := font.RenderMap(false)
    if err != nil {
        log.Error("fail render font map: %s",err)
        return log.NewError("fail render font map: %s",err)
    }
    err = grid.texture.LoadRGBA(rgba)
    if err != nil {
        log.Error("fail load font map: %s",err)
        return log.NewError("fail to load font map: %s",err)
    }
    grid.texture.TexImage()
	
	grid.ScheduleRefresh()
    return nil
}





func (grid *Grid) autoScale(camera *gfx.Camera, font *gfx.Font) float32 {

	fontRatio := font.Ratio()
	screenRatio := camera.Ratio()
	
	ratio := screenRatio / fontRatio
	
	scaleWidth  :=  ratio * 2. / float32(grid.width) 
	scaleHeight :=          2. / float32(grid.height ) 
	
	if scaleWidth < scaleHeight { 
		return scaleWidth
	} else { 
		return scaleHeight
	}
	
	return float32(1.0)	
}


//func (grid *Grid) autoWidth(camera *gfx.Camera, font *gfx.Font) {
//	h := grid.Height
//	var cfg = make(GridConfig)
//	cfg.SetHeight(h)
//	cfg.autoWidth(camera, font)
//	
//}







func (grid *Grid) Configure(config *GridConfig, camera *gfx.Camera, font *gfx.Font) {
    if config == nil { return }


    log.Debug("%s configure %s",grid.Desc(),config.Desc())


	config.autoWidth(camera.Ratio(),font.Ratio())

    {
        changed := false
        if config.GetSetWidth() && config.GetWidth() != 0 && uint(config.GetWidth()) != grid.width {
            grid.width = uint( config.GetWidth() )
            changed = true
        } 
    
        if config.GetSetHeight() && config.GetHeight() != 0 && uint(config.GetHeight()) != grid.height { 
            grid.height = uint( config.GetHeight() )
            changed = true
        }
    
        if config.GetSetBuffer() && uint(config.GetBuffer()) != grid.buffer {
            grid.buffer = uint( config.GetBuffer() )
            changed = true
        }
        
        if changed {
            grid.lineBuffer.Resize(grid.height,grid.buffer)
            grid.termBuffer.Resize(grid.width,grid.height)   
        }
    }
        
    if config.GetSetTerminal() {
        grid.terminal = config.GetTerminal()
    }

    if true {  //optimize!!
        log.Debug("%s rendermap %s",grid.Desc(),font.Desc())
        grid.RenderMap(font)
//        grid.texture.TexImage()
    }

    {
        if config.GetSetDownward() { 
            grid.downward = config.GetDownward() 
        }	
	    if config.GetSetSpeed() { 
            grid.lineBuffer.SetSpeed( float32(config.GetSpeed() ) )
    	}
    	
    	if config.GetSetAdaptive() {
            grid.lineBuffer.Adaptive = config.GetAdaptive()
        }
        if config.GetSetDrop() {
            grid.lineBuffer.Drop = config.GetDrop()
        }

        if config.GetSetSmooth() {
            grid.lineBuffer.Smooth = config.GetSmooth()
        }


    }

    
    {
	    changed := false
		vert,frag := grid.vert, grid.frag
		if config.GetSetVert() { changed = true; grid.vert = config.GetVert() }
		if config.GetSetFrag() { changed = true; grid.frag = config.GetFrag() }
		if changed {
			err := grid.LoadShaders()    
			if err != nil {
				grid.vert = vert
				grid.frag = frag
			}
		}
    }
    

    if config.GetSetFill() {
        fillStr := grid.fill( config.GetFill() )
        grid.lineBuffer.Fill( fillStr )
        grid.termBuffer.Fill( fillStr )

    }
	
	grid.ScheduleRefresh()

}

func (grid *Grid) DescBuffer() string {
    if grid.terminal {
        return grid.termBuffer.Desc()    
    } else {
        return grid.lineBuffer.Desc()
    }    
}
    

func (grid *Grid) DumpBuffer() string {
    if grid.terminal {
        return grid.termBuffer.Dump()    
    } else {
        return grid.lineBuffer.Dump(grid.width)
    }    
}


func NewGrid(lineBuffer *LineBuffer, termBuffer *TermBuffer) *Grid {
    ret:= &Grid{}    
    ret.width =  uint(GridDefaults.GetWidth())
    ret.height = uint(GridDefaults.GetHeight())

    ret.downward = GridDefaults.GetDownward()

    ret.terminal = GridDefaults.GetTerminal()
    ret.buffer = uint(GridDefaults.GetBuffer())

    ret.vert  = GridDefaults.GetVert()
    ret.frag  = GridDefaults.GetFrag()
    
    ret.refreshChan = make( chan bool, 1 )
    ret.lineBuffer = lineBuffer
    ret.termBuffer = termBuffer
    
    return ret
}

func (grid *Grid) Desc() string { 
    return grid.Config().Desc()
}

func (grid *Grid) Config() *GridConfig {
    return &GridConfig{
        SetWidth: true,   Width: uint64(grid.width),
        SetHeight: true, Height: uint64(grid.height),
        
        SetSpeed: true, Speed: float64(grid.lineBuffer.Speed()),
        SetDownward: true, Downward: grid.downward,
        SetAdaptive: true, Adaptive: grid.lineBuffer.Adaptive,
        
        SetTerminal: true, Terminal: grid.terminal,
        SetBuffer: true, Buffer: uint64(grid.buffer),
        
        SetVert: true, Vert: grid.vert,
        SetFrag: true, Frag: grid.frag,
        
    }
}


