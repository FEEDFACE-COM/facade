
package grid

import(
    "fmt"
    conf "../conf"
    log "../log"
)

type Grid struct {
    width uint
    height uint
    *Buffer   
}



func (grid *Grid) Render() {
    //opengl stuff    
}

func (grid *Grid) Queue(text string) {
    grid.Buffer.Queue(text, 1.0)
}

func (grid *Grid) Configure(config *conf.GridConfig) {
    log.Debug("configure grid: %s",config.Describe())
    
    if config.Width != grid.width {
        grid.width = config.Width    
    }
    
    if config.Height != grid.height {
        grid.height = config.Height
        grid.Buffer.Configure(config)    
    }
}

func NewGrid() *Grid {
    ret := &Grid{}
    ret.Configure(conf.NewGridConfig())
    ret.Buffer = NewBuffer(ret.height)
    return ret
}

func (grid *Grid) Describe() string {
    ret := fmt.Sprintf("grid[%dx%d] %s",grid.width,grid.height,grid.Buffer.Describe())
    return ret
}


