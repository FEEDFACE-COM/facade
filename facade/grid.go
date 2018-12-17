
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

    buffer *gfx.Buffer
    
    texture *gfx.Texture
    program *gfx.Program
    object *gfx.Object
    data []float32
    
    scroller *gfx.Scroller
    
    state GridState
    
    empty *gfx.Text
        
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

    
    if grid.checkRefresh() {

	    grid.generateData(font)
	    
	}
    
    
    gl.ActiveTexture(gl.TEXTURE0)
    
    grid.program.UseProgram(debug)
    grid.object.BindBuffer()
    
    
    
    tileCount := mgl32.Vec2{ float32(grid.state.Width), float32(grid.state.Height) }
    grid.program.Uniform2fv(gfx.TILECOUNT, 1, &tileCount[0] );
    
    tileSize := mgl32.Vec2{ font.MaxSize().W/font.MaxSize().H, font.MaxSize().H/font.MaxSize().H }
    grid.program.Uniform2fv(gfx.TILESIZE, 1, &tileSize[0] );

    clocknow := float32( gfx.NOW() )
    grid.program.Uniform1fv(gfx.CLOCKNOW, 1, &clocknow )
        
    grid.scroller.Uniform(grid.program, grid.state.Downward)
    camera.Uniform(grid.program)
    
    grid.texture.Uniform(grid.program)

    { 
        dw := float32(0.0); 
        if grid.state.Downward { dw = 1.0 }
        grid.program.Uniform1f(gfx.DOWNWARD, dw)
    }
    
    
    scale := float32( 1.0 )
    scale = grid.autoScale(camera,font)

    var trans = float32(-0.5)
    if ( grid.state.Downward ) {
        trans *= -1.
    } 

    model := mgl32.Ident4()
    model = model.Mul4( mgl32.Scale3D(scale,scale,scale) )
//	model = model.Mul4( mgl32.Translate3D(0.0,trans,0.0) )
    grid.program.UniformMatrix4fv(gfx.MODEL, 1, &model[0] )
    
    grid.program.VertexAttribPointer(gfx.VERTEX,    3, (3+2+2)*4, (0)*4 )
	grid.program.VertexAttribPointer(gfx.TEXCOORD,  2, (3+2+2)*4, (3)*4 )
    grid.program.VertexAttribPointer(gfx.TILECOORD, 2, (3+2+2)*4, (3+2)*4 )
    
    
    count := int32(grid.state.Width*(grid.state.Height+1))
	offset := 0
	if grid.state.Downward { //need to skip first row
		offset = len( grid.buffer.Tail(0).Text ) 
	}
	

    if !debug || debug {    
	    off := offset
	    grid.program.SetDebug(false)
        grid.texture.BindTexture()
        gl.DrawArrays(gl.TRIANGLES, int32(off * 2*3), count*(2*3)  )
	    grid.program.SetDebug(debug)
    }

    if debug {
        gl.LineWidth(3.0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
        w,h := int(grid.state.Width), int(grid.state.Height)
        from := 0
        to := h+1
        if ! grid.state.Downward {
	    	from = -1;
	    	to = h;    
	    }
        for r:=from; r<to; r++ {
            var line = grid.lineForRow(r)
            for c:=0; c<w; c++ {

				if line == nil { break }
				if int(c) >= len(line.Text) { break }

				
				off := r*w + c
				off += w
				gl.DrawArrays(gl.LINE_STRIP, int32(off * 2*3), int32(2*3) )

            }
        }
        
        
    }
//    if verbose { log.Debug( grid.scroller.Desc() ) }
}


func (grid *Grid) Height() uint { return grid.state.Height }



func (grid *Grid) lineForRow(row int) *gfx.Text {
	r := uint(row)
	
//	if r >= grid.state.Height { 
//		r= grid.state.Height-1 
//	}

	if r>= grid.state.Height {
		return grid.empty
	}

	if r < 0 {
		return grid.empty
	}

	if grid.state.Downward {
		return grid.buffer.Tail(r)
	} else {
		return grid.buffer.Head(r)
	}	
	return grid.empty
}


func (grid *Grid) Fill(font *gfx.Font, fill string) {
    
    switch fill {
    
    	//todo: cheeck widht, switch different titles
    	//also, clear!
        case "title": 
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
        
        
        case "grid":
            w,h := int(grid.state.Width), int(grid.state.Height)
//            for r:=0; r<h-1; r++ {
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
            w,h := int(grid.state.Width), int(grid.state.Height)
            alpha := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`~!@#$^&*()-_=+[{]}|;:',<.>/?"
            s := 0
            for r:=0; r<h; r++ {
                line := alpha[ s%len(alpha) : min(s+w,len(alpha)-1) ]
                grid.Queue(line,font)
                s += 1
            }
            
            
        case "clear":
        	h := int(grid.state.Height)
        	for r:=0; r<h; r++ {
	        	grid.Queue("",font)
	        }

    }    
}


func min(a,b int) int { if a < b { return a; }; return b; }


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



const DEBUG_DATA = false

func (grid *Grid) generateData(font *gfx.Font) {
    grid.data = []float32{}
    tmp := fmt.Sprintf("generate %s %s",grid.Desc(),font.Desc())
    w,h := int(grid.state.Width), int(grid.state.Height+1)
    for r:=-1; r<h; r++ {
        tmp += "\n"
        y:= -1 * (r-h/2)

        var line = grid.lineForRow(r)
        for c:=0; c<w; c++ {
//            x:= c-w/2 + (1-w%2)
            x:= c-w/2 + (w%2)
            
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

//    grid.autoWidth(grid.state,camera,font)
    

    log.Debug("init %s",grid.Desc())
    grid.texture.Init()
    grid.RenderMap(font)
//    grid.texture.TexImage()
    

    grid.object.Init()
    grid.LoadShaders()

	grid.empty.RenderTexture(font)
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






func (grid *Grid) Queue(text string, font *gfx.Font) {
    newText := gfx.NewText(text)
    fun := func() { 
	    grid.empty = gfx.NewText("") 
	    grid.ScheduleRefresh()
	    log.Debug("empty funned: %s",grid.Desc())
	}
    if grid.scroller.Once(fun) {
	    log.Debug("empty primed: %s",grid.Desc())
		tmp := grid.buffer.Head(0)
//	    if grid.state.Downward {
//		    tmp = grid.buffer.Tail(0)
//		}
	    if tmp == nil {
			grid.empty = gfx.NewText("")
		} else {
			grid.empty = gfx.NewText( tmp.Text )
		}
	} 
	grid.buffer.Queue( newText )
	grid.ScheduleRefresh()
    
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

    if width,ok := config.Width(); ok { 
	    grid.state.Width = width 
	} 

    if height,ok := config.Height(); ok && height != grid.state.Height { 
	    grid.state.Height = height 
        log.Debug("resize %s",grid.buffer.Desc())
        grid.buffer.Resize(height)    
    }

    if true {  //optimize!!
        log.Debug("rendermap %s",font.Desc())
        grid.RenderMap(font)
//        grid.texture.TexImage()
		
		grid.empty.RenderTexture(font)
    }

	{
		if tmp,ok := config.Downward(); ok { grid.state.Downward = tmp }	
	}

	{
		changed := false
//	    scroll,speed := grid.state.Scroll, grid.state.Speed
    	if tmp,ok := config.Scroll(); ok { changed = true; grid.state.Scroll = tmp }
	    if tmp,ok := config.Speed(); ok { changed = true; grid.state.Speed = tmp }    
    	if changed {
	        grid.scroller.SetScrollSpeed(grid.state.Scroll, float32(grid.state.Speed))
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
    

    if fill,ok := config.Fill(); ok {
        grid.Fill(font,fill)
    }
	
	grid.ScheduleRefresh()

}





func NewGrid(config *GridConfig) *Grid {
    ret := &Grid{}
    ret.state = GridDefaults
    ret.state.ApplyConfig(config)
    ret.refreshChan = make( chan bool, 1 )
    ret.buffer = gfx.NewBuffer(ret.state.Height)
    ret.program = gfx.GetProgram("grid")
    ret.object = gfx.NewObject("grid")
    ret.texture = gfx.NewTexture("grid")
    ret.scroller = gfx.NewScroller(ret.state.Scroll,float32(ret.state.Speed))
    ret.empty = gfx.NewText("")
    return ret
}

func (grid *Grid) Desc() string { return grid.state.Desc()  }
func (grid *Grid) Dump() string { return grid.buffer.Dump() }


