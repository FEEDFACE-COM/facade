
package font

import (
)


type Font struct {
    Config
    //texture data etc    
}


func NewFont() *Font {
    return &Font{Config: *NewConfig()}    
}


