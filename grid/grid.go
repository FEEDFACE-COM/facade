
package grid

import(
    "fmt"
    log "../log"
)

type Grid struct {
    Config 
    Buffer   
}



func (grid *Grid) Render() {
    //opengl stuff    
}

func (grid *Grid) Queue(text string) {
    grid.Buffer.Queue(text, 1.0)
}

func (grid *Grid) Configure(config *Config) {
    log.Debug("configure grid: %s",config.Describe())
    
    if config.Width != grid.Config.Width {
        grid.Config.Width = config.Width    
    }
    
    if config.Height != grid.Config.Height {
        grid.Config.Height = config.Height
        grid.Buffer.Configure(config)    
    }
}

func NewGrid() *Grid {
    config := NewConfig()
    buffer := NewBuffer(config.Height)
    return &Grid{Config: *config, Buffer: *buffer}
}

func (grid *Grid) Describe() string {
    ret := fmt.Sprintf("%s %s",grid.Config.Describe(),grid.Buffer.Describe())
    return ret
}


