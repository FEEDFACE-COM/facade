
package grid

import(
    "fmt"
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


func NewGrid() *Grid {
    config := NewConfig()
    buffer := NewBuffer(config.Height)
    return &Grid{Config: *config, Buffer: *buffer}
}

func (grid *Grid) Describe() string {
    ret := fmt.Sprintf("%s %s",grid.Config.Describe(),grid.Buffer.Describe())
    return ret
}