
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



type Grid struct {

    lineBuffer   *LineBuffer
    termBuffer   *TermBuffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    
    state GridState
        
    refreshChan chan bool
}





const DEBUG_GRID = false


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
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.state.Width), float32(grid.state.Height) }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
    tileSize := mgl32.Vec2{ font.MaxSize().W/font.MaxSize().H, font.MaxSize().H/font.MaxSize().H }
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );
    
    tileOffset := mgl32.Vec2{-1., 0.0}
    if grid.state.Width % 2 == 0 { //even columns
        tileOffset[0] = 0.5
    }
    if grid.state.Height % 2 == 0 { //even rows
        tileOffset[1] = -0.5
    }
    if grid.state.Downward && ! grid.state.Term {
        tileOffset[1] += 1.    
    }
    grid.program.Uniform2fv(gfx.TILEOFFSET, 1, &tileOffset[0] );


    cursorPos := mgl32.Vec2{-1., -1.}
    if grid.state.Term {
        x,y := grid.termBuffer.GetCursor()
        cursorPos[0] = float32(x)
        cursorPos[1] = float32(y)
    }
    grid.program.Uniform2fv(gfx.CURSORPOS, 1, &cursorPos[0] );

    clocknow := float32( gfx.NOW() )
    grid.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow )

    scroller := float32(0.0)
    if ! grid.state.Term {
        scroller = -1. * float32( grid.lineBuffer.GetScroller() )
        if grid.state.Downward {
            scroller *= -1.
        }
    }
    grid.program.Uniform1f(gfx.SCROLLER,scroller)

    camera.Uniform(grid.program)
    
    grid.texture.Uniform(grid.program)

//    { 
//        dw := float32(0.0); 
//        if grid.state.Downward { dw = 1.0 }
//        grid.program.Uniform1f(gfx.DOWNWARD, dw)
//    }
    
    
    scale := float32( 1.0 )
    scale = grid.autoScale(camera,font)

//    var trans = float32(-0.5)
//    if ( grid.state.Downward ) {
//        trans *= -1.
//    } 

    model := mgl32.Ident4()
    model = model.Mul4( mgl32.Scale3D(scale,scale,scale) )
//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    grid.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+2+2)*4, (0)*4 )
	grid.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+2+2)*4, (3)*4 )
    grid.program.VertexAttribPointer(gfx.TILECOORD, 2, (3+2+2+2)*4, (3+2)*4 )
    grid.program.VertexAttribPointer(gfx.GRIDCOORD, 2, (3+2+2+2)*4, (3+2+2)*4 )
    
    
    count := int32(grid.state.Width*(grid.state.Height+1))
	offset := int32(0)

//    if grid.state.Downward { // need to skip first row
//        if grid.ringBuffer.Tail(0) != nil {
//            offset = int32(grid.state.Width)
//        }
//    }

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
        for r:=0; r<int(grid.state.Height+1); r++ {
            for c:=0; c<int(grid.state.Width); c++ {
                gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(1*2*3))
                off += int32(1)          
            }
//    	   gl.DrawArrays(gl.LINE_STRIP,int32(off*2*3), int32(grid.state.Width*2*3) )    
//    	   off += int32(grid.state.Width)
        }	  
        
    }
}


func (grid *Grid) Height() uint { return grid.state.Height }







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
            w,h := int(grid.state.Width), int(grid.state.Height)
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
            w,h := int(grid.state.Width), int(grid.state.Height)
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
        	h := int(grid.state.Height)
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

    th := 1./float32(gfx.GlyphRows)
    tw := glyphSize.W / ( maxGlyphSize.W * float32(gfx.GlyphCols) )

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
    if DEBUG_GRID { log.Debug("generate %s %s",grid.Desc(),font.Desc()) }
    
    width, height := int(grid.state.Width), int(grid.state.Height)
    
    for r:=0; r<=height; r++ {
        y := -1 * (r - height/2)
        
        row := r
        if grid.state.Downward && ! grid.state.Term {
            row = height - r    
        }
        
        line := Line("")
        if grid.state.Term {
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
                X: float32(glyphCoord.X) / (gfx.GlyphCols),
                Y: float32(glyphCoord.Y) / (gfx.GlyphRows),
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
    
    cols := byte(gfx.GlyphCols)

    col := chr % cols
    row := chr / cols
    return gfx.Coord{
        X: int(col),
        Y: int(row),
    }
}


func (grid *Grid) Init(camera *gfx.Camera, font *gfx.Font) {

//    grid.autoWidth(grid.state,camera,font)
    

    log.Debug("init %s",grid.Desc())
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
    err = grid.program.GetCompileShaders("grid/",grid.state.Vert,grid.state.Frag)
    if err != nil { return log.NewError("fail load grid shaders: %s",err) }
    err = grid.program.LinkProgram(); 
    if err != nil { return log.NewError("fail link grid program: %v",err) }
    return nil
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
    grid.texture.TexImage()
	
	grid.ScheduleRefresh()
    return nil
}





func (grid *Grid) autoScale(camera *gfx.Camera, font *gfx.Font) float32 {

	fontRatio := font.Ratio()
	screenRatio := camera.Ratio()
	
	ratio := screenRatio / fontRatio
	
	scaleWidth  :=  ratio * 2. / float32(grid.state.Width) 
	scaleHeight :=          2. / float32(grid.state.Height ) 
	
	if scaleWidth < scaleHeight { 
		return scaleWidth
	} else { 
		return scaleHeight
	}
	
	return float32(1.0)	
}


func (grid *Grid) autoWidth(camera *gfx.Camera, font *gfx.Font) {
	h := grid.state.Height
	var cfg = make(GridConfig)
	cfg.SetHeight(h)
	cfg.autoWidth(camera, font)
	
}

func (config *GridConfig) autoWidth(camera *gfx.Camera, font *gfx.Font) {
    if width,ok := config.Width(); !ok || width == 0 { // no width given
        height,ok := config.Height()
        if !ok { 
	        log.Debug("autowidth fail: %s",config.Desc())
	        return
	    }
        w := camera.Ratio() / font.Ratio() * float32(height)
        if height == 1 { w = 5. }
        config.SetWidth( uint(w) )
        log.Debug("autowidth %s",config.Desc())
    } 
}






func (grid *Grid) Configure(config *GridConfig, camera *gfx.Camera, font *gfx.Font) {
    if config == nil { return }


    log.Debug("config %s",config.Desc())


	config.autoWidth(camera,font)

    if width,ok := config.Width(); ok && width != 0 { 
	    grid.state.Width = width 
        grid.termBuffer.Resize(grid.state.Width,grid.state.Height)   
	} 

    if height,ok := config.Height(); ok && height != 0 && height != grid.state.Height { 
	    grid.state.Height = height 
        grid.lineBuffer.Resize(grid.state.Height,grid.state.BufLen)
        grid.termBuffer.Resize(grid.state.Width,grid.state.Height)   
    }

    if buflen,ok := config.BufLen(); ok && buflen != 0 && buflen != grid.state.BufLen {
        grid.state.BufLen = buflen
        grid.lineBuffer.Resize(grid.state.Height,grid.state.BufLen)
    }
    
    if term,ok := config.Term(); ok && term != grid.state.Term {
        grid.state.Term = term
    }

    if true {  //optimize!!
        log.Debug("rendermap %s",font.Desc())
        grid.RenderMap(font)
//        grid.texture.TexImage()
		
//		grid.empty.RenderTexture(font)
    }

	{
		if tmp,ok := config.Downward(); ok { grid.state.Downward = tmp }	
	}

//	{
//		changed := false
//	    scroll,speed := grid.state.Scroll, grid.state.Speed
//    	if tmp,ok := config.Scroll(); ok { changed = true; grid.state.Scroll = tmp }
//	    if tmp,ok := config.Speed(); ok { changed = true; grid.state.Speed = tmp }    
//    	if changed {
//	        grid.scroller.SetScrollSpeed(grid.state.Scroll, float32(grid.state.Speed))
//    	}
//    }

	{
		changed := false
	    if tmp,ok := config.Speed(); ok { changed = true; grid.state.Speed = tmp }    
    	if changed {
	        grid.lineBuffer.Speed = float32(grid.state.Speed)
    	}
    }
    
    {
	    changed := false
		vert,frag := grid.state.Vert, grid.state.Frag
		if tmp,ok := config.Vert(); ok { changed = true; grid.state.Vert = tmp }
		if tmp,ok := config.Frag(); ok { changed = true; grid.state.Frag = tmp }
		if changed {
			err := grid.LoadShaders()    
			if false && err != nil {
				grid.state.Vert = vert
				grid.state.Frag = frag
			}
		}
    }
    

    if fillName,ok := config.Fill(); ok {
        
        fillStr := grid.fill( fillName )

        grid.lineBuffer.Fill( fillStr )
        grid.termBuffer.Fill( fillStr )

    }
	
	grid.ScheduleRefresh()

}





func NewGrid(config *GridConfig, lineBuffer *LineBuffer, termBuffer *TermBuffer) *Grid {
    ret := &Grid{}
    ret.state = GridDefaults
    ret.state.ApplyConfig(config)
    ret.refreshChan = make( chan bool, 1 )
    ret.lineBuffer = lineBuffer
    ret.termBuffer = termBuffer
    ret.program = gfx.GetProgram("grid")
    ret.object = gfx.NewObject("grid")
    ret.texture = gfx.NewTexture("grid")
    return ret
}

func (grid *Grid) Desc() string { return grid.state.Desc()  }

func (grid *Grid) Dump() string { 
    if grid.state.Term {
        return grid.termBuffer.Dump() 
    } else { 
        return grid.lineBuffer.Dump(grid.state.Width)
    }
}


