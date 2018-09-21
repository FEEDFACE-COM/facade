
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
    buffer Buffer
}

type gridItem struct {
    text string
}

func (item *gridItem) Desc() string { return item.text }

func (item *gridItem) Close() { log.Debug("close grid[%s]",item.text) }

func (item *gridItem) Bind(uint32)  { return  }

func (grid *Grid) Render(camera *gfx.Camera) {
    gl.ClearColor(0xff,0x0,0x0,1.0)
}

func (grid *Grid) Queue(text string) {
    newItem := gridItem{text: text}
    grid.buffer.Queue( &newItem )
}


func (grid *Grid) Init(camera *gfx.Camera) {
    
}

func (grid *Grid) Configure(config *conf.GridConfig) {
    if config == nil {
        return
    }
    log.Debug("configure grid: %s",config.Desc())
    
    if config.Width != grid.width {
        grid.width = config.Width    
    }
    
    if config.Height != grid.height {
        grid.height = config.Height
        grid.buffer.Resize(config.Height)    
    }
}

func NewGrid(config *conf.GridConfig) *Grid {
    if config == nil {
        config = conf.NewGridConfig()
    }
    ret := &Grid{width: config.Width, height: config.Height}
    ret.buffer = NewBuffer(config.Height)
    return ret
}

func (grid *Grid) Desc() string {
    ret := fmt.Sprintf("grid[%dx%d]",grid.width,grid.height)
    ret += "\n>> "
    for i:=uint(0);i<grid.height;i++ {
        item := grid.buffer.Item(i)
        if item != nil { ret += (*item).Desc() }
        ret += ","
    }
    ret += "\n<< "
    for i:=uint(grid.height);i>0;i-- {
        item := grid.buffer.Item(i-1)
        if item != nil { ret += (*item).Desc() }
        ret += ","
    }
    ret += "\n"
    return ret
}

func (grid *Grid) Dump() string {
    return grid.buffer.Dump()
}


